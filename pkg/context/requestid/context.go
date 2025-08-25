package requestid

import (
	"context"
	"fmt"
)

type ctxRequestIDMarker struct{}

type ctxRequestID struct {
	requestID string
}

var (
	ctxRequestIDKey = &ctxRequestIDMarker{}
)

// Extract takes the call-scoped requestID from the context.
// If the ctxRequestID  wasn't used, an error is returned.
func Extract(ctx context.Context) (string, error) {
	r, ok := ctx.Value(ctxRequestIDKey).(*ctxRequestID)
	if !ok || r == nil {
		return "", fmt.Errorf("unable to get the requestID")
	}

	return r.requestID, nil
}

// ToContext adds the sdkrequestid.RequestID to the context for extraction later.
// Returning the new context that has been created.
func ToContext(ctx context.Context, requestID string) context.Context {
	r := &ctxRequestID{
		requestID: requestID,
	}
	return context.WithValue(ctx, ctxRequestIDKey, r)
}
