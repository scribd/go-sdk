package sqs

import (
	"context"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"

	"github.com/scribd/go-sdk/pkg/pubsub"
	"github.com/scribd/go-sdk/pkg/pubsub/pool"
)

type (
	SQSClient interface {
		ReceiveMessage(
			ctx context.Context,
			params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
	}

	Subscriber struct {
		client      SQSClient
		queueURL    string
		handler     MsgHandler
		maxMessages int
		pool        *pool.Pool
		waitTime    time.Duration
		wg          sync.WaitGroup // tracks active handlers
		stopCh      chan struct{}
	}

	SubscriberConfig struct {
		SQSClient  *sqs.Client
		MsgHandler MsgHandler
		SQSConfig  pubsub.SQS
	}

	MsgHandler func(msg types.Message)
)

const (
	defaultNumWorkers = 1
	maxWaitTime       = 20 * time.Second
)

func NewSubscriber(c SubscriberConfig) *Subscriber {
	workers := c.SQSConfig.Subscriber.Workers
	if workers == 0 {
		workers = defaultNumWorkers
	}

	waitTime := c.SQSConfig.Subscriber.WaitTime
	if waitTime > maxWaitTime {
		waitTime = maxWaitTime
	}

	return &Subscriber{
		client:      c.SQSClient,
		handler:     c.MsgHandler,
		maxMessages: c.SQSConfig.Subscriber.MaxMessages,
		queueURL:    c.SQSConfig.Subscriber.QueueURL,
		pool:        pool.New(workers),
		waitTime:    waitTime,
		stopCh:      make(chan struct{}),
	}
}

func (s *Subscriber) Subscribe(ctx context.Context) chan error {
	ch := make(chan error)

	req := &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(s.queueURL),
		MaxNumberOfMessages:   int32(s.maxMessages),
		MessageAttributeNames: []string{"All"},
		MessageSystemAttributeNames: []types.MessageSystemAttributeName{
			types.MessageSystemAttributeNameAll,
		},
	}
	if s.waitTime > 0 {
		req.WaitTimeSeconds = int32(s.waitTime.Seconds())
	}

	go func() {
		defer close(ch)

		for {
			select {
			case <-s.stopCh:
				return
			default:
				response, err := s.client.ReceiveMessage(ctx, req)
				if err != nil {
					ch <- err

					return
				}

				for _, message := range response.Messages {
					s.wg.Add(1)
					s.pool.Schedule(func() {
						s.handler(message)
						s.wg.Done()
					})
				}
			}
		}
	}()

	return ch
}

func (s *Subscriber) Unsubscribe() error {
	close(s.stopCh)
	s.wg.Wait()

	return nil
}
