package middleware

import (
	"net/http"

	sdkmetricscontext "git.lo/microservices/sdk/go-sdk/pkg/context/metrics"
	sdkmetrics "git.lo/microservices/sdk/go-sdk/pkg/metrics"
)

// MetricsMiddleware wraps an instantiated Metrics client that will be
// injected in the request context.
type MetricsMiddleware struct {
	Metrics sdkmetrics.Metrics
}

// NewMetricsMiddleware is a constructor used to build a MetricsMiddleware.
func NewMetricsMiddleware(metrics sdkmetrics.Metrics) MetricsMiddleware {
	return MetricsMiddleware{
		Metrics: metrics,
	}
}

// Handler implements the middlewares.Handlerer interface: it returns a
// http.Handler to be mounted as middleware. The handler injects the Metrics
// client to the request context.
func (sm MetricsMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := sdkmetricscontext.ToContext(r.Context(), sm.Metrics)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
