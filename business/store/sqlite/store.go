package sqlite

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

func ConnectDB(source string) (*sql.DB, error) {
	database, err := sql.Open("sqlite3", source)
	if err != nil {
		logrus.Fatalf("%v", err)
		return nil, err
	}
	return database, nil
}
