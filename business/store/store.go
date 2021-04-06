package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var (
	errEmptyUser   = errors.New("empty user")
	errEmptyUserID = errors.New("empty user id")
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

// DB implements store interface
type DB struct {
	Conn *sql.DB
}

// NewDB creates a new instance if the DB
func NewDB(database *sql.DB) *DB {
	return &DB{
		Conn: database,
	}
}

// Store have CRUD functions for user management
type Store interface {
	Insert(ctx context.Context, u *User, uid string) error
	Read(ctx context.Context, email string) (*User, error)
	Authenticate(email, password string) (bool, error)
	Delete(email string) (int64, error)
}

// Insert creates a new user during registration
func (d *DB) Insert(ctx context.Context, u *User, uid string) error {
	if d.Conn == nil {
		return errEmptyDBConnection
	}

	if u == nil {
		return errEmptyUser
	}

	if uid == "" {
		return errEmptyUserID
	}

	insert, err := d.Conn.Prepare(`INSERT INTO identity_users 
(id, first_name, last_name,password,
 email, company, post_code, terms) 
 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		logrus.Errorf("failed to prepare user insert: %v", err)
		return err
	}

	_, err = insert.ExecContext(ctx, uid, u.FirstName, u.LastName,
		u.Password, u.Email, u.Company, u.PostCode, u.Terms)
	if err != nil {
		logrus.Errorf("failed to insert user data: %v", err)
		return err
	}

	return nil
}

// Read will fetch data from db for a user as per the email
// will return nil if user is not found
func (d *DB) Read(ctx context.Context, email string) (*User, error) {
	if d.Conn == nil {
		return nil, errEmptyDBConnection
	}

	fetch, err := d.Conn.Prepare(
		"SELECT first_name, last_name,company, post_code FROM identity_users where email = ?")
	if err != nil {
		return nil, err
	}

	rows, err := fetch.QueryContext(ctx, email)
	if err != nil {
		return nil, err
	}

	user := &User{}
	err = rows.Scan(&user.FirstName, &user.LastName, &user.PostCode, &user.Company)
	if errors.Is(err, sql.ErrNoRows) {
		logrus.Infof("user not found :: %s", email)
		return nil, nil
	}

	return user, nil
}

// Authenticate checks the validity of a given password for an email
func (d *DB) Authenticate(email, password string) (bool, error) {
	login, err := d.Conn.Prepare(`SELECT  password FROM
										    identity_users where email = ?`)
	if err != nil {
		logrus.Fatalf("%v", err)
		return false, err
	}
	rows := login.QueryRow(email)
	var hashedPass string
	err = rows.Scan(&hashedPass)
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

// Delete removes an email from db
func (d *DB) Delete(email string) (int64, error) {
	remove, err := d.Conn.Prepare(`DELETE  FROM identity_users where email = ?`)
	if err != nil {
		logrus.Fatalf("%v", err)
		return 0, err
	}
	result, err := remove.Exec(email)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
