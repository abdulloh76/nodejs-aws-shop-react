package main

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfront"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfrontorigins"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3deployment"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type S3CloudFrontStackProps struct {
	awscdk.StackProps
}

func NewS3CloudFrontStack(scope constructs.Construct, id string, props *S3CloudFrontStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	shopBucket := awss3.NewBucket(stack, jsii.String("abdulloh76-aws-shop-react"), &awss3.BucketProps{
		RemovalPolicy:     awscdk.RemovalPolicy_DESTROY,
		BlockPublicAccess: awss3.BlockPublicAccess_BLOCK_ALL(),
		AutoDeleteObjects: jsii.Bool(true),
	})

	oai := awscloudfront.NewOriginAccessIdentity(stack, jsii.String("OAI_NEW"), &awscloudfront.OriginAccessIdentityProps{
		Comment: jsii.String(fmt.Sprintf("OAI for %s", *shopBucket.BucketName())),
	})

	distribution := awscloudfront.NewDistribution(stack, jsii.String("aws-shop-cloudfront"), &awscloudfront.DistributionProps{
		DefaultBehavior: &awscloudfront.BehaviorOptions{
			Origin: awscloudfrontorigins.NewS3Origin(shopBucket, &awscloudfrontorigins.S3OriginProps{
				OriginAccessIdentity: oai,
			}),
			ViewerProtocolPolicy: awscloudfront.ViewerProtocolPolicy_REDIRECT_TO_HTTPS,
		},
		DefaultRootObject: jsii.String("index.html"),
	})

	shopBucket.AddToResourcePolicy(
		awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
			Actions:   jsii.Strings("s3:GetObject"),
			Resources: jsii.Strings(fmt.Sprintf("%s/*", *shopBucket.BucketArn())),
			Principals: &[]awsiam.IPrincipal{
				awsiam.NewCanonicalUserPrincipal(oai.CloudFrontOriginAccessIdentityS3CanonicalUserId()),
			},
		}),
	)

	awss3deployment.NewBucketDeployment(
		stack,
		jsii.String("DeploymentWithInvalidation"),
		&awss3deployment.BucketDeploymentProps{
			Sources: &[]awss3deployment.ISource{
				awss3deployment.Source_Asset(jsii.String("../dist"), &awss3assets.AssetOptions{}),
			},
			DestinationBucket: shopBucket,
			Distribution:      distribution,
			DistributionPaths: jsii.Strings("/*"),
		})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewS3CloudFrontStack(app, "S3CloudFrontStack", &S3CloudFrontStackProps{})

	app.Synth(nil)
}
