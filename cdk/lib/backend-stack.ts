import * as cdk from "@aws-cdk/core";
import * as lambda from "@aws-cdk/aws-lambda-go";
import * as dynamodb from "@aws-cdk/aws-dynamodb";
import { CorsHttpMethod, HttpApi, HttpMethod } from "@aws-cdk/aws-apigatewayv2";
import { LambdaProxyIntegration } from "@aws-cdk/aws-apigatewayv2-integrations";

export class BackendStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
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
      },
    });

    table.grantReadWriteData(fatLambda);

    const httpApi = new HttpApi(this, "api-gateway", {
      corsPreflight: {
        allowMethods: [
          CorsHttpMethod.GET,
          CorsHttpMethod.HEAD,
          CorsHttpMethod.OPTIONS,
          CorsHttpMethod.POST,
        ],
        allowOrigins: ["http://localhost:3000"],
        maxAge: cdk.Duration.minutes(10),
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
