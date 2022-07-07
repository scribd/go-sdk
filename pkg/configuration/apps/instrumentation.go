package apps

import (
	"fmt"

	"github.com/scribd/go-sdk/pkg/configuration/builder"
)

// Instrumentation configuration struct.
type Instrumentation struct {
	Environment    string
	Enabled        bool   `mapstructure:"enabled"`
	ServiceVersion string `mapstructure:"service_version"`
	// Enable Profiler Code Hostspots feature
	CodeHotspotsEnabled bool `mapstructure:"code_hotspots_enabled"`
}

// FetchConfig fetches data from builder into it's caller.
func (ins *Instrumentation) FetchConfig(b builder.Builder) error {
	config, err := b.Build()
	if err != nil {
		return fmt.Errorf("unable to build for instrumentation. err: %w", err)
	}

	if err = config.Unmarshal(ins); err != nil {
		return fmt.Errorf("Unable to decode into Instrumentation struct: %w", err)
	}

	ins.Environment = config.String("ENV")

	return nil
}
