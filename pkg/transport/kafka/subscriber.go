package kafka

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport"
	"github.com/go-kit/log"
	"github.com/twmb/franz-go/pkg/kgo"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	kafkasdk "github.com/scribd/go-sdk/pkg/instrumentation/kafka"
	sdkkafka "github.com/scribd/go-sdk/pkg/pubsub/kafka"
)

// Subscriber wraps an endpoint and provides a handler for kafka messages.
type Subscriber struct {
	e            endpoint.Endpoint
	dec          DecodeRequestFunc
	before       []RequestFunc
	after        []SubscriberResponseFunc
	finalizer    []SubscriberFinalizerFunc
	errorHandler transport.ErrorHandler
	errorEncoder ErrorEncoder
}

// NewSubscriber constructs a new subscriber provides a handler for kafka messages.
func NewSubscriber(
	e endpoint.Endpoint,
	dec DecodeRequestFunc,
	opts ...SubscriberOption,
) *Subscriber {
	c := &Subscriber{
		e:            e,
		dec:          dec,
		errorEncoder: DefaultErrorEncoder,
		errorHandler: transport.NewLogErrorHandler(log.NewNopLogger()),
	}
	for _, opt := range opts {
		opt(c)
	}

	return c
}

// SubscriberOption sets an optional parameter for subscribers.
type SubscriberOption func(consumer *Subscriber)

// SubscriberBefore functions are executed on the subscriber message object
// before the request is decoded.
func SubscriberBefore(before ...RequestFunc) SubscriberOption {
	return func(c *Subscriber) {
		c.before = append(c.before, before...)
	}
}

// SubscriberAfter functions are executed on the subscriber reply after the
// endpoint is invoked, but before anything is published to the reply.
func SubscriberAfter(after ...SubscriberResponseFunc) SubscriberOption {
	return func(c *Subscriber) {
		c.after = append(c.after, after...)
	}
}

// SubscriberErrorEncoder is used to encode errors to the subscriber reply
// whenever they're encountered in the processing of a request. Clients can
// use this to provide custom error formatting. By default,
// errors will be published with the DefaultErrorEncoder.
func SubscriberErrorEncoder(ee ErrorEncoder) SubscriberOption {
	return func(s *Subscriber) { s.errorEncoder = ee }
}

// SubscriberErrorHandler is used to handle non-terminal errors. By default, non-terminal errors
// are ignored. This is intended as a diagnostic measure.
func SubscriberErrorHandler(errorHandler transport.ErrorHandler) SubscriberOption {
	return func(c *Subscriber) {
		c.errorHandler = errorHandler
	}
}

// SubscriberFinalizer is executed at the end of every message processing.
// By default, no finalizer is registered.
func SubscriberFinalizer(f ...SubscriberFinalizerFunc) SubscriberOption {
	return func(c *Subscriber) {
		c.finalizer = append(c.finalizer, f...)
	}
}

// ServeMsg provides kafka.MsgHandler.
func (s Subscriber) ServeMsg(h Handler) sdkkafka.MsgHandler {
	return func(msg *kgo.Record) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if len(s.finalizer) > 0 {
			defer func() {
				for _, f := range s.finalizer {
					f(ctx, msg)
				}
			}()
		}

		for _, f := range s.before {
			ctx = f(ctx, msg)
		}

		request, err := s.dec(ctx, msg)
		if err != nil {
			s.errorEncoder(ctx, err, msg, h)
			s.errorHandler.Handle(ctx, err)
			return
		}

		response, err := s.e(ctx, request)
		if err != nil {
			s.errorEncoder(ctx, err, msg, h)
			s.errorHandler.Handle(ctx, err)
			return
		}

		for _, f := range s.after {
			ctx = f(ctx, response)
		}
	}
}

// SubscriberFinalizerFunc can be used to perform work at the end of message processing,
// after the response has been constructed. The principal
// intended use is for request logging.
type SubscriberFinalizerFunc func(ctx context.Context, msg *kgo.Record)

// ErrorEncoder is responsible for encoding an error to the subscriber reply.
// Users are encouraged to use custom ErrorEncoders to encode errors to
// their replies, and will likely want to pass and check for their own error
// types.
type ErrorEncoder func(ctx context.Context,
	err error, msg *kgo.Record, h Handler)

// DefaultErrorEncoder simply ignores the message.
func DefaultErrorEncoder(ctx context.Context,
	err error, msg *kgo.Record, h Handler) {
}

// NewInstrumentedSubscriber constructs a new subscriber provides a handler for kafka messages.
// It also instruments the subscriber with datadog tracing.
func NewInstrumentedSubscriber(e endpoint.Endpoint, dec DecodeRequestFunc, opts ...SubscriberOption) *Subscriber {
	options := []SubscriberOption{
		SubscriberBefore(startMessageHandlerTrace),
		SubscriberFinalizer(finishMessageHandlerTrace),
	}

	options = append(options, opts...)

	return NewSubscriber(e, dec, options...)
}

func startMessageHandlerTrace(ctx context.Context, msg *kgo.Record) context.Context {
	if spanctx, err := tracer.Extract(kafkasdk.NewMessageCarrier(msg)); err == nil {
		span := tracer.StartSpan("kafka.msghandler", tracer.ChildOf(spanctx))

		ctx = tracer.ContextWithSpan(ctx, span)
	}

	return ctx
}

func finishMessageHandlerTrace(ctx context.Context, msg *kgo.Record) {
	if span, ok := tracer.SpanFromContext(ctx); ok {
		span.Finish()
	}
}
