package database

import (
	"github.com/jinzhu/gorm"
	// Imports required gorm MySQL dialect.
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// NewConnection returns a new Gorm database connection.
func NewConnection(config *Config) (*gorm.DB, error) {
	connectionDetails := NewConnectionDetails(config)
	db, err := gorm.Open(connectionDetails.Dialect, connectionDetails.String())
	if err == nil {
		db.DB().SetMaxIdleConns(connectionDetails.Pool)
	}

	return db, err
}
