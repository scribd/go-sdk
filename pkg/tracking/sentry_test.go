package tracking

import (
	"testing"

	"github.com/scribd/go-sdk/pkg/configuration/apps"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const testSentryDSN = "https://aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@0000000.ingest.sentry.io/0000000"

func TestSentryHookLevels(t *testing.T) {
	config := apps.Tracking{}
	hook, err := NewSentryHook(config)
	assert.NoError(t, err)
	assert.Equal(t, []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel}, hook.Levels())
}

func TestSentryHookDsn(t *testing.T) {
	config := apps.Tracking{SentryDSN: testSentryDSN}
	_, err := NewSentryHook(config)
	assert.NoError(t, err)
}

func TestSentryHookErrorOnInvalidDsn(t *testing.T) {
	config := apps.Tracking{SentryDSN: "invalidSentryDSN"}
	_, err := NewSentryHook(config)
	assert.Error(t, err)
}

func TestSentryHookManualTag(t *testing.T) {
	config := apps.Tracking{SentryDSN: testSentryDSN}
	hook, err := NewSentryHook(config)
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
