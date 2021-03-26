package store

import (
	"database/sql"
	"github.com/riyadennis/identity-server/business/store/mysql"
)

// Connect opens DB connection using evn vars set during application initialisation
func Connect() (*sql.DB, error) {
	db, err := mysql.ConnectDB()
	if err != nil {
		return nil, err
	}

	return db, nil
}
