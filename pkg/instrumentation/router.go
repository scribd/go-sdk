package instrumentation

import (
	"fmt"

	ddmux "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
	tracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const (
	datadogServiceSuffix = "app"
)

// Tracer is a "controller" to a ddtrace.tracer.
//
// Tracer is not exactly a "wrapper" because the tracer is a
// private/global entity in the tracer library and it's not directly
// accessible.
// - https://godoc.org/gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer#Start
// - https://github.com/DataDog/dd-trace-go/blob/v1.19.0/ddtrace/tracer/tracer.go
//
// Tracer specifies an implementation of the Datadog tracer which allows
// starting and propagating spans.
type Tracer struct {
	Enabled     bool
	Environment string
	Options     []tracer.StartOption
}

// New tracer returns a new tracer with the giver configuration and an optional
// list of ddtrace's tracer.StartOptions.
func NewTracer(config *Config, options ...tracer.StartOption) *Tracer {
	options = append(
		options,
		tracer.WithGlobalTag("env", config.environment),
	)

	return &Tracer{
		Enabled:     config.Enabled,
		Environment: config.environment,
		Options:     options,
	}
}

// Start starts the current tracer.
func (t *Tracer) Start() {
	if !t.Enabled {
		return
	}

	tracer.Start(t.Options...)
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
func (t *Tracer) Router(appName string) *ddmux.Router {
	serviceName := fmt.Sprintf("%s-%s", appName, datadogServiceSuffix)
	return ddmux.NewRouter(ddmux.WithServiceName(serviceName))
}
