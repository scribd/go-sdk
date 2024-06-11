package instrumentation

import (
	"context"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestTraceDatabase(t *testing.T) {
	const (
		mockOperationName = "test"
		mockDBName        = "test_db"
	)

	dbFile := path.Join(t.TempDir(), mockDBName)
	db, err := gorm.Open(sqlite.Open(dbFile))
	if err != nil {
		t.Fatalf("Failed to open DB: %s", err)
	}

	mockTracer := mocktracer.Start()
	defer mockTracer.Stop()

	mockSpan := tracer.StartSpan(mockOperationName)
	defer mockSpan.Finish()

	mockContext := tracer.ContextWithSpan(context.Background(), mockSpan)
	db = TraceDatabase(mockContext, db)

	_, ok := tracer.SpanFromContext(db.Statement.Context)
	if !ok {
		require.True(t, ok)
	}

	openSpans := mockTracer.OpenSpans()
	require.Len(t, openSpans, 1)
	require.Equal(t, openSpans[0].OperationName(), mockOperationName)
}
