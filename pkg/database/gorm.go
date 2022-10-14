package database

import (
	"strconv"
	"time"

	"github.com/DATA-DOG/go-txdb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const testEnv = "test"

// NewConnection returns a new Gorm database connection.
func NewConnection(config *Config, environment string) (*gorm.DB, error) {
	connectionDetails := NewConnectionDetails(config)

	connectionString := connectionDetails.String()
	driverName := connectionDetails.Dialect

	// Register the test driver and mock driver name and connection string in test environment.
	if environment == testEnv {
		// Using time.Now() as a unique identifier for the test database so that we can call NewConnection()
		// multiple times without getting an error.
		testDriverName := strconv.Itoa(int(time.Now().UnixNano()))

		txdb.Register(testDriverName, connectionDetails.Dialect, connectionString)
		driverName = testDriverName
		connectionString = testDriverName
	}

	dialector := mysql.New(mysql.Config{
		DSN:        connectionString,
		DriverName: driverName,
	})

	db, err := gorm.Open(dialector)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(config.Pool)

	return db, nil
}
