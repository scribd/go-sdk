package middleware

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/scribd/go-sdk/pkg/context/requestid"
)

const (
	// RequestIDHeader is the header key for carrying the RequestID.
	RequestIDHeader string = "X-Request-Id"
)

// RequestIDMiddleware propagates or sets a request ID to the incoming request
// setting the value both in the HTTP header and in the request context.
type RequestIDMiddleware struct{}

// NewRequestIDMiddleware creates a new RequestIDMiddleware.
func NewRequestIDMiddleware() RequestIDMiddleware {
	return RequestIDMiddleware{}
}

// Handler implements the middlewares.Handlerer interface. Returns an
// http.Hanlder to inject the RequestID.
func (rm RequestIDMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(RequestIDHeader)

		if requestID == "" {
			if uuidObject, err := uuid.NewRandom(); err == nil {
				requestID = uuidObject.String()
			}
		}

		ctx := requestid.ToContext(r.Context(), requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
