package database

import (
	"fmt"
	"strings"

	"github.com/scribd/go-sdk/pkg/configuration/apps"
)

// ConnectionDetails represents database connection details.
type ConnectionDetails struct {
	Dialect  string
	Username string
	Password string
	Host     string
	Port     int
	Database string
	Encoding string
	Timeout  string
	Pool     int
}

// NewConnectionDetails creates a new ConnectionDetails struct from a DB configuration.
func NewConnectionDetails(config apps.Database) ConnectionDetails {
	return ConnectionDetails{
		Dialect:  "mysql",
		Username: config.Username,
		Password: config.Password,
		Host:     config.Host,
		Port:     config.Port,
		Database: config.Database,
		Encoding: "utf8mb4_unicode_ci",
		Timeout:  config.Timeout,
		Pool:     config.Pool,
	}
}

// String builds a connection string from a database Config.
func (cd ConnectionDetails) String() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?%s",
		cd.Username,
		cd.Password,
		cd.Host,
		cd.Port,
		cd.Database,
		cd.opts(),
	)
}

// StringWithoutDB builds a connection string from a database Config without the database.
func (cd ConnectionDetails) StringWithoutDB() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/?%s",
		cd.Username,
		cd.Password,
		cd.Host,
		cd.Port,
		cd.opts(),
	)
}

func (cd ConnectionDetails) opts() string {
	if cd.Timeout == "" {
		cd.Timeout = "1s"
	}

	options := map[string]string{
		"charset":   "utf8",
		"parseTime": "True",
		"loc":       "Local",
		"timeout":   cd.Timeout,
	}

	var opts []string

	for key, value := range options {
		opts = append(opts, fmt.Sprintf("%s=%s", key, value))
	}

	return strings.Join(opts, "&")
}
