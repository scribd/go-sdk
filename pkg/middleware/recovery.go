package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	sdkloggercontext "github.com/scribd/go-sdk/pkg/context/logger"
)

// RecoveryMiddleware is a middleware that recovers from panics and logs them.
type RecoveryMiddleware struct{}

// NewRecoveryMiddleware is a constructor used to build a RecoveryMiddleware.
// IMPORTANT: This middleware should be the last one in the middleware chain.
func NewRecoveryMiddleware() RecoveryMiddleware {
	return RecoveryMiddleware{}
}

// Handler implements the middlewares.Handlerer interface: it returns a
// http.Handler to be mounted as middleware. The Handler recovers from a panic,
// sends a sentry event, sends fatal error log and halts the service.
func (rm RecoveryMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				l, err := sdkloggercontext.Extract(r.Context())
				if err != nil {
					debug.PrintStack()
					log.Printf("logger not found in context: %v\n", err)
					log.Fatalf("http: panic serving URI %s: %v", r.URL.RequestURI(), rec)
				}

				l.WithError(fmt.Errorf("%v", rec)).
					Fatalf("http: panic serving URI %s: %v", r.URL.RequestURI(), rec)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
