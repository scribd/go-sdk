package database

import (
	"github.com/DATA-DOG/go-txdb"
	"github.com/jinzhu/gorm"

	// Imports required gorm MySQL dialect.
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// NewConnection returns a new Gorm database connection.
func NewConnection(config *Config, environment string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	connectionDetails := NewConnectionDetails(config)

	switch environment {
	case "test":
		txdb.Register("txdb", connectionDetails.Dialect, connectionDetails.String())
		db, err = gorm.Open(connectionDetails.Dialect, "txdb", "tx_1")
	default:
		db, err = gorm.Open(connectionDetails.Dialect, connectionDetails.String())
		if err == nil {
			db.DB().SetMaxIdleConns(connectionDetails.Pool)
		}
	}

	return db, err
}
