import * as cdk from "@aws-cdk/core";
import * as lambda from "@aws-cdk/aws-lambda-go";
import * as dynamodb from "@aws-cdk/aws-dynamodb";
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

    const roomLambda = new lambda.GoFunction(this, "handler", {
      entry: "../backend/go/cmd/room",
      environment: {
        table: table.tableName,
        region: this.region,
      },
    });

    table.grantReadWriteData(roomLambda);

    const httpApi = new HttpApi(this, "api-gateway", {});

    httpApi.addRoutes({
      path: "/room/{id}",
      methods: [HttpMethod.ANY],
      integration: new LambdaProxyIntegration({
        handler: roomLambda,
      }),
    });
  }
}
