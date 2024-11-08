package sqs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// DecodeRequestFunc extracts a user-domain request object from
// an SQS message object. It is designed to be used in Subscribers.
type DecodeRequestFunc func(context.Context, types.Message) (request interface{}, err error)

// EncodeRequestFunc encodes the passed payload object into
// an SQS message object. It is designed to be used in Publishers.
type EncodeRequestFunc func(context.Context, *sqs.SendMessageInput, interface{}) error

// EncodeResponseFunc encodes the passed response object to
// an SQS message object. It is designed to be used in Subscribers.
type EncodeResponseFunc func(context.Context, *sqs.SendMessageInput, interface{}) error

// DecodeResponseFunc extracts a user-domain response object from
// an SQS message object. It is designed to be used in Publishers.
type DecodeResponseFunc func(context.Context, types.Message) (response interface{}, err error)
