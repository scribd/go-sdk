package aws

import (
	"fmt"
	"os"
	"strings"

	cbuilder "github.com/scribd/go-sdk/internal/pkg/configuration/builder"
)

type (
	HTTPClient struct {
		// MaxIdleConns, if non-zero, controls the maximum idle
		// (keep-alive) connections to keep per-host.
		MaxIdleConns int `mapstructure:"max_idle_conns"`
	}

	StaticConfig struct {
		// AccessKeyID is the AWS access key ID.
		AccessKeyID string `mapstructure:"access_key_id"`
		// SecretAccessKey is the AWS secret access key.
		SecretAccessKey string `mapstructure:"secret_access_key"`
		// SessionToken is the AWS session token.
		SessionToken string `mapstructure:"session_token"`
	}

	AssumeRoleConfig struct {
		// ARN is the ARN of the role to assume.
		ARN string `mapstructure:"arn"`
	}

	CredentialsConfig struct {
		// Static is the configuration for static credentials.
		Static StaticConfig `mapstructure:"static"`
		// AssumeRole is the configuration for assuming a role.
		AssumeRole AssumeRoleConfig `mapstructure:"assume_role"`
	}

	AWSConfig struct {
		// Region is the region to send requests to.
		Region string `mapstructure:"region"`
		// HTTPClient is the configuration for the HttpClient AWS the SDK's API clients will use to invoke HTTP requests.
		HTTPClient HTTPClient `mapstructure:"http_client"`
	}

	S3Config struct {
		// Region is the region to send requests to.
		Region string `mapstructure:"region"`
		// HTTPClient is the configuration for the HttpClient AWS the SDK's API clients will use to invoke HTTP requests.
		HTTPClient HTTPClient `mapstructure:"http_client"`
		// Credentials is the configuration for the AWS credentials.
		Credentials CredentialsConfig `mapstructure:"credentials"`
	}

	SagemakerRuntimeConfig struct {
		// Region is the region to send requests to.
		Region string `mapstructure:"region"`
		// HTTPClient is the configuration for the HttpClient AWS the SDK's API clients will use to invoke HTTP requests.
		HTTPClient HTTPClient `mapstructure:"http_client"`
		// Credentials is the configuration for the AWS credentials.
		Credentials CredentialsConfig `mapstructure:"credentials"`
	}

	SFNConfig struct {
		// Region is the region to send requests to.
		Region string `mapstructure:"region"`
		// HTTPClient is the configuration for the HttpClient AWS the SDK's API clients will use to invoke HTTP requests.
		HTTPClient HTTPClient `mapstructure:"http_client"`
		// Credentials is the configuration for the AWS credentials.
		Credentials CredentialsConfig `mapstructure:"credentials"`
	}

	SQSConfig struct {
		// Region is the region to send requests to.
		Region string `mapstructure:"region"`
		// HTTPClient is the configuration for the HttpClient AWS the SDK's API clients will use to invoke HTTP requests.
		HTTPClient HTTPClient `mapstructure:"http_client"`
		// Credentials is the configuration for the AWS credentials.
		Credentials CredentialsConfig `mapstructure:"credentials"`
	}

	Config struct {
		// AWSConfig is the configuration for the AWS SDK.
		AWSConfig AWSConfig `mapstructure:"config"`
		// CredentialsConfig is the configuration for the AWS credentials.
		CredentialsConfig CredentialsConfig `mapstructure:"credentials"`
		// S3 is the configuration for the S3 clients.
		S3 map[string]S3Config `mapstructure:"s3"`
		// SagemakerRuntime is the configuration for the SagemakerRuntime clients.
		SagemakerRuntime map[string]SagemakerRuntimeConfig `mapstructure:"sagemakerruntime"`
		// SFN is the configuration for the Sfn clients.
		SFN map[string]SFNConfig `mapstructure:"sfn"`
		// SQS is the configuration for the Sqs clients.
		SQS map[string]SQSConfig `mapstructure:"sqs"`
	}
)

func NewConfig() (*Config, error) {
	config := &Config{}
	viperBuilder := cbuilder.New("aws")

	appName := strings.ReplaceAll(os.Getenv("APP_SETTINGS_NAME"), "-", "_")
	viperBuilder.SetDefault("aws", fmt.Sprintf("%s_%s", appName, os.Getenv("APP_ENV")))

	vConf, err := viperBuilder.Build()
	if err != nil {
		return config, err
	}

	if err = vConf.Unmarshal(config); err != nil {
		return config, fmt.Errorf("unable to decode into struct: %s", err.Error())
	}

	return config, nil
}
