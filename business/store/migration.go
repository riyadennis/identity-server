package store

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
)

// Migrate runs migration on the db specified in the connection
// will create all the tables in the migrations folder
func Migrate(db *sql.DB, dbName, basePath string) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	if dbName == "" {
		return errors.New("no database set in  config")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+basePath+"migrations",
		dbName,
		driver)
	if err != nil {
		return err
	}

	_ = m.Up()

	return nil
}
