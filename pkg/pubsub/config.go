package pubsub

import (
	"fmt"
	"strings"
	"time"

	cbuilder "github.com/scribd/go-sdk/internal/pkg/configuration/builder"
)

type (
	SASLMechanism int

	Config struct {
		// Kafka contains the configuration for Kafka client
		Kafka Kafka `mapstructure:"kafka"`
	}

	// Publisher contains the publisher specific configuration
	Publisher struct {
		// MaxAttempts represents the maximum number of times
		// the client will try to send message again in case of failure
		MaxAttempts int `mapstructure:"max_attempts"`
		// WriteTimeout the maximum amount of time the client will wait for message to be written to Kafka topic
		WriteTimeout time.Duration `mapstructure:"write_timeout"`
		// Topic the Kafka topic name to publish messages to
		Topic string `mapstructure:"topic"`
		// Enabled whether the publisher is enabled or not
		Enabled bool `mapstructure:"enabled"`
		// MetricsEnabled controls if metrics publishing is enabled or not
		MetricsEnabled bool `mapstructure:"metrics_enabled"`
	}

	Subscriber struct {
		// Topic the Kafka topic name to retrieve messages from
		Topic string `mapstructure:"topic"`
		// GroupId the Kafka consumer group id
		GroupId string `mapstructure:"group_id"`
		// Enabled whether the subscriber id enabled or not
		Enabled bool `mapstructure:"enabled"`
		// MetricsEnabled controls if metrics publishing is enabled or not
		MetricsEnabled bool `mapstructure:"metrics_enabled"`
		// AutoCommit controls if the subscriber should auto commit messages
		AutoCommit AutoCommit `mapstructure:"auto_commit"`
		// Workers controls the number of workers that will be used to process messages
		Workers int `mapstructure:"workers"`
	}

	AutoCommit struct {
		// Enabled whether the auto commit is enabled or not
		Enabled bool `mapstructure:"enabled"`
	}

	TLS struct {
		// Enabled whether the TLS connection is enabled or not
		Enabled bool `mapstructure:"enabled"`

		// Ca Root CA certificate
		Ca string `mapstructure:"ca"`
		// Cert is a PEM certificate string
		Cert string `mapstructure:"cert_pem"`
		// CertKey is a PEM key certificate string
		CertKey string `mapstructure:"cert_pem_key"`
		// Passphrase is used in case the private key needs to be decrypted
		Passphrase string `mapstructure:"passphrase"`
		// InsecureSkipVerify whether to skip TLS verification or not
		InsecureSkipVerify bool `mapstructure:"insecure_skip_verify"`
	}

	SASL struct {
		// Enabled whether the SASL connection is enabled or not
		Enabled bool `mapstructure:"enabled"`

		// Mechanism is a string representation of the SASL mechanism
		// Currently, only "plain" and "aws_msk_iam" are supported
		Mechanism string `mapstructure:"mechanism"`
		// The username to authenticate Kafka requests
		Username string `mapstructure:"username"`
		// The password to authenticate Kafka requests
		Password string `mapstructure:"password"`
		// AWSMskIam AWS MSK IAM configuration
		// To learn more visit AWS MSK documentation
		// https://docs.aws.amazon.com/msk/latest/developerguide/iam-access-control.html
		AWSMskIam SASLAwsMskIam `mapstructure:"aws_msk_iam"`
	}

	SASLAwsMskIam struct {
		// AWS MSK IAM access key to authenticate AWS MSK requests
		AccessKey string `mapstructure:"access_key"`
		// AWS MSK IAM secret key to authenticate AWS MSK requests
		SecretKey string `mapstructure:"secret_key"`

		// SessionToken is used to authenticate AWS MSK requests via AWS STS service
		// For more information see https://docs.aws.amazon.com/STS/latest/APIReference/welcome.html
		SessionToken string `mapstructure:"session_token"`
		// The client's user agent string
		UserAgent string `mapstructure:"user_agent"`
		// If provided, this role will be used to establish connection to AWS MSK ignoring the static credentials
		AssumableRole string `mapstructure:"role"`
		// Will be passed to AWS STS when assuming the role
		SessionName string `mapstructure:"session_name"`
	}

	Kafka struct {
		// List of Kafka broker URLs to connect to
		BrokerUrls []string `mapstructure:"broker_urls"`
		// Client identification
		ClientId string `mapstructure:"client_id"`
		// Cert is a PEM certificate string
		// Deprecated: use TLS configuration instead
		Cert string `mapstructure:"cert_pem"`
		// CertKey is a PEM key certificate string
		// Deprecated: use TLS configuration instead
		CertKey string `mapstructure:"cert_pem_key"`
		// Security protocol to use for authentication purposes.
		// Deprecated: use TLS and/or SASL configuration instead
		SecurityProtocol string `mapstructure:"security_protocol"`
		// Publisher specific configuration
		Publisher Publisher `mapstructure:"publisher"`
		// Subscriber specific configuration
		Subscriber Subscriber `mapstructure:"subscriber"`
		// Whether to skip SSL verification or not
		// Deprecated: use TLS configuration instead
		SSLVerificationEnabled bool `mapstructure:"ssl_verification_enabled"`

		// TLS configuration
		TLS TLS `mapstructure:"tls"`
		// SASL configuration
		SASL SASL `mapstructure:"sasl"`

		// MetricsEnabled controls if metrics publishing is enabled or not
		MetricsEnabled bool `mapstructure:"metrics_enabled"`
	}
)

const (
	Unknown SASLMechanism = iota
	Plain
	AWSMskIam

	saslMechanismPlainString = "plain"

	saslMechanismAWsMskIamString = "aws_msk_iam"
)

var (
	_stringToSASLMechanism = map[string]SASLMechanism{
		saslMechanismPlainString:     Plain,
		saslMechanismAWsMskIamString: AWSMskIam,
	}
)

// NewConfig returns a new Config instance.
func NewConfig() (*Config, error) {
	config := &Config{}
	viperBuilder := cbuilder.New("pubsub")

	viperBuilder.SetDefault("kafka.subscriber.auto_commit.enabled", true)

	vConf, err := viperBuilder.Build()
	if err != nil {
		return config, err
	}

	if err = vConf.Unmarshal(config); err != nil {
		return config, fmt.Errorf("unable to decode into struct: %s", err.Error())
	}

	if err = config.validate(); err != nil {
		return config, err
	}

	config.Kafka.BrokerUrls = vConf.GetStringSlice("kafka.broker_urls")

	return config, nil
}

func (c *Config) validate() error {
	if c.Kafka.SASL.Enabled && c.Kafka.SASLMechanism() == Unknown {
		var allowedMechanisms []string
		for k := range _stringToSASLMechanism {
			allowedMechanisms = append(allowedMechanisms, k)
		}

		return fmt.Errorf(
			"%s mechanism provided, but following mechanisms are allowed: %s",
			c.Kafka.SASL.Mechanism,
			strings.Join(allowedMechanisms, ","),
		)
	}

	return nil
}

func (k Kafka) SASLMechanism() SASLMechanism {
	if v, ok := _stringToSASLMechanism[k.SASL.Mechanism]; ok {
		return v
	}

	return Unknown
}
