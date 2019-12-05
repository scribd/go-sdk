package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	sdklogger "git.lo/microservices/sdk/go-sdk/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testingHandlerBody = "Testing Handler Body"

func testingHandler(t *testing.T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(testingHandlerBody))
		require.Nil(t, err)
	})
}

func logger(t *testing.T, b *bytes.Buffer) sdklogger.Logger {
	t.Helper()

	// Inject this "owned" buffer as Output in the logger wrapped by
	// the loggingMiddleware under test.

	config := &sdklogger.Config{
		ConsoleEnabled:    true,
		ConsoleJSONFormat: true,
		ConsoleLevel:      "info",
		FileEnabled:       false,
	}
	l, err := sdklogger.NewBuilder(config).BuildTestLogger(b)
	require.Nil(t, err)

	return l
}

func TestPreviousRequestDoesNotAffectNewOne(t *testing.T) {
	handler := testingHandler(t)

	recorder := httptest.NewRecorder()

	req1, err := http.NewRequest("GET", "http://example.com?param1=42", nil)
	require.Nil(t, err)

	req2, err := http.NewRequest("GET", "http://example.com?param2=100500", nil)
	require.Nil(t, err)

	var buffer bytes.Buffer
	loggingMiddleware := NewLoggingMiddleware(logger(t, &buffer))

	loggingMiddleware.Handler(handler).ServeHTTP(recorder, req1)

	var fields1 map[string]interface{}
	err = json.Unmarshal(buffer.Bytes(), &fields1)
	require.Nil(t, err)

	var http1 map[string]interface{} = (fields1["http"]).(map[string]interface{})
	var params1 map[string]interface{} = (http1["request_params"]).(map[string]interface{})

	assert.NotNil(t, params1["param1"])

	buffer.Reset()

	loggingMiddleware.Handler(handler).ServeHTTP(recorder, req2)

	var fields2 map[string]interface{}
	err = json.Unmarshal(buffer.Bytes(), &fields2)
	require.Nil(t, err)

	var http2 map[string]interface{} = (fields2["http"]).(map[string]interface{})
	var params2 map[string]interface{} = (http2["request_params"]).(map[string]interface{})

	assert.Nil(t, params2["param1"])
	assert.NotNil(t, params2["param2"])
}

func TestOutputStructuredContentFromMiddleware(t *testing.T) {
	handler := testingHandler(t)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	require.Nil(t, err)

	recorder := httptest.NewRecorder()

	// Inject this "owned" buffer as Output in the logger wrapped by
	// the loggingMiddleware under test.
	var buffer bytes.Buffer

	loggingMiddleware := NewLoggingMiddleware(logger(t, &buffer))
	loggingMiddleware.Handler(handler).ServeHTTP(recorder, req)

	expectedBody := testingHandlerBody
	actualBody := recorder.Body.String()
	assert.Equal(t, expectedBody, actualBody)

	expectedCode := http.StatusOK
	actualCode := recorder.Code
	assert.Equal(t, expectedCode, actualCode)

	var fields map[string]interface{}
	err = json.Unmarshal(buffer.Bytes(), &fields)
	require.Nil(t, err)

	assertions := func(fields map[string]interface{}) {
		assert.NotEmpty(t, fields["message"])
		assert.Equal(t, "info", fields["level"])
		assert.NotEmpty(t, fields["timestamp"])
		assert.NotEmpty(t, fields["http"])

		var http map[string]interface{} = (fields["http"]).(map[string]interface{})

		assert.NotNil(t, http["remote_addr"])
		assert.NotEmpty(t, http["request_fullpath"])
		assert.NotNil(t, http["request_id"])
		assert.NotNil(t, http["request_ip"])
		assert.NotEmpty(t, http["request_method"])
		assert.NotNil(t, http["request_params"])
		assert.NotNil(t, http["request_path"])
		assert.NotNil(t, http["request_user_agent"])
		assert.NotEmpty(t, http["response_status"])
		assert.NotNil(t, http["response_time_total_ms"])
	}

	assertions(fields)
}
