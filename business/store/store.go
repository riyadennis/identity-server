package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	errEmptyUser = errors.New("empty user")
)

// Store have CRUD functions for user management
type Store interface {
	Insert(ctx context.Context, u *UserRequest) (*UserResource, error)
	Read(ctx context.Context, email string) (*UserResource, error)
	Delete(id string) (int64, error)
}

// UserRequest hold data from registration request body
type UserRequest struct {
	ID        string `jsonapi:"primary,user"`
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
	ID        string `jsonapi:"primary,user"`
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

// Insert creates a new user during registration
func (d *DB) Insert(ctx context.Context, u *UserRequest) (*UserResource, error) {
	if d.Conn == nil {
		return nil, errEmptyDBConnection
	}

	if u == nil {
		return nil, errEmptyUser
	}

	id := uuid.New().String()

	insert, err := d.Conn.Prepare(`INSERT INTO identity_users 
(id, first_name, last_name,password,
 email, company, post_code, terms) 
 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		logrus.Errorf("failed to prepare user insert: %v", err)
		return nil, err
	}

	_, err = insert.ExecContext(ctx, id, u.FirstName, u.LastName,
		u.Password, u.Email, u.Company, u.PostCode, u.Terms)
	if err != nil {
		logrus.Errorf("failed to insert user data: %v", err)
		return nil, err
	}

	return d.Retrieve(ctx, id)
}

// Retrieve will fetch data from db for a user as per the id
// will return nil if user is not found
func (d *DB) Retrieve(ctx context.Context, id string) (*UserResource, error) {
	if d.Conn == nil {
		return nil, errEmptyDBConnection
	}

	fetch, err := d.Conn.Prepare(
		`SELECT
       first_name,
       last_name,
       email,
       company,
       post_code,
       created_at,
       updated_at
		FROM identity_users 
		where id = ? limit 1`)
	if err != nil {
		return nil, err
	}

	row := fetch.QueryRowContext(ctx, id)
	user := &UserResource{}
	err = row.Scan(
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Company,
		&user.PostCode,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logrus.Infof("user not found :: %s", id)
			return nil, nil
		}

		return nil, err
	}
	user.ID = id
	return user, nil
}

// Read will fetch data from db for a user as per the email
// will return nil if user is not found
func (d *DB) Read(ctx context.Context, email string) (*UserResource, error) {
	query := `SELECT id,
       first_name,
       last_name,
       email,
       company,
       post_code,
       created_at,
       updated_at
		FROM identity_users 
		where email = ?`

	rows, err := d.Conn.QueryContext(ctx, query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := &UserResource{}
	for rows.Next() {
		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Company,
			&user.PostCode,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logrus.Infof("user not found :: %s", email)
				return nil, nil
			}

			return nil, err
		}
	}

	return user, nil
}

// Delete removes a user from db as per the ID
func (d *DB) Delete(id string) (int64, error) {
	remove, err := d.Conn.Prepare(`DELETE  FROM identity_users WHERE id = ?`)
	if err != nil {
		logrus.Fatalf("%v", err)
		return 0, err
	}

	result, err := remove.Exec(id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
