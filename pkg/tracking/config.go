package tracking

import (
	"fmt"

	cbuilder "git.lo/microservices/sdk/go-sdk/internal/pkg/configuration/builder"
)

var timeout = 200

// Config stores the configuration for the tracking.
type Config struct {
	SentryDSN     string `mapstructure:"dsn"`
	SentryTimeout int    `mapstructure:"timeout"`
}

// NewConfig returns a new TrackingConfig instance
func NewConfig() (*Config, error) {
	config := &Config{}

	viperBuilder := cbuilder.New("sentry")

	viperBuilder.SetDefault("timeout", timeout)

	vConf, err := viperBuilder.Build()
	if err != nil {
		return config, err
	}

	if err = vConf.Unmarshal(config); err != nil {
		return config, fmt.Errorf("Unable to decode into struct: %s", err.Error())
	}

	return config, nil
}
