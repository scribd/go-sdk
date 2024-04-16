package kafka

import (
	"context"

	"github.com/google/uuid"
	"github.com/twmb/franz-go/pkg/kgo"

	sdkloggercontext "github.com/scribd/go-sdk/pkg/context/logger"
	sdkmetricscontext "github.com/scribd/go-sdk/pkg/context/metrics"
	sdkrequestidcontext "github.com/scribd/go-sdk/pkg/context/requestid"
	sdkinstrumentation "github.com/scribd/go-sdk/pkg/instrumentation"
	sdklogger "github.com/scribd/go-sdk/pkg/logger"
	sdkmetrics "github.com/scribd/go-sdk/pkg/metrics"
)

// RequestFunc may take information from a Kafka message and put it into a
// request context. In Subscribers, RequestFuncs are executed prior to invoking the
// endpoint.
type RequestFunc func(ctx context.Context, msg *kgo.Record) context.Context

// SubscriberResponseFunc may take information from a request context and use it to
// manipulate a Publisher. SubscriberResponseFuncs are only executed in
// consumers, after invoking the endpoint but prior to publishing a reply.
type SubscriberResponseFunc func(ctx context.Context, response interface{}) context.Context

// PublisherResponseFunc may take information from a request context.
// PublisherResponseFunc are only executed in producers, after a request has been produced.
type PublisherResponseFunc func(ctx context.Context, msg *kgo.Record) context.Context

// SetMetrics returns RequestFunc that sets the Metrics client to the request context.
func SetMetrics(m sdkmetrics.Metrics) RequestFunc {
	return func(ctx context.Context, msg *kgo.Record) context.Context {
		return sdkmetricscontext.ToContext(ctx, m)
	}
}

// SetRequestID returns RequestFunc that sets RequestID to the request context if not previously set.
func SetRequestID() RequestFunc {
	return func(ctx context.Context, msg *kgo.Record) context.Context {
		_, err := sdkrequestidcontext.Extract(ctx)
		if err != nil {
			if uuidObject, err := uuid.NewRandom(); err == nil {
				requestID := uuidObject.String()
				return sdkrequestidcontext.ToContext(ctx, requestID)
			}
		}

		return ctx
	}
}

// SetLogger returns RequestFunc that sets SDK Logger to the request context.
// It will also try to setup context values to the logger fields.
func SetLogger(l sdklogger.Logger) RequestFunc {
	return func(ctx context.Context, msg *kgo.Record) context.Context {
		logContext := sdkinstrumentation.TraceLogs(ctx)

		requestID, err := sdkrequestidcontext.Extract(ctx)
		if err != nil {
			l.WithFields(sdklogger.Fields{
				"error": err.Error(),
			}).Tracef("Could not retrieve request id from the context")
		}

		logger := l.WithFields(sdklogger.Fields{
			"pubsub": sdklogger.Fields{
				"request_id": requestID,
			},
			"dd": sdklogger.Fields{
				"trace_id": logContext.TraceID,
				"span_id":  logContext.SpanID,
			},
		})

		return sdkloggercontext.ToContext(ctx, logger)
	}
}
