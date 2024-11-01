package sqs

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"

	"github.com/scribd/go-sdk/pkg/pubsub/pool"
)

type (
	mockSQSClient struct {
		msgs []types.Message
	}
)

func (m *mockSQSClient) ReceiveMessage(
	ctx context.Context,
	params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {
	return &sqs.ReceiveMessageOutput{
		Messages: m.msgs,
	}, nil
}

func Test_Subscriber_Subscribe(t *testing.T) {
	t.Run("all subscribers finished", func(t *testing.T) {
		var nHandlers int64 // atomic
		var executedTimes int64

		c := make(chan int, 6)

		sub := &Subscriber{
			pool:   pool.New(2),
			stopCh: make(chan struct{}),
			client: &mockSQSClient{
				msgs: []types.Message{
					{
						Body: aws.String("1"),
					},
					{
						Body: aws.String("2"),
					},
					{
						Body: aws.String("3"),
					},
					{
						Body: aws.String("4"),
					},
					{
						Body: aws.String("5"),
					},
					{
						Body: aws.String("6"),
					},
				},
			},
			handler: func(msg types.Message) {
				c <- 0

				atomic.AddInt64(&nHandlers, 1)
				defer atomic.AddInt64(&nHandlers, -1)
				atomic.AddInt64(&executedTimes, 1)

				time.Sleep(time.Millisecond * 10)
			},
		}

		_ = sub.Subscribe(context.Background())
		// Make sure all goroutines have started.
		for i := 0; i < cap(c); i++ {
			<-c
		}

		err := sub.Unsubscribe()
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if got := atomic.LoadInt64(&nHandlers); got != 0 {
			t.Errorf("expected 0, got %d", got)
		}

		if got := atomic.LoadInt64(&executedTimes); got != 6 {
			t.Errorf("expected 6, got %d", got)
		}
	})
}
