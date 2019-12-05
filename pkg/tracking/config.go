package tracking

import (
	"fmt"
	"time"

	cbuilder "git.lo/microservices/sdk/go-sdk/internal/pkg/configuration/builder"
)

var timeout = 200 * time.Millisecond

// Config stores the configuration for the tracking.
type Config struct {
	SentryDSN     string        `mapstructure:"dsn"`
	SentryTimeout time.Duration `mapstructure:"timeout"`
}

// NewConfig returns a new TrackingConfig instance
func NewConfig() (*Config, error) {
	config := &Config{}

	viperBuilder := cbuilder.New("sentry")
	vConf, err := viperBuilder.Build()
	if err != nil {
		return config, err
	}

	vConf.SetDefault("timeout", timeout)

	if err = vConf.Unmarshal(config); err != nil {
		return config, fmt.Errorf("Unable to decode into struct: %s", err.Error())
	}

	return config, nil
}
