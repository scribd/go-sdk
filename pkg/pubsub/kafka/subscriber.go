package kafka

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/twmb/franz-go/pkg/kerr"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kmsg"

	sdkkafka "github.com/scribd/go-sdk/pkg/instrumentation/kafka"
	sdklogger "github.com/scribd/go-sdk/pkg/logger"
	"github.com/scribd/go-sdk/pkg/pubsub/pool"
)

type (
	Subscriber struct {
		logger             sdklogger.Logger
		consumer           *sdkkafka.Client
		autoCommitDisabled bool
		mu                 sync.Mutex
		consumers          map[string]map[int32]pconsumer
		handler            func(rec *kgo.Record)
		numWorkers         int
		maxRecords         int
	}
	MsgHandler func(msg *kgo.Record)
)

const (
	subscriberServiceNameSuffix = "pubsub-subscriber"

	defaultMaxRecords = 10000
)

// NewSubscriber is a tiny wrapper around the sdk kafka.Client and provides API to subscribe to a kafka topic.
func NewSubscriber(c Config, opts ...kgo.Opt) (*Subscriber, error) {
	serviceName := fmt.Sprintf("%s-%s", c.ApplicationName, subscriberServiceNameSuffix)

	cfg, err := newConfig(c, opts...)
	if err != nil {
		return nil, err
	}

	autoCommitDisabled := !c.KafkaConfig.Subscriber.AutoCommit.Enabled
	if autoCommitDisabled {
		cfg = append(cfg, kgo.DisableAutoCommit())
	}
	if c.KafkaConfig.Subscriber.BlockRebalance {
		cfg = append(cfg, kgo.BlockRebalanceOnPoll())
	}

	s := &Subscriber{
		logger:             c.Logger,
		mu:                 sync.Mutex{},
		consumers:          make(map[string]map[int32]pconsumer),
		numWorkers:         c.KafkaConfig.Subscriber.Workers,
		handler:            c.MsgHandler,
		autoCommitDisabled: autoCommitDisabled,
		maxRecords:         c.KafkaConfig.Subscriber.MaxRecords,
	}

	if s.maxRecords == 0 {
		s.maxRecords = defaultMaxRecords
	}

	cfg = append(cfg, []kgo.Opt{
		kgo.ConsumerGroup(c.KafkaConfig.Subscriber.GroupId),
		kgo.ConsumeTopics(c.KafkaConfig.Subscriber.Topic),
		kgo.OnPartitionsLost(s.lost),
		kgo.OnPartitionsRevoked(s.revoked),
		kgo.OnPartitionsAssigned(s.assigned),
	}...)

	client, err := sdkkafka.NewClient(cfg, sdkkafka.WithServiceName(serviceName))
	if err != nil {
		return nil, err
	}
	s.consumer = client

	return s, nil
}

// Subscribe subscribes to a configured topic and reads messages.
// Returns unbuffered channel to inspect possible errors.
func (s *Subscriber) Subscribe(ctx context.Context) chan error {
	ch := make(chan error)

	go func() {
		defer close(ch)

		for {
			fetches := s.consumer.PollRecords(ctx, s.maxRecords)
			if fetches.IsClientClosed() {
				return
			}

			if fetches.Err() != nil {
				var containsFatalErr bool

				fetches.EachError(func(_ string, _ int32, err error) {
					if !containsFatalErr {
						containsFatalErr = isFatalFetchError(err)
					}
					ch <- err
				})

				if containsFatalErr {
					return
				}
			}

			fetches.EachTopic(func(t kgo.FetchTopic) {
				s.mu.Lock()
				tconsumers := s.consumers[t.Topic]
				s.mu.Unlock()

				if tconsumers == nil {
					return
				}
				t.EachPartition(func(p kgo.FetchPartition) {
					pc, ok := tconsumers[p.Partition]
					if !ok {
						return
					}
					select {
					case pc.recs <- s.consumer.WrapFetchPartition(ctx, p):
					case <-pc.quit:
					}
				})
			})
			if kgoClient, ok := s.consumer.KafkaClient.(*kgo.Client); ok {
				// this call does nothing in case rebalance is not blocked
				kgoClient.AllowRebalance()
			}
		}
	}()

	return ch
}

func (s *Subscriber) assigned(_ context.Context, cl *kgo.Client, assigned map[string][]int32) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for topic, partitions := range assigned {
		if s.consumers[topic] == nil {
			s.consumers[topic] = make(map[int32]pconsumer)
		}
		for _, partition := range partitions {
			pc := pconsumer{
				quit: make(chan struct{}),
				recs: make(chan *sdkkafka.FetchPartition),
				pool: pool.New(s.numWorkers),
				done: make(chan struct{}),
			}
			s.consumers[topic][partition] = pc
			go pc.consume(cl, s.logger, s.autoCommitDisabled, s.handler)
		}
	}
}

func (s *Subscriber) stopConsumers(lost map[string][]int32) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var wg sync.WaitGroup
	defer wg.Wait()

	for topic, partitions := range lost {
		ptopics := s.consumers[topic]
		for _, partition := range partitions {
			pc := ptopics[partition]
			delete(ptopics, partition)

			if len(ptopics) == 0 {
				delete(s.consumers, topic)
			}
			close(pc.quit)
			wg.Add(1)
			go func() { <-pc.done; wg.Done() }()
		}
	}
}

func (s *Subscriber) lost(_ context.Context, _ *kgo.Client, lost map[string][]int32) {
	s.stopConsumers(lost)
}

func (s *Subscriber) revoked(ctx context.Context, cl *kgo.Client, lost map[string][]int32) {
	s.stopConsumers(lost)
	if !s.autoCommitDisabled {
		cl.CommitOffsetsSync(ctx, cl.MarkedOffsets(),
			func(cl *kgo.Client, _ *kmsg.OffsetCommitRequest, _ *kmsg.OffsetCommitResponse, err error) {
				if err != nil {
					s.logger.WithError(err).Errorf("Revoke commit failed")
				}
			},
		)
	}
}

func (s *Subscriber) Unsubscribe() error {
	s.consumer.Close()

	return nil
}

func isFatalFetchError(err error) bool {
	var kafkaErr *kerr.Error

	return errors.As(err, &kafkaErr)
}
