package instrumentation

import (
	"context"
	"testing"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/mocktracer"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"github.com/stretchr/testify/assert"
)

func TestTraceLogs(t *testing.T) {
	traceIDValue := uint64(1)
	spanIDValue := uint64(2)
	operationName := "TestOperation"
	ctx := context.Background()

	mt := mocktracer.Start()
	defer mt.Stop()

	var span *tracer.Span

	// Create the ParentSpan.
	span, ctx = tracer.StartSpanFromContext(
		ctx,
		operationName,
		tracer.WithSpanID(traceIDValue),
		tracer.ServiceName("TestService"),
	)
	defer span.Finish()

	parentTraceID := span.Context().TraceID()

	// Create the ChildSpan. Using `WithSpanID` it force-set the
	// SpanID, rather than use a random number. If no Parent
	// SpanContext is provided is present, then this will also set
	// the TraceID to the same value˜.
	span, ctx = tracer.StartSpanFromContext(
		ctx,
		operationName,
		tracer.WithSpanID(spanIDValue),
		tracer.ServiceName("TestService"),
	)
	defer span.Finish()

	lc := TraceLogs(ctx)

	assert.Equal(t, parentTraceID, lc.TraceID)
	assert.Equal(t, spanIDValue, lc.SpanID)
}
