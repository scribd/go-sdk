package instrumentation

import (
	"context"
	"testing"

	assert "github.com/stretchr/testify/assert"
	mocktracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func TestTraceLogs(t *testing.T) {
	traceIDValue := uint64(1)
	spanIDValue := uint64(2)
	operationName := "TestOperation"
	ctx := context.Background()

	mt := mocktracer.Start()
	defer mt.Stop()

	var span ddtracer.Span

	// Create the ParentSpan.
	span, ctx = ddtracer.StartSpanFromContext(
		ctx,
		operationName,
		ddtracer.WithSpanID(traceIDValue),
		ddtracer.ServiceName("TestService"),
	)
	defer span.Finish()

	// Create the ChildSpan. Using `WithSpanID` it force-set the
	// SpanID, rather than use a random number. If no Parent
	// SpanContext is provided is present, then this will also set
	// the TraceID to the same value˜.
	span, ctx = ddtracer.StartSpanFromContext(
		ctx,
		operationName,
		ddtracer.WithSpanID(spanIDValue),
		ddtracer.ServiceName("TestService"),
	)
	defer span.Finish()

	lc := TraceLogs(ctx)

	assert.Equal(t, traceIDValue, lc.TraceID)
	assert.Equal(t, spanIDValue, lc.SpanID)
}
