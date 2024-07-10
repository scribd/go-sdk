package database

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/DATA-DOG/go-txdb"
	mysqldriver "github.com/go-sql-driver/mysql"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	gormtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorm.io/gorm.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const testEnv = "test"

// NewConnection returns a new instrumented Gorm database connection.
func NewConnection(config *Config, environment, appName string) (*gorm.DB, error) {
	connectionDetails := NewConnectionDetails(config)
	connectionString := connectionDetails.String()
	driverName := connectionDetails.Dialect

	var d driver.Driver
	d = &mysqldriver.MySQLDriver{}

	// Register the test driver.
	if environment == testEnv {
		connector := txdb.New(driverName, connectionString)
		d = connector.Driver()
	}

	serviceName := fmt.Sprintf("%s-mysql", appName)

	sqltrace.Register(driverName, d, sqltrace.WithServiceName(serviceName))
	sqlDB, err := sqltrace.Open(driverName, connectionString)
	if err != nil {
		return nil, err
	}

	databasePoolSettings(sqlDB, config)

	dialector := mysql.New(mysql.Config{Conn: sqlDB})
	db, err := gormtrace.Open(dialector, nil, gormtrace.WithServiceName(serviceName))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func databasePoolSettings(db *sql.DB, config *Config) {
	db.SetMaxIdleConns(config.Pool)
	db.SetMaxOpenConns(config.MaxOpenConnections)
	db.SetConnMaxIdleTime(config.ConnectionMaxIdleTime)
	db.SetConnMaxLifetime(config.ConnectionMaxLifetime)
}
