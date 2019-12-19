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

func TestOutputStructuredContentFromMiddleware(t *testing.T) {
	handler := testingHandler(t)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	require.Nil(t, err)

	recorder := httptest.NewRecorder()

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

	loggingMiddleware := NewLoggingMiddleware(l)
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
		assert.NotNil(t, http["request_id"])
		assert.NotNil(t, http["request_ip"])
		assert.NotEmpty(t, http["request_method"])
		assert.NotNil(t, http["request_path"])
		assert.NotEmpty(t, http["request_fullpath"])
		assert.NotNil(t, http["request_params"])
		assert.NotNil(t, http["request_user_agent"])
		assert.NotEmpty(t, http["response_status"])
		assert.NotNil(t, http["response_time_total_ms"])
	}

	assertions(fields)
}
