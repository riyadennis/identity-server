package store

import (
	"database/sql"
	"errors"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
)

// Migrate runs migration on the db specified in the connection
// will create all the tables in the migrations folder
func Migrate(db *sql.DB) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}
	if os.Getenv("MYSQL_DATABASE") == "" {
		return errors.New("no database set in .env")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		os.Getenv("MYSQL_DATABASE"),
		driver)
	if err != nil {
		return err
	}

	_ = m.Up()

	return nil
}
