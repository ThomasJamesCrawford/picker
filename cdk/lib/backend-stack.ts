import * as cdk from "@aws-cdk/core";
import * as lambda from "@aws-cdk/aws-lambda-go";
import * as dynamodb from "@aws-cdk/aws-dynamodb";
import * as iam from "@aws-cdk/aws-iam";
import { CorsHttpMethod, HttpApi, HttpMethod } from "@aws-cdk/aws-apigatewayv2";
import { LambdaProxyIntegration } from "@aws-cdk/aws-apigatewayv2-integrations";
import { CloudFrontWebDistribution } from "@aws-cdk/aws-cloudfront";

interface BackendStackProps extends cdk.StackProps {
  distribution: CloudFrontWebDistribution;
}

const SSM_BASE_PATH = "/picker/";

export class BackendStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props: BackendStackProps) {
    super(scope, id, props);

    const table = new dynamodb.Table(this, "picker-table", {
      partitionKey: { name: "PK", type: dynamodb.AttributeType.STRING },
      sortKey: { name: "SK", type: dynamodb.AttributeType.STRING },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
    });

    table.addGlobalSecondaryIndex({
      partitionKey: { name: "GSI1PK", type: dynamodb.AttributeType.STRING },
      sortKey: { name: "GSI1SK", type: dynamodb.AttributeType.STRING },
      indexName: "GSI1",
    });

    const fatLambda = new lambda.GoFunction(this, "handler", {
      entry: "../backend/go/cmd/fat-lambda",
      environment: {
        table: table.tableName,
        region: this.region,
        GIN_MODE: "release",
        ssm_path: SSM_BASE_PATH,
      },
    });

    fatLambda.role?.attachInlinePolicy(
      new iam.Policy(this, "lambda-ssm", {
        statements: [
          new iam.PolicyStatement({
            effect: iam.Effect.ALLOW,
            actions: ["ssm:GetParametersByPath"],
            resources: [
              `arn:aws:ssm:${this.region}:${this.account}:parameter${SSM_BASE_PATH}`,
            ],
          }),
        ],
      })
    );

    table.grantReadWriteData(fatLambda);

    const httpApi = new HttpApi(this, "api-gateway", {
      corsPreflight: {
        allowHeaders: ["*"],
        allowMethods: [CorsHttpMethod.ANY],
        allowOrigins: [
          "http://localhost:3000",
          `https://${props.distribution.distributionDomainName}`,
        ],
        maxAge: cdk.Duration.minutes(10),
        allowCredentials: true,
      },
    });

    httpApi.addRoutes({
      path: "/{proxy+}",
      methods: [HttpMethod.ANY],
      integration: new LambdaProxyIntegration({
        handler: fatLambda,
      }),
    });

    new cdk.CfnOutput(this, "httpGateway", {
      value: httpApi.apiEndpoint,
      description: "This needs to go into the frontend .env",
    });
  }
}
