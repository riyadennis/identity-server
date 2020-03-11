package sqlM

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

func ConnectDB() (*sql.DB, error) {
	database, err := sql.Open("mysql",
		"root:root@tcp(mysql-development:3306)/identity_db")
	if err != nil {
		logrus.Fatalf("%v", err)
		return nil, err
	}
	return database, nil
}
