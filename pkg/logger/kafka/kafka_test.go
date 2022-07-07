package kafka

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/scribd/go-sdk/pkg/configuration/apps"
	"github.com/scribd/go-sdk/pkg/logger"
)

func TestNewKafkaLogger(t *testing.T) {
	tests := []struct {
		name          string
		opt           Opt
		level         logger.Level
		expectedLevel kgo.LogLevel
	}{
		{
			name:          "No level provided",
			level:         logger.Error,
			expectedLevel: kgo.LogLevelError,
		},
		{
			name:          "Panic level provided",
			level:         logger.Error,
			opt:           WithLevel(logger.Panic),
			expectedLevel: kgo.LogLevelError,
		},
		{
			name:          "Fatal level provided",
			level:         logger.Error,
			opt:           WithLevel(logger.Fatal),
			expectedLevel: kgo.LogLevelError,
		},
		{
			name:          "Error level provided",
			level:         logger.Error,
			opt:           WithLevel(logger.Error),
			expectedLevel: kgo.LogLevelError,
		},
		{
			name:          "Warn level provided",
			level:         logger.Warn,
			opt:           WithLevel(logger.Warn),
			expectedLevel: kgo.LogLevelWarn,
		},
		{
			name:          "Info level provided",
			level:         logger.Info,
			opt:           WithLevel(logger.Info),
			expectedLevel: kgo.LogLevelInfo,
		},
		{
			name:          "Debug level provided",
			level:         logger.Debug,
			opt:           WithLevel(logger.Debug),
			expectedLevel: kgo.LogLevelDebug,
		},
		{
			name:          "Trace level provided",
			level:         logger.Trace,
			opt:           WithLevel(logger.Trace),
			expectedLevel: kgo.LogLevelDebug,
		},
		{
			name:  "WithLevelFn provided",
			level: logger.Info,
			opt: WithLevelFn(func() kgo.LogLevel {
				return kgo.LogLevelInfo
			}),
			expectedLevel: kgo.LogLevelInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fields logger.Fields
			var buffer bytes.Buffer

			b := logger.NewBuilder(apps.Logger{
				ConsoleEnabled:    true,
				ConsoleJSONFormat: true,
				ConsoleLevel:      string(tt.level),
			})
			l, err := b.BuildTestLogger(&buffer)
			require.Nil(t, err)

			var opts []Opt
			if tt.opt != nil {
				opts = []Opt{tt.opt}
			}
			kl := NewKafkaLogger(l, opts...)

			lev := kl.Level()
			assert.Equal(t, tt.expectedLevel, lev)

			kl.Log(tt.expectedLevel, "test", "key", 1)

			err = json.Unmarshal(buffer.Bytes(), &fields)
			assert.Nil(t, err)

			assert.Equal(t, float64(1), fields["key"])
			assert.Equal(t, fields["message"], "test")
			assert.NotEmpty(t, fields["timestamp"])
		})
	}
}
