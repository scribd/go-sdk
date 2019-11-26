package logger

import (
	"bytes"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	firstKey       = "firstKey"
	firstValue     = "firstValue"
	logKey         = "logKey"
	messageContent = "test message"
	secondKey      = "second_key"
	secondValue    = "second_value"
	withJSON       = true
	withoutJSON    = false
)

func logAndAssertContent(
	t *testing.T,
	config *Config,
	log func(Logger),
	withExpectedContent bool,
) {
	var buffer bytes.Buffer

	b := NewBuilder(config)
	lLogger, err := b.BuildTestLogger(&buffer)
	require.Nil(t, err)

	log(lLogger)

	content := buffer.String()
	withActualContent := content != ""
	assert.Equal(t, withExpectedContent, withActualContent)
}

func logAndAssertJSONFields(
	t *testing.T,
	config *Config,
	log func(Logger),
	assertions func(fields Fields),
) {
	var buffer bytes.Buffer
	var fields Fields

	b := NewBuilder(config)
	lLogger, err := b.BuildTestLogger(&buffer)
	require.Nil(t, err)

	log(lLogger)

	err = json.Unmarshal(buffer.Bytes(), &fields)
	assert.Nil(t, err)

	assertions(fields)
}

func logAndAssertTextFields(
	t *testing.T,
	config *Config,
	log func(Logger),
	assertions func(fields map[string]string),
) {
	var buffer bytes.Buffer

	b := NewBuilder(config)
	lLogger, err := b.BuildTestLogger(&buffer)
	require.Nil(t, err)

	log(lLogger)

	fields := make(map[string]string)
	for _, kv := range strings.Split(strings.TrimRight(buffer.String(), "\n"), " ") {
		if !strings.Contains(kv, "=") {
			continue
		}
		kvArr := strings.Split(kv, "=")
		key := strings.TrimSpace(kvArr[0])
		val := kvArr[1]
		if kvArr[1][0] == '"' {
			var err error
			val, err = strconv.Unquote(val)
			require.NoError(t, err)
		}
		fields[key] = val
	}

	assertions(fields)
}

func logConfigForTest(withJSONFormat bool) *Config {
	return &Config{
		ConsoleEnabled:    true,
		ConsoleJSONFormat: withJSONFormat,
		ConsoleLevel:      "trace",
	}
}

// This test is proving the following:
// - the config enables JSON formatter and the output is a parseable JSON;
// - the config asks for level "trace" and above and the logger has output;
// - the formatter in the SDK is customized and the logger correctly use those
//   key fields;
// - the formatter disables the logrus standard "msg" key in the field and
//   the logger correctly doesn't show it;
func TestInfoLevelWithJSONFields(t *testing.T) {
	logAndAssertJSONFields(
		t,
		logConfigForTest(withJSON),
		func(log Logger) {
			log.Infof(messageContent)
		},
		func(fields Fields) {
			assert.Nil(t, fields["msg"])
			assert.Equal(t, "info", fields["level"])
			assert.NotEmpty(t, fields[fieldKeyTime])
			assert.Equal(t, messageContent, fields[fieldKeyMsg])
		},
	)
}

func TestInfoLevelWithTextFields(t *testing.T) {
	// Using an underscore as separator to simplify the parser
	// in `logAndAssertTextFields`.
	logAndAssertTextFields(
		t,
		logConfigForTest(withoutJSON),
		func(log Logger) {
			log.Infof("messageContent")
		},
		func(fields map[string]string) {
			assert.Empty(t, fields["msg"])
			assert.Equal(t, "info", fields["level"])
			assert.NotEmpty(t, fields[fieldKeyTime])
			assert.Equal(t, "messageContent", fields[fieldKeyMsg])
		},
	)
}

func TestLevelConfiguration(t *testing.T) {
	t.Run("RunningWithLoggerInTestEnvironment", func(t *testing.T) {
		expected := "test"
		actual := os.Getenv("APP_ENV")
		assert.Equal(t, expected, actual)
	})

	testCases := []struct {
		name                string
		config              *Config
		log                 func(log Logger)
		withExpectedContent bool
	}{
		{
			name: "WhenConfigLevelIsTraceDebugIsLogged",
			config: &Config{
				ConsoleEnabled:    true,
				ConsoleJSONFormat: withJSON,
				ConsoleLevel:      "trace",
			},
			log: func(l Logger) {
				l.Debugf(messageContent)
			},
			withExpectedContent: true,
		},
		{
			name: "WhenConfigLevelIsWarnInfoIsNotLogged",
			config: &Config{
				ConsoleEnabled:    true,
				ConsoleJSONFormat: withJSON,
				ConsoleLevel:      "warn",
			},
			log: func(l Logger) {
				l.Infof(messageContent)
			},
			withExpectedContent: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logAndAssertContent(
				t,
				tc.config,
				tc.log,
				tc.withExpectedContent,
			)
		})
	}
}

var logFunc = func(log Logger) {
	log.Infof(messageContent)
}

func buildTestLoggerWithFields(t *testing.T, buffer *bytes.Buffer, fields Fields) Logger {
	b := NewBuilder(logConfigForTest(withJSON))
	lLogger, err := b.BuildTestLogger(buffer)
	require.Nil(t, err)

	lLogger = lLogger.WithFields(fields)

	return lLogger
}

func testMergeFieldsBeforeAndAfter(
	t *testing.T,
	initialFields Fields,
	initialAssertions func(fields Fields),
	finalFields Fields,
	finalAssertions func(fields Fields),
) {
	// Build the logger with the given fields.
	var buffer bytes.Buffer
	lLogger := buildTestLoggerWithFields(t, &buffer, initialFields)
	// Run the logger.
	logFunc(lLogger)
	// Verify the assertions.
	actualFields := Fields{}
	err := json.Unmarshal(buffer.Bytes(), &actualFields)
	assert.Nil(t, err)
	initialAssertions(actualFields)

	// Add the fields to the previous logger.
	lLogger = lLogger.WithFields(finalFields)

	// Run again the logger with the new fields.
	buffer.Reset()
	logFunc(lLogger)

	// Verify the assertions
	actualFields = Fields{}
	err = json.Unmarshal(buffer.Bytes(), &actualFields)
	assert.Nil(t, err)
	finalAssertions(actualFields)
}

func TestMergeFields(t *testing.T) {
	// Add a field to the logger.
	initialFields := Fields{
		logKey: Fields{
			firstKey: firstValue,
		},
	}

	// Assert that the initial set of fields is present.
	initialAssertions := func(fields Fields) {
		assert.Nil(t, fields["msg"])
		assert.Equal(t, "info", fields["level"])
		assert.NotEmpty(t, fields[fieldKeyTime])
		assert.Equal(t, messageContent, fields[fieldKeyMsg])
		assert.Equal(t, firstValue, (fields[logKey]).(map[string]interface{})[firstKey])
	}

	// Add a new field that should be merged in the initial set of fields.
	// Also, override the existing value of the `firstField`.
	finalFields := Fields{
		logKey: Fields{
			secondKey: secondValue,
		},
	}

	// Assert that the new field is added to the initial set of fields.
	finalAssertions := func(fields Fields) {
		assert.Nil(t, fields["msg"])
		assert.Equal(t, "info", fields["level"])
		assert.NotEmpty(t, fields[fieldKeyTime])
		assert.Equal(t, messageContent, fields[fieldKeyMsg])
		assert.Equal(t, firstValue, (fields[logKey]).(map[string]interface{})[firstKey])
		assert.Equal(t, secondValue, (fields[logKey]).(map[string]interface{})[secondKey])
	}

	testMergeFieldsBeforeAndAfter(
		t,
		initialFields,
		initialAssertions,
		finalFields,
		finalAssertions,
	)
}

func TestMergeAndOverrideFields(t *testing.T) {
	// Add a field to the logger.
	initialFields := Fields{
		logKey: Fields{
			firstKey: firstValue,
		},
	}

	// Assert that the initial set of fields is present.
	initialAssertions := func(fields Fields) {
		assert.Nil(t, fields["msg"])
		assert.Equal(t, "info", fields["level"])
		assert.NotEmpty(t, fields[fieldKeyTime])
		assert.Equal(t, messageContent, fields[fieldKeyMsg])
		assert.Equal(t, firstValue, (fields[logKey]).(map[string]interface{})[firstKey])
	}

	// Add a new field that should be merged in the initial set of fields.
	// Also, override the existing value of the `firstField`.
	firstValueUpdated := "first_value_updated"
	finalFields := Fields{
		logKey: Fields{
			firstKey:  firstValueUpdated,
			secondKey: secondValue,
		},
	}

	// Assert that the new field is added to the initial set of fields
	// and that the previous value is overwritten.
	finalAssertions := func(fields Fields) {
		assert.Nil(t, fields["msg"])
		assert.Equal(t, "info", fields["level"])
		assert.NotEmpty(t, fields[fieldKeyTime])
		assert.Equal(t, messageContent, fields[fieldKeyMsg])
		assert.Equal(t, firstValueUpdated, (fields[logKey]).(map[string]interface{})[firstKey])
		assert.Equal(t, secondValue, (fields[logKey]).(map[string]interface{})[secondKey])
	}

	testMergeFieldsBeforeAndAfter(
		t,
		initialFields,
		initialAssertions,
		finalFields,
		finalAssertions,
	)
}
