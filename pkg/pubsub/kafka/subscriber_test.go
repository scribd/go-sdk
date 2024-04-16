package kafka

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/twmb/franz-go/pkg/kerr"
	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/scribd/go-sdk/pkg/instrumentation/kafka"
	sdklogger "github.com/scribd/go-sdk/pkg/logger"
)

type mockKafkaClient struct {
	fetches chan kgo.Fetches
}

func (c *mockKafkaClient) Produce(ctx context.Context, r *kgo.Record, promise func(*kgo.Record, error)) {
}

func (c *mockKafkaClient) ProduceSync(ctx context.Context, rs ...*kgo.Record) kgo.ProduceResults {
	return nil
}

func (c *mockKafkaClient) PollRecords(ctx context.Context, num int) kgo.Fetches {
	return <-c.fetches
}

func (c *mockKafkaClient) Flush(ctx context.Context) error {
	return nil
}

func (c *mockKafkaClient) Close() {
	close(c.fetches)
}

func TestSubscriber_Subscribe(t *testing.T) {
	tests := []struct {
		name           string
		fetches        kgo.Fetches
		wantErr        bool
		msgHandler     func(record *kgo.Record)
		wantTerminated bool
		expectedErrCnt int
	}{
		{
			name: "success",
			fetches: kgo.Fetches{
				{
					Topics: []kgo.FetchTopic{
						{
							Topic: "test",
							Partitions: []kgo.FetchPartition{
								{
									Records: []*kgo.Record{{Key: []byte("test")}},
								},
							},
						},
					},
				},
			},
			msgHandler: func(record *kgo.Record) {
				assert.Equal(t, []byte("test"), record.Key)
			},
		},
		{
			name: "non-fatal error",
			fetches: kgo.Fetches{
				{
					Topics: []kgo.FetchTopic{
						{
							Topic: "test",
							Partitions: []kgo.FetchPartition{
								{
									Err: errors.New("test"),
								},
							},
						},
					},
				},
			},
			expectedErrCnt: 1,
			wantErr:        true,
		},
		{
			name: "non-fatal errors",
			fetches: kgo.Fetches{
				{
					Topics: []kgo.FetchTopic{
						{
							Topic: "test",
							Partitions: []kgo.FetchPartition{
								{
									Err: errors.New("test"),
								},
								{
									Err: errors.New("test2"),
								},
							},
						},
					},
				},
			},
			expectedErrCnt: 2,
		},
		{
			name: "fatal error (terminate subscriber)",
			fetches: kgo.Fetches{
				{
					Topics: []kgo.FetchTopic{
						{
							Topic: "test",
							Partitions: []kgo.FetchPartition{
								{
									Err: kerr.BrokerIDNotRegistered,
								},
							},
						},
					},
				},
			},
			expectedErrCnt: 1,
			wantErr:        true,
			wantTerminated: true,
		},
		{
			name: "fatal error (closed client)",
			fetches: kgo.Fetches{
				{
					Topics: []kgo.FetchTopic{
						{
							Topic: "test",
							Partitions: []kgo.FetchPartition{
								{
									Partition: 0,
									Err:       kgo.ErrClientClosed,
								},
							},
						},
					},
				},
			},
			wantTerminated: true,
		},
		{
			name: "fatal and non-fatal errors",
			fetches: kgo.Fetches{
				{
					Topics: []kgo.FetchTopic{
						{
							Topic: "test",
							Partitions: []kgo.FetchPartition{
								{
									Err: errors.New("test"),
								},
								{
									Err: kerr.BrokerIDNotRegistered,
								},
							},
						},
					},
				},
			},
			wantErr:        true,
			expectedErrCnt: 2,
			wantTerminated: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fetchesChan := make(chan kgo.Fetches, 1)
			fetchesChan <- tt.fetches

			var wg sync.WaitGroup
			wg.Add(1)

			s := &Subscriber{
				mu:         sync.Mutex{},
				consumers:  make(map[string]map[int32]pconsumer),
				numWorkers: 1,
				handler: func(rec *kgo.Record) {
					defer wg.Done()

					tt.msgHandler(rec)
				},
				consumer: kafka.WrapClient(&kafka.Client{
					KafkaClient: &mockKafkaClient{
						fetches: fetchesChan,
					},
				}),
			}

			s.assigned(context.Background(), nil, getAssigns(1))

			errCh := s.Subscribe(context.Background())
			if tt.msgHandler != nil {
				wg.Wait()
			}

			if tt.wantErr {
				i := 0
				for err := range errCh {
					i++
					assert.NotNil(t, err)
					if i == tt.expectedErrCnt && !tt.wantTerminated {
						break
					}
				}
				assert.Equal(t, tt.expectedErrCnt, i)
			}

			if tt.wantTerminated {
				_, ok := <-errCh
				assert.False(t, ok)

				err := s.Unsubscribe()
				assert.Nil(t, err)

				_, ok = <-fetchesChan
				assert.False(t, ok)
			}
		})
	}
}

func TestSubscriber_ConcurrentSubscribe(t *testing.T) {
	tests := []struct {
		name          string
		numWorkers    int
		numPartitions int
	}{
		{
			name:          "1 worker, 1 partition",
			numWorkers:    1,
			numPartitions: 1,
		}, {
			name:          "1 worker, 10 partitions",
			numWorkers:    1,
			numPartitions: 10,
		}, {
			name:          "10 workers, 10 partitions",
			numWorkers:    10,
			numPartitions: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fetches := generateFetches(tt.numPartitions)
			fetchesChan := make(chan kgo.Fetches, len(fetches))
			fetchesChan <- fetches

			var wg sync.WaitGroup
			wg.Add(100)

			var i uint64

			s := &Subscriber{
				mu:         sync.Mutex{},
				consumers:  make(map[string]map[int32]pconsumer),
				numWorkers: tt.numWorkers,
				handler: func(rec *kgo.Record) {
					atomic.AddUint64(&i, 1)
					wg.Done()
				},
				consumer: kafka.WrapClient(&kafka.Client{
					KafkaClient: &mockKafkaClient{
						fetches: fetchesChan,
					},
				}),
			}

			s.assigned(context.Background(), nil, getAssigns(tt.numPartitions))

			s.Subscribe(context.Background())

			wg.Wait()

			assert.Equal(t, 100, int(i))
		})
	}
}

func TestSubscriber_RevokePartition(t *testing.T) {
	fetchesChan := make(chan kgo.Fetches, 1)

	fetch := kgo.Fetches{
		{
			Topics: []kgo.FetchTopic{
				{
					Topic: "test",
					Partitions: []kgo.FetchPartition{
						{
							Records: []*kgo.Record{{Key: []byte("test")}},
						},
					},
				},
			},
		},
	}

	fetchesChan <- fetch

	recv := make(chan struct{}, 1)

	l, err := logger()

	assert.Nil(t, err)

	s := &Subscriber{
		mu:         sync.Mutex{},
		consumers:  make(map[string]map[int32]pconsumer),
		logger:     l,
		numWorkers: 1,
		handler: func(rec *kgo.Record) {
			recv <- struct{}{}
		},
		consumer: kafka.WrapClient(&kafka.Client{
			KafkaClient: &mockKafkaClient{
				fetches: fetchesChan,
			},
		}),
	}

	s.assigned(context.Background(), &kgo.Client{}, getAssigns(1))

	s.Subscribe(context.Background())

	select {
	case <-recv:
	case <-time.After(5 * time.Millisecond):
		t.Fatal("expected to receive a message")
	}

	// make sure we stop the partition consumer
	s.revoked(context.Background(), &kgo.Client{}, getAssigns(1))

	fetchesChan <- fetch

	select {
	case <-recv:
		t.Fatal("expected to consumer to be stopped")
	case <-time.After(5 * time.Millisecond):
	}

	err = s.Unsubscribe()
	assert.Nil(t, err)
}

func getAssigns(numPartitions int) map[string][]int32 {
	assigns := make(map[string][]int32)
	for i := 0; i < numPartitions; i++ {
		assigns["test"] = append(assigns["test"], int32(i))
	}
	return assigns
}

func generateFetches(partitions int) kgo.Fetches {
	fetches := kgo.Fetches{
		{
			Topics: []kgo.FetchTopic{
				{
					Topic: "test",
				},
			},
		},
	}

	numRecords := 100 / partitions
	for i := 0; i < partitions; i++ {
		records := make([]*kgo.Record, numRecords)
		for j := 0; j < numRecords; j++ {
			records[j] = &kgo.Record{
				Value: []byte(fmt.Sprintf("test %d", j)),
			}
		}

		fetches[0].Topics[0].Partitions = append(fetches[0].Topics[0].Partitions, kgo.FetchPartition{
			Partition: int32(i),
			Records:   records,
		})
	}

	return fetches
}

func logger() (sdklogger.Logger, error) {
	var buffer bytes.Buffer
	return sdklogger.NewBuilder(
		&sdklogger.Config{
			ConsoleEnabled:    true,
			ConsoleJSONFormat: true,
			ConsoleLevel:      "info",
			FileEnabled:       false,
		}).BuildTestLogger(&buffer)
}
