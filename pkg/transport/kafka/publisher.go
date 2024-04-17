package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/twmb/franz-go/pkg/kgo"
)

const (
	defaultPublisherTimeout = 10 * time.Second
)

// Publisher wraps single Kafka topic for message publishing
// and implements endpoint.Endpoint.
type Publisher struct {
	handler   Handler
	topic     string
	enc       EncodeRequestFunc
	dec       DecodeResponseFunc
	before    []RequestFunc
	after     []PublisherResponseFunc
	deliverer Deliverer
	timeout   time.Duration
}

// NewPublisher constructs a new publisher for a single Kafka topic,
// which implements endpoint.Endpoint.
func NewPublisher(
	handler Handler,
	topic string,
	enc EncodeRequestFunc,
	dec DecodeResponseFunc,
	options ...PublisherOption,
) *Publisher {
	p := &Publisher{
		handler:   handler,
		topic:     topic,
		deliverer: SyncDeliverer,
		enc:       enc,
		dec:       dec,
		timeout:   defaultPublisherTimeout,
	}
	for _, opt := range options {
		opt(p)
	}

	return p
}

// PublisherOption sets an optional parameter for publishers.
type PublisherOption func(publisher *Publisher)

// PublisherBefore sets the RequestFuncs that are applied to the outgoing publisher
// request before it's invoked.
func PublisherBefore(before ...RequestFunc) PublisherOption {
	return func(p *Publisher) {
		p.before = append(p.before, before...)
	}
}

// PublisherAfter adds one or more PublisherResponseFuncs, which are applied to the
// context after successful message publishing.
// This is useful for context-manipulation operations.
func PublisherAfter(after ...PublisherResponseFunc) PublisherOption {
	return func(p *Publisher) {
		p.after = append(p.after, after...)
	}
}

// PublisherDeliverer sets the deliverer function that the Publisher invokes.
func PublisherDeliverer(deliverer Deliverer) PublisherOption {
	return func(p *Publisher) { p.deliverer = deliverer }
}

// PublisherTimeout sets the available timeout for a kafka request.
func PublisherTimeout(timeout time.Duration) PublisherOption {
	return func(p *Publisher) { p.timeout = timeout }
}

// Endpoint returns a usable endpoint that invokes message publishing.
func (p Publisher) Endpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		ctx, cancel := context.WithTimeout(ctx, p.timeout)
		defer cancel()

		msg := &kgo.Record{
			Topic: p.topic,
		}

		if err := p.enc(ctx, msg, request); err != nil {
			return nil, err
		}

		for _, f := range p.before {
			ctx = f(ctx, msg)
		}

		event, err := p.deliverer(ctx, p, msg)
		if err != nil {
			return nil, err
		}

		for _, f := range p.after {
			ctx = f(ctx, event)
		}

		response, err := p.dec(ctx, event)
		if err != nil {
			return nil, err
		}

		return response, nil
	}
}

// Deliverer is invoked by the Publisher to publish the specified Message, and to
// retrieve the appropriate response Event object.
type Deliverer func(
	context.Context,
	Publisher,
	*kgo.Record,
) (*kgo.Record, error)

// SyncDeliverer is a deliverer that publishes the specified message
// and returns the first  object.
// If the context times out while waiting for a reply, an error will be returned.
func SyncDeliverer(ctx context.Context, pub Publisher, msg *kgo.Record) (*kgo.Record, error) {
	results := pub.handler.ProduceSync(ctx, msg)

	if len(results) > 0 && results[0].Err != nil {
		return nil, results[0].Err
	}

	return results[0].Record, nil
}

// AsyncDeliverer delivers the supplied message and
// returns a nil response.
//
// When using this deliverer please ensure that the supplied DecodeResponseFunc and
// PublisherResponseFunc are able to handle nil-type responses.
//
// AsyncDeliverer will produce the message with the context detached due to the fact that actual
// message producing is called asynchronously (another goroutine) and at that time original context might be
// already canceled causing the producer to fail. The detached context will include values attached to the original
// context, but deadline and cancel will be reset. To provide a context for asynchronous deliverer please
// use AsyncDelivererCtx function instead.
func AsyncDeliverer(ctx context.Context, pub Publisher, msg *kgo.Record) (*kgo.Record, error) {
	pub.handler.Produce(detach{ctx: ctx}, msg, nil)

	return nil, nil
}

// AsyncDelivererCtx delivers the supplied message and
// returns a nil response.
//
// When using this deliverer please ensure that the supplied DecodeResponseFunc and
// PublisherResponseFunc are able to handle nil-type responses.
func AsyncDelivererCtx(ctx context.Context, pub Publisher, msg *kgo.Record) (*kgo.Record, error) {
	pub.handler.Produce(ctx, msg, nil)

	return nil, nil
}

// EncodeJSONRequest is an EncodeRequestFunc that serializes the request as a
// JSON object to the Message value.
// Many services can use it as a sensible default.
func EncodeJSONRequest(_ context.Context, msg *kgo.Record, request interface{}) error {
	rawJSON, err := json.Marshal(request)
	if err != nil {
		return err
	}

	msg.Value = rawJSON

	return nil
}

// Handler is a handler interface to make testing possible.
// It is highly recommended to use *kafka.Producer as the interface implementation.
type Handler interface {
	Produce(ctx context.Context, rec *kgo.Record, fn func(record *kgo.Record, err error))
	ProduceSync(ctx context.Context, rs ...*kgo.Record) kgo.ProduceResults
}

type detach struct {
	ctx context.Context
}

func (d detach) Deadline() (time.Time, bool) {
	return time.Time{}, false
}
func (d detach) Done() <-chan struct{} {
	return nil
}
func (d detach) Err() error {
	return nil
}

func (d detach) Value(key interface{}) interface{} {
	return d.ctx.Value(key)
}
