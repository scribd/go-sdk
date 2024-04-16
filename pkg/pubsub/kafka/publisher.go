package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"

	sdkkafka "github.com/scribd/go-sdk/pkg/instrumentation/kafka"
)

type (
	Publisher struct {
		producer *sdkkafka.Client
	}
)

const (
	defaultFlushTimeout = time.Second * 10

	publisherServiceNameSuffix = "pubsub-publisher"
)

// NewPublisher is a tiny wrapper around the go-sdk kafka.Client and provides API to Publish kafka messages.
func NewPublisher(c Config, opts ...kgo.Opt) (*Publisher, error) {
	serviceName := fmt.Sprintf("%s-%s", c.ApplicationName, publisherServiceNameSuffix)

	cfg, err := newConfig(c, opts...)
	if err != nil {
		return nil, err
	}

	cfg = append(cfg, []kgo.Opt{
		kgo.ProduceRequestTimeout(c.KafkaConfig.Publisher.WriteTimeout),
		kgo.RecordRetries(c.KafkaConfig.Publisher.MaxAttempts),
	}...)

	producer, err := sdkkafka.NewClient(cfg, sdkkafka.WithServiceName(serviceName))
	if err != nil {
		return nil, err
	}

	return &Publisher{producer: producer}, nil
}

// Publish publishes kgo.Record message.
func (p *Publisher) Publish(ctx context.Context, rec *kgo.Record, fn func(record *kgo.Record, err error)) {
	p.producer.Produce(ctx, rec, fn)
}

// Produce is an alias to Publish to satisfy kafka go-kit transport.
func (p *Publisher) Produce(ctx context.Context, rec *kgo.Record, fn func(record *kgo.Record, err error)) {
	p.Publish(ctx, rec, fn)
}

// ProduceSync publishes kgo.Record messages synchronously.
func (p *Publisher) ProduceSync(ctx context.Context, rs ...*kgo.Record) kgo.ProduceResults {
	return p.producer.ProduceSync(ctx, rs...)
}

// GetKafkaProducer returns underlying kafka.Producer for fine-grained tuning purposes.
func (p *Publisher) GetKafkaProducer() *sdkkafka.Client {
	return p.producer
}

// Stop flushes and waits for outstanding messages and requests to complete delivery.
// It also closes a Producer instance.
func (p *Publisher) Stop(ctx context.Context) error {
	if _, deadlineSet := ctx.Deadline(); !deadlineSet {
		timeoutCtx, cancel := context.WithTimeout(ctx, defaultFlushTimeout)
		defer cancel()

		ctx = timeoutCtx
	}

	err := p.producer.Flush(ctx)
	if err != nil {
		return err
	}

	p.producer.Close()

	return nil
}
