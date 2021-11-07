import * as cdk from "@aws-cdk/core";
import * as s3 from "@aws-cdk/aws-s3";
import * as s3deploy from "@aws-cdk/aws-s3-deployment";
import * as cloudfront from "@aws-cdk/aws-cloudfront";
import * as route53 from "@aws-cdk/aws-route53";
import * as targets from "@aws-cdk/aws-route53-targets";
import * as acm from "@aws-cdk/aws-certificatemanager";
import { HttpApi } from "@aws-cdk/aws-apigatewayv2";
import { StringParameter } from "@aws-cdk/aws-ssm";
import { PriceClass, ViewerCertificate } from "@aws-cdk/aws-cloudfront";
import { APP_NAME, DOMAIN_NAME } from "./shared-parameters";
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

    const hostedZone = route53.PublicHostedZone.fromLookup(this, "hostedZone", {
      domainName: DOMAIN_NAME,
    });

    const certificate = new acm.DnsValidatedCertificate(
      this,
      "acmCertificate",
      {
        domainName: DOMAIN_NAME,
        hostedZone,
        region: "us-east-1",
      }
    );

    const distribution = new cloudfront.CloudFrontWebDistribution(
      this,
      "distribution",
      {
        viewerCertificate: ViewerCertificate.fromAcmCertificate(certificate),
        /**
         * SPA routing will 404, needs to be handled client side
         *
         * API Gateway will get hit by this too (probably doesn't matter)
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

    new route53.AaaaRecord(this, "cloudfrontDistributionAlias", {
      zone: hostedZone,
      target: route53.RecordTarget.fromAlias(
        new targets.CloudFrontTarget(distribution)
      ),
    });

    new route53.ARecord(this, "cloudfrontDistributionAliasIPv4", {
      zone: hostedZone,
      target: route53.RecordTarget.fromAlias(
        new targets.CloudFrontTarget(distribution)
      ),
    });

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
