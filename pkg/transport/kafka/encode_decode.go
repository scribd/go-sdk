package kafka

import (
	"context"

	"github.com/twmb/franz-go/pkg/kgo"
)

// DecodeRequestFunc extracts a user-domain request object from
// an Kafka message. It is designed to be used in Kafka Subscribers.
type DecodeRequestFunc func(ctx context.Context, msg *kgo.Record) (request interface{}, err error)

// EncodeRequestFunc encodes the passed request object into
// an Kafka message object. It is designed to be used in Kafka Publishers.
type EncodeRequestFunc func(context.Context, *kgo.Record, interface{}) error

// EncodeResponseFunc encodes the passed response object into
// a Kafka message object. It is designed to be used in Kafka Subscribers.
type EncodeResponseFunc func(context.Context, *kgo.Record, interface{}) error

// DecodeResponseFunc extracts a user-domain response object from kafka
// response object. It's designed to be used in kafka publisher, for publisher-side
// endpoints. One straightforward DecodeResponseFunc could be something that
// JSON decodes from the response payload to the concrete response type.
type DecodeResponseFunc func(context.Context, *kgo.Record) (response interface{}, err error)
