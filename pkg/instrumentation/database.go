package instrumentation

import (
	"context"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gorm.io/gorm"
)

// TraceDatabase fetches the span from the context and injects it to a new
// database session's statement context.
func TraceDatabase(ctx context.Context, db *gorm.DB) *gorm.DB {
	if ctx == nil {
		return db
	}

	parentSpan, _ := tracer.SpanFromContext(ctx)
	return db.WithContext(tracer.ContextWithSpan(ctx, parentSpan))
}
