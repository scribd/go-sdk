package pubsub

import (
	"fmt"
	"time"

	cbuilder "github.com/scribd/go-sdk/internal/pkg/configuration/builder"
)

type (
	Config struct {
		Kafka Kafka `mapstructure:"kafka"`
	}

	Publisher struct {
		MaxAttempts  int           `mapstructure:"max_attempts"`
		WriteTimeout time.Duration `mapstructure:"write_timeout"`
		Topic        string        `mapstructure:"topic"`
		Enabled      bool          `mapstructure:"enabled"`
	}

	Subscriber struct {
		Topic   string `mapstructure:"topic"`
		GroupId string `mapstructure:"group_id"`
		Enabled bool   `mapstructure:"enabled"`
	}

	Kafka struct {
		BrokerUrls             []string   `mapstructure:"broker_urls"`
		ClientId               string     `mapstructure:"client_id"`
		Cert                   string     `mapstructure:"cert_pem"`
		CertKey                string     `mapstructure:"cert_pem_key"`
		SecurityProtocol       string     `mapstructure:"security_protocol"`
		Publisher              Publisher  `mapstructure:"publisher"`
		Subscriber             Subscriber `mapstructure:"subscriber"`
		SSLVerificationEnabled bool       `mapstructure:"ssl_verification_enabled"`
	}
)

// NewConfig returns a new Config instance.
func NewConfig() (*Config, error) {
	config := &Config{}
	viperBuilder := cbuilder.New("pubsub")

	vConf, err := viperBuilder.Build()
	if err != nil {
		return config, err
	}

	if err = vConf.Unmarshal(config); err != nil {
		return config, fmt.Errorf("unable to decode into struct: %s", err.Error())
	}

	config.Kafka.BrokerUrls = vConf.GetStringSlice("kafka.broker_urls")

	return config, nil
}
