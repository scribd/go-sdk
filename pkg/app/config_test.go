package app

import "testing"

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
