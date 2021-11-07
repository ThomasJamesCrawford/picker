#!/usr/bin/env node
import "source-map-support/register";
import * as cdk from "@aws-cdk/core";
import { BackendStack } from "../lib/backend-stack";
import { FrontendStack } from "../lib/frontend-stack";
import { SharedParametersStack } from "../lib/shared-parameters";

const app = new cdk.App();

const env = {
  account: process.env.CDK_DEFAULT_ACCOUNT,
  region: process.env.CDK_DEFAULT_REGION,
};

const sharedParameters = new SharedParametersStack(app, "shared-parameters", {
  env,
});

const backend = new BackendStack(app, "backend", {
  sessionCookie: sharedParameters.sessionCookie,
  env,
});

new FrontendStack(app, "frontend", {
  httpApi: backend.httpApi,
  sessionCookie: sharedParameters.sessionCookie,
  env,
});
