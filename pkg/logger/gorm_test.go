package logger

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGormLogger(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		cfg      Config
		isEmpty  bool
		expected map[string]interface{}
	}{
		{
			name:    "Empty on fatal log level",
			input:   []interface{}{"sql", "test/test", time.Second, "select * from test;", []interface{}{}, int64(10)},
			isEmpty: true,
			cfg: Config{
				ConsoleEnabled:    true,
				ConsoleJSONFormat: true,
				ConsoleLevel:      "fatal",
			},
		},
		{
			name:    "Empty on error log level",
			input:   []interface{}{"sql", "test/test", time.Second, "select * from test;", []interface{}{}, int64(10)},
			isEmpty: true,
			cfg: Config{
				ConsoleEnabled:    true,
				ConsoleJSONFormat: true,
				ConsoleLevel:      "error",
			},
		},
		{
			name:    "Empty on warning log level",
			input:   []interface{}{"sql", "test/test", time.Second, "select * from test;", []interface{}{}, int64(10)},
			isEmpty: true,
			cfg: Config{
				ConsoleEnabled:    true,
				ConsoleJSONFormat: true,
				ConsoleLevel:      "warning",
			},
		},
		{
			name:    "Empty on info log level",
			input:   []interface{}{"sql", "test/test", time.Second, "select * from test;", []interface{}{}, int64(10)},
			isEmpty: true,
			cfg: Config{
				ConsoleEnabled:    true,
				ConsoleJSONFormat: true,
				ConsoleLevel:      "info",
			},
		},
		{
			name:  "Print database log",
			input: []interface{}{"sql", "test/test", time.Second, "select * from test;", []interface{}{}, int64(10)},
			cfg: Config{
				ConsoleEnabled:    true,
				ConsoleJSONFormat: true,
				ConsoleLevel:      "debug",
			},
			expected: map[string]interface{}{
				"duration":      float64(time.Second),
				"affected_rows": float64(10),
				"file_location": "test/test",
				"values":        "",
			},
		},
		{
			name:  "Empty on non sql input in debug mode",
			input: []interface{}{"test string", "test/test", time.Second, "select * from test;", []interface{}{}, int64(10)},
			cfg: Config{
				ConsoleEnabled:    true,
				ConsoleJSONFormat: true,
				ConsoleLevel:      "debug",
			},
			isEmpty: true,
		},
		{
			name:  "Print log in trace mode on non sql input",
			input: []interface{}{"test string", "test/test", time.Second, "select * from test;", []interface{}{}, int64(10)},
			cfg: Config{
				ConsoleEnabled:    true,
				ConsoleJSONFormat: true,
				ConsoleLevel:      "trace",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fields Fields
			var buffer bytes.Buffer

			b := NewBuilder(&tt.cfg)
			l, err := b.BuildTestLogger(&buffer)
			require.Nil(t, err)

			gl := NewGormLogger(l)
			gl.Print(tt.input...)

			if tt.isEmpty {
				assert.Empty(t, buffer.Bytes())
			} else {
				err = json.Unmarshal(buffer.Bytes(), &fields)
				assert.Nil(t, err)

				if tt.expected != nil {
					assert.Equal(t, tt.expected, fields["sql"])
				} else {
					assert.NotEmpty(t, fields["message"])
				}
			}
		})
	}
}
