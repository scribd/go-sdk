package sqs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"

	sdkloggercontext "github.com/scribd/go-sdk/pkg/context/logger"
	sdkmetricscontext "github.com/scribd/go-sdk/pkg/context/metrics"
	sdkrequestidcontext "github.com/scribd/go-sdk/pkg/context/requestid"
	sdkinstrumentation "github.com/scribd/go-sdk/pkg/instrumentation"
	sdklogger "github.com/scribd/go-sdk/pkg/logger"
	sdkmetrics "github.com/scribd/go-sdk/pkg/metrics"
)

// SubscriberRequestFunc may take information from a subscriber request result and
// put it into a request context. In Subscribers, RequestFuncs are executed prior
// to invoking the endpoint.
// use cases eg. in Subscriber : extract message information into context.
type SubscriberRequestFunc func(
	ctx context.Context, cancel context.CancelFunc, message types.Message) context.Context

// PublisherRequestFunc may take information from a producer request and put it into a
// request context, or add some informations to SendMessageInput. In Publishers,
// RequestFuncs are executed prior to publishing the message but after encoding.
// use cases eg. in Publisher : enforce some message attributes to SendMessageInput.
type PublisherRequestFunc func(ctx context.Context, input *sqs.SendMessageInput) context.Context

// SubscriberResponseFunc may take information from a request context and use it to
// manipulate a Publisher. SubscriberResponseFunc are only executed in
// subscriber, after invoking the endpoint.
// use cases eg. : Pipe information from request message, delete msg from queue, etc.
type SubscriberResponseFunc func(
	ctx context.Context, cancel context.CancelFunc, message types.Message, resp interface{}) context.Context

// PublisherResponseFunc may take information from an sqs.SendMessageOutput and
// fetch response using the Client. SQS is not req-reply out-of-the-box. Responses need to be fetched.
// PublisherResponseFunc are only executed in producers, after a request has been made,
// but prior to its response being decoded. So this is the perfect place to fetch actual response.
type PublisherResponseFunc func(
	context.Context, SQSPublisher, *sqs.SendMessageOutput) (context.Context, types.Message, error)

// SetPublisherLogger returns PublisherRequestFunc that sets SDK Logger to the request context.
// It will also try to setup context values to the logger fields.
func SetPublisherLogger(l sdklogger.Logger) PublisherRequestFunc {
	return func(ctx context.Context, input *sqs.SendMessageInput) context.Context {
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

// SetSubscriberLogger returns SubscriberRequestFunc that sets SDK Logger to the request context.
// It will also try to setup context values to the logger fields.
func SetSubscriberLogger(l sdklogger.Logger) SubscriberRequestFunc {
	return func(ctx context.Context, cancel context.CancelFunc, message types.Message) context.Context {
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

// SetPublisherMetrics returns PublisherRequestFunc that sets the Metrics client to the request context.
func SetPublisherMetrics(m sdkmetrics.Metrics) PublisherRequestFunc {
	return func(ctx context.Context, input *sqs.SendMessageInput) context.Context {
		return sdkmetricscontext.ToContext(ctx, m)
	}
}

// SetSubscriberMetrics returns SubscriberRequestFunc that sets the Metrics client to the request context.
func SetSubscriberMetrics(m sdkmetrics.Metrics) SubscriberRequestFunc {
	return func(ctx context.Context, cancel context.CancelFunc, message types.Message) context.Context {
		return sdkmetricscontext.ToContext(ctx, m)
	}
}
