package kafka

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twmb/franz-go/pkg/kgo"

	sdkloggercontext "github.com/scribd/go-sdk/pkg/context/logger"
	sdkrequestidcontext "github.com/scribd/go-sdk/pkg/context/requestid"
	sdklogger "github.com/scribd/go-sdk/pkg/logger"
)

const (
	errString = "err"
)

type (
	mockHandler struct {
		producePromise func(ctx context.Context, record *kgo.Record, err error)
		deliverAfter   time.Duration
	}

	testReq struct {
		A int `json:"a"`
	}
	testRes struct {
		A int `json:"a"`
	}
)

func (h *mockHandler) ProduceSync(ctx context.Context, rs ...*kgo.Record) kgo.ProduceResults {
	var results kgo.ProduceResults
	done := make(chan struct{})

	h.Produce(ctx, rs[0], func(record *kgo.Record, err error) {
		results = append(results, kgo.ProduceResult{Record: record, Err: err})
		close(done)
	})

	<-done
	return results
}

func (h *mockHandler) Produce(ctx context.Context, rec *kgo.Record, promise func(record *kgo.Record, err error)) {
	fn := func(rec *kgo.Record, err error) {
		if promise != nil {
			promise(rec, err)
		}
		if h.producePromise != nil {
			h.producePromise(ctx, rec, err)
		}
	}

	go func() {
		select {
		case <-ctx.Done():
			fn(rec, ctx.Err())
		case <-time.After(h.deliverAfter):
			if fn != nil {
				fn(rec, nil)
			}
		}
	}()
}

// TestBadEncode tests if encode errors are handled properly.
func TestBadEncode(t *testing.T) {
	h := &mockHandler{}
	pub := NewPublisher(
		h,
		"test",
		func(context.Context, *kgo.Record, interface{}) error { return errors.New(errString) },
		func(context.Context, *kgo.Record) (response interface{}, err error) { return struct{}{}, nil },
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
	if want, have := errString, err.Error(); want != have {
		t.Errorf("want %s, have %s", want, have)
	}
}

// TestBadDecode tests if decode errors are handled properly.
func TestBadDecode(t *testing.T) {
	h := &mockHandler{}

	pub := NewPublisher(
		h,
		"test",
		func(context.Context, *kgo.Record, interface{}) error { return nil },
		func(context.Context, *kgo.Record) (response interface{}, err error) {
			return struct{}{}, errors.New(errString)
		},
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
	if want, have := errString, err.Error(); want != have {
		t.Errorf("want %s, have %s", want, have)
	}
}

// TestPublisherTimeout ensures that the publisher timeout mechanism works.
func TestPublisherTimeout(t *testing.T) {
	h := &mockHandler{
		deliverAfter: time.Second,
	}

	pub := NewPublisher(
		h,
		"test",
		func(context.Context, *kgo.Record, interface{}) error { return nil },
		func(context.Context, *kgo.Record) (response interface{}, err error) {
			return struct{}{}, nil
		},
		PublisherTimeout(50*time.Millisecond),
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
		t.Fatal("timed out waiting for result")
	}

	if err == nil {
		t.Error("expected error")
	}
	if want, have := context.DeadlineExceeded.Error(), err.Error(); want != have {
		t.Errorf("want %s, have %s", want, have)
	}
}

func TestSuccessfulPublisher(t *testing.T) {
	mockReq := testReq{1}
	mockRes := testRes{
		A: 1,
	}
	_, err := json.Marshal(mockRes)
	if err != nil {
		t.Fatal(err)
	}
	h := &mockHandler{}

	pub := NewPublisher(
		h,
		"test",
		testReqEncoder,
		testResMessageDecoder,
	)
	var res testRes
	var ok bool
	resChan := make(chan interface{}, 1)
	errChan := make(chan error, 1)
	go func() {
		res, pubErr := pub.Endpoint()(context.Background(), mockReq)
		if pubErr != nil {
			errChan <- pubErr
		} else {
			resChan <- res
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
	if want, have := mockRes.A, res.A; want != have {
		t.Errorf("want %d, have %d", want, have)
	}
}

// TestSendAndForgetPublisher tests that the AsyncDeliverer is working
func TestAsyncPublisher(t *testing.T) {
	finishChan := make(chan struct{})
	contextKey, contextValue := "test", "test"
	h := &mockHandler{
		producePromise: func(ctx context.Context, record *kgo.Record, err error) {
			if err != nil {
				t.Fatal(err)
			}
			if want, have := contextValue, ctx.Value(contextKey); want != have {
				t.Errorf("want %s, have %s", want, have)
			}
			finishChan <- struct{}{}
		},
	}

	pub := NewPublisher(
		h,
		"test",
		func(context.Context, *kgo.Record, interface{}) error { return nil },
		func(ctx context.Context, rec *kgo.Record) (response interface{}, err error) {
			val := ctx.Value(contextValue).(string)
			assert.Equal(t, contextValue, val)

			return struct{}{}, nil
		},
		PublisherDeliverer(AsyncDeliverer),
		PublisherTimeout(50*time.Millisecond),
	)

	var err error
	errChan := make(chan error)
	go func() {
		ctx := context.WithValue(context.Background(), contextKey, contextValue)
		ctx, cancel := context.WithCancel(ctx)
		cancel()
		_, pubErr := pub.Endpoint()(ctx, struct{}{})
		if pubErr != nil {
			errChan <- pubErr
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

func TestSetRequestID(t *testing.T) {
	h := &mockHandler{}

	mockReq := testReq{1}

	pub := NewPublisher(
		h,
		"test",
		func(context.Context, *kgo.Record, interface{}) error { return nil },
		func(context.Context, *kgo.Record) (response interface{}, err error) {
			return struct{}{}, nil
		},
		PublisherBefore(SetRequestID()),
		PublisherDeliverer(func(ctx context.Context, publisher Publisher, message *kgo.Record) (*kgo.Record, error) {
			r, err := sdkrequestidcontext.Extract(ctx)
			require.NotNil(t, r)
			require.Nil(t, err)

			return nil, nil
		}),
	)

	_, err := pub.Endpoint()(context.Background(), mockReq)
	require.Nil(t, err)
}

func TestSetLogger(t *testing.T) {
	h := &mockHandler{}

	mockReq := testReq{1}

	// Inject this "owned" buffer as Output in the logger wrapped by
	// the loggingMiddleware under test.
	var buffer bytes.Buffer

	config := &sdklogger.Config{
		ConsoleEnabled:    true,
		ConsoleJSONFormat: true,
		ConsoleLevel:      "info",
		FileEnabled:       false,
	}
	l, err := sdklogger.NewBuilder(config).BuildTestLogger(&buffer)
	require.Nil(t, err)

	pub := NewPublisher(
		h,
		"test",
		func(context.Context, *kgo.Record, interface{}) error { return nil },
		func(context.Context, *kgo.Record) (response interface{}, err error) {
			return struct{}{}, nil
		},
		PublisherBefore(SetLogger(l)),
		PublisherDeliverer(func(ctx context.Context, publisher Publisher, message *kgo.Record) (*kgo.Record, error) {
			l, ctxErr := sdkloggercontext.Extract(ctx)
			require.NotNil(t, l)
			require.Nil(t, ctxErr)

			l.Infof("test")

			return nil, nil
		}),
		PublisherAfter(func(ctx context.Context, ev *kgo.Record) context.Context {
			l, ctxErr := sdkloggercontext.Extract(ctx)
			require.NotNil(t, l)
			require.Nil(t, ctxErr)

			return ctx
		}),
	)

	_, err = pub.Endpoint()(context.Background(), mockReq)
	require.Nil(t, err)

	var fields map[string]interface{}
	err = json.Unmarshal(buffer.Bytes(), &fields)
	require.Nil(t, err)

	assert.NotEmpty(t, fields["pubsub"])
	assert.NotEmpty(t, fields["dd"])

	var pubsub = (fields["pubsub"]).(map[string]interface{})
	assert.NotNil(t, pubsub["request_id"])
}

func testReqEncoder(_ context.Context, m *kgo.Record, request interface{}) error {
	req, ok := request.(testReq)
	if !ok {
		return errors.New("type assertion failure")
	}
	b, err := json.Marshal(req)
	if err != nil {
		return err
	}
	m.Value = b
	return nil
}

func testResMessageDecoder(_ context.Context, m *kgo.Record) (interface{}, error) {
	return testResDecoder(m.Value)
}

func testResDecoder(b []byte) (interface{}, error) {
	var obj testRes
	err := json.Unmarshal(b, &obj)
	return obj, err
}
