import * as cdk from "@aws-cdk/core";
import * as ssm from "@aws-cdk/aws-ssm";
import { StringParameter } from "@aws-cdk/aws-ssm";

export const APP_NAME = "pickr";

export const SSM_BASE_PATH = `/${APP_NAME}/`;

export class SharedParametersStack extends cdk.Stack {
  public sessionCookie: StringParameter;

  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    this.sessionCookie = new ssm.StringParameter(this, "session-cookie", {
      stringValue: `${APP_NAME}-session`,
    });
  }
}
