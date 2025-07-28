package store

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang-migrate/migrate/database/mysql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDBInsertSuccess(t *testing.T) {
	scenarios := []struct {
		name        string
		db          *DB
		user        *UserRequest
		uid         string
		expectedErr error
	}{
		{
			name:        "empty connection",
			db:          &DB{},
			expectedErr: errEmptyDBConnection,
		},
		{
			name:        "empty user",
			db:          &DB{Conn: &sql.DB{}},
			user:        nil,
			expectedErr: errEmptyUser,
		},
		{
			name: "success",
			db: func() *DB {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectPrepare("INSERT INTO identity_users").ExpectExec().
					WithArgs(sqlmock.AnyArg(), "John", "Doe", "check", "john.doe@test.com", "Arctura", "12345", true).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectPrepare(regexp.QuoteMeta("SELECT first_name, last_name, email, company, post_code, created_at, updated_at FROM identity_users where id = ? limit 1")).
					ExpectQuery().
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"first_name", "last_name", "email", "company", "post_code", "created_at", "updated_at"}).
						AddRow("John", "Doe", "john.doe@test.com", "Arctura", "12345", time.Now(), time.Now()))
				return NewDB(conn)
			}(),
			user: &UserRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@test.com",
				Password:  "check",
				Company:   "Arctura",
				PostCode:  "12345",
				Terms:     true,
			},
			uid:         uuid.NewString(),
			expectedErr: nil,
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			re, err := sc.db.Insert(context.Background(), sc.user)
			if !errors.Is(err, sc.expectedErr) {
				t.Fatalf("unexpected error, wanted %v, got %v", sc.expectedErr, err)
			}
			if re != nil {
				assert.Equal(t, re.Email, sc.user.Email)
			}
		})
	}
}

func TestDBInsertExecuteFail(t *testing.T) {
	conn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("mock initailisation failed: %v", err)
	}
	uid := uuid.New().String()
	user := &UserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@test.com",
		Password:  "check",
		Terms:     true,
	}
	mock.ExpectPrepare("INSERT INTO identity_users").ExpectExec().
		WithArgs(uid, user.FirstName, user.LastName, user.Password,
			user.Email, user.Company, user.PostCode, user.Terms).
		WillReturnResult(nil)
	_, err = NewDB(conn).Insert(context.Background(), user)
	if err == nil {
		t.Fail()
	}
}

func TestDBInsertPrepareFail(t *testing.T) {
	conn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("mock initailisation failed: %v", err)
	}

	user := &UserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@test.com",
		Password:  "check",
		Terms:     true,
	}
	mock.ExpectPrepare("INSERT INTO identity_users").WillReturnError(mysql.ErrNoDatabaseName)
	_, err = NewDB(conn).Insert(context.Background(), user)
	if err == nil {
		t.Fail()
	}
}
