package logger

import (
	"bytes"
	"testing"

	"git.lo/microservices/sdk/go-sdk/pkg/tracking"

	"github.com/stretchr/testify/assert"
)

var testConfig = &Config{
	ConsoleEnabled:    true,
	ConsoleJSONFormat: withJSON,
	ConsoleLevel:      "trace",
}

func TestNewBuilder(t *testing.T) {
	testCases := []struct {
		name   string
		config *Config
	}{
		{
			name:   "NewBuilderWithConfigFileSetValues",
			config: testConfig,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := NewBuilder(tc.config)

			assert.Equal(t, tc.config.ConsoleEnabled, b.config.ConsoleEnabled)
			assert.Equal(t, tc.config.ConsoleJSONFormat, b.config.ConsoleJSONFormat)
			assert.Equal(t, tc.config.ConsoleLevel, b.config.ConsoleLevel)

			assert.Empty(t, b.config.FileEnabled)
			assert.Empty(t, b.config.FileJSONFormat)
			assert.Empty(t, b.config.FileLevel)
			assert.Empty(t, b.config.FileLocation)
			assert.Empty(t, b.config.FileName)
		})
	}
}

func TestSetFields(t *testing.T) {
	testCases := []struct {
		name   string
		fields Fields
	}{
		{
			name: "SetFieldsWithValues",
			fields: Fields{
				"role": "test",
			},
		},
		{
			name:   "SetFieldsWithoutValues",
			fields: Fields{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := NewBuilder(testConfig).SetFields(tc.fields)

			assert.Equal(t, b.fields, tc.fields)
		})
	}
}

func TestSetTracking(t *testing.T) {
	testCases := []struct {
		name           string
		trackingConfig *tracking.Config
	}{
		{
			name: "SetTrackingIsTrueWithATrackingConfig",
			trackingConfig: &tracking.Config{
				SentryDSN: "https://key@sentry.io/project",
			},
		},
		{
			name:           "SetTrackingIsFalseWithoutATrackingConfig",
			trackingConfig: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := NewBuilder(testConfig).SetTracking(tc.trackingConfig)

			assert.Equal(t, b.trackingConfig, tc.trackingConfig)
		})
	}
}

func TestBuild(t *testing.T) {
	testCases := []struct {
		name           string
		trackingConfig *tracking.Config
		fields         Fields
	}{
		{
			name: "WithATrackingConfigAndFieldsItBuilds",
			trackingConfig: &tracking.Config{
				SentryDSN: "https://thealphanumericsentrydns00000000@a012345.ingest.sentry.io/0000000",
			},
			fields: Fields{
				"role": "test",
			},
		},
		{
			name:           "WithoutATrackingConfigAndEmptyFieldsItBuilds",
			fields:         Fields{},
			trackingConfig: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := NewBuilder(testConfig).SetFields(tc.fields)

			if tc.trackingConfig != nil {
				b.SetTracking(tc.trackingConfig)
			}

			_, err := b.Build()

			assert.Nil(t, err)
			assert.Equal(t, tc.fields, b.fields)
		})
	}
}

func TestBuildTestLogger(t *testing.T) {
	testCases := []struct {
		name string
	}{
		{
			name: "WhenBuildingATestLoggerIsNotTracking",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			_, err := NewBuilder(testConfig).BuildTestLogger(&out)

			assert.Nil(t, err)
		})
	}
}
