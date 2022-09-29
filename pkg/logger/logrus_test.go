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
	withJSON    bool = true
	withoutJSON bool = false
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
//   - the config enables JSON formatter and the output is a parseable JSON;
//   - the config asks for level "trace" and above and the logger has output;
//   - the formatter in the SDK is customized and the logger correctly use those
//     key fields;
//   - the formatter disables the logrus standard "msg" key in the field and
//     the logger correctly doesn't show it;
func TestInfoLevelWithJSONFields(t *testing.T) {
	messageContent := "test message"
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
	messageContent := "test_message"
	logAndAssertTextFields(
		t,
		logConfigForTest(withoutJSON),
		func(log Logger) {
			log.Infof(messageContent)
		},
		func(fields map[string]string) {
			assert.Empty(t, fields["msg"])
			assert.Equal(t, "info", fields["level"])
			assert.NotEmpty(t, fields[fieldKeyTime])
			assert.Equal(t, messageContent, fields[fieldKeyMsg])
		},
	)
}

func TestLevelConfiguration(t *testing.T) {
	t.Run("RunningWithLoggerInTestEnvironment", func(t *testing.T) {
		expected := "test"
		actual := os.Getenv("APP_ENV")
		assert.Equal(t, expected, actual)
	})

	messageContent := "test message"
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
