package instrumentation

import (
	"fmt"

	cbuilder "github.com/scribd/go-sdk/internal/pkg/configuration/builder"
)

type Config struct {
	environment string
	Enabled     bool `mapstructure:"enabled"`
	// Enable Profiler Code Hostspots feature
	CodeHotspotsEnabled bool `mapstructure:"code_hotspots_enabled"`
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

	config.environment = vConf.GetString("ENV")

	return config, nil
}
