#!/usr/bin/env node
import "source-map-support/register";
import * as cdk from "@aws-cdk/core";
import { BackendStack } from "../lib/backend-stack";
import { FrontendStack } from "../lib/frontend-stack";

const app = new cdk.App();

const frontend = new FrontendStack(app, "frontend");

new BackendStack(app, "backend", {
  distribution: frontend.distribution,
}).addDependency(frontend);
