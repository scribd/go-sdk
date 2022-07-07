package logger

import (
	"bytes"

	"github.com/scribd/go-sdk/pkg/configuration/apps"
)

// Builder is a Logger builder.
type Builder struct {
	config         apps.Logger
	fields         Fields
	trackingConfig apps.Tracking
}

// NewBuilder initializes a Logger builder with the given configuration.
func NewBuilder(config apps.Logger) *Builder {
	return &Builder{
		config: config,
	}
}

// SetFields adds Fields to the Builder.
func (b *Builder) SetFields(fields Fields) *Builder {
	b.fields = fields
	return b
}

// SetTracking sets the error reporting configuration.
func (b *Builder) SetTracking(trackingConfig apps.Tracking) *Builder {
	b.trackingConfig = trackingConfig
	return b
}

// Build applies the given configuration and returns a Logger instance.
func (b *Builder) Build() (Logger, error) {
	lLogrus, err := newLogrusLogger(b.config)
	if err != nil {
		return nil, err
	}

	logrusEntry := logrusLogEntry{
		entry: lLogrus.WithFields(convertToLogrusFields(b.fields)),
	}

	if b.trackingConfig.SentryDSN != "" {
		if err := logrusEntry.setTracking(b.trackingConfig); err != nil {
			return nil, err
		}
	}

	return &logrusEntry, nil
}

// BuildTestLogger returns a Logger instance that will write into the bytes buffer
// passed as parameter.
// BuildTestLogger is only for testing.
func (b *Builder) BuildTestLogger(out *bytes.Buffer) (Logger, error) {
	lLogrus, err := newTestLogrusLogger(b.config, out)
	if err != nil {
		return nil, err
	}

	return &logrusLogEntry{
		entry: lLogrus.WithFields(convertToLogrusFields(b.fields)),
	}, nil
}
