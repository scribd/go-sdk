package instrumentation

import (
	"time"

	"github.com/DataDog/dd-trace-go/v2/profiler"
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
	serviceName := globalServiceName(config.ServiceName)

	options = append(
		options,
		profiler.WithService(serviceName),
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
		// Profiler runs when instrumentation is enabled and profiling is not
		// opted out. ProfilerEnabled defaults to true (see NewConfig), so this
		// preserves the previous behavior (profiler on wherever Enabled is on)
		// while allowing a service to keep tracing but disable the profiler via
		// profiler_enabled: false.
		enabled: config.Enabled && config.ProfilerEnabled,
		start:   profiler.Start,
		stop:    profiler.Stop,
		options: options,
	}
}
