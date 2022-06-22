package kafka

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/scribd/go-sdk/pkg/metrics"
)

type (
	commonMetrics struct {
		metrics    metrics.Metrics
		sampleRate float64

		sampleRatesPerMetric map[string]float64
	}

	// Opt applies options to the common metrics.
	Opt interface {
		apply(m *commonMetrics)
	}

	opt struct{ fn func(m *commonMetrics) }

	BrokerMetrics struct {
		commonMetrics
	}

	ProducerMetrics struct {
		commonMetrics
	}

	ConsumerMetrics struct {
		commonMetrics
	}
)

var (
	_ kgo.HookBrokerConnect           = new(BrokerMetrics)
	_ kgo.HookBrokerDisconnect        = new(BrokerMetrics)
	_ kgo.HookBrokerWrite             = new(BrokerMetrics)
	_ kgo.HookBrokerRead              = new(BrokerMetrics)
	_ kgo.HookBrokerThrottle          = new(BrokerMetrics)
	_ kgo.HookProduceBatchWritten     = new(ProducerMetrics)
	_ kgo.HookProduceRecordBuffered   = new(ProducerMetrics)
	_ kgo.HookProduceRecordUnbuffered = new(ProducerMetrics)
	_ kgo.HookFetchBatchRead          = new(ConsumerMetrics)
	_ kgo.HookFetchRecordBuffered     = new(ConsumerMetrics)
	_ kgo.HookFetchRecordUnbuffered   = new(ConsumerMetrics)
)

const (
	defaultSampleRate = 1.0
)

func (o opt) apply(m *commonMetrics) { o.fn(m) }

// WithSampleRate sets the sample rate which will be used when publishing metrics
func WithSampleRate(rate float64) Opt {
	return opt{fn: func(m *commonMetrics) {
		m.sampleRate = rate
	}}
}

// WithSampleRatesPerMetric sets the sample rate per metric name which will be used when publishing metrics
func WithSampleRatesPerMetric(ratesPerMetric map[string]float64) Opt {
	return opt{fn: func(m *commonMetrics) {
		m.sampleRatesPerMetric = ratesPerMetric
	}}
}

func (m *commonMetrics) Incr(name string, tags []string) {
	err := m.metrics.Incr(name, tags, m.rate(name))
	if err != nil {
		// ignore error
		return
	}
}

func (m *commonMetrics) Count(name string, value int64, tags []string) {
	err := m.metrics.Count(name, value, tags, m.rate(name))
	if err != nil {
		// ignore error
		return
	}
}

func (m *commonMetrics) TimeInMilliseconds(name string, value float64, tags []string) {
	err := m.metrics.TimeInMilliseconds(name, value, tags, m.rate(name))
	if err != nil {
		// ignore error
		return
	}
}

func (m *commonMetrics) Histogram(name string, value float64, tags []string) {
	err := m.metrics.Histogram(name, value, tags, m.rate(name))
	if err != nil {
		// ignore error
		return
	}
}

func (m *commonMetrics) rate(metricName string) float64 {
	r := m.sampleRate

	if metricRate, ok := m.sampleRatesPerMetric[metricName]; ok {
		r = metricRate
	}

	// no sample rate provided
	if r == 0 {
		r = defaultSampleRate
	}

	return r
}

func NewConsumerMetrics(m metrics.Metrics, opts ...Opt) *ConsumerMetrics {
	cm := commonMetrics{metrics: m}

	for _, opt := range opts {
		opt.apply(&cm)
	}

	return &ConsumerMetrics{commonMetrics: cm}
}

func (c *ConsumerMetrics) OnFetchRecordUnbuffered(r *kgo.Record, polled bool) {
	c.Incr(
		"kafka_client.consumer.records_unbuffered",
		[]string{},
	)
}

func (c *ConsumerMetrics) OnFetchRecordBuffered(record *kgo.Record) {
	c.Incr(
		"kafka_client.consumer.records_buffered",
		[]string{},
	)
}

func (c *ConsumerMetrics) OnFetchBatchRead(
	meta kgo.BrokerMetadata,
	topic string, partition int32,
	metrics kgo.FetchBatchMetrics) {
	node := strconv.Itoa(int(meta.NodeID))
	c.Count(
		"kafka_client.consumer.fetch_bytes_uncompressed_total",
		int64(metrics.UncompressedBytes),
		mapToStatsdTags(map[string]string{
			"node":      node,
			"topic":     topic,
			"partition": strconv.Itoa(int(partition)),
		}),
	)
	c.Count(
		"kafka_client.consumer.fetch_bytes_compressed_total",
		int64(metrics.CompressedBytes),
		mapToStatsdTags(map[string]string{
			"node":      node,
			"topic":     topic,
			"partition": strconv.Itoa(int(partition)),
		}),
	)
}

func NewProducerMetrics(m metrics.Metrics, opts ...Opt) *ProducerMetrics {
	cm := commonMetrics{metrics: m}

	for _, opt := range opts {
		opt.apply(&cm)
	}

	return &ProducerMetrics{commonMetrics: cm}
}

func (p *ProducerMetrics) OnProduceRecordUnbuffered(record *kgo.Record, err error) {
	if err != nil {
		p.Incr(
			"kafka_client.producer.records_error",
			[]string{},
		)
		return
	}

	p.Incr(
		"kafka_client.producer.records_unbuffered",
		[]string{},
	)
}

func (p *ProducerMetrics) OnProduceRecordBuffered(record *kgo.Record) {
	p.Incr(
		"kafka_client.producer.records_buffered",
		[]string{},
	)
}

func (p *ProducerMetrics) OnProduceBatchWritten(
	meta kgo.BrokerMetadata,
	topic string,
	partition int32,
	metrics kgo.ProduceBatchMetrics) {
	node := strconv.Itoa(int(meta.NodeID))
	p.Count(
		"kafka_client.producer.produce_bytes_uncompressed_total",
		int64(metrics.UncompressedBytes),
		mapToStatsdTags(map[string]string{
			"node":      node,
			"topic":     topic,
			"partition": strconv.Itoa(int(partition)),
		}),
	)
	p.Count(
		"kafka_client.producer.produce_bytes_compressed_total",
		int64(metrics.CompressedBytes),
		mapToStatsdTags(map[string]string{
			"node":      node,
			"topic":     topic,
			"partition": strconv.Itoa(int(partition)),
		}),
	)
}

func NewBrokerMetrics(m metrics.Metrics, opts ...Opt) *BrokerMetrics {
	cm := commonMetrics{metrics: m}

	for _, opt := range opts {
		opt.apply(&cm)
	}

	return &BrokerMetrics{commonMetrics: cm}
}

func (b *BrokerMetrics) OnBrokerThrottle(
	meta kgo.BrokerMetadata,
	throttleInterval time.Duration,
	throttledAfterResponse bool) {
	node := strconv.Itoa(int(meta.NodeID))
	b.Histogram(
		"kafka_client.broker.throttle_latency",
		throttleInterval.Seconds(),
		mapToStatsdTags(map[string]string{
			"node": node,
		}),
	)
}

func (b *BrokerMetrics) OnBrokerRead(
	meta kgo.BrokerMetadata,
	key int16,
	bytesRead int,
	readWait,
	timeToRead time.Duration,
	err error) {
	node := strconv.Itoa(int(meta.NodeID))

	if err != nil {
		b.Incr(
			"kafka_client.broker.read_errors_total",
			mapToStatsdTags(map[string]string{
				"node": node,
			}))
		return
	}

	b.Count(
		"kafka_client.broker.read_bytes_total",
		int64(bytesRead),
		mapToStatsdTags(map[string]string{
			"node": node,
		}))
	b.TimeInMilliseconds(
		"kafka_client.broker.read_wait_latency",
		float64(readWait.Milliseconds()),
		mapToStatsdTags(map[string]string{
			"node": node,
		}))
	b.TimeInMilliseconds(
		"kafka_client.broker.read_latency",
		float64(timeToRead.Milliseconds()),
		mapToStatsdTags(map[string]string{
			"node": node,
		}),
	)
}

func (b *BrokerMetrics) OnBrokerWrite(
	meta kgo.BrokerMetadata,
	key int16,
	bytesWritten int,
	writeWait,
	timeToWrite time.Duration,
	err error) {
	node := strconv.Itoa(int(meta.NodeID))

	if err != nil {
		b.Incr(
			"kafka_client.broker.write_errors_total",
			mapToStatsdTags(map[string]string{
				"node": node,
			}))
		return
	}

	b.Count(
		"kafka_client.broker.write_bytes_total",
		int64(bytesWritten),
		mapToStatsdTags(map[string]string{
			"node": node,
		}))
	b.TimeInMilliseconds(
		"kafka_client.broker.write_wait_latency",
		float64(writeWait.Milliseconds()),
		mapToStatsdTags(map[string]string{
			"node": node,
		}))
	b.TimeInMilliseconds(
		"kafka_client.broker.write_latency",
		float64(timeToWrite.Milliseconds()),
		mapToStatsdTags(map[string]string{
			"node": node,
		}),
	)
}

func (b *BrokerMetrics) OnBrokerDisconnect(meta kgo.BrokerMetadata, conn net.Conn) {
	node := strconv.Itoa(int(meta.NodeID))
	b.Incr(
		"kafka_client.broker.disconnects_total",
		mapToStatsdTags(map[string]string{
			"node": node,
		}))
}

func (b *BrokerMetrics) OnBrokerConnect(meta kgo.BrokerMetadata, dialDur time.Duration, conn net.Conn, err error) {
	node := strconv.Itoa(int(meta.NodeID))

	if err != nil {
		b.Incr(
			"kafka_client.broker.connect_errors_total",
			mapToStatsdTags(map[string]string{
				"node": node,
			}))
		return
	}

	b.Incr(
		"kafka_client.broker.connects_total",
		mapToStatsdTags(map[string]string{
			"node": node,
		}))
}

func mapToStatsdTags(m map[string]string) []string {
	var tags []string
	for k, v := range m {
		tags = append(tags, fmt.Sprintf("%s:%s", k, v))
	}

	return tags
}
