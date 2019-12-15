package sqlite

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
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
	fetch, err := database.Prepare(`SELECT first_name, last_name, 
											company, post_code FROM
										    identity_users where email = ?`)
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
	uid := uuid.New()
	_, err := id.InsertNew.Exec(uid, u.FirstName, u.LastName, u.Email, u.Company, u.PostCode, u.Terms)
	if err != nil {
		logrus.Errorf("failed to insert user data :: %v", err)
		return err
	}
	return nil
}

func (id *IdentityDB) Read(email string) (*store.User, error) {
	rows := id.Fetch.QueryRow(email)
	var u *store.User

	var fname, lname, post, company string
	err := rows.Scan(&fname, &lname, &post, &company)
	if err != nil {
		logrus.Errorf("%v", err)
	}
	if fname != "" {
		u = &store.User{}
		u.FirstName = fname
		u.LastName = lname
		u.Company = company
		u.PostCode = post
	}
	return u, nil
}
