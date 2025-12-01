package instrumentation

import (
	"context"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
)

// LogContext represents a log state that can be used to collerate logs
// emitted across the request cycle.
type LogContext struct {
	// TraceID returns the trace ID that this context is carrying.
	TraceID string

	// SpanID returns the span ID that this context is carrying.
	SpanID uint64
}

// TraceLogs extracts and returns a LogContext for Logs.
func TraceLogs(ctx context.Context) *LogContext {
	span, _ := tracer.SpanFromContext(ctx)

	return &LogContext{
		TraceID: span.Context().TraceID(),
		SpanID:  span.Context().SpanID(),
	}
}
