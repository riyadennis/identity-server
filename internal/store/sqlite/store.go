package sqlite

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

func Setup(source string) error {
	database, err := sql.Open("sqlite3", source)
	if err != nil {
		logrus.Fatalf("%v", err)
		return err
	}
	prepare, err := database.Prepare(`CREATE TABLE IF NOT EXISTS 
											identity_users (id TEXT PRIMARY KEY, 
											first_name  TEXT, 
											last_name TEXT, 
											email TEXT, 
											password TEXT, 
											company TEXT, 
											post_code TEXT,
											terms INTEGER)`)
	if err != nil {
		logrus.Fatalf("%v", err)
		return err
	}
	_, err = prepare.Exec()
	if err != nil {
		logrus.Fatalf("%v", err)
		return err
	}
	return nil
}

func ConnectDB(source string) (*sql.DB, error) {
	database, err := sql.Open("sqlite3", source)
	if err != nil {
		logrus.Fatalf("%v", err)
		return nil, err
	}
	return database, nil
}
