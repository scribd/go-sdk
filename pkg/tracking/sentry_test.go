package tracking

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const (
	testSentryDSNValid   = "https://aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@0000000.ingest.sentry.io/0000000"
	testSentryDSNInvalid = "invalidSentryDSN"
)

func TestNewSentryHook(t *testing.T) {
	hook, err := NewSentryHook(&Config{})
	assert.NoError(t, err)

	logger := newMockLogger(hook)

	assert.Empty(t, sentry.LastEventID(),
		"eventID must be empty without calling sentry hook levels")

	logger.Errorf("sample message")
	assert.Empty(t, sentry.LastEventID(),
		"eventID must be empty without calling WithError")

	logger.WithError(errors.New("sample error")).Errorf("sample message")
	assert.NotEmpty(t, sentry.LastEventID(),
		"last eventID must be set as expected")
}

func TestSentryHookLevels(t *testing.T) {
	config := Config{}
	hook, err := NewSentryHook(&config)
	assert.NoError(t, err)
	assert.Equal(t, []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel}, hook.Levels())
}

func TestSentryHookDsn(t *testing.T) {
	config := Config{SentryDSN: testSentryDSNValid}
	_, err := NewSentryHook(&config)
	assert.NoError(t, err)
}

func TestSentryHookErrorOnInvalidDsn(t *testing.T) {
	config := Config{SentryDSN: testSentryDSNInvalid}
	_, err := NewSentryHook(&config)
	assert.Error(t, err)
}

func TestSentryHookManualTag(t *testing.T) {
	config := Config{SentryDSN: testSentryDSNValid}
	hook, err := NewSentryHook(&config)
	assert.NoError(t, err)

	key := "testKey"
	value := "testValue"
	hook.tags[key] = value

	assert.NotNil(t, hook.tags[key])
	assert.Equal(t, value, hook.tags[key])

	entry := logrus.Entry{}
	err = hook.Fire(&entry)

	assert.NoError(t, err)
	assert.Equal(t, value, hook.tags[key])
}

func newMockLogger(hook *Hook) *logrus.Logger {
	logger := logrus.New()

	logger.SetOutput(ioutil.Discard)
	logger.Hooks.Add(hook)

	return logger
}
