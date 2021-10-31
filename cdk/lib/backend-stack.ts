import * as cdk from "@aws-cdk/core";
import * as lambda from "@aws-cdk/aws-lambda-go";
import * as dynamodb from "@aws-cdk/aws-dynamodb";
import * as iam from "@aws-cdk/aws-iam";
import { HttpApi, HttpMethod } from "@aws-cdk/aws-apigatewayv2";
import { LambdaProxyIntegration } from "@aws-cdk/aws-apigatewayv2-integrations";
import { SSM_BASE_PATH } from "./shared-parameters";
import { StringParameter } from "@aws-cdk/aws-ssm";

interface BackendStackProps extends cdk.StackProps {
  sessionCookie: StringParameter;
}

export class BackendStack extends cdk.Stack {
  public httpApi: HttpApi;

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
        session_cookie: props.sessionCookie.stringValue,
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

    /**
     * This is behind a CloudFront distribution
     *
     * TODO Needs an authorizer to only accept requests from CloudFront
     */
    const httpApi = new HttpApi(this, "api-gateway");

    httpApi.addRoutes({
      path: "/api/{proxy+}",
      methods: [HttpMethod.ANY],
      integration: new LambdaProxyIntegration({
        handler: fatLambda,
      }),
    });

    this.httpApi = httpApi;

    new cdk.CfnOutput(this, "httpGateway", {
      value: httpApi.apiEndpoint,
    });
  }
}
