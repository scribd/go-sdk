package kafka

import (
	"fmt"
	"net"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/twmb/franz-go/pkg/kgo"
)

type (
	mockMetrics struct {
		incr []incrArgs
		cnt  []cntArgs
		ts   []tsArgs
		hs   []hsArgs
	}

	incrArgs struct {
		name string
		tags []string
		rate float64
	}

	cntArgs struct {
		name string
		val  int64
		tags []string
		rate float64
	}

	hsArgs struct {
		name string
		val  float64
		tags []string
		rate float64
	}

	tsArgs struct {
		name string
		val  float64
		tags []string
		rate float64
	}
)

func (m *mockMetrics) Gauge(name string, value float64, tags []string, rate float64) error {
	return nil
}

func (m *mockMetrics) Count(name string, value int64, tags []string, rate float64) error {
	// so we have consistent sorted layout of the slice
	sort.Strings(tags)

	m.cnt = append(m.cnt, cntArgs{
		name: name,
		val:  value,
		tags: tags,
		rate: rate,
	})

	return nil
}

func (m *mockMetrics) Histogram(name string, value float64, tags []string, rate float64) error {
	// so we have consistent sorted layout of the slice
	sort.Strings(tags)

	m.hs = append(m.hs, hsArgs{
		name: name,
		val:  value,
		tags: tags,
		rate: rate,
	})

	return nil
}

func (m *mockMetrics) Distribution(name string, value float64, tags []string, rate float64) error {
	return nil
}

func (m *mockMetrics) Decr(name string, tags []string, rate float64) error {
	return nil
}

func (m *mockMetrics) Incr(name string, tags []string, rate float64) error {
	// so we have consistent sorted layout of the slice
	sort.Strings(tags)

	m.incr = append(m.incr, incrArgs{
		name: name,
		tags: tags,
		rate: rate,
	})

	return nil
}

func (m *mockMetrics) Set(name string, value string, tags []string, rate float64) error {
	return nil
}

func (m *mockMetrics) Timing(name string, value time.Duration, tags []string, rate float64) error {
	return nil
}

func (m *mockMetrics) TimeInMilliseconds(name string, value float64, tags []string, rate float64) error {
	// so we have consistent sorted layout of the slice
	sort.Strings(tags)

	m.ts = append(m.ts, tsArgs{
		name: name,
		val:  value,
		tags: tags,
		rate: rate,
	})

	return nil
}

func (m *mockMetrics) SimpleEvent(title, text string) error {
	return nil
}

func (m *mockMetrics) Close() error {
	return nil
}

func (m *mockMetrics) Flush() error {
	return nil
}

func (m *mockMetrics) SetWriteTimeout(d time.Duration) error {
	return nil
}

func TestNewBrokerMetrics(t *testing.T) {
	tests := []struct {
		name       string
		sampleRate float64

		sampleRatesPerMetric map[string]float64
	}{
		{
			name: "no sample rate specified",
		},
		{
			name:       "sample rate specified",
			sampleRate: 0.5,
		},
		{
			name: "sample rate specified for metric name",
			sampleRatesPerMetric: map[string]float64{
				"kafka_client.broker.read_errors_total":  0.3,
				"kafka_client.broker.write_errors_total": 0.1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// dummy net.Conn to pass down to hooks
			n, _ := net.Pipe()

			ms := &mockMetrics{}

			bm := NewBrokerMetrics(ms, WithSampleRate(tt.sampleRate),
				WithSampleRatesPerMetric(tt.sampleRatesPerMetric))

			bm.OnBrokerThrottle(
				kgo.BrokerMetadata{
					NodeID: 1,
				},
				time.Second,
				false,
			)
			bm.OnBrokerRead(
				kgo.BrokerMetadata{
					NodeID: 1,
				},
				int16(1),
				10,
				time.Millisecond,
				time.Millisecond,
				fmt.Errorf("test"),
			)
			bm.OnBrokerRead(
				kgo.BrokerMetadata{
					NodeID: 1,
				},
				int16(1),
				10,
				time.Millisecond,
				time.Millisecond,
				nil,
			)
			bm.OnBrokerWrite(
				kgo.BrokerMetadata{
					NodeID: 1,
				},
				int16(1),
				10,
				time.Millisecond,
				time.Millisecond,
				fmt.Errorf("test"),
			)
			bm.OnBrokerWrite(
				kgo.BrokerMetadata{
					NodeID: 1,
				},
				int16(1),
				10,
				time.Millisecond,
				time.Millisecond,
				nil,
			)
			bm.OnBrokerDisconnect(
				kgo.BrokerMetadata{
					NodeID: 1,
				},
				n,
			)
			bm.OnBrokerConnect(
				kgo.BrokerMetadata{
					NodeID: 1,
				},
				time.Second,
				n,
				fmt.Errorf("test"),
			)
			bm.OnBrokerConnect(
				kgo.BrokerMetadata{
					NodeID: 1,
				},
				time.Second,
				n,
				nil,
			)

			assert.Equal(t, []incrArgs{
				{
					name: "kafka_client.broker.read_errors_total",
					tags: []string{"node:1"},
					rate: getExpectedSampleRate(
						tt.sampleRatesPerMetric,
						tt.sampleRate,
						"kafka_client.broker.read_errors_total",
					),
				},
				{
					name: "kafka_client.broker.write_errors_total",
					tags: []string{"node:1"},
					rate: getExpectedSampleRate(
						tt.sampleRatesPerMetric,
						tt.sampleRate,
						"kafka_client.broker.write_errors_total",
					),
				},
				{
					name: "kafka_client.broker.disconnects_total",
					tags: []string{"node:1"},
					rate: getExpectedSampleRate(
						tt.sampleRatesPerMetric,
						tt.sampleRate,
						"kafka_client.broker.disconnects_total",
					),
				},
				{
					name: "kafka_client.broker.connect_errors_total",
					tags: []string{"node:1"},
					rate: getExpectedSampleRate(
						tt.sampleRatesPerMetric,
						tt.sampleRate,
						"kafka_client.broker.connect_errors_total",
					),
				},
				{
					name: "kafka_client.broker.connects_total",
					tags: []string{"node:1"},
					rate: getExpectedSampleRate(
						tt.sampleRatesPerMetric,
						tt.sampleRate,
						"kafka_client.broker.connects_total",
					),
				},
			}, ms.incr)
			assert.Equal(t, []cntArgs{
				{
					name: "kafka_client.broker.read_bytes_total",
					val:  10,
					tags: []string{"node:1"},
					rate: getExpectedSampleRate(
						tt.sampleRatesPerMetric,
						tt.sampleRate,
						"kafka_client.broker.read_bytes_total",
					),
				},
				{
					name: "kafka_client.broker.write_bytes_total",
					val:  10,
					tags: []string{"node:1"},
					rate: getExpectedSampleRate(
						tt.sampleRatesPerMetric,
						tt.sampleRate,
						"kafka_client.broker.write_bytes_total",
					),
				},
			}, ms.cnt)
			assert.Equal(t, []hsArgs{
				{
					name: "kafka_client.broker.throttle_latency",
					val:  1,
					tags: []string{"node:1"},
					rate: getExpectedSampleRate(
						tt.sampleRatesPerMetric,
						tt.sampleRate,
						"kafka_client.broker.throttle_latency",
					),
				},
			}, ms.hs)
			assert.Equal(t, []tsArgs{
				{
					name: "kafka_client.broker.read_wait_latency",
					val:  1,
					tags: []string{"node:1"},
					rate: getExpectedSampleRate(
						tt.sampleRatesPerMetric,
						tt.sampleRate,
						"kafka_client.broker.read_wait_latency",
					),
				},
				{
					name: "kafka_client.broker.read_latency",
					val:  1,
					tags: []string{"node:1"},
					rate: getExpectedSampleRate(
						tt.sampleRatesPerMetric,
						tt.sampleRate,
						"kafka_client.broker.read_latency",
					),
				},
				{
					name: "kafka_client.broker.write_wait_latency",
					val:  1,
					tags: []string{"node:1"},
					rate: getExpectedSampleRate(
						tt.sampleRatesPerMetric,
						tt.sampleRate,
						"kafka_client.broker.write_wait_latency",
					),
				},
				{
					name: "kafka_client.broker.write_latency",
					val:  1,
					tags: []string{"node:1"},
					rate: getExpectedSampleRate(
						tt.sampleRatesPerMetric,
						tt.sampleRate,
						"kafka_client.broker.write_latency",
					),
				},
			}, ms.ts)
		})
	}
}

func TestNewConsumerMetrics(t *testing.T) {
	tests := []struct {
		name       string
		sampleRate float64

		sampleRatesPerMetric map[string]float64
	}{
		{
			name: "no sample rate specified",
		},
		{
			name:       "sample rate specified",
			sampleRate: 0.5,
		},
		{
			name: "sample rate specified for metric name",
			sampleRatesPerMetric: map[string]float64{
				"kafka_client.consumer.records_unbuffered": 0.3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &mockMetrics{}

			cm := NewConsumerMetrics(
				ms, WithSampleRate(tt.sampleRate),
				WithSampleRatesPerMetric(tt.sampleRatesPerMetric))

			cm.OnFetchRecordUnbuffered(&kgo.Record{}, false)
			cm.OnFetchRecordBuffered(&kgo.Record{})
			cm.OnFetchBatchRead(
				kgo.BrokerMetadata{
					NodeID: 1,
				},
				"test",
				int32(1),
				kgo.FetchBatchMetrics{
					CompressedBytes:   10,
					UncompressedBytes: 100,
				},
			)
			assert.Equal(t, []incrArgs{
				{
					name: "kafka_client.consumer.records_unbuffered",
					tags: []string{},
					rate: getExpectedSampleRate(
						tt.sampleRatesPerMetric,
						tt.sampleRate,
						"kafka_client.consumer.records_unbuffered",
					),
				},
				{
					name: "kafka_client.consumer.records_buffered",
					tags: []string{},
					rate: getExpectedSampleRate(
						tt.sampleRatesPerMetric,
						tt.sampleRate,
						"kafka_client.consumer.records_buffered",
					),
				},
			}, ms.incr)
			assert.Equal(t, []cntArgs{
				{
					name: "kafka_client.consumer.fetch_bytes_uncompressed_total",
					val:  100, tags: []string{"node:1", "partition:1", "topic:test"},
					rate: getExpectedSampleRate(
						tt.sampleRatesPerMetric,
						tt.sampleRate,
						"kafka_client.consumer.fetch_bytes_uncompressed_total",
					),
				},
				{
					name: "kafka_client.consumer.fetch_bytes_compressed_total",
					val:  10, tags: []string{"node:1", "partition:1", "topic:test"},
					rate: getExpectedSampleRate(
						tt.sampleRatesPerMetric,
						tt.sampleRate,
						"kafka_client.consumer.fetch_bytes_compressed_total",
					),
				},
			}, ms.cnt)
			assert.Empty(t, ms.hs)
			assert.Empty(t, ms.ts)
		})
	}
}

func TestNewProducerMetrics(t *testing.T) {
	tests := []struct {
		name       string
		sampleRate float64

		sampleRatesPerMetric map[string]float64
	}{
		{
			name: "no sample rate specified",
		},
		{
			name:       "sample rate specified",
			sampleRate: 0.5,
		},
		{
			name: "sample rate specified for metric name",
			sampleRatesPerMetric: map[string]float64{
				"kafka_client.producer.records_error": 0.3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ms := &mockMetrics{}

			cm := NewProducerMetrics(ms, WithSampleRate(tt.sampleRate),
				WithSampleRatesPerMetric(tt.sampleRatesPerMetric))

			cm.OnProduceRecordUnbuffered(&kgo.Record{}, fmt.Errorf("test"))
			cm.OnProduceRecordUnbuffered(&kgo.Record{}, nil)
			cm.OnProduceRecordBuffered(&kgo.Record{})
			cm.OnProduceBatchWritten(
				kgo.BrokerMetadata{
					NodeID: 1,
				},
				"test",
				int32(1),
				kgo.ProduceBatchMetrics{
					CompressedBytes:   10,
					UncompressedBytes: 100,
				},
			)
			assert.Equal(t, []incrArgs{
				{
					name: "kafka_client.producer.records_error",
					tags: []string{},
					rate: getExpectedSampleRate(
						tt.sampleRatesPerMetric,
						tt.sampleRate,
						"kafka_client.producer.records_error",
					),
				},
				{
					name: "kafka_client.producer.records_unbuffered",
					tags: []string{},
					rate: getExpectedSampleRate(
						tt.sampleRatesPerMetric,
						tt.sampleRate,
						"kafka_client.producer.records_unbuffered",
					),
				},
				{
					name: "kafka_client.producer.records_buffered",
					tags: []string{},
					rate: getExpectedSampleRate(
						tt.sampleRatesPerMetric,
						tt.sampleRate,
						"kafka_client.producer.records_buffered",
					),
				},
			}, ms.incr)
			assert.Equal(t, []cntArgs{
				{
					name: "kafka_client.producer.produce_bytes_uncompressed_total",
					val:  100,
					tags: []string{"node:1", "partition:1", "topic:test"},
					rate: getExpectedSampleRate(
						tt.sampleRatesPerMetric,
						tt.sampleRate,
						"kafka_client.producer.produce_bytes_uncompressed_total",
					),
				},
				{
					name: "kafka_client.producer.produce_bytes_compressed_total",
					val:  10,
					tags: []string{"node:1", "partition:1", "topic:test"},
					rate: getExpectedSampleRate(
						tt.sampleRatesPerMetric,
						tt.sampleRate,
						"kafka_client.producer.produce_bytes_compressed_total",
					),
				},
			}, ms.cnt)
			assert.Empty(t, ms.hs)
			assert.Empty(t, ms.ts)
		})
	}
}

func getExpectedSampleRate(sampleRates map[string]float64, rate float64, metric string) float64 {
	if r, ok := sampleRates[metric]; ok {
		return r
	}

	if rate != 0 {
		return rate
	}

	return defaultSampleRate
}
