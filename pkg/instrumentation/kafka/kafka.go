package kafka

import (
	"context"
	"math"

	"github.com/twmb/franz-go/pkg/kgo"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type (
	KafkaClient interface {
		Produce(ctx context.Context, r *kgo.Record, promise func(*kgo.Record, error))
		ProduceSync(ctx context.Context, rs ...*kgo.Record) kgo.ProduceResults
		PollRecords(ctx context.Context, num int) kgo.Fetches
		Flush(ctx context.Context) error
		Close()
	}

	Client struct {
		KafkaClient
		cfg  *config
		prev ddtrace.Span
	}

	FetchesRecordIter struct {
		*kgo.FetchesRecordIter
		ctx    context.Context
		client *Client
	}
)

// NewClient calls kgo.NewClient and wraps the resulting Client.
func NewClient(conf []kgo.Opt, opts ...Option) (*Client, error) {
	c, err := kgo.NewClient(conf...)
	if err != nil {
		return nil, err
	}

	return WrapClient(c, opts...), nil
}

// WrapClient wraps kgo.Client so that produced and consumed messages are traced.
func WrapClient(c KafkaClient, opts ...Option) *Client {
	return &Client{
		KafkaClient: c,
		cfg:         newConfig(opts...),
	}
}

func (c *Client) startProducerSpan(ctx context.Context, msg *kgo.Record) ddtrace.Span {
	opts := []tracer.StartSpanOption{
		tracer.ServiceName(c.cfg.producerServiceName),
		tracer.ResourceName("Produce Topic " + msg.Topic),
		tracer.SpanType(ext.SpanTypeMessageProducer),
	}

	if !math.IsNaN(c.cfg.analyticsRate) {
		opts = append(opts, tracer.Tag(ext.EventSampleRate, c.cfg.analyticsRate))
	}

	carrier := NewMessageCarrier(msg)
	if spanctx, err := tracer.Extract(carrier); err == nil {
		opts = append(opts, tracer.ChildOf(spanctx))
	}

	span, _ := tracer.StartSpanFromContext(ctx, "kafka.produce", opts...)
	err := tracer.Inject(span.Context(), carrier)
	// ignoring the error because carrier implements the interface
	if err != nil {
		return span
	}

	return span
}

func (c *Client) startConsumerSpan(ctx context.Context, msg *kgo.Record) ddtrace.Span {
	opts := []tracer.StartSpanOption{
		tracer.ServiceName(c.cfg.consumerServiceName),
		tracer.ResourceName("Consume Topic " + msg.Topic),
		tracer.SpanType(ext.SpanTypeMessageConsumer),
		tracer.Tag("partition", msg.Partition),
		tracer.Tag("offset", msg.Offset),
		tracer.Measured(),
	}

	if !math.IsNaN(c.cfg.analyticsRate) {
		opts = append(opts, tracer.Tag(ext.EventSampleRate, c.cfg.analyticsRate))
	}

	carrier := NewMessageCarrier(msg)
	if spanctx, err := tracer.Extract(carrier); err == nil {
		opts = append(opts, tracer.ChildOf(spanctx))
	}

	span, _ := tracer.StartSpanFromContext(ctx, "kafka.consume", opts...)
	err := tracer.Inject(span.Context(), carrier)
	// ignoring the error because carrier implements the interface
	if err != nil {
		return span
	}

	return span
}

// Produce calls the underlying *kgo.Client.Produce, the request will be traced.
// This function is used for producing message asynchronously.
func (c *Client) Produce(ctx context.Context, msg *kgo.Record, fn func(record *kgo.Record, err error)) {
	span := c.startProducerSpan(ctx, msg)

	c.KafkaClient.Produce(ctx, msg, fn)

	finishSpan(span, msg.Partition, msg.Offset, nil)
}

// ProduceSync calls the underlying *kgo.Client.ProduceSync and traces all results.
func (c *Client) ProduceSync(ctx context.Context, msgs ...*kgo.Record) kgo.ProduceResults {
	spans := make([]ddtrace.Span, len(msgs))
	for i := range msgs {
		spans[i] = c.startProducerSpan(ctx, msgs[i])
	}

	results := c.KafkaClient.ProduceSync(ctx, msgs...)
	for i, span := range spans {
		finishSpan(span, msgs[i].Partition, msgs[i].Offset, results[i].Err)
	}

	return results
}

// WrapFetchesRecordIter wraps the kgo.FetchesRecordIter and links it to the Client.
func (c *Client) WrapFetchesRecordIter(ctx context.Context, i *kgo.FetchesRecordIter) *FetchesRecordIter {
	return &FetchesRecordIter{
		FetchesRecordIter: i,

		ctx:    ctx,
		client: c,
	}
}

// Close calls the underlying *kgo.Close and finishes the remaining span.
func (c *Client) Close() {
	c.KafkaClient.Close()

	if c.prev != nil {
		c.prev.Finish()
		c.prev = nil
	}
}

// Next calls underlying kgo.FetchesRecordIter.Next and traces the message.
func (i *FetchesRecordIter) Next() *kgo.Record {
	if i.client.prev != nil {
		i.client.prev.Finish()
		i.client.prev = nil
	}

	msg := i.FetchesRecordIter.Next()
	i.client.prev = i.client.startConsumerSpan(i.ctx, msg)

	return msg
}

func finishSpan(span ddtrace.Span, partition int32, offset int64, err error) {
	span.SetTag("partition", partition)
	span.SetTag("offset", offset)
	span.Finish(tracer.WithError(err))
}
