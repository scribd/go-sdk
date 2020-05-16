package tracking

import (
	"fmt"
	"os"

	cbuilder "git.lo/microservices/sdk/go-sdk/internal/pkg/configuration/builder"
)

// Config stores the configuration for the tracking.
type Config struct {
	SentryDSN string `mapstructure:"dsn"`

	environment string
	release     string
	serverName  string
}

// NewConfig returns a new TrackingConfig instance.
func NewConfig() (*Config, error) {
	config := &Config{}
	config.environment = os.Getenv("APP_ENV")
	config.release = os.Getenv("APP_VERSION")
	config.serverName = os.Getenv("APP_SERVER_NAME")

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
