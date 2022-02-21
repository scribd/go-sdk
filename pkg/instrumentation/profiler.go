package instrumentation

import (
	"time"

	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

// Profiler wraps DataDog profiles exporter.
type Profiler struct {
	enabled bool
	start   func(options ...profiler.Option) error
	stop    func()
	options []profiler.Option
}

// Start calls DD profiler with options set during Profiler construction.
func (p *Profiler) Start() error {
	if !p.enabled {
		return nil
	}

	return p.start(p.options...)
}

// Stop DataDog profiles exporter.
func (p *Profiler) Stop() {
	p.stop()
}

// NewProfiler constructs new profiler with options.
// You can include common options like: profiler.WithService(appName), profiler.WithVersion(version).
func NewProfiler(config *Config, options ...profiler.Option) *Profiler {
	options = append(
		options,
		profiler.WithEnv(config.environment),
		profiler.WithVersion(config.ServiceVersion),
	)

	if config.CodeHotspotsEnabled {
		options = append(
			options,
			profiler.CPUDuration(time.Minute),
			profiler.WithPeriod(time.Minute),
		)
	}

	return &Profiler{
		enabled: config.Enabled,
		start:   profiler.Start,
		stop:    profiler.Stop,
		options: options,
	}
}
