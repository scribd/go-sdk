package sqs

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport"
	"github.com/go-kit/log"
)

type (
	SQSClient interface {
		ChangeMessageVisibility(
			ctx context.Context,
			input *sqs.ChangeMessageVisibilityInput,
			optFns ...func(opts *sqs.Options)) (*sqs.ChangeMessageVisibilityOutput, error)
		DeleteMessage(
			ctx context.Context,
			input *sqs.DeleteMessageInput,
			optFns ...func(opts *sqs.Options)) (*sqs.DeleteMessageOutput, error)
	}

	// Subscriber wraps an endpoint and provides a handler for SQS messages.
	Subscriber struct {
		sqsClient    SQSClient
		e            endpoint.Endpoint
		dec          DecodeRequestFunc
		enc          EncodeResponseFunc
		queueURL     string
		before       []SubscriberRequestFunc
		after        []SubscriberResponseFunc
		errorEncoder ErrorEncoder
		finalizer    []SubscriberFinalizerFunc
		errorHandler transport.ErrorHandler
	}
)

// NewSubscriber constructs a new Subscriber, which provides a ServeMessage method
// and message handlers that wrap the provided endpoint.
func NewSubscriber(
	sqsClient SQSClient,
	e endpoint.Endpoint,
	dec DecodeRequestFunc,
	enc EncodeResponseFunc,
	queueURL string,
	options ...SubscriberOption,
) *Subscriber {
	s := &Subscriber{
		sqsClient:    sqsClient,
		e:            e,
		dec:          dec,
		enc:          enc,
		queueURL:     queueURL,
		errorEncoder: DefaultErrorEncoder,
		errorHandler: transport.NewLogErrorHandler(log.NewNopLogger()),
	}
	for _, option := range options {
		option(s)
	}
	return s
}

// SubscriberOption sets an optional parameter for subscribers.
type SubscriberOption func(*Subscriber)

// SubscriberBefore functions are executed on the producer request object before the
// request is decoded.
func SubscriberBefore(before ...SubscriberRequestFunc) SubscriberOption {
	return func(s *Subscriber) { s.before = append(s.before, before...) }
}

// SubscriberAfter functions are executed on the subscriber reply after the
// endpoint is invoked.
func SubscriberAfter(after ...SubscriberResponseFunc) SubscriberOption {
	return func(s *Subscriber) { s.after = append(s.after, after...) }
}

// SubscriberErrorEncoder is used to encode errors to the subscriber reply
// whenever they're encountered in the processing of a request. Clients can
// use this to provide custom error formatting. By default,
// errors will be published with the DefaultErrorEncoder.
func SubscriberErrorEncoder(ee ErrorEncoder) SubscriberOption {
	return func(s *Subscriber) { s.errorEncoder = ee }
}

// SubscriberErrorHandler is used to handle non-terminal errors. By default, non-terminal errors
// are ignored. This is intended as a diagnostic measure. Finer-grained control
// of error handling, including logging in more detail, should be performed in a
// custom SubscriberErrorEncoder which has access to the context.
func SubscriberErrorHandler(errorHandler transport.ErrorHandler) SubscriberOption {
	return func(s *Subscriber) { s.errorHandler = errorHandler }
}

// SubscriberFinalizer is executed once all the received SQS messages are done being processed.
// By default, no finalizer is registered.
func SubscriberFinalizer(f ...SubscriberFinalizerFunc) SubscriberOption {
	return func(s *Subscriber) { s.finalizer = f }
}

// SubscriberSetContextTimeout returns a SubscriberOption that sets the context timeout.
func SubscriberSetContextTimeout(timeout time.Duration) SubscriberOption {
	return func(s *Subscriber) {
		before := func(ctx context.Context, cancel context.CancelFunc, msg types.Message) context.Context {
			newCtx, newCancel := context.WithTimeout(ctx, timeout)
			defer newCancel()

			return newCtx
		}
		s.before = append(s.before, before)
	}

}

// SubscriberDeleteMessageBefore returns a SubscriberOption that appends a function
// that delete the message from queue to the list of subscriber's before functions.
func SubscriberDeleteMessageBefore() SubscriberOption {
	return func(s *Subscriber) {
		deleteBefore := func(ctx context.Context, cancel context.CancelFunc, msg types.Message) context.Context {
			if err := deleteMessage(ctx, s.sqsClient, s.queueURL, msg); err != nil {
				s.errorHandler.Handle(ctx, err)
				s.errorEncoder(ctx, err, msg, s.sqsClient)
				cancel()
			}
			return ctx
		}
		s.before = append(s.before, deleteBefore)
	}
}

// SubscriberDeleteMessageAfter returns a SubscriberOption that appends a function
// that delete a message from queue to the list of subscriber's after functions.
func SubscriberDeleteMessageAfter() SubscriberOption {
	return func(s *Subscriber) {
		deleteAfter := func(
			ctx context.Context, cancel context.CancelFunc, msg types.Message, _ interface{}) context.Context {
			if err := deleteMessage(ctx, s.sqsClient, s.queueURL, msg); err != nil {
				s.errorHandler.Handle(ctx, err)
				s.errorEncoder(ctx, err, msg, s.sqsClient)
				cancel()
			}
			return ctx
		}
		s.after = append(s.after, deleteAfter)
	}
}

// ServeMessage serves an SQS message.
func (s Subscriber) ServeMessage(ctx context.Context) func(msg types.Message) error {
	return func(msg types.Message) error {
		newCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		if len(s.finalizer) > 0 {
			defer func() {
				for _, f := range s.finalizer {
					f(newCtx, msg)
				}
			}()
		}

		for _, f := range s.before {
			newCtx = f(newCtx, cancel, msg)
		}

		req, err := s.dec(newCtx, msg)
		if err != nil {
			s.errorHandler.Handle(newCtx, err)
			s.errorEncoder(newCtx, err, msg, s.sqsClient)
			return err
		}

		response, err := s.e(newCtx, req)
		if err != nil {
			s.errorHandler.Handle(newCtx, err)
			s.errorEncoder(newCtx, err, msg, s.sqsClient)
			return err
		}

		for _, f := range s.after {
			newCtx = f(newCtx, cancel, msg, response)
		}

		return nil
	}
}

// ErrorEncoder is responsible for encoding an error to the subscriber's reply.
// Users are encouraged to use custom ErrorEncoders to encode errors to
// their replies, and will likely want to pass and check for their own error
// types.
type ErrorEncoder func(ctx context.Context, err error, req types.Message, sqsClient SQSClient)

// SubscriberFinalizerFunc can be used to perform work at the end of a request
// from a producer, after the response has been written to the producer. The
// principal intended use is for request logging.
// Can also be used to delete messages once fully proccessed.
type SubscriberFinalizerFunc func(ctx context.Context, msg types.Message)

// DefaultErrorEncoder simply ignores the message.
func DefaultErrorEncoder(context.Context, error, types.Message, SQSClient) {
}

// SubscriberNackMessageErrorEncoder can be used to perform an immediate nack on the message.
func SubscriberNackMessageErrorEncoder() SubscriberOption {
	return func(s *Subscriber) {
		nackErrorHandler := func(ctx context.Context, err error, msg types.Message, sqsClient SQSClient) {
			_, sqsErr := sqsClient.ChangeMessageVisibility(ctx, &sqs.ChangeMessageVisibilityInput{
				QueueUrl:          &s.queueURL,
				ReceiptHandle:     msg.ReceiptHandle,
				VisibilityTimeout: 1,
			})
			if sqsErr != nil {
				s.errorHandler.Handle(ctx, sqsErr)
			}
		}
		s.errorEncoder = nackErrorHandler
	}
}

func deleteMessage(ctx context.Context, sqsClient SQSClient, queueURL string, msg types.Message) error {
	_, err := sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &queueURL,
		ReceiptHandle: msg.ReceiptHandle,
	})
	return err
}

// EncodeJSONResponse marshals response as json and loads it into an sqs.SendMessageInput MessageBody.
func EncodeJSONResponse(_ context.Context, input *sqs.SendMessageInput, response interface{}) error {
	payload, err := json.Marshal(response)
	if err != nil {
		return err
	}
	input.MessageBody = aws.String(string(payload))
	return nil
}
