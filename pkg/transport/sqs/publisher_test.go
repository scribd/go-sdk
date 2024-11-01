package sqs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

const (
	queueURL = "someURL"
)

type testReq struct {
	Squadron int `json:"s"`
}

type testRes struct {
	Squadron int    `json:"s"`
	Name     string `json:"n"`
}

var names = map[int]string{
	424: "tiger",
	426: "thunderbird",
	429: "bison",
	436: "tusker",
	437: "husky",
}

// mockClient is a mock of SQS Handler.
type mockClient struct {
	SQSPublisher
	SQSClient
	err               error
	sendOutputChan    chan types.Message
	receiveOutputChan chan *sqs.ReceiveMessageOutput
	sendMsgID         string
	deleteError       error
}

func (mock *mockClient) Publish(ctx context.Context, input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	if input != nil && input.MessageBody != nil && *input.MessageBody != "" {
		go func() {
			mock.receiveOutputChan <- &sqs.ReceiveMessageOutput{
				Messages: []types.Message{
					{
						MessageAttributes: input.MessageAttributes,
						Body:              input.MessageBody,
						MessageId:         aws.String(mock.sendMsgID),
					},
				},
			}
		}()
		return &sqs.SendMessageOutput{MessageId: aws.String(mock.sendMsgID)}, nil
	}
	// Add logic to allow context errors.
	for {
		select {
		case d := <-mock.sendOutputChan:
			return &sqs.SendMessageOutput{MessageId: d.MessageId}, mock.err
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

// TestBadEncode tests if encode errors are handled properly.
func TestBadEncode(t *testing.T) {
	mock := &mockClient{
		sendOutputChan: make(chan types.Message),
	}
	pub := NewPublisher(
		mock,
		queueURL,
		func(context.Context, *sqs.SendMessageInput, interface{}) error { return errors.New("err!") },
		func(context.Context, types.Message) (response interface{}, err error) { return struct{}{}, nil },
	)
	errChan := make(chan error, 1)
	var err error
	go func() {
		_, pubErr := pub.Endpoint()(context.Background(), struct{}{})
		errChan <- pubErr

	}()
	select {
	case err = <-errChan:
		break

	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for result")
	}
	if err == nil {
		t.Error("expected error")
	}
	if want, have := "err!", err.Error(); want != have {
		t.Errorf("want %s, have %s", want, have)
	}
}

// TestBadDecode tests if decode errors are handled properly.
func TestBadDecode(t *testing.T) {
	mock := &mockClient{
		sendOutputChan: make(chan types.Message),
	}
	go func() {
		mock.sendOutputChan <- types.Message{
			MessageId: aws.String("someMsgID"),
		}
	}()

	pub := NewPublisher(
		mock,
		queueURL,
		func(context.Context, *sqs.SendMessageInput, interface{}) error { return nil },
		func(context.Context, types.Message) (response interface{}, err error) {
			return struct{}{}, errors.New("err!")
		},
		PublisherAfter(func(
			ctx context.Context, _ SQSPublisher, msg *sqs.SendMessageOutput) (context.Context, types.Message, error) {
			// Set the actual response for the request.
			return ctx, types.Message{Body: aws.String("someMsgContent")}, nil
		}),
	)

	var err error
	errChan := make(chan error, 1)
	go func() {
		_, pubErr := pub.Endpoint()(context.Background(), struct{}{})
		errChan <- pubErr
	}()

	select {
	case err = <-errChan:
		break

	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for result")
	}

	if err == nil {
		t.Error("expected error")
	}
	if want, have := "err!", err.Error(); want != have {
		t.Errorf("want %s, have %s", want, have)
	}
}

// TestSuccessfulPublisher ensures that the producer mechanisms work.
func TestSuccessfulPublisher(t *testing.T) {
	mockReq := testReq{437}
	mockRes := testRes{
		Squadron: mockReq.Squadron,
		Name:     names[mockReq.Squadron],
	}
	b, err := json.Marshal(mockRes)
	if err != nil {
		t.Fatal(err)
	}
	mock := &mockClient{
		sendOutputChan: make(chan types.Message),
		sendMsgID:      "someMsgID",
	}
	go func() {
		mock.sendOutputChan <- types.Message{
			MessageId: aws.String("someMsgID"),
		}
	}()

	pub := NewPublisher(
		mock,
		queueURL,
		EncodeJSONRequest,
		func(_ context.Context, msg types.Message) (interface{}, error) {
			response := testRes{}
			if unmarshallErr := json.Unmarshal([]byte(*msg.Body), &response); unmarshallErr != nil {
				return nil, unmarshallErr
			}
			return response, nil
		},
		PublisherAfter(func(
			ctx context.Context, _ SQSPublisher, msg *sqs.SendMessageOutput) (context.Context, types.Message, error) {
			// Sets the actual response for the request.
			if *msg.MessageId == "someMsgID" {
				return ctx, types.Message{Body: aws.String(string(b))}, nil
			}
			return nil, types.Message{}, fmt.Errorf("Did not receive expected SendMessageOutput")
		}),
	)
	var res testRes
	var ok bool
	resChan := make(chan interface{}, 1)
	errChan := make(chan error, 1)
	go func() {
		r, pubErr := pub.Endpoint()(context.Background(), mockReq)
		if pubErr != nil {
			errChan <- pubErr
		} else {
			resChan <- r
		}
	}()

	select {
	case response := <-resChan:
		res, ok = response.(testRes)
		if !ok {
			t.Error("failed to assert endpoint response type")
		}
		break

	case err = <-errChan:
		break

	case <-time.After(100 * time.Millisecond):
		t.Fatal("timed out waiting for result")
	}

	if err != nil {
		t.Fatal(err)
	}
	if want, have := mockRes.Name, res.Name; want != have {
		t.Errorf("want %s, have %s", want, have)
	}
}

// TestSuccessfulPublisherNoResponse ensures that the producer response mechanism works.
func TestSuccessfulPublisherNoResponse(t *testing.T) {
	mock := &mockClient{
		sendOutputChan:    make(chan types.Message),
		receiveOutputChan: make(chan *sqs.ReceiveMessageOutput),
		sendMsgID:         "someMsgID",
	}

	pub := NewPublisher(
		mock,
		queueURL,
		EncodeJSONRequest,
		NoResponseDecode,
	)
	var err error
	errChan := make(chan error, 1)
	finishChan := make(chan bool, 1)
	go func() {
		_, pubErr := pub.Endpoint()(context.Background(), struct{}{})
		if pubErr != nil {
			errChan <- pubErr
		} else {
			finishChan <- true
		}
	}()

	select {
	case <-finishChan:
		break
	case err = <-errChan:
		t.Errorf("unexpected error %s", err)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timed out waiting for result")
	}
}

// TestPublisherWithBefore adds a PublisherBefore function that adds responseQueueURL to context,
// and another on that adds it as a message attribute to outgoing message.
// This test ensures that setting multiple before functions work as expected
// and that SetPublisherResponseQueueURL works as expected.
func TestPublisherWithBefore(t *testing.T) {
	mock := &mockClient{
		sendOutputChan:    make(chan types.Message),
		receiveOutputChan: make(chan *sqs.ReceiveMessageOutput),
		sendMsgID:         "someMsgID",
	}

	responseQueueURL := "someOtherURL"
	pub := NewPublisher(
		mock,
		queueURL,
		EncodeJSONRequest,
		NoResponseDecode,
		PublisherBefore(SetPublisherResponseQueueURL(responseQueueURL)),
		PublisherBefore(func(c context.Context, s *sqs.SendMessageInput) context.Context {
			responseQueueURL := c.Value(ContextKeyResponseQueueURL).(string)
			if s.MessageAttributes == nil {
				s.MessageAttributes = make(map[string]types.MessageAttributeValue)
			}
			s.MessageAttributes["responseQueueURL"] = types.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: &responseQueueURL,
			}
			return c
		}),
	)
	var err error
	errChan := make(chan error, 1)
	go func() {
		_, pubErr := pub.Endpoint()(context.Background(), struct{}{})
		if pubErr != nil {
			errChan <- pubErr
		}
	}()

	want := types.MessageAttributeValue{
		DataType:    aws.String("String"),
		StringValue: &responseQueueURL,
	}

	select {
	case receiveOutput := <-mock.receiveOutputChan:
		if len(receiveOutput.Messages) != 1 {
			t.Errorf("published %d messages instead of 1", len(receiveOutput.Messages))
		}
		if have, exists := receiveOutput.Messages[0].MessageAttributes["responseQueueURL"]; !exists {
			t.Errorf("expected MessageAttributes responseQueueURL not found")
		} else if *have.StringValue != responseQueueURL || *have.DataType != "String" {
			t.Errorf("want %v, have %v", want, have)
		}
		break
	case err = <-errChan:
		t.Errorf("unexpected error %s", err)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timed out waiting for result")
	}
}
