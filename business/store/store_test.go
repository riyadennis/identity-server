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

func TestPing(t *testing.T) {
	t.Run("ping success", func(t *testing.T) {
		conn, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		assert.NoError(t, err)
		mock.ExpectPing()
		db := &MYSQL{Conn: conn}
		assert.NoError(t, db.Ping())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ping failure", func(t *testing.T) {
		conn, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		assert.NoError(t, err)
		mock.ExpectPing().WillReturnError(errors.New("connection refused"))
		db := &MYSQL{Conn: conn}
		assert.EqualError(t, db.Ping(), "connection refused")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

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
				mock.ExpectPrepare(regexp.QuoteMeta(RetrieveQuery)).
					ExpectQuery().
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"first_name", "last_name", "email", "company", "post_code", "created_at", "updated_at", "role"}).
						AddRow("John", "Doe", "john.doe@test.com", "Arctura", "12345", time.Now(), time.Now(), "user"))
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
				mock.ExpectPrepare(regexp.QuoteMeta(RetrieveQuery)).
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
				mock.ExpectPrepare(regexp.QuoteMeta(RetrieveQuery)).
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
				mock.ExpectPrepare(regexp.QuoteMeta(RetrieveQuery)).
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
				mock.ExpectPrepare(regexp.QuoteMeta(RetrieveQuery)).
					ExpectQuery().
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"first_name", "last_name", "email", "company", "post_code", "created_at", "updated_at", "role"}).
							AddRow("john", "doe", "john.doe@gmail.com", "Arctura", "12345", "2024-01-01", "2024-01-01", "user"))
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
				Role:      "user",
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

func TestDB_UpdateRole(t *testing.T) {
	scenarios := []struct {
		name        string
		db          *MYSQL
		userID      string
		role        string
		expectedErr error
	}{
		{
			name: "prepare failed",
			db: func() *MYSQL {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectPrepare(`UPDATE identity_users SET role`).
					WillReturnError(errors.New("error"))
				return NewDB(conn)
			}(),
			expectedErr: errors.New("error"),
		},
		{
			name: "exec failed",
			db: func() *MYSQL {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectPrepare(`UPDATE identity_users SET role`).
					ExpectExec().WillReturnError(errors.New("exec error"))
				return NewDB(conn)
			}(),
			userID:      "user-123",
			role:        "admin",
			expectedErr: errors.New("exec error"),
		},
		{
			name: "user not found",
			db: func() *MYSQL {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectPrepare(`UPDATE identity_users SET role`).
					ExpectExec().WillReturnResult(sqlmock.NewResult(0, 0))
				return NewDB(conn)
			}(),
			userID:      "nonexistent",
			role:        "admin",
			expectedErr: errors.New("user not found"),
		},
		{
			name: "success",
			db: func() *MYSQL {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectPrepare(`UPDATE identity_users SET role`).
					ExpectExec().WillReturnResult(sqlmock.NewResult(0, 1))
				return NewDB(conn)
			}(),
			userID: "user-123",
			role:   "admin",
		},
	}
	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			err := sc.db.UpdateRole(context.Background(), sc.userID, sc.role)
			assert.Equal(t, sc.expectedErr, err)
		})
	}
}

func TestDB_ListAll(t *testing.T) {
	scenarios := []struct {
		name           string
		db             *MYSQL
		expectedUsers  []*User
		expectedErr    error
		expectedErrMsg string
	}{
		{
			name: "query failed",
			db: func() *MYSQL {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectQuery(`SELECT id, first_name`).
					WillReturnError(errors.New("query error"))
				return NewDB(conn)
			}(),
			expectedErr: errors.New("query error"),
		},
		{
			name: "scan error",
			db: func() *MYSQL {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectQuery(`SELECT id, first_name`).
					WillReturnRows(sqlmock.NewRows(
						[]string{"id", "first_name", "last_name", "email", "company", "post_code", "role", "created_at", "updated_at"}).
						AddRow(nil, nil, nil, nil, nil, nil, nil, nil, nil))
				return NewDB(conn)
			}(),
			expectedErrMsg: `sql: Scan error on column index 0, name "id": converting NULL to string is unsupported`,
		},
		{
			name: "empty result",
			db: func() *MYSQL {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectQuery(`SELECT id, first_name`).
					WillReturnRows(sqlmock.NewRows(
						[]string{"id", "first_name", "last_name", "email", "company", "post_code", "role", "created_at", "updated_at"}))
				return NewDB(conn)
			}(),
			expectedUsers: nil,
		},
		{
			name: "success",
			db: func() *MYSQL {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectQuery(`SELECT id, first_name`).
					WillReturnRows(sqlmock.NewRows(
						[]string{"id", "first_name", "last_name", "email", "company", "post_code", "role", "created_at", "updated_at"}).
						AddRow("1", "John", "Doe", "john@test.com", "Acme", "12345", "user", "2024-01-01", "2024-01-01").
						AddRow("2", "Jane", "Doe", "jane@test.com", "Acme", "12345", "admin", "2024-01-01", "2024-01-01"))
				return NewDB(conn)
			}(),
			expectedUsers: []*User{
				{ID: "1", FirstName: "John", LastName: "Doe", Email: "john@test.com", Company: "Acme", PostCode: "12345", Role: "user", CreatedAt: "2024-01-01", UpdatedAt: "2024-01-01"},
				{ID: "2", FirstName: "Jane", LastName: "Doe", Email: "jane@test.com", Company: "Acme", PostCode: "12345", Role: "admin", CreatedAt: "2024-01-01", UpdatedAt: "2024-01-01"},
			},
		},
	}
	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			users, err := sc.db.ListAll(context.Background())
			if sc.expectedErrMsg != "" {
				assert.EqualError(t, err, sc.expectedErrMsg)
			} else {
				assert.Equal(t, sc.expectedErr, err)
			}
			assert.Equal(t, sc.expectedUsers, users)
		})
	}
}

func TestDB_ListByRole(t *testing.T) {
	scenarios := []struct {
		name           string
		db             *MYSQL
		role           string
		expectedUsers  []*User
		expectedErr    error
		expectedErrMsg string
	}{
		{
			name: "query failed",
			db: func() *MYSQL {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectQuery(`SELECT id, first_name`).
					WillReturnError(errors.New("query error"))
				return NewDB(conn)
			}(),
			role:        "admin",
			expectedErr: errors.New("query error"),
		},
		{
			name: "scan error",
			db: func() *MYSQL {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectQuery(`SELECT id, first_name`).
					WillReturnRows(sqlmock.NewRows(
						[]string{"id", "first_name", "last_name", "email", "company", "post_code", "role", "created_at", "updated_at"}).
						AddRow(nil, nil, nil, nil, nil, nil, nil, nil, nil))
				return NewDB(conn)
			}(),
			role:           "admin",
			expectedErrMsg: `sql: Scan error on column index 0, name "id": converting NULL to string is unsupported`,
		},
		{
			name: "empty result",
			db: func() *MYSQL {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectQuery(`SELECT id, first_name`).
					WillReturnRows(sqlmock.NewRows(
						[]string{"id", "first_name", "last_name", "email", "company", "post_code", "role", "created_at", "updated_at"}))
				return NewDB(conn)
			}(),
			role:          "admin",
			expectedUsers: nil,
		},
		{
			name: "success",
			db: func() *MYSQL {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectQuery(`SELECT id, first_name`).
					WithArgs("admin").
					WillReturnRows(sqlmock.NewRows(
						[]string{"id", "first_name", "last_name", "email", "company", "post_code", "role", "created_at", "updated_at"}).
						AddRow("1", "John", "Doe", "john@test.com", "Acme", "12345", "admin", "2024-01-01", "2024-01-01"))
				return NewDB(conn)
			}(),
			role: "admin",
			expectedUsers: []*User{
				{ID: "1", FirstName: "John", LastName: "Doe", Email: "john@test.com", Company: "Acme", PostCode: "12345", Role: "admin", CreatedAt: "2024-01-01", UpdatedAt: "2024-01-01"},
			},
		},
	}
	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			users, err := sc.db.ListByRole(context.Background(), sc.role)
			if sc.expectedErrMsg != "" {
				assert.EqualError(t, err, sc.expectedErrMsg)
			} else {
				assert.Equal(t, sc.expectedErr, err)
			}
			assert.Equal(t, sc.expectedUsers, users)
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
