package sqlM

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/riyadennis/identity-server/internal/store"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	Db        *sql.DB
	InsertNew *sql.Stmt
	Fetch     *sql.Stmt
	Login     *sql.Stmt
}

func ConnectDB() (*sql.DB, error) {
	database, err := sql.Open("mysql",
		"root:root@tcp(mysql-development:3306)/identity_db")
	if err != nil {
		logrus.Fatalf("%v", err)
		return nil, err
	}
	return database, nil
}

func PrepareDB(database *sql.DB) (*DB, error) {
	insert, err := database.Prepare(`INSERT INTO identity_users 
(id, first_name, last_name,password,
 email, company, post_code, terms) 
 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Fatalf("%v", err)
		return nil, err
	}
	fetch, err := database.Prepare(`SELECT first_name, last_name, 
											company, post_code FROM
										    identity_users where email = ?`)
	if err != nil {
		log.Fatalf("%v", err)
		return nil, err
	}
	login, err := database.Prepare(`SELECT first_name, last_name, password FROM
										    identity_users where email = ?`)
	if err != nil {
		log.Fatalf("%v", err)
		return nil, err
	}
	return &DB{
		Db:        database,
		InsertNew: insert,
		Fetch:     fetch,
		Login:     login,
	}, nil
}

func (id *DB) Insert(u *store.User) error {
	uid := uuid.New()
	_, err := id.InsertNew.Exec(uid, u.FirstName, u.LastName,
		u.Password, u.Email, u.Company, u.PostCode, u.Terms)
	if err != nil {
		logrus.Errorf("failed to insert user data :: %v", err)
		return err
	}
	return nil
}

func (id *DB) Read(email string) (*store.User, error) {
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

func (id *DB) Authenticate(email, password string) (bool, error) {
	rows := id.Login.QueryRow(email)
	var hashedPass string
	err := rows.Scan(&hashedPass)
	if err != nil {
		logrus.Errorf("%v", err)
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password))
	if err != nil {
		logrus.Errorf("invalid password :: %v", err)
		return false, err
	}
	return true, nil
}
