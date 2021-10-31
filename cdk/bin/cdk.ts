#!/usr/bin/env node
import "source-map-support/register";
import * as cdk from "@aws-cdk/core";
import { BackendStack } from "../lib/backend-stack";
import { FrontendStack } from "../lib/frontend-stack";
import { SharedParametersStack } from "../lib/shared-parameters";

const app = new cdk.App();

const sharedParameters = new SharedParametersStack(app, "shared-parameters");

const backend = new BackendStack(app, "backend", {
  sessionCookie: sharedParameters.sessionCookie,
});

backend.addDependency(sharedParameters);

const frontend = new FrontendStack(app, "frontend", {
  httpApi: backend.httpApi,
  sessionCookie: sharedParameters.sessionCookie,
});

frontend.addDependency(backend);
frontend.addDependency(sharedParameters);
