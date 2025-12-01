package database

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/DATA-DOG/go-txdb"
	sqltrace "github.com/DataDog/dd-trace-go/contrib/database/sql/v2"
	gormtrace "github.com/DataDog/dd-trace-go/contrib/gorm.io/gorm.v1/v2"
	mysqldriver "github.com/go-sql-driver/mysql"
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

	sqltrace.Register(driverName, d, sqltrace.WithService(serviceName))
	sqlDB, err := sqltrace.Open(driverName, connectionString)
	if err != nil {
		return nil, err
	}

	databasePoolSettings(sqlDB, config)

	dialector := mysql.New(mysql.Config{Conn: sqlDB})

	gormConfig := &gorm.Config{}
	if config.DisableDefaultGormTransaction {
		gormConfig.SkipDefaultTransaction = true
	}
	if config.CachePreparedStatements {
		gormConfig.PrepareStmt = true
	}

	db, err := gormtrace.Open(dialector, gormConfig, gormtrace.WithService(serviceName))
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
