package database

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	// Imports required gorm MySQL dialect.
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// NewConnection returns a new Gorm database connection.
func NewConnection(config *Config) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", connectionString(config))
	if err == nil {
		db.DB().SetMaxIdleConns(config.Pool)
	}

	return db, err
}

func connectionString(config *Config) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		options(config.Timeout),
	)
}

func options(timeout string) string {
	if timeout == "" {
		timeout = "1s"
	}

	options := map[string]string{
		"charset":   "utf8",
		"parseTime": "True",
		"loc":       "Local",
		"timeout":   timeout,
	}

	var opts []string

	for key, value := range options {
		opts = append(opts, fmt.Sprintf("%s=%s", key, value))
	}

	return strings.Join(opts, "&")
}
