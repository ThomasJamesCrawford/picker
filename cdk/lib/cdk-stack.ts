import * as cdk from "@aws-cdk/core";
import * as lambda from "@aws-cdk/aws-lambda-go";
import * as dynamodb from "@aws-cdk/aws-dynamodb";
import * as s3 from "@aws-cdk/aws-s3";
import * as s3deploy from "@aws-cdk/aws-s3-deployment";
import * as cloudfront from "@aws-cdk/aws-cloudfront";
import * as origins from "@aws-cdk/aws-cloudfront-origins";
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
      publicReadAccess: true,
    });

    new s3deploy.BucketDeployment(this, "static-svelte-website-deployment", {
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
                "pnpm run build",
                "cp -r /asset-input/build/* /asset-output/",
              ].join(" && "),
            ],
          },
        }),
      ],
      destinationBucket: svelteBucket,
    });

    const distribution = new cloudfront.Distribution(this, "distribution", {
      defaultBehavior: { origin: new origins.S3Origin(svelteBucket) },
    });

    new cdk.CfnOutput(this, "cloudfrontDistribution", {
      value: distribution.domainName,
    });

    new cdk.CfnOutput(this, "httpGateway", {
      value: httpApi.apiEndpoint,
    });
  }
}
