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
		db          *MYSQL
		user        *User
		uid         string
		expectedErr error
	}{
		{
			name:        "empty connection",
			db:          &MYSQL{},
			expectedErr: errEmptyDBConnection,
		},
		{
			name:        "empty user",
			db:          &MYSQL{Conn: &sql.DB{}},
			user:        nil,
			expectedErr: errEmptyUser,
		},
		{
			name: "success",
			db: func() *MYSQL {
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
			user: &User{
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
	user := &User{
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

	user := &User{
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

func TestDBRetrieve(t *testing.T) {
	scenarios := []struct {
		name        string
		db          *MYSQL
		user        *User
		uid         string
		expectedErr error
	}{
		{
			name:        "empty connection",
			db:          &MYSQL{},
			expectedErr: errEmptyDBConnection,
		},
		{
			name: "prepare failed",
			db: func() *MYSQL {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectPrepare(regexp.QuoteMeta("SELECT first_name, last_name, email, company, post_code, created_at, updated_at FROM identity_users where id = ? limit 1")).
					WillReturnError(errors.New("error"))
				return NewDB(conn)
			}(),
			expectedErr: errors.New("error"),
		},
		{
			name: "sql failed",
			db: func() *MYSQL {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectPrepare(regexp.QuoteMeta("SELECT first_name, last_name, email, company, post_code, created_at, updated_at FROM identity_users where id = ? limit 1")).
					ExpectQuery().
					WithArgs(sqlmock.AnyArg()).
					WillReturnError(errors.New("error"))
				return NewDB(conn)
			}(),
			expectedErr: errors.New("error"),
		},
		{
			name: "user not found",
			db: func() *MYSQL {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectPrepare(regexp.QuoteMeta("SELECT first_name, last_name, email, company, post_code, created_at, updated_at FROM identity_users where id = ? limit 1")).
					ExpectQuery().
					WithArgs(sqlmock.AnyArg()).
					WillReturnError(sql.ErrNoRows)
				return NewDB(conn)
			}(),
		},
		{
			name: "success",
			db: func() *MYSQL {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectPrepare(regexp.QuoteMeta("SELECT first_name, last_name, email, company, post_code, created_at, updated_at FROM identity_users where id = ? limit 1")).
					ExpectQuery().
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"first_name", "last_name", "email", "company", "post_code", "created_at", "updated_at"}).
							AddRow("john", "doe", "john.doe@gmail.com", "Arctura", "12345", "2024-01-01", "2024-01-01"))
				return NewDB(conn)
			}(),
			user: &User{
				FirstName: "john",
				LastName:  "doe",
				Email:     "john.doe@gmail.com",
				Company:   "Arctura",
				PostCode:  "12345",
				CreatedAt: "2024-01-01",
				UpdatedAt: "2024-01-01",
			},
		},
	}
	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			user, err := sc.db.Retrieve(context.Background(), sc.uid)
			assert.Equal(t, sc.expectedErr, err)
			assert.Equal(t, sc.user, user)
		})
	}
}

func TestDB_Read(t *testing.T) {
	scenarios := []struct {
		name        string
		db          *MYSQL
		user        *User
		email       string
		expectedErr error
	}{
		{
			name:        "empty connection",
			db:          &MYSQL{},
			expectedErr: errEmptyDBConnection,
		},
		{
			name: "query failed",
			db: func() *MYSQL {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectQuery(ReadQuery).
					WillReturnError(errors.New("error"))
				return NewDB(conn)
			}(),
			expectedErr: errors.New("error"),
		},
		{
			name: "invalid data in MYSQL",
			db: func() *MYSQL {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectQuery(ReadQuery).
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows(
						[]string{"id", "first_name", "last_name", "email", "company", "post_code", "created_at", "updated_at"}).
						AddRow(nil, nil, nil, nil, nil, nil, nil, nil))
				return NewDB(conn)
			}(),
			expectedErr: errInvalidDataInDB,
		},
		{
			name: "success",
			db: func() *MYSQL {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectQuery(ReadQuery).
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows(
						[]string{"id", "first_name", "last_name", "email", "company", "post_code", "created_at", "updated_at"}).
						AddRow(123, "john", "doe", "john.doe@gmail.com", "Arctura", "12345", "2024-01-01", "2024-01-01"))
				return NewDB(conn)
			}(),
			user: &User{
				ID:        "123",
				FirstName: "john",
				LastName:  "doe",
				Email:     "john.doe@gmail.com",
				Company:   "Arctura",
				PostCode:  "12345",
				CreatedAt: "2024-01-01",
				UpdatedAt: "2024-01-01",
			},
		},
	}
	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			user, err := sc.db.Read(context.Background(), sc.email)
			assert.Equal(t, sc.expectedErr, err)
			assert.Equal(t, sc.user, user)
		})
	}
}

func TestDB_Delete(t *testing.T) {
	scenarios := []struct {
		name                 string
		db                   *MYSQL
		id                   string
		expectedErr          error
		expectedRowsAffected int64
	}{
		{
			name: "deletion prepare failed",
			db: func() *MYSQL {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectPrepare(`DELETE  FROM identity_users WHERE id = ?`).
					WillReturnError(errors.New("error"))
				return NewDB(conn)
			}(),
			expectedErr: errors.New("error"),
		},
		{
			name: "deletion execution failed",
			db: func() *MYSQL {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectPrepare(`DELETE  FROM identity_users WHERE id = ?`).
					ExpectExec().WillReturnError(errors.New("error"))
				return NewDB(conn)
			}(),
			expectedErr: errors.New("error"),
		},
		{
			name: "success",
			db: func() *MYSQL {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectPrepare(`DELETE  FROM identity_users WHERE id = ?`).
					ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))
				return NewDB(conn)
			}(),
			expectedRowsAffected: 1,
		},
	}
	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			rowsAffected, err := sc.db.Delete(sc.id)
			assert.Equal(t, sc.expectedErr, err)
			assert.Equal(t, sc.expectedRowsAffected, rowsAffected)
		})
	}
}
