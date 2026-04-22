package logger

import (
	"maps"

	"github.com/sirupsen/logrus"
)

// Fields is the struct that that stores key/value pairs for structured logs.
type Fields map[string]any

// Set sets a key/value pair in the Fields map.
func (f *Fields) Set(key string, value any) {
	map[string]any(*f)[key] = value
}

func convertToLogrusFields(fields Fields) logrus.Fields {
	logrusFields := logrus.Fields{}
	maps.Copy(logrusFields, fields)
	return logrusFields
}
