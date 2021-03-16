package store

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"

	"github.com/riyadennis/identity-server/business/store/mysql"
	"github.com/riyadennis/identity-server/business/store/sqlite"
)

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
	db, err := mysql.ConnectDB()
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
