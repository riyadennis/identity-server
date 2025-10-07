package store

import (
	"errors"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"github.com/riyadennis/identity-server/business"
)

func TestAuthenticate(t *testing.T) {
	testcases := []struct {
		name           string
		db             *DB
		expectedResult bool
		expectedError  error
	}{
		{
			name: "prepare failed",
			db: func() *DB {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectPrepare(authQuery).
					WillReturnError(errors.New("error"))
				return NewDB(conn)
			}(),
			expectedError:  errors.New("error"),
			expectedResult: false,
		},
		{
			name: "scan failed",
			db: func() *DB {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectPrepare(authQuery).
					ExpectQuery().
					WithArgs(sqlmock.AnyArg()).
					WillReturnError(errors.New("error"))
				return NewDB(conn)
			}(),
			expectedError:  errors.New("error"),
			expectedResult: false,
		},
		{
			name: "invalid password in DB",
			db: func() *DB {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectPrepare(authQuery).
					ExpectQuery().
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"password"}).AddRow("pass"))
				return NewDB(conn)
			}(),
			expectedError:  bcrypt.ErrHashTooShort,
			expectedResult: false,
		},
		{
			name: "valid password in DB",
			db: func() *DB {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				password, err := business.EncryptPassword("pass")
				assert.NoError(t, err)
				mock.ExpectPrepare(authQuery).
					ExpectQuery().
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"password"}).AddRow(password))
				return NewDB(conn)
			}(),
			expectedResult: true,
		},
	}
	os.Setenv("MYSQL_USERNAME", "root")
	os.Setenv("MYSQL_PASSWORD", "root")
	os.Setenv("MYSQL_HOST", "localhost")
	os.Setenv("MYSQL_DATABASE", "identity")
	os.Setenv("MYSQL_PORT", "80")
	for _, testCase := range testcases {
		t.Run(testCase.name, func(t *testing.T) {
			authenticated, err := testCase.db.Authenticate("email", "pass")
			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedResult, authenticated)
		})
	}

}
