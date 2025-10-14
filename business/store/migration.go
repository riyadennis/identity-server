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
	errInvalidDataInDB         = errors.New("invalid data in db")
)

// Migrate runs migration on the db specified in the connection
// will create all the tables in the migrations folder
func Migrate(db *sql.DB, dbName, migrationPath string) error {
	if db == nil {
		return errEmptyDBConnection
	}
	if dbName == "" {
		return errEmptyDatabaseName
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{
		DatabaseName: dbName,
	})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationPath,
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
