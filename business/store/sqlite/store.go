package sqlite

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

// ConnectDB opens a connection to sqlite
// used mainly for tests
func ConnectDB(source string) (*sql.DB, error) {
	database, err := sql.Open("sqlite3", source)
	if err != nil {
		logrus.Fatalf("%v", err)
		return nil, err
	}
	return database, nil
}
