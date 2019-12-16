package sqlite

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"

	"github.com/google/uuid"
	"github.com/riyadennis/identity-server/internal/store"
	"github.com/sirupsen/logrus"
)

type LiteDB struct {
	Db        *sql.DB
	InsertNew *sql.Stmt
	Fetch     *sql.Stmt
	Login     *sql.Stmt
}

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

func PrepareDB(source string) *LiteDB {
	database, err := sql.Open("sqlite3", source)
	if err != nil {
		logrus.Fatalf("%v", err)
		return nil
	}
	insert, err := database.Prepare(`INSERT INTO identity_users 
(id, first_name, last_name,password,
 email, company, post_code, terms) 
 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Fatalf("%v", err)
	}
	fetch, err := database.Prepare(`SELECT first_name, last_name, 
											company, post_code FROM
										    identity_users where email = ?`)
	if err != nil {
		log.Fatalf("%v", err)
	}
	login, err := database.Prepare(`SELECT first_name, last_name, password FROM
										    identity_users where email = ?`)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return &LiteDB{
		Db:        database,
		InsertNew: insert,
		Fetch:     fetch,
		Login:     login,
	}
}

func (id *LiteDB) Insert(u *store.User) error {
	uid := uuid.New()
	_, err := id.InsertNew.Exec(uid, u.FirstName, u.LastName,
		u.Password, u.Email, u.Company, u.PostCode, u.Terms)
	if err != nil {
		logrus.Errorf("failed to insert user data :: %v", err)
		return err
	}
	return nil
}

func (id *LiteDB) Read(email string) (*store.User, error) {
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

func (id *LiteDB) Authenticate(email, password string) (string, error) {
	rows := id.Login.QueryRow(email)
	var fname, lname, hashedPass string
	err := rows.Scan(&fname, &lname, &hashedPass)
	if err != nil {
		logrus.Errorf("%v", err)
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password))
	if err != nil {
		logrus.Errorf("invalid password :: %v", err)
		return "", err
	}
	if fname == "" && lname == "" {
		return "", nil
	}

	return fmt.Sprintf("%s %s", fname, lname), nil
}
