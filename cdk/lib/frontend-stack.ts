import * as cdk from "@aws-cdk/core";
import * as s3 from "@aws-cdk/aws-s3";
import * as s3deploy from "@aws-cdk/aws-s3-deployment";
import * as cloudfront from "@aws-cdk/aws-cloudfront";
import { HttpApi } from "@aws-cdk/aws-apigatewayv2";
import { StringParameter } from "@aws-cdk/aws-ssm";
import { PriceClass } from "@aws-cdk/aws-cloudfront";
import { APP_NAME } from "./shared-parameters";
import { BlockPublicAccess } from "@aws-cdk/aws-s3";

interface FrontendStackProps extends cdk.StackProps {
  httpApi: HttpApi;
  sessionCookie: StringParameter;
}

export class FrontendStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props: FrontendStackProps) {
    super(scope, id, props);

    const svelteBucket = new s3.Bucket(this, "svelte-bucket", {
      websiteIndexDocument: "index.html",
      websiteErrorDocument: "index.html",
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      blockPublicAccess: BlockPublicAccess.BLOCK_ALL,
    });

    /**
     * Grant CloudFront access to S3 without making S3 public
     */
    const cloudfrontOAI = new cloudfront.OriginAccessIdentity(this, "oai");
    svelteBucket.grantRead(cloudfrontOAI.grantPrincipal);

    const distribution = new cloudfront.CloudFrontWebDistribution(
      this,
      "distribution",
      {
        /**
         * SPA routing will 404, needs to be handled client side
         *
         * API Gateway will get hit by this too
         */
        errorConfigurations: [
          {
            errorCode: 404,
            responseCode: 200,
            responsePagePath: "/index.html",
          },
        ],
        priceClass: PriceClass.PRICE_CLASS_ALL,
        originConfigs: [
          // Serve the S3 bucket
          {
            s3OriginSource: {
              s3BucketSource: svelteBucket,
              originAccessIdentity: cloudfrontOAI,
            },
            behaviors: [{ isDefaultBehavior: true }],
          },
          // Proxy the API Gateway
          {
            behaviors: [
              {
                pathPattern: "/api/*",
                allowedMethods: cloudfront.CloudFrontAllowedMethods.ALL,
                /**
                 * Block CloudFront from caching any requests
                 *
                 * maxTtl could be set higher to allow cache-control headers set from applications to cause CloudFront caching
                 */
                maxTtl: cdk.Duration.seconds(0),
                minTtl: cdk.Duration.seconds(0),
                defaultTtl: cdk.Duration.seconds(0),
                forwardedValues: {
                  queryString: true,
                  cookies: {
                    forward: "whitelist",
                    whitelistedNames: [props.sessionCookie.stringValue],
                  },
                },
              },
            ],
            customOriginSource: {
              domainName: cdk.Fn.parseDomainName(props.httpApi.apiEndpoint),
            },
          },
        ],
      }
    );

    new s3deploy.BucketDeployment(this, "static-svelte-website-deployment", {
      /**
       * Invalidate the index.html on every deploy
       */
      distribution,
      distributionPaths: ["/index.html"],
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
                `pnpm install`,
                `VITE_API_URL=/api VITE_APP_NAME=${APP_NAME} pnpm run build`,
                "cp -r /asset-input/build/* /asset-output/",
              ].join(" && "),
            ],
          },
        }),
      ],
      destinationBucket: svelteBucket,
    });

    new cdk.CfnOutput(this, "cloudfrontDistribution", {
      value: `https://${distribution.distributionDomainName}`,
    });
  }
}
