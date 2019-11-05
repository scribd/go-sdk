package database

import (
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestNewConnectionDetails(t *testing.T) {
	config, err := NewConfig()
	if err != nil {
		t.Logf("Cannot create NewConfig: %s", err.Error())
	}
	details := NewConnectionDetails(config)

	assert.Equal(t, details.Dialect, "mysql")
	assert.Equal(t, details.Username, config.Username)
	assert.Equal(t, details.Password, config.Password)
	assert.Equal(t, details.Host, config.Host)
	assert.Equal(t, details.Port, config.Port)
	assert.Equal(t, details.Database, config.Database)
	assert.Equal(t, details.Encoding, "utf8mb4_unicode_ci")
	assert.Equal(t, details.Timeout, config.Timeout)
	assert.Equal(t, details.Pool, config.Pool)
}

func TestString(t *testing.T) {
	cases := []struct {
		name             string
		config           *Config
		connectionString string
		optionsString    string
	}{
		{
			name: "WithAllAttributesPresent",
			config: &Config{
				Host:     "192.168.1.1",
				Port:     8080,
				Username: "john",
				Password: "doe",
				Database: "microlith",
				Timeout:  "10s",
			},
			connectionString: "john:doe@tcp(192.168.1.1:8080)/microlith",
			optionsString:    "timeout=10s",
		},
		{
			name: "WithOneAttributeBlank",
			config: &Config{
				Host:     "192.168.1.1",
				Port:     8080,
				Username: "john",
				Password: "",
				Database: "microlith",
				Timeout:  "10s",
			},
			connectionString: "john:@tcp(192.168.1.1:8080)/microlith",
			optionsString:    "timeout=10s",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			details := NewConnectionDetails(c.config)
			got := details.String()

			assert.Contains(t, got, c.connectionString)
			assert.Contains(t, got, c.optionsString)
		})
	}
}

func TestStringWithoutDB(t *testing.T) {
	cases := []struct {
		name             string
		config           *Config
		connectionString string
		optionsString    string
	}{
		{
			name: "WithAllAttributesPresent",
			config: &Config{
				Host:     "192.168.1.1",
				Port:     8080,
				Username: "john",
				Password: "doe",
				Database: "microlith",
				Timeout:  "10s",
			},
			connectionString: "john:doe@tcp(192.168.1.1:8080)/",
			optionsString:    "timeout=10s",
		},
		{
			name: "WithOneAttributeBlank",
			config: &Config{
				Host:     "192.168.1.1",
				Port:     8080,
				Username: "john",
				Password: "",
				Database: "microlith",
				Timeout:  "10s",
			},
			connectionString: "john:@tcp(192.168.1.1:8080)/",
			optionsString:    "timeout=10s",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			details := NewConnectionDetails(c.config)
			got := details.StringWithoutDB()

			assert.Contains(t, got, c.connectionString)
			assert.Contains(t, got, c.optionsString)
		})
	}
}

func TestOpts(t *testing.T) {
	cases := []struct {
		name      string
		details   ConnectionDetails
		timeout   string
		charset   string
		parseTime string
		loc       string
	}{
		{
			name: "WithPresentTimeout",
			details: ConnectionDetails{
				Timeout: "100s",
			},
			timeout:   "timeout=100s",
			charset:   "charset=utf8",
			parseTime: "parseTime=True",
			loc:       "loc=Local",
		},
		{
			name: "WithBlankTimeout",
			details: ConnectionDetails{
				Timeout: "",
			},
			timeout:   "timeout=1s",
			charset:   "charset=utf8",
			parseTime: "parseTime=True",
			loc:       "loc=Local",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.details.opts()

			assert.Contains(t, got, c.timeout)
			assert.Contains(t, got, c.charset)
			assert.Contains(t, got, c.parseTime)
			assert.Contains(t, got, c.loc)
		})
	}
}
