package store

import (
	"errors"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	"github.com/sirupsen/logrus"
	"os"
)

func Migrate() error {
	driver, err := mysql.WithInstance(GetStore(), &mysql.Config{})
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

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logrus.Info("no migration to apply")
		} else {
			return err
		}
	}

	return nil
}
