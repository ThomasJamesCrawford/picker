import * as cdk from "@aws-cdk/core";
import * as lambda from "@aws-cdk/aws-lambda-go";

export class CdkStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    new lambda.GoFunction(this, "handler", {
      entry: "../../backend/go/cmd/",
    });
  }
}
