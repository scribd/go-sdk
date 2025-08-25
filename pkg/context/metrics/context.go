package metrics

import (
	"context"
	"fmt"

	sdkmetrics "github.com/scribd/go-sdk/pkg/metrics"
)

type ctxMetricsMarker struct{}

type ctxMetrics struct {
	metrics sdkmetrics.Metrics
}

var (
	ctxMetricsKey = &ctxMetricsMarker{}
)

// Extract takes the call-scoped sdkmetrics.Metrics from the context.
// If the ctxMetrics  wasn't used, an error is returned.
func Extract(ctx context.Context) (sdkmetrics.Metrics, error) {
	m, ok := ctx.Value(ctxMetricsKey).(*ctxMetrics)
	if !ok || m == nil {
		return nil, fmt.Errorf("unable to get the metrics")
	}

	return m.metrics, nil
}

// ToContext adds the sdkmetrics.Metrics to the context for extraction later.
// Returning the new context that has been created.
func ToContext(ctx context.Context, metrics sdkmetrics.Metrics) context.Context {
	m := &ctxMetrics{
		metrics: metrics,
	}
	return context.WithValue(ctx, ctxMetricsKey, m)
}
