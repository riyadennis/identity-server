package store

import (
	"database/sql"
	"errors"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	"github.com/sirupsen/logrus"
)

var (
	errEmptyDBConnection       = errors.New("empty database connection")
	errEmptyDatabaseName       = errors.New("no database set in  config")
	errMigrationInitialisation = errors.New("failed to initialise migration")
)

// Migrate runs migration on the db specified in the connection
// will create all the tables in the migrations folder
func Migrate(db *sql.DB, dbName, basePath string) error {
	if db == nil {
		return errEmptyDBConnection
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	if dbName == "" {
		return errEmptyDatabaseName
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+basePath+"migrations",
		dbName,
		driver)
	if err != nil {
		logrus.Errorf("failed tp initialise migration: %v", err)
		return errMigrationInitialisation
	}

	err = m.Up()
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
	}

	return nil
}
