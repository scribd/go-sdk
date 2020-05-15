package tracking

import (
	"fmt"
	"os"

	cbuilder "git.lo/microservices/sdk/go-sdk/internal/pkg/configuration/builder"
)

// Config stores the configuration for the tracking.
type Config struct {
	environment string
	Release     string `mapstructure:"release"`
	SentryDSN   string `mapstructure:"dsn"`
	ServerName  string `mapstructure:"servername"`
}

// NewConfig returns a new TrackingConfig instance
func NewConfig() (*Config, error) {
	config := &Config{}
	config.environment = os.Getenv("APP_ENV")

	viperBuilder := cbuilder.New("sentry")

	vConf, err := viperBuilder.Build()
	if err != nil {
		return config, err
	}

	if err = vConf.Unmarshal(config); err != nil {
		return config, fmt.Errorf("Unable to decode into struct: %s", err.Error())
	}

	return config, nil
}
