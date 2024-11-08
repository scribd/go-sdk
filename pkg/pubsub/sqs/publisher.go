package sqs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type (
	Publisher struct {
		client   *sqs.Client
		queueURL string
	}
)

func NewPublisher(sqsClient *sqs.Client, queueURL string) *Publisher {
	return &Publisher{
		client:   sqsClient,
		queueURL: queueURL,
	}
}

func (p *Publisher) Publish(ctx context.Context, msg *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	if msg.QueueUrl == nil {
		msg.QueueUrl = &p.queueURL
	}

	return p.client.SendMessage(ctx, msg)
}
