package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var (
	errEmptyUser = errors.New("empty user")
)

// UserRequest hold data from registration request body
type UserRequest struct {
	FirstName string `jsonapi:"attr,first_name"`
	LastName  string `jsonapi:"attr,last_name"`
	Email     string `jsonapi:"attr,email"`
	Password  string `jsonapi:"attr,password"`
	Company   string `jsonapi:"attr,company"`
	PostCode  string `jsonapi:"attr,post_code"`
	Terms     bool   `jsonapi:"attr,terms"`
}

// UserResource hold data about user in the database
type UserResource struct {
	ID        string `jsonapi:"attr,id"`
	FirstName string `jsonapi:"attr,first_name"`
	LastName  string `jsonapi:"attr,last_name"`
	Email     string `jsonapi:"attr,email"`
	Password  string `jsonapi:"attr,password"`
	Company   string `jsonapi:"attr,company"`
	PostCode  string `jsonapi:"attr,post_code"`
	Terms     bool   `jsonapi:"attr,terms"`
	CreatedAt string `jsonapi:"attr,created_at"`
	UpdatedAt string `jsonapi:"attr,updated_at"`
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
	Insert(ctx context.Context, u *UserRequest) error
	Read(ctx context.Context, email string) (*UserResource, error)
	Authenticate(email, password string) (bool, error)
	Delete(email string) (int64, error)
}

// Insert creates a new user during registration
func (d *DB) Insert(ctx context.Context, u *UserRequest) error {
	if d.Conn == nil {
		return errEmptyDBConnection
	}

	if u == nil {
		return errEmptyUser
	}

	id := uuid.New().String()

	insert, err := d.Conn.Prepare(`INSERT INTO identity_users 
(id, first_name, last_name,password,
 email, company, post_code, terms) 
 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		logrus.Errorf("failed to prepare user insert: %v", err)
		return err
	}

	_, err = insert.ExecContext(ctx, id, u.FirstName, u.LastName,
		u.Password, u.Email, u.Company, u.PostCode, u.Terms)
	if err != nil {
		logrus.Errorf("failed to insert user data: %v", err)
		return err
	}

	return nil
}

// Retrieve will fetch data from db for a user as per the id
// will return nil if user is not found
func (d *DB) Retrieve(ctx context.Context, email string) (*UserResource, error) {
	if d.Conn == nil {
		return nil, errEmptyDBConnection
	}

	fetch, err := d.Conn.Prepare(
		`SELECT id,
       first_name,
       last_name,
       company,
       post_code,
       created_at,
       updated_at,
		FROM identity_users 
		where email = ?`)
	if err != nil {
		return nil, err
	}

	rows, err := fetch.QueryContext(ctx, email)
	if err != nil {
		return nil, err
	}

	user := &UserResource{}
	err = rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.PostCode,
		&user.Company,
		&user.PostCode,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		logrus.Infof("user not found :: %s", email)
		return nil, nil
	}

	return user, nil
}

// Read will fetch data from db for a user as per the email
// will return nil if user is not found
func (d *DB) Read(ctx context.Context, email string) (*UserResource, error) {
	if d.Conn == nil {
		return nil, errEmptyDBConnection
	}

	fetch, err := d.Conn.Prepare(
		`SELECT id,
       first_name,
       last_name,
       company,
       post_code,
       created_at,
       updated_at
		FROM identity_users 
		where email = ?`)
	if err != nil {
		return nil, err
	}

	rows, err := fetch.QueryContext(ctx, email)
	if err != nil {
		return nil, err
	}

	user := &UserResource{}
	err = rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.PostCode,
		&user.Company,
		&user.PostCode,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
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
