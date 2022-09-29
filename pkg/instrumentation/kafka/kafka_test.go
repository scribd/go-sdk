package kafka

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twmb/franz-go/pkg/kgo"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
)

type (
	mockKafkaClient struct {
		*Client
		produceError bool
		ch           chan kgo.Fetch
	}
)

func (c *mockKafkaClient) Produce(ctx context.Context, r *kgo.Record, promise func(*kgo.Record, error)) {
	c.ch <- kgo.Fetch{Topics: []kgo.FetchTopic{{
		Partitions: []kgo.FetchPartition{{
			Records: []*kgo.Record{r},
		}},
	}}}
}

func (c *mockKafkaClient) ProduceSync(ctx context.Context, rs ...*kgo.Record) kgo.ProduceResults {
	res := kgo.ProduceResult{}
	res.Record = rs[0]

	if c.produceError {
		res.Err = fmt.Errorf("error")
	}

	return kgo.ProduceResults{res}
}

func (c *mockKafkaClient) PollRecords(ctx context.Context, num int) kgo.Fetches {
	fetch := <-c.ch

	return kgo.Fetches{fetch}
}

func (c *mockKafkaClient) Close() {
	close(c.ch)
}

func TestNewClient(t *testing.T) {
	mt := mocktracer.Start()
	defer mt.Stop()

	tests := []struct {
		name string
		err  bool
		fn   func(c *Client)
	}{
		{
			name: "produce",
			fn: func(c *Client) {
				c.Produce(context.Background(), &kgo.Record{Topic: "test"}, nil)
			},
		},
		{
			name: "produce sync",
			fn: func(c *Client) {
				c.ProduceSync(context.Background(), &kgo.Record{Topic: "test"})
			},
		},
		{
			name: "produce sync with error",
			fn: func(c *Client) {
				c.ProduceSync(context.Background(), &kgo.Record{Topic: "test"})
			},
			err: true,
		},
		{
			name: "consume",
			fn: func(c *Client) {
				ctx := context.Background()

				c.Produce(ctx, &kgo.Record{Topic: "test"}, nil)

				fetches := c.KafkaClient.PollRecords(context.Background(), 1)

				iter := c.WrapFetchesRecordIter(ctx, fetches.RecordIter())
				for !iter.Done() {
					iter.Next()
				}
			},
		},
	}

	for _, tt := range tests {
		c := WrapClient(&mockKafkaClient{
			ch:           make(chan kgo.Fetch, 1),
			produceError: tt.err,
		}, WithAnalyticsRate(0.1))

		tt.fn(c)
		c.Close()
	}

	spans := mt.FinishedSpans()

	assert.Len(t, spans, 5)

	// produce
	for i := 0; i < 4; i++ {
		s := spans[i]
		assert.Equal(t, "kafka.produce", s.OperationName())
		assert.Equal(t, "kafka", s.Tag(ext.ServiceName))
		assert.Equal(t, "Produce Topic test", s.Tag(ext.ResourceName))
		assert.Equal(t, 0.1, s.Tag(ext.EventSampleRate))
		assert.Equal(t, "queue", s.Tag(ext.SpanType))
		assert.Equal(t, int32(0), s.Tag("partition"))
	}

	s1 := spans[4] // consume
	assert.Equal(t, "kafka.consume", s1.OperationName())
	assert.Equal(t, "kafka", s1.Tag(ext.ServiceName))
	assert.Equal(t, "Consume Topic test", s1.Tag(ext.ResourceName))
	assert.Equal(t, "queue", s1.Tag(ext.SpanType))
	assert.Equal(t, int32(0), s1.Tag("partition"))
}
