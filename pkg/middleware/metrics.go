package middleware

import (
	"context"
	"net/http"

	"git.lo/microservices/sdk/go-sdk/pkg/metrics"
)

const (
	// Metrics is the context key for the carrying the Metrics client.
	Metrics = "Metrics"
)

// MetricsMiddleware wraps an instantiated Metrics client that will be
// injected in the request context.
type MetricsMiddleware struct {
	Metrics metrics.Metrics
}

// NewMetricsMiddleware is a constructor used to build a MetricsMiddleware.
func NewMetricsMiddleware(metrics metrics.Metrics) MetricsMiddleware {
	return MetricsMiddleware{
		Metrics: metrics,
	}
}

// Handler implements the middlewares.Handlerer interface: it returns a
// http.Handler to be mounted as middleware. The handler injects the Metrics
// client to the request context.
func (sm MetricsMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), Metrics, sm.Metrics)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
