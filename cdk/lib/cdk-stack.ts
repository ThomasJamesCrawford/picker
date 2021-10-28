import * as cdk from "@aws-cdk/core";
import * as lambda from "@aws-cdk/aws-lambda-go";
import * as dynamodb from "@aws-cdk/aws-dynamodb";
import * as s3 from "@aws-cdk/aws-s3";
import * as s3deploy from "@aws-cdk/aws-s3-deployment";
import * as cloudfront from "@aws-cdk/aws-cloudfront";
import { HttpApi, HttpMethod } from "@aws-cdk/aws-apigatewayv2";
import { LambdaProxyIntegration } from "@aws-cdk/aws-apigatewayv2-integrations";

export class CdkStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const table = new dynamodb.Table(this, "picker-table", {
      partitionKey: { name: "PK", type: dynamodb.AttributeType.STRING },
      sortKey: { name: "SK", type: dynamodb.AttributeType.STRING },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
    });

    const fatLambda = new lambda.GoFunction(this, "handler", {
      entry: "../backend/go/cmd/fat-lambda",
      environment: {
        table: table.tableName,
        region: this.region,
        GIN_MODE: "release",
      },
    });

    table.grantReadWriteData(fatLambda);

    const httpApi = new HttpApi(this, "api-gateway", {});

    httpApi.addRoutes({
      path: "/{proxy+}",
      methods: [HttpMethod.ANY],
      integration: new LambdaProxyIntegration({
        handler: fatLambda,
      }),
    });

    const svelteBucket = new s3.Bucket(this, "svelte-bucket", {
      websiteIndexDocument: "index.html",
    });

    /**
     * Grant CloudFront access to S3 without making S3 public
     */
    const cloudfrontOAI = new cloudfront.OriginAccessIdentity(this, "oai");
    svelteBucket.grantRead(cloudfrontOAI.grantPrincipal);

    const distribution = new cloudfront.CloudFrontWebDistribution(
      this,
      "distribution",
      {
        originConfigs: [
          {
            s3OriginSource: {
              s3BucketSource: svelteBucket,
              originAccessIdentity: cloudfrontOAI,
            },
            behaviors: [{ isDefaultBehavior: true }],
          },
        ],
      }
    );

    new s3deploy.BucketDeployment(this, "static-svelte-website-deployment", {
      distribution,
      distributionPaths: ["/index.html"],
      sources: [
        s3deploy.Source.asset("../frontend", {
          bundling: {
            image: cdk.DockerImage.fromBuild("../frontend", {
              file: "Dockerfile.pnpm",
            }),
            command: [
              "bash",
              "-c",
              [
                "pnpm install",
                `VITE_API_URL=${httpApi.apiEndpoint} pnpm run build`,
                "cp -r /asset-input/build/* /asset-output/",
              ].join(" && "),
            ],
          },
        }),
      ],
      destinationBucket: svelteBucket,
    });

    new cdk.CfnOutput(this, "cloudfrontDistribution", {
      value: `https://${distribution.domainName}`,
    });

    new cdk.CfnOutput(this, "httpGateway", {
      value: httpApi.apiEndpoint,
    });
  }
}
