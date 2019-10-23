package logger

import (
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestFields(t *testing.T) {
	testCases := []struct {
		key      string
		value    string
		initial  Fields
		expected Fields
		name     string
	}{
		{
			initial: Fields{},
			key:     "key",
			value:   "newValue",
			expected: Fields{
				"key": "newValue",
			},
			name: "SetANewKey",
		},
		{
			initial: Fields{
				"key": "oldValue",
			},
			key:   "key",
			value: "newValue",
			expected: Fields{
				"key": "newValue",
			},
			name: "UpdatesAnExistingKey",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.initial
			actual.Set(tc.key, tc.value)
			assert.Equal(t, actual, tc.expected)
		})
	}
}
