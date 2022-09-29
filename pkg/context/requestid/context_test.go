package requestid

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtract(t *testing.T) {
	tests := []struct {
		name          string
		ctxSet        func(ctx context.Context) context.Context
		expected      string
		expectedError error
	}{
		{
			name: "Context without request id",
			ctxSet: func(ctx context.Context) context.Context {
				return ctx
			},
			expectedError: fmt.Errorf("Unable to get the requestID"),
		},
		{
			name: "Context contains request id",
			ctxSet: func(ctx context.Context) context.Context {
				return ToContext(ctx, "test")
			},
			expected: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			resultCtx := tt.ctxSet(ctx)

			requestId, err := Extract(resultCtx)

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.Equal(t, tt.expected, requestId)
			}
		})
	}
}
