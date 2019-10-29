package logger

import (
	"github.com/sirupsen/logrus"
)

// Fields is the struct that that stores key/value pairs for structured logs.
type Fields map[string]interface{}

// Set sets a key/value pair in the Fields map.
func (f *Fields) Set(key string, value interface{}) {
	map[string]interface{}(*f)[key] = value
}

func convertToLogrusFields(fields Fields) logrus.Fields {
	logrusFields := logrus.Fields{}
	for index, val := range fields {
		logrusFields[index] = val
	}
	return logrusFields
}
