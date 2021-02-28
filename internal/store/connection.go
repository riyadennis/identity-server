package store

import (
	"database/sql"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/riyadennis/identity-server/internal/store/sqlM"
	"github.com/riyadennis/identity-server/internal/store/sqlite"
)

// SetStore set the store interface
func SetStore(db *sql.DB) {
	SQLDB = db
}

// GetStore fetch store interface
func GetStore() *sql.DB {
	return SQLDB
}

func Connect() (*sql.DB, error) {
	var db *sql.DB
	var err error

	db, err = connectMysql()
	if err != nil {
		return nil, err
	}

	if os.Getenv("ENV") == "test" {
		db, err = connectSQLite()
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func connectMysql() (*sql.DB, error) {
	db, err := sqlM.ConnectDB()
	if err != nil {
		return nil, err
	}
	logrus.Infof("MYSQL db details %v", db.Stats())
	return db, nil
}

func connectSQLite() (*sql.DB, error) {
	db, err := sqlite.ConnectDB(viper.GetString("source"))
	if err != nil {
		return nil, err
	}

	logrus.Infof("SQLite db details %#v", db.Stats().MaxOpenConnections)
	return db, nil
}
