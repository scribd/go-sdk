package instrumentation

import (
	"fmt"

	ddmux "github.com/DataDog/dd-trace-go/contrib/gorilla/mux/v2"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
)

const (
	datadogServiceSuffix = "app"
)

// Tracer is a "controller" to a ddtrace.tracer.
//
// Tracer is not exactly a "wrapper" because the tracer is a
// private/global entity in the tracer library and it's not directly
// accessible.
//
// Tracer specifies an implementation of the Datadog tracer which allows
// starting and propagating spans.
type Tracer struct {
	Enabled     bool
	Environment string
	Options     []tracer.StartOption

	globalServiceName string
}

// NewTracer returns a new tracer with the giver configuration and an optional
// list of ddtrace's tracer.StartOptions.
//
// NewTracer assigns universal the version of the service that is running, and will be applied to all spans,
// regardless of whether span service name and config service name match.
func NewTracer(config *Config, options ...tracer.StartOption) *Tracer {
	serviceName := globalServiceName(config.ServiceName)

	options = append(
		options,
		tracer.WithService(serviceName),
		tracer.WithEnv(config.environment),
		tracer.WithUniversalVersion(config.ServiceVersion),
	)

	if config.RuntimeMetricsEnabled {
		options = append(options, tracer.WithRuntimeMetrics())
	}

	if config.CodeHotspotsEnabled {
		options = append(
			options,
			tracer.WithProfilerCodeHotspots(true),
			tracer.WithProfilerEndpoints(true),
		)
	}

	return &Tracer{
		Enabled:           config.Enabled,
		Environment:       config.environment,
		globalServiceName: serviceName,
		Options:           options,
	}
}

// Start starts the current tracer.
func (t *Tracer) Start() error {
	if !t.Enabled {
		return nil
	}

	return tracer.Start(t.Options...)
}

// Stop stops the current tracer.
func (t *Tracer) Stop() {
	if !t.Enabled {
		return
	}

	tracer.Stop()
}

// Router returns an instrumented-mux-compatible router instance traced
// with the global tracer.
//
// Returning a Router is part of the Tracer API to ensure a single entry-point
// for the instrumentation features.
func (t *Tracer) Router() *ddmux.Router {
	return ddmux.NewRouter(ddmux.WithService(t.globalServiceName))
}

func globalServiceName(serviceName string) string {
	return fmt.Sprintf("%s-%s", serviceName, datadogServiceSuffix)
}
