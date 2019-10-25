package database

import (
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		timeout   string
		charset   string
		parseTime string
		loc       string
	}{
		{
			name:      "WithPresentTimeout",
			input:     "100s",
			timeout:   "timeout=100s",
			charset:   "charset=utf8",
			parseTime: "parseTime=True",
			loc:       "loc=Local",
		},
		{
			name:      "WithBlankTimeout",
			input:     "",
			timeout:   "timeout=1s",
			charset:   "charset=utf8",
			parseTime: "parseTime=True",
			loc:       "loc=Local",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := options(c.input)

			assert.Contains(t, got, c.timeout)
			assert.Contains(t, got, c.charset)
			assert.Contains(t, got, c.parseTime)
			assert.Contains(t, got, c.loc)
		})
	}
}

func TestConnectionString(t *testing.T) {
	cases := []struct {
		name             string
		username         string
		password         string
		host             string
		port             int
		database         string
		timeout          string
		connectionString string
		optionsString    string
	}{
		{
			name:             "WithAllAttributesPresent",
			username:         "john",
			password:         "doe",
			host:             "192.168.1.1",
			port:             8080,
			database:         "microlith",
			timeout:          "10s",
			connectionString: "john:doe@tcp(192.168.1.1:8080)/microlith",
			optionsString:    "timeout=10s",
		},
		{
			name:             "WithOneAttributeBlank",
			username:         "john",
			password:         "",
			host:             "192.168.1.1",
			port:             8080,
			database:         "microlith",
			timeout:          "10s",
			connectionString: "john:@tcp(192.168.1.1:8080)/microlith",
			optionsString:    "timeout=10s",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			config := &Config{
				Username: c.username,
				Password: c.password,
				Host:     c.host,
				Port:     c.port,
				Database: c.database,
				Timeout:  c.timeout,
			}
			got := connectionString(config)

			assert.Contains(t, got, c.connectionString)
			assert.Contains(t, got, c.optionsString)
		})
	}
}
