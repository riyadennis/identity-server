package sqlM

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

func ConnectDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s):%s/%s",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_ROOT_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"))
	database, err := sql.Open("mysql", connStr)
	if err != nil {
		logrus.Fatalf("%v", err)
		return nil, err
	}
	return database, nil
}
