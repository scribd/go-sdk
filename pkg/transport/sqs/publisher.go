package sqs

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/go-kit/kit/endpoint"
)

type contextKey int

const (
	// ContextKeyResponseQueueURL is the context key that allows fetching
	// and setting the response queue URL from and into context.
	ContextKeyResponseQueueURL contextKey = iota
)

type (
	SQSPublisher interface {
		Publish(ctx context.Context, message *sqs.SendMessageInput) (*sqs.SendMessageOutput, error)
	}

	// Publisher wraps an Publisher client, and provides a method that
	// implements endpoint.Endpoint.
	Publisher struct {
		Handler  SQSPublisher
		queueURL string
		enc      EncodeRequestFunc
		dec      DecodeResponseFunc
		before   []PublisherRequestFunc
		after    []PublisherResponseFunc
	}
)

// NewPublisher constructs a usable Publisher for a single remote method.
func NewPublisher(
	handler SQSPublisher,
	queueURL string,
	enc EncodeRequestFunc,
	dec DecodeResponseFunc,
	options ...PublisherOption,
) *Publisher {
	p := &Publisher{
		Handler:  handler,
		queueURL: queueURL,
		enc:      enc,
		dec:      dec,
	}
	for _, option := range options {
		option(p)
	}
	return p
}

// PublisherOption sets an optional parameter for clients.
type PublisherOption func(*Publisher)

// PublisherBefore sets the RequestFuncs that are applied to the outgoing SQS
// request before it's invoked.
func PublisherBefore(before ...PublisherRequestFunc) PublisherOption {
	return func(p *Publisher) { p.before = append(p.before, before...) }
}

// PublisherAfter sets the ClientResponseFuncs applied to the incoming SQS
// request prior to it being decoded. This is useful for obtaining the response
// and adding any information onto the context prior to decoding.
func PublisherAfter(after ...PublisherResponseFunc) PublisherOption {
	return func(p *Publisher) { p.after = append(p.after, after...) }
}

// SetPublisherResponseQueueURL can be used as a before function to add
// provided url as responseQueueURL in context.
func SetPublisherResponseQueueURL(url string) PublisherRequestFunc {
	return func(ctx context.Context, _ *sqs.SendMessageInput) context.Context {
		return context.WithValue(ctx, ContextKeyResponseQueueURL, url)
	}
}

// Endpoint returns a usable endpoint that invokes the remote endpoint.
func (p Publisher) Endpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		msgInput := sqs.SendMessageInput{
			QueueUrl: &p.queueURL,
		}
		if err := p.enc(ctx, &msgInput, request); err != nil {
			return nil, err
		}

		for _, f := range p.before {
			ctx = f(ctx, &msgInput)
		}

		output, err := p.Handler.Publish(ctx, &msgInput)
		if err != nil {
			return nil, err
		}

		var responseMsg types.Message
		for _, f := range p.after {
			ctx, responseMsg, err = f(ctx, p.Handler, output)
			if err != nil {
				return nil, err
			}
		}

		response, err := p.dec(ctx, responseMsg)
		if err != nil {
			return nil, err
		}

		return response, nil
	}
}

// EncodeJSONRequest is an EncodeRequestFunc that serializes the request as a
// JSON object and loads it as the MessageBody of the sqs.SendMessageInput.
// This can be enough for most JSON over SQS communications.
func EncodeJSONRequest(_ context.Context, msg *sqs.SendMessageInput, request interface{}) error {
	b, err := json.Marshal(request)
	if err != nil {
		return err
	}

	msg.MessageBody = aws.String(string(b))

	return nil
}

// NoResponseDecode is a DecodeResponseFunc that can be used when no response is needed.
// It returns nil value and nil error.
func NoResponseDecode(_ context.Context, _ types.Message) (interface{}, error) {
	return nil, nil
}
