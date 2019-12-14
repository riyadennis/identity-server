package sqlite

import (
	"database/sql"
	"errors"
	"log"

	"github.com/riyadennis/identity-server/internal/store"
	"github.com/sirupsen/logrus"
)

type IdentityDB struct {
	Db        *sql.DB
	InsertNew *sql.Stmt
	Fetch     *sql.Stmt
}

func Setup() error {
	database, err := sql.Open("sqlite3", "./identity.db")
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

func PrepareDB() *IdentityDB {
	database, err := sql.Open("sqlite3", "./identity.db")
	if err != nil {
		logrus.Fatalf("%v", err)
		return nil
	}
	insert, err := database.Prepare(`INSERT INTO identity_users 
(id, first_name, last_name,
 email, company, post_code, terms) 
 VALUES (?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Fatalf("%v", err)
	}
	fetch, err := database.Prepare("SELECT first_name,last_name, company  FROM identity_users where email = ? and password = ?")
	if err != nil {
		log.Fatalf("%v", err)
	}
	return &IdentityDB{
		Db:        database,
		InsertNew: insert,
		Fetch:     fetch,
	}
}

func (id *IdentityDB) Insert(u *store.User) error {
	_, err := id.InsertNew.Exec(1, u.FirstName, u.LastName, u.Email, u.Company, u.PostCode, u.Terms)
	if err != nil {
		logrus.Errorf("failed to insert user data :: %v", err)
		return err
	}
	return nil
}

func (id *IdentityDB) Read(email, password string) (*store.User, error) {
	rows, err := id.Fetch.Query(email, password)
	if err != nil {
		logrus.Errorf("failed to fetch user data :: %v", err)
		return nil, err
	}
	defer rows.Close()
	var (
		fname, lname, company string
	)

	for rows.Next() {
		u := &store.User{}
		err := rows.Scan(&fname, &lname, &company)
		if err != nil {
			logrus.Errorf("failed to   fetch row:: %v", err)
			return nil, err
		}
		u.FirstName = fname
		u.LastName = lname
		u.Company = company
		return u, nil
	}
	return nil, errors.New("no user found")
}
