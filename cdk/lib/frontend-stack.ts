import * as cdk from "@aws-cdk/core";
import * as s3 from "@aws-cdk/aws-s3";
import * as s3deploy from "@aws-cdk/aws-s3-deployment";
import * as cloudfront from "@aws-cdk/aws-cloudfront";

export class FrontendStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const svelteBucket = new s3.Bucket(this, "svelte-bucket", {
      websiteIndexDocument: "index.html",
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
        originConfigs: [
          {
            s3OriginSource: {
              s3BucketSource: svelteBucket,
              originAccessIdentity: cloudfrontOAI,
            },
            behaviors: [{ isDefaultBehavior: true }],
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
                "pnpm run build",
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
