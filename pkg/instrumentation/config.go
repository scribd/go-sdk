package instrumentation

import (
	"fmt"

	cbuilder "github.com/scribd/go-sdk/internal/pkg/configuration/builder"
)

type Config struct {
	environment string
	Enabled     bool `mapstructure:"enabled"`

	ServiceName    string `mapstructure:"service_name"`
	ServiceVersion string `mapstructure:"service_version"`
	// ProfilerEnabled gates the continuous profiler. It defaults to true (see
	// NewConfig), so the profiler stays on wherever instrumentation is enabled —
	// this is opt-out. Set it to false (e.g. `APP_DATADOG_PROFILER_ENABLED=false`)
	// to run tracing without the profiler, for instance to avoid per-host
	// profiling usage in non-production environments.
	ProfilerEnabled bool `mapstructure:"profiler_enabled"`
	// Enable Profiler Code Hotspots feature
	CodeHotspotsEnabled bool `mapstructure:"code_hotspots_enabled"`
	// Enable runtime metrics.
	RuntimeMetricsEnabled bool `mapstructure:"runtime_metrics_enabled"`
}

// NewConfig returns a new ServerConfig instance.
func NewConfig() (*Config, error) {
	config := &Config{}
	// The profiler is opt-out: default profiler_enabled to true so upgrading
	// go-sdk does not silently disable profiling for existing services.
	viperBuilder := cbuilder.New("datadog").SetDefault("profiler_enabled", true)

	vConf, err := viperBuilder.Build()
	if err != nil {
		return config, err
	}

	if err = vConf.Unmarshal(config); err != nil {
		return config, fmt.Errorf("unable to decode into struct: %s", err.Error())
	}

	config.environment = vConf.GetString("ENV")

	return config, nil
}
