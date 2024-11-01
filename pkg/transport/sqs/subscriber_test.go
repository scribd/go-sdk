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
	testErrMessage = "err!"
)

var (
	errTypeAssertion = errors.New("type assertion error")
)

func (mock *mockClient) ReceiveMessage(
	ctx context.Context, input *sqs.ReceiveMessageInput,
) (*sqs.ReceiveMessageOutput, error) {
	// Add logic to allow context errors.
	for {
		select {
		case d := <-mock.receiveOutputChan:
			return d, mock.err
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (mock *mockClient) DeleteMessage(
	ctx context.Context, input *sqs.DeleteMessageInput,
	optFns ...func(options *sqs.Options)) (*sqs.DeleteMessageOutput, error) {
	return nil, mock.deleteError
}

// TestSubscriberDeleteBefore checks if deleteMessage is set properly using subscriber options.
func TestSubscriberDeleteBefore(t *testing.T) {
	mock := &mockClient{
		sendOutputChan:    make(chan types.Message),
		receiveOutputChan: make(chan *sqs.ReceiveMessageOutput),
		deleteError:       fmt.Errorf("delete err!"),
	}
	errEncoder := SubscriberErrorEncoder(func(
		ctx context.Context, err error, req types.Message, sqsClient SQSClient) {
		publishError := sqsError{
			Err:   err.Error(),
			MsgID: *req.MessageId,
		}
		payload, err := json.Marshal(publishError)
		if err != nil {
			t.Fatal(err)
		}

		publisher := sqsClient.(*mockClient)
		_, err = publisher.Publish(ctx, &sqs.SendMessageInput{
			MessageBody: aws.String(string(payload)),
		})
		if err != nil {
			t.Fatal(err)
		}
	})
	subscriber := NewSubscriber(mock,
		func(context.Context, interface{}) (interface{}, error) { return struct{}{}, nil },
		func(context.Context, types.Message) (interface{}, error) { return nil, nil },
		func(context.Context, *sqs.SendMessageInput, interface{}) error { return nil },
		queueURL,
		errEncoder,
		SubscriberDeleteMessageBefore(),
	)

	err := subscriber.ServeMessage(context.Background())(types.Message{
		Body:      aws.String("MessageBody"),
		MessageId: aws.String("fakeMsgID"),
	})
	if err != nil {
		t.Fatal(err)
	}

	var receiveOutput *sqs.ReceiveMessageOutput
	select {
	case receiveOutput = <-mock.receiveOutputChan:
		break

	case <-time.After(200 * time.Millisecond):
		t.Fatal("Timed out waiting for publishing")
	}
	res, err := decodeSubscriberError(receiveOutput)
	if err != nil {
		t.Fatal(err)
	}
	if want, have := "delete err!", res.Err; want != have {
		t.Errorf("want %s, have %s", want, have)
	}
}

// TestSubscriberBadDecode checks if decoder errors are handled properly.
func TestSubscriberBadDecode(t *testing.T) {
	mock := &mockClient{
		sendOutputChan:    make(chan types.Message),
		receiveOutputChan: make(chan *sqs.ReceiveMessageOutput),
	}
	errEncoder := SubscriberErrorEncoder(func(
		ctx context.Context, err error, req types.Message, sqsClient SQSClient) {
		publishError := sqsError{
			Err:   err.Error(),
			MsgID: *req.MessageId,
		}
		payload, err := json.Marshal(publishError)
		if err != nil {
			t.Fatal(err)
		}

		publisher := sqsClient.(*mockClient)
		_, err = publisher.Publish(ctx, &sqs.SendMessageInput{
			MessageBody: aws.String(string(payload)),
		})
		if err != nil {
			t.Fatal(err)
		}
	})
	subscriber := NewSubscriber(mock,
		func(context.Context, interface{}) (interface{}, error) { return struct{}{}, nil },
		func(context.Context, types.Message) (interface{}, error) { return nil, errors.New(testErrMessage) },
		func(context.Context, *sqs.SendMessageInput, interface{}) error { return nil },
		queueURL,
		errEncoder,
	)

	err := subscriber.ServeMessage(context.Background())(types.Message{
		Body:      aws.String("MessageBody"),
		MessageId: aws.String("fakeMsgID"),
	})
	if err == nil {
		t.Errorf("expected error")
	}

	var receiveOutput *sqs.ReceiveMessageOutput
	select {
	case receiveOutput = <-mock.receiveOutputChan:
		break

	case <-time.After(200 * time.Millisecond):
		t.Fatal("Timed out waiting for publishing")
	}
	res, err := decodeSubscriberError(receiveOutput)
	if err != nil {
		t.Fatal(err)
	}
	if want, have := testErrMessage, res.Err; want != have {
		t.Errorf("want %s, have %s", want, have)
	}
}

// TestSubscriberBadEndpoint checks if endpoint errors are handled properly.
func TestSubscriberBadEndpoint(t *testing.T) {
	mock := &mockClient{
		sendOutputChan:    make(chan types.Message),
		receiveOutputChan: make(chan *sqs.ReceiveMessageOutput),
	}
	errEncoder := SubscriberErrorEncoder(func(
		ctx context.Context, err error, req types.Message, sqsClient SQSClient) {
		publishError := sqsError{
			Err:   err.Error(),
			MsgID: *req.MessageId,
		}
		payload, err := json.Marshal(publishError)
		if err != nil {
			t.Fatal(err)
		}

		publisher := sqsClient.(*mockClient)
		_, err = publisher.Publish(ctx, &sqs.SendMessageInput{
			MessageBody: aws.String(string(payload)),
		})
		if err != nil {
			t.Fatal(err)
		}
	})
	subscriber := NewSubscriber(mock,
		func(context.Context, interface{}) (interface{}, error) { return struct{}{}, errors.New(testErrMessage) },
		func(context.Context, types.Message) (interface{}, error) { return nil, nil },
		func(context.Context, *sqs.SendMessageInput, interface{}) error { return nil },
		queueURL,
		errEncoder,
	)

	err := subscriber.ServeMessage(context.Background())(types.Message{
		Body:      aws.String("MessageBody"),
		MessageId: aws.String("fakeMsgID"),
	})
	if err == nil {
		t.Errorf("expected error")
	}

	var receiveOutput *sqs.ReceiveMessageOutput
	select {
	case receiveOutput = <-mock.receiveOutputChan:
		break

	case <-time.After(200 * time.Millisecond):
		t.Fatal("Timed out waiting for publishing")
	}
	res, err := decodeSubscriberError(receiveOutput)
	if err != nil {
		t.Fatal(err)
	}
	if want, have := testErrMessage, res.Err; want != have {
		t.Errorf("want %s, have %s", want, have)
	}
}

// TestSubscriberSuccess checks if subscriber responds correctly to message.
func TestSubscriberSuccess(t *testing.T) {
	obj := testReq{
		Squadron: 436,
	}
	b, err := json.Marshal(obj)
	if err != nil {
		t.Fatal(err)
	}
	mock := &mockClient{
		sendOutputChan:    make(chan types.Message),
		receiveOutputChan: make(chan *sqs.ReceiveMessageOutput),
	}
	subscriber := NewSubscriber(mock,
		testEndpoint,
		testReqDecoderfunc,
		EncodeJSONResponse,
		queueURL,
		SubscriberAfter(func(
			ctx context.Context, cancel context.CancelFunc, msg types.Message, resp interface{}) context.Context {
			_, err = mock.Publish(context.Background(), &sqs.SendMessageInput{
				MessageBody: msg.Body,
			})
			if err != nil {
				t.Fatal(err)
			}

			return ctx
		}),
	)

	err = subscriber.ServeMessage(context.Background())(types.Message{
		Body:      aws.String(string(b)),
		MessageId: aws.String("fakeMsgID"),
	})
	if err != nil {
		t.Fatal(err)
	}

	var receiveOutput *sqs.ReceiveMessageOutput
	select {
	case receiveOutput = <-mock.receiveOutputChan:
		break

	case <-time.After(200 * time.Millisecond):
		t.Fatal("Timed out waiting for publishing")
	}
	res, err := decodeResponse(receiveOutput)
	if err != nil {
		t.Fatal(err)
	}
	want := testRes{
		Squadron: 436,
	}
	if have := res; want != have {
		t.Errorf("want %v, have %v", want, have)
	}
}

// TestSubscriberSuccessNoReply checks if subscriber processes correctly message
// without sending response.
func TestSubscriberSuccessNoReply(t *testing.T) {
	obj := testReq{
		Squadron: 436,
	}
	b, err := json.Marshal(obj)
	if err != nil {
		t.Fatal(err)
	}
	mock := &mockClient{
		sendOutputChan:    make(chan types.Message),
		receiveOutputChan: make(chan *sqs.ReceiveMessageOutput),
	}
	subscriber := NewSubscriber(mock,
		testEndpoint,
		testReqDecoderfunc,
		EncodeJSONResponse,
		queueURL,
	)

	err = subscriber.ServeMessage(context.Background())(types.Message{
		Body:      aws.String(string(b)),
		MessageId: aws.String("fakeMsgID"),
	})
	if err != nil {
		t.Fatal(err)
	}

	var receiveOutput *sqs.ReceiveMessageOutput
	select {
	case receiveOutput = <-mock.receiveOutputChan:
		t.Errorf("received output when none was expected, have %v", receiveOutput)
		return

	case <-time.After(200 * time.Millisecond):
		// As expected, we did not receive any response from subscriber.
		return
	}
}

// TestSubscriberAfter checks if subscriber after is called as expected.
// Here after is used to transfer some info from received message in response.
func TestSubscriberAfter(t *testing.T) {
	obj1 := testReq{
		Squadron: 436,
	}
	b1, err := json.Marshal(obj1)
	if err != nil {
		t.Fatal(err)
	}
	mock := &mockClient{
		sendOutputChan:    make(chan types.Message),
		receiveOutputChan: make(chan *sqs.ReceiveMessageOutput),
	}
	correlationID := "test"
	msg := types.Message{
		Body:      aws.String(string(b1)),
		MessageId: aws.String("fakeMsgID1"),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"correlationID": {
				DataType:    aws.String("String"),
				StringValue: &correlationID,
			},
		},
	}
	subscriber := NewSubscriber(mock,
		testEndpoint,
		testReqDecoderfunc,
		EncodeJSONResponse,
		queueURL,
		SubscriberAfter(func(
			ctx context.Context, cancel context.CancelFunc, msg types.Message, resp interface{}) context.Context {
			_, pubErr := mock.Publish(ctx, &sqs.SendMessageInput{
				MessageBody:       msg.Body,
				MessageAttributes: msg.MessageAttributes,
			})
			if pubErr != nil {
				t.Fatal(pubErr)
			}

			return ctx
		}),
	)
	ctx := context.Background()
	err = subscriber.ServeMessage(ctx)(msg)
	if err != nil {
		t.Fatal(err)
	}

	var receiveOutput *sqs.ReceiveMessageOutput
	select {
	case receiveOutput = <-mock.receiveOutputChan:
		break

	case <-time.After(200 * time.Millisecond):
		t.Fatal("Timed out waiting for publishing")
	}
	if len(receiveOutput.Messages) != 1 {
		t.Errorf("received %d messages instead of 1", len(receiveOutput.Messages))
	}
	if correlationIDAttribute, exists := receiveOutput.Messages[0].MessageAttributes["correlationID"]; exists {
		if have := correlationIDAttribute.StringValue; *have != correlationID {
			t.Errorf("have %s, want %s", *have, correlationID)
		}
	} else {
		t.Errorf("expected message attribute with key correlationID in response, but it was not found")
	}
}

type sqsError struct {
	Err   string `json:"err"`
	MsgID string `json:"msgID"`
}

func decodeSubscriberError(receiveOutput *sqs.ReceiveMessageOutput) (sqsError, error) {
	receivedError := sqsError{}
	err := json.Unmarshal([]byte(*receiveOutput.Messages[0].Body), &receivedError)
	return receivedError, err
}

func testEndpoint(ctx context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(testReq)
	if !ok {
		return nil, errTypeAssertion
	}
	name, prs := names[req.Squadron]
	if !prs {
		return nil, errors.New("unknown squadron name")
	}
	res := testRes{
		Squadron: req.Squadron,
		Name:     name,
	}
	return res, nil
}

func testReqDecoderfunc(_ context.Context, msg types.Message) (interface{}, error) {
	var obj testReq
	err := json.Unmarshal([]byte(*msg.Body), &obj)
	return obj, err
}

func decodeResponse(receiveOutput *sqs.ReceiveMessageOutput) (interface{}, error) {
	if len(receiveOutput.Messages) != 1 {
		return nil, fmt.Errorf("Error : received %d messages instead of 1", len(receiveOutput.Messages))
	}
	resp := testRes{}
	err := json.Unmarshal([]byte(*receiveOutput.Messages[0].Body), &resp)
	return resp, err
}
