package database

import (
	"fmt"
	"strconv"
	"time"

	"github.com/DATA-DOG/go-txdb"
	gormtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorm.io/gorm.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

const testEnv = "test"

// NewConnection returns a new instrumented Gorm database connection.
func NewConnection(config *Config, environment, appName string) (*gorm.DB, error) {
	serviceName := fmt.Sprintf("%s-mysql", appName)
	dialector := getDialectorFromConfig(config, environment)

	db, err := gormtrace.Open(dialector, nil, gormtrace.WithServiceName(serviceName))
	if err != nil {
		return nil, err
	}
	if len(config.DBs) > 0 {
		if err := db.Use(dbresolver.Register(
			getDbResolverConfig(config, environment),
		)); err != nil {
			return nil, err
		}
	}

	if err := databasePoolSettings(db, config); err != nil {
		return nil, err
	}

	return db, nil
}

func getDbResolverConfig(config *Config, env string) dbresolver.Config {
	resolverCfg := dbresolver.Config{}
	for _, dbConfig := range config.DBs {
		if dbConfig.Replica {
			resolverCfg.Replicas = []gorm.Dialector{getDialectorFromConfig(&dbConfig, env)}
		} else {
			resolverCfg.Sources = []gorm.Dialector{getDialectorFromConfig(&dbConfig, env)}
		}
	}

	return resolverCfg
}

func getDialectorFromConfig(config *Config, environment string) gorm.Dialector {
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

	return mysql.New(mysql.Config{
		DSN:        connectionString,
		DriverName: driverName,
	})
}

func databasePoolSettings(gormDB *gorm.DB, config *Config) error {
	db, err := gormDB.DB()
	if err != nil {
		return err
	}

	db.SetMaxIdleConns(config.Pool)
	db.SetMaxOpenConns(config.MaxOpenConnections)
	db.SetConnMaxIdleTime(config.ConnectionMaxIdleTime)
	db.SetConnMaxLifetime(config.ConnectionMaxLifetime)

	return nil
}
