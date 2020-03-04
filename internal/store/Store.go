package store

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// User hold information needed to complete user registration
type User struct {
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	Company          string `json:"company"`
	PostCode         string `json:"post_code"`
	Terms            bool   `json:"terms"`
	RegistrationDate string
}

type DB struct {
	InsertNew *sql.Stmt
	Fetch     *sql.Stmt
	Login     *sql.Stmt
	Remove    *sql.Stmt
}

func PrepareDB(database *sql.DB) (*DB, error) {
	insert, err := database.Prepare(`INSERT INTO identity_users 
(id, first_name, last_name,password,
 email, company, post_code, terms) 
 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		logrus.Fatalf("%v", err)
		return nil, err
	}
	fetch, err := database.Prepare(`SELECT first_name, last_name, 
											company, post_code FROM
										    identity_users where email = ?`)
	if err != nil {
		logrus.Fatalf("%v", err)
		return nil, err
	}
	login, err := database.Prepare(`SELECT  password FROM
										    identity_users where email = ?`)
	if err != nil {
		logrus.Fatalf("%v", err)
		return nil, err
	}
	delete, err := database.Prepare(`DELETE  FROM identity_users where email = ?`)
	if err != nil {
		logrus.Fatalf("%v", err)
		return nil, err
	}
	return &DB{
		InsertNew: insert,
		Fetch:     fetch,
		Login:     login,
		Remove:    delete,
	}, nil
}

//TODO to consolidate these functions between two database type
type Store interface {
	Insert(u *User) error
	Read(email string) (*User, error)
	Authenticate(email, password string) (bool, error)
	Delete(email string) (bool, error)
}

func (id *DB) Insert(u *User) error {
	uid := uuid.New()
	_, err := id.InsertNew.Exec(uid, u.FirstName, u.LastName,
		u.Password, u.Email, u.Company, u.PostCode, u.Terms)
	if err != nil {
		logrus.Errorf("failed to insert user data :: %v", err)
		return err
	}
	return nil
}

func (id *DB) Read(email string) (*User, error) {
	rows := id.Fetch.QueryRow(email)

	var fname, lname, post, company string
	err := rows.Scan(&fname, &lname, &post, &company)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Infof("user not found :: %s", email)
			return nil, nil
		}
		return nil, err
	}
	u := &User{
		FirstName: fname,
		LastName:  lname,
		Company:   company,
		PostCode:  post,
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

func (id *DB) Delete(email string) (bool, error) {
	u, err := id.Read(email)
	if err != nil {
		return false, err
	}
	if u == nil {
		return false, errors.New("user not found")
	}
	_, err = id.Remove.Query(email)
	if err != nil {
		return false, err
	}
	return true, nil
}
