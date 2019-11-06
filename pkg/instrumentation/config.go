package instrumentation

import (
	"fmt"

	cbuilder "git.lo/microservices/sdk/go-sdk/internal/pkg/configuration/builder"
)

type Config struct {
	Enabled bool `mapstructure:"enabled"`
}

// NewConfig returns a new ServerConfig instance.
func NewConfig() (*Config, error) {
	config := &Config{}
	viperBuilder := cbuilder.New("datadog")

	vConf, err := viperBuilder.Build()
	if err != nil {
		return config, err
	}

	if err = vConf.Unmarshal(config); err != nil {
		return config, fmt.Errorf("Unable to decode into struct: %s", err.Error())
	}

	return config, nil
}
