package aws

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/stretchr/testify/assert"
)

func TestBuilder(t *testing.T) {
	httpClientConfig := HTTPClient{
		MaxIdleConns: 100,
	}
	staticCredentialsConfig := CredentialsConfig{
		Static: StaticConfig{
			AccessKeyID:     "test",
			SecretAccessKey: "test",
			SessionToken:    "test",
		},
	}
	assumeRoleConfig := CredentialsConfig{
		AssumeRole: AssumeRoleConfig{
			ARN: "test",
		},
	}

	cfg := &Config{
		AWSConfig: AWSConfig{
			Region: "us-east-2",
		},
		S3: map[string]S3Config{
			"default": {
				HTTPClient: httpClientConfig,
			},
			"test": {
				Region:      "us-east-1",
				Credentials: staticCredentialsConfig,
			},
			"test2": {
				Region:      "us-west-2",
				Credentials: assumeRoleConfig,
			},
		},
		SagemakerRuntime: map[string]SagemakerRuntimeConfig{
			"default": {
				HTTPClient: httpClientConfig,
			},
			"test": {
				Region:      "us-east-1",
				Credentials: staticCredentialsConfig,
			},
			"test2": {
				Region:      "us-west-2",
				Credentials: assumeRoleConfig,
			},
		},
		SFN: map[string]SFNConfig{
			"default": {
				HTTPClient: httpClientConfig,
			},
			"test": {
				Region:      "us-east-1",
				Credentials: staticCredentialsConfig,
			},
			"test2": {
				Region:      "us-west-2",
				Credentials: assumeRoleConfig,
			},
		},
		SQS: map[string]SQSConfig{
			"default": {
				HTTPClient: httpClientConfig,
			},
			"test": {
				Region:      "us-east-1",
				Credentials: staticCredentialsConfig,
			},
			"test2": {
				Region:      "us-west-2",
				Credentials: assumeRoleConfig,
			},
		},
	}
	builder := NewBuilder(cfg)

	c, err := builder.LoadConfig(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, c)

	tests := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{
			name: "S3 service with HTTP client settings",
			fn: func(t *testing.T) {
				s3svc, err := builder.NewS3Service(c, "default")
				assert.NoError(t, err)
				assert.NotNil(t, s3svc)

				opts := s3svc.Options()
				assert.Equal(t, opts.Region, "us-east-2")
				assert.NotNil(t, opts.HTTPClient)

				httpTransport := opts.HTTPClient.(*http.Client).Transport.(*http.Transport)
				assert.Equal(t, httpTransport.MaxIdleConnsPerHost, 100)
			},
		}, {
			name: "S3 service with static credentials",
			fn: func(t *testing.T) {
				s3svc, err := builder.NewS3Service(c, "test")
				assert.NoError(t, err)
				assert.NotNil(t, s3svc)

				opts := s3svc.Options()
				assert.Equal(t, opts.Region, "us-east-1")
				assert.NotNil(t, opts.Credentials)

				creds := opts.Credentials.(credentials.StaticCredentialsProvider)
				assert.Equal(t, creds.Value.AccessKeyID, "test")
				assert.Equal(t, creds.Value.SecretAccessKey, "test")
				assert.Equal(t, creds.Value.SessionToken, "test")
			},
		},
		{
			name: "S3 service with assume role credentials",
			fn: func(t *testing.T) {
				s3svc, err := builder.NewS3Service(c, "test2")
				assert.NoError(t, err)
				assert.NotNil(t, s3svc)

				opts := s3svc.Options()
				assert.Equal(t, opts.Region, "us-west-2")
				assert.NotNil(t, opts.Credentials)

				creds := opts.Credentials.(*aws.CredentialsCache)
				assert.NotNil(t, creds)
			},
		},
		{
			name: "non-existent S3 service",
			fn: func(t *testing.T) {
				s3svc, err := builder.NewS3Service(c, "non-existent")
				assert.Error(t, err)
				assert.Nil(t, s3svc)
			},
		},
		{
			name: "SagemakerRuntime service with HTTP client settings",
			fn: func(t *testing.T) {
				sagemakerRuntimeSvc, err := builder.NewSagemakerRuntimeService(c, "default")
				assert.NoError(t, err)
				assert.NotNil(t, sagemakerRuntimeSvc)

				opts := sagemakerRuntimeSvc.Options()
				assert.Equal(t, opts.Region, "us-east-2")
				assert.NotNil(t, opts.HTTPClient)

				httpTransport := opts.HTTPClient.(*http.Client).Transport.(*http.Transport)
				assert.Equal(t, httpTransport.MaxIdleConnsPerHost, 100)
			},
		},
		{
			name: "SagemakerRuntime service with static credentials",
			fn: func(t *testing.T) {
				sagemakerRuntimeSvc, err := builder.NewSagemakerRuntimeService(c, "test")
				assert.NoError(t, err)
				assert.NotNil(t, sagemakerRuntimeSvc)

				opts := sagemakerRuntimeSvc.Options()
				assert.Equal(t, opts.Region, "us-east-1")
				assert.NotNil(t, opts.Credentials)

				creds := opts.Credentials.(credentials.StaticCredentialsProvider)
				assert.Equal(t, creds.Value.AccessKeyID, "test")
				assert.Equal(t, creds.Value.SecretAccessKey, "test")
				assert.Equal(t, creds.Value.SessionToken, "test")
			},
		},
		{
			name: "SagemakerRuntime service with assume role credentials",
			fn: func(t *testing.T) {
				sagemakerRuntimeSvc, err := builder.NewSagemakerRuntimeService(c, "test2")
				assert.NoError(t, err)
				assert.NotNil(t, sagemakerRuntimeSvc)

				opts := sagemakerRuntimeSvc.Options()
				assert.Equal(t, opts.Region, "us-west-2")
				assert.NotNil(t, opts.Credentials)

				creds := opts.Credentials.(*aws.CredentialsCache)
				assert.NotNil(t, creds)
			},
		},
		{
			name: "non-existent SagemakerRuntime service",
			fn: func(t *testing.T) {
				sagemakerRuntimeSvc, err := builder.NewSagemakerRuntimeService(c, "non-existent")
				assert.Error(t, err)
				assert.Nil(t, sagemakerRuntimeSvc)
			},
		},
		{
			name: "Sfn service with HTTP client settings",
			fn: func(t *testing.T) {
				sfnSvc, err := builder.NewSFNService(c, "default")

				// SFN service does not have access to options for now
				assert.NoError(t, err)
				assert.NotNil(t, sfnSvc)
			},
		},
		{
			name: "Sfn service with static credentials",
			fn: func(t *testing.T) {
				sfnSvc, err := builder.NewSFNService(c, "test")

				// SFN service does not have access to options for now
				assert.NoError(t, err)
				assert.NotNil(t, sfnSvc)
			},
		},
		{
			name: "Sfn service with assume role credentials",
			fn: func(t *testing.T) {
				sfnSvc, err := builder.NewSFNService(c, "test2")

				// SFN service does not have access to options for now
				assert.NoError(t, err)
				assert.NotNil(t, sfnSvc)
			},
		},
		{
			name: "non-existent Sfn service",
			fn: func(t *testing.T) {
				sfnSvc, err := builder.NewSFNService(c, "non-existent")
				assert.Error(t, err)
				assert.Nil(t, sfnSvc)
			},
		},
		{
			name: "Sqs service with HTTP client settings",
			fn: func(t *testing.T) {
				sqsSvc, err := builder.NewSQSService(c, "default")

				// SQS service does not have access to options for now
				assert.NoError(t, err)
				assert.NotNil(t, sqsSvc)
			},
		},
		{
			name: "Sqs service with static credentials",
			fn: func(t *testing.T) {
				sqsSvc, err := builder.NewSQSService(c, "test")

				// SQS service does not have access to options for now
				assert.NoError(t, err)
				assert.NotNil(t, sqsSvc)
			},
		},
		{
			name: "Sqs service with assume role credentials",
			fn: func(t *testing.T) {
				sqsSvc, err := builder.NewSQSService(c, "test2")

				// SQS service does not have access to options for now
				assert.NoError(t, err)
				assert.NotNil(t, sqsSvc)
			},
		},
		{
			name: "non-existent Sqs service",
			fn: func(t *testing.T) {
				sqsSvc, err := builder.NewSQSService(c, "non-existent")
				assert.Error(t, err)
				assert.Nil(t, sqsSvc)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tc.fn(t)
		})
	}
}
