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
	Insert(ctx context.Context, u *User) (*User, error)
	Read(ctx context.Context, email string) (*User, error)
	Retrieve(ctx context.Context, id string) (*User, error)
	Delete(id string) (int64, error)
	Ping() error
	UpdateRole(ctx context.Context, userID string, role string) error
	ListByRole(ctx context.Context, role string) ([]*User, error)
	ListAll(ctx context.Context) ([]*User, error)
}

// User holds data from the registration request body
type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Company   string `json:"company"`
	PostCode  string `json:"post_code"`
	Terms     bool   `json:"terms"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// MYSQL implements store interface
type MYSQL struct {
	Conn *sql.DB
}

// NewDB creates a new instance if the MYSQL
func NewDB(database *sql.DB) *MYSQL {
	return &MYSQL{
		Conn: database,
	}
}

// Insert creates a new user during registration
func (m *MYSQL) Insert(ctx context.Context, u *User) (*User, error) {
	if m.Conn == nil {
		return nil, errEmptyDBConnection
	}

	if u == nil {
		return nil, errEmptyUser
	}

	id := uuid.New().String()

	insert, err := m.Conn.Prepare(`INSERT INTO identity_users
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

	return m.Retrieve(ctx, id)
}

// Retrieve will fetch data from auth for a user as per the id
// will return nil if the user is not found
func (m *MYSQL) Retrieve(ctx context.Context, id string) (*User, error) {
	if m.Conn == nil {
		return nil, errEmptyDBConnection
	}

	fetch, err := m.Conn.Prepare(RetrieveQuery)
	if err != nil {
		return nil, err
	}

	row := fetch.QueryRowContext(ctx, id)
	user := &User{}
	err = row.Scan(
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Company,
		&user.PostCode,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Role,
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

var RetrieveQuery = `SELECT first_name, last_name, email, company, post_code, created_at, updated_at, role FROM identity_users where id = ? limit 1`

var ReadQuery = `SELECT id,
       first_name,
       last_name,
       email,
       company,
       post_code,
       created_at,
       updated_at
		FROM identity_users
		where email = ?`

// Read will fetch data from auth for a user as per the email
// will return nil if the user is not found
func (m *MYSQL) Read(ctx context.Context, email string) (*User, error) {
	if m.Conn == nil {
		return nil, errEmptyDBConnection
	}

	rows, err := m.Conn.QueryContext(ctx, ReadQuery, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := &User{}
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
			logrus.Errorf("failed to read user data: %v", err)
			return nil, errInvalidDataInDB
		}
	}

	return user, nil
}

// Delete removes a user from auth as per the ID
func (m *MYSQL) Delete(id string) (int64, error) {
	remove, err := m.Conn.Prepare(`DELETE  FROM identity_users WHERE id = ?`)
	if err != nil {
		logrus.Errorf("user deletion failed with error: %v", err)
		return 0, err
	}

	result, err := remove.Exec(id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (m *MYSQL) Ping() error {
	return m.Conn.Ping()
}

// UpdateRole sets the role for a user identified by userID.
func (m *MYSQL) UpdateRole(ctx context.Context, userID string, role string) error {
	stmt, err := m.Conn.Prepare(`UPDATE identity_users SET role = ? WHERE id = ?`)
	if err != nil {
		return err
	}
	result, err := stmt.ExecContext(ctx, role, userID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("user not found")
	}
	return nil
}

// ListAll returns all registered users.
func (m *MYSQL) ListAll(ctx context.Context) ([]*User, error) {
	rows, err := m.Conn.QueryContext(ctx,
		`SELECT id, first_name, last_name, email, company, post_code, role, created_at, updated_at
		 FROM identity_users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		u := &User{}
		if err := rows.Scan(
			&u.ID, &u.FirstName, &u.LastName, &u.Email,
			&u.Company, &u.PostCode, &u.Role, &u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// ListByRole returns all users with the given role.
func (m *MYSQL) ListByRole(ctx context.Context, role string) ([]*User, error) {
	rows, err := m.Conn.QueryContext(ctx,
		`SELECT id, first_name, last_name, email, company, post_code, role, created_at, updated_at
		 FROM identity_users WHERE role = ?`, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		u := &User{}
		if err := rows.Scan(
			&u.ID, &u.FirstName, &u.LastName, &u.Email,
			&u.Company, &u.PostCode, &u.Role, &u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
