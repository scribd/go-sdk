package aws

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sagemakerruntime"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type (
	Builder struct {
		config *Config
	}
)

func NewBuilder(c *Config) *Builder {
	return &Builder{config: c}
}

func (b *Builder) LoadConfig(
	ctx context.Context, opts ...func(options *awscfg.LoadOptions) error) (aws.Config, error) {
	defaultOpts := []func(options *awscfg.LoadOptions) error{
		awscfg.WithRegion(b.config.AWSConfig.Region),
		awscfg.WithHTTPClient(createHttpClient(&b.config.AWSConfig.HTTPClient)),
	}

	opts = append(defaultOpts, opts...)

	return awscfg.LoadDefaultConfig(ctx, opts...)
}

func (b *Builder) NewS3Service(
	cfg aws.Config, serviceName string, opts ...func(options *s3.Options)) (*s3.Client, error) {
	s3Cfg, ok := b.config.S3[serviceName]
	if !ok {
		return nil, fmt.Errorf("s3 config for service %s not found", serviceName)
	}
	defaultOpts := []func(*s3.Options){
		func(options *s3.Options) {
			if s3Cfg.Region != "" {
				options.Region = s3Cfg.Region
			}
			if credentialsSet(&s3Cfg.Credentials) {
				options.Credentials = getCredentialsProvider(cfg, &s3Cfg.Credentials)
			}
			if s3Cfg.HTTPClient.MaxIdleConns > 0 {
				options.HTTPClient = createHttpClient(&s3Cfg.HTTPClient)
			}
		},
	}

	opts = append(defaultOpts, opts...)

	return s3.NewFromConfig(cfg, opts...), nil
}

func (b *Builder) NewSagemakerRuntimeService(
	cfg aws.Config,
	serviceName string, opts ...func(options *sagemakerruntime.Options)) (*sagemakerruntime.Client, error) {
	sagemakerRuntimeCfg, ok := b.config.SagemakerRuntime[serviceName]
	if !ok {
		return nil, fmt.Errorf("sagemaker runtime config for service %s not found", serviceName)
	}
	defaultOpts := []func(*sagemakerruntime.Options){
		func(options *sagemakerruntime.Options) {
			if sagemakerRuntimeCfg.Region != "" {
				options.Region = sagemakerRuntimeCfg.Region
			}
			if credentialsSet(&sagemakerRuntimeCfg.Credentials) {
				options.Credentials = getCredentialsProvider(cfg, &sagemakerRuntimeCfg.Credentials)
			}
			if sagemakerRuntimeCfg.HTTPClient.MaxIdleConns > 0 {
				options.HTTPClient = createHttpClient(&sagemakerRuntimeCfg.HTTPClient)
			}
		},
	}

	opts = append(defaultOpts, opts...)

	return sagemakerruntime.NewFromConfig(cfg, opts...), nil
}

func (b *Builder) NewSFNService(
	cfg aws.Config,
	serviceName string, opts ...func(options *sfn.Options)) (*sfn.Client, error) {
	sfnCfg, ok := b.config.SFN[serviceName]
	if !ok {
		return nil, fmt.Errorf("sfn config for service %s not found", serviceName)
	}
	defaultOpts := []func(*sfn.Options){
		func(options *sfn.Options) {
			if sfnCfg.Region != "" {
				options.Region = sfnCfg.Region
			}
			if credentialsSet(&sfnCfg.Credentials) {
				options.Credentials = getCredentialsProvider(cfg, &sfnCfg.Credentials)
			}
			if sfnCfg.HTTPClient.MaxIdleConns > 0 {
				options.HTTPClient = createHttpClient(&sfnCfg.HTTPClient)
			}
		},
	}

	opts = append(defaultOpts, opts...)

	return sfn.NewFromConfig(cfg, opts...), nil
}

func (b *Builder) NewSQSService(
	cfg aws.Config, serviceName string, opts ...func(options *sqs.Options)) (*sqs.Client, error) {
	sqsCfg, ok := b.config.SQS[serviceName]
	if !ok {
		return nil, fmt.Errorf("sqs config for service %s not found", serviceName)
	}
	defaultOpts := []func(*sqs.Options){
		func(options *sqs.Options) {
			if sqsCfg.Region != "" {
				options.Region = sqsCfg.Region
			}
			if credentialsSet(&sqsCfg.Credentials) {
				options.Credentials = getCredentialsProvider(cfg, &sqsCfg.Credentials)
			}
			if sqsCfg.HTTPClient.MaxIdleConns > 0 {
				options.HTTPClient = createHttpClient(&sqsCfg.HTTPClient)
			}
		},
	}

	opts = append(defaultOpts, opts...)

	return sqs.NewFromConfig(cfg, opts...), nil
}

func createHttpClient(cfg *HTTPClient) *http.Client {
	defaultRoundTripper := http.DefaultTransport
	defaultTransport := defaultRoundTripper.(*http.Transport)

	httpTransport := defaultTransport.Clone()
	httpTransport.MaxIdleConnsPerHost = cfg.MaxIdleConns

	return &http.Client{Transport: httpTransport}
}

func credentialsSet(cfg *CredentialsConfig) bool {
	return cfg.Static.AccessKeyID != "" &&
		cfg.Static.SecretAccessKey != "" ||
		cfg.AssumeRole.ARN != ""
}

func getCredentialsProvider(
	awsConf aws.Config, cfg *CredentialsConfig) aws.CredentialsProvider {
	if cfg.AssumeRole.ARN != "" {
		return aws.NewCredentialsCache(
			stscreds.NewAssumeRoleProvider(
				sts.NewFromConfig(awsConf),
				cfg.AssumeRole.ARN))
	}
	return credentials.NewStaticCredentialsProvider(
		cfg.Static.AccessKeyID,
		cfg.Static.SecretAccessKey,
		cfg.Static.SessionToken)
}
