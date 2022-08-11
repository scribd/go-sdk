package instrumentation

import (
	"fmt"

	cbuilder "github.com/scribd/go-sdk/internal/pkg/configuration/builder"
)

type Config struct {
	environment string
	// TODO must be application name to match pkg/metrics service name.
	// should be fixed in future.
	ServiceName string

	Enabled        bool   `mapstructure:"enabled"`
	ServiceVersion string `mapstructure:"service_version"`
	// Enable Profiler Code Hotspots feature
	CodeHotspotsEnabled bool `mapstructure:"code_hotspots_enabled"`
	// Enable runtime metrics.
	RuntimeMetricsEnabled bool `mapstructure:"runtime_metrics_enabled"`
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
