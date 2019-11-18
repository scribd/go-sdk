package app

import (
	"strings"
	"testing"
	"time"
)

func TestNewConfig(t *testing.T) {
	cases := []struct {
		name       string
		configName string
		want       bool
	}{
		{
			name:       "ValidFile",
			configName: "valid",
			want:       false,
		},
		{
			name:       "InvalidFile",
			configName: "invalid",
			want:       true,
		},
		{
			name:       "NonExistentFile",
			configName: "nonexsistent",
			want:       true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := NewConfig("testdata", c.configName)
			got := err != nil

			if c.want != got {
				t.Errorf("Expected to get %v, got %v", c.want, got)
			}
		})
	}
}

func TestGetters(t *testing.T) {
	var (
		expectedBool        = true
		expectedFloat64     = 100.59
		expectedInt         = 42
		expectedString      = "foo"
		expectedStringSlice = strings.Join([]string{"foo", "bar"}, ",")
		expectedTime, _     = time.Parse(time.RFC3339, "2009-11-10T23:00:00Z")
	)

	cfg, err := NewConfig("testdata", "getters")
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if actual := cfg.Bool("bool"); actual != expectedBool {
		t.Errorf("Got: %t, expected: %t", actual, expectedBool)
	}

	if actual := cfg.Float64("float64"); actual != expectedFloat64 {
		t.Errorf("Got: %f, expected: %f", actual, expectedFloat64)
	}

	if actual := cfg.Int("int"); actual != expectedInt {
		t.Errorf("Got: %d, expected: %d", actual, expectedInt)
	}

	if actual := cfg.String("string"); actual != expectedString {
		t.Errorf("Got: %s, expected: %s", actual, expectedString)
	}

	actualStrSlice := cfg.StringSlice("string_slice")

	if strings.Join(actualStrSlice, ",") != expectedStringSlice {
		t.Errorf("Got slice: %s, expected: %s", actualStrSlice, expectedStringSlice)
	}

	if actual := cfg.Time("time"); actual != expectedTime {
		t.Errorf("Got time: %s, expected time: %s", actual, expectedTime)
	}
}
