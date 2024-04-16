package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

// TestSubscriberBadDecode checks if decoder errors are handled properly.
func TestSubscriberBadDecode(t *testing.T) {
	errCh := make(chan error, 1)

	sub := NewSubscriber(
		func(context.Context, interface{}) (interface{}, error) { return struct{}{}, nil },
		func(context.Context, *kgo.Record) (interface{}, error) { return nil, errors.New("err!") },
		SubscriberErrorEncoder(createTestErrorEncoder(errCh)),
	)

	sub.ServeMsg(nil)(&kgo.Record{})

	var err error
	select {
	case err = <-errCh:
		break

	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for error")
	}
	if want, have := "err!", err.Error(); want != have {
		t.Errorf("want %s, have %s", want, have)
	}
}

// TestSubscriberBadEndpoint checks if endpoint errors are handled properly.
func TestSubscriberBadEndpoint(t *testing.T) {
	errCh := make(chan error, 1)

	sub := NewSubscriber(
		func(context.Context, interface{}) (interface{}, error) { return nil, errors.New("err!") },
		func(context.Context, *kgo.Record) (interface{}, error) { return struct{}{}, nil },
		SubscriberErrorEncoder(createTestErrorEncoder(errCh)),
	)

	sub.ServeMsg(nil)(&kgo.Record{})

	var err error
	select {
	case err = <-errCh:
		break

	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for error")
	}
	if want, have := "err!", err.Error(); want != have {
		t.Errorf("want %s, have %s", want, have)
	}
}

func TestSubscriberSuccess(t *testing.T) {
	obj := testReq{A: 1}
	b, err := json.Marshal(obj)
	if err != nil {
		t.Fatal(err)
	}

	sub := NewSubscriber(
		testEndpoint,
		testReqDecoder,
		SubscriberAfter(func(ctx context.Context, response interface{}) context.Context {
			res := response.(testRes)
			if res.A != 2 {
				t.Errorf("got wrong result: %d", res.A)
			}

			return ctx
		}),
	)

	// golangci-lint v1.44.2 in linux reports that variable 'b' declared but unused
	// in case of usage of self-invoked function. Probably need an investigation
	handler := sub.ServeMsg(nil)
	handler(&kgo.Record{
		Value: b,
	})
}

func createTestErrorEncoder(ch chan error) ErrorEncoder {
	return func(ctx context.Context, err error, msg *kgo.Record, h Handler) {
		ch <- err
	}
}

func testReqDecoder(_ context.Context, m *kgo.Record) (interface{}, error) {
	var obj testReq
	err := json.Unmarshal(m.Value, &obj)
	return obj, err
}

func testEndpoint(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(testReq)
	if !ok {
		return nil, errors.New("type assertion error")
	}

	res := testRes{
		A: req.A + 1,
	}
	return res, nil
}
