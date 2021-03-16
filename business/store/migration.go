package store

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	"os"
)

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
