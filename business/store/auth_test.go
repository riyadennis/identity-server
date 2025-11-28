package store

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var (
	testExpiry = time.Now().Add(time.Hour * 24)
)

func TestAuthenticate(t *testing.T) {
	testcases := []struct {
		name           string
		db             *Auth
		expectedResult bool
		expectedError  error
	}{
		{
			name:           "prepare failed",
			db:             prepareFailedAuth(t, authQuery),
			expectedError:  errors.New("error"),
			expectedResult: false,
		},
		{
			name: "scan failed",
			db: func() *Auth {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectPrepare(regexp.QuoteMeta(authQuery)).
					ExpectQuery().
					WithArgs(sqlmock.AnyArg()).
					WillReturnError(errors.New("error"))
				return &Auth{
					Conn: conn,
				}
			}(),
			expectedError:  errors.New("error"),
			expectedResult: false,
		},
		{
			name: "invalid password in DB",
			db: func() *Auth {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectPrepare(regexp.QuoteMeta(authQuery)).
					ExpectQuery().
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"password"}).AddRow("pass"))
				return &Auth{
					Conn:   conn,
					Logger: logrus.New(),
				}
			}(),
			expectedError:  bcrypt.ErrHashTooShort,
			expectedResult: false,
		},
		{
			name: "valid password in DB",
			db: func() *Auth {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				enPass, err := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.MinCost)
				assert.NoError(t, err)
				mock.ExpectPrepare(regexp.QuoteMeta(authQuery)).
					ExpectQuery().
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"password"}).AddRow(string(enPass)))
				return &Auth{
					Conn: conn,
				}
			}(),
			expectedResult: true,
		},
	}
	for _, testCase := range testcases {
		t.Run(testCase.name, func(t *testing.T) {
			authenticated, err := testCase.db.Authenticate("email", "pass")
			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedResult, authenticated)
		})
	}

}

func TestAuth_FetchLoginToken(t *testing.T) {
	testcases := []struct {
		name           string
		auth           *Auth
		expectedResult *TokenRecord
		expectedError  error
	}{
		{
			name:          "prepare failed",
			auth:          prepareFailedAuth(t, tokenQuery),
			expectedError: errors.New("error"),
		},
		{
			name: "query failed",
			auth: func() *Auth {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectPrepare(regexp.QuoteMeta(tokenQuery)).
					ExpectQuery().WithArgs(sqlmock.AnyArg()).
					WillReturnError(errors.New("error"))
				return &Auth{
					Conn: conn,
				}
			}(),
			expectedError: errors.New("error"),
		},
		{
			name: "no record",
			auth: func() *Auth {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectPrepare(regexp.QuoteMeta(tokenQuery)).
					ExpectQuery().WithArgs(sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows(
						[]string{"id", "token", "ttl", "expiry", "last_used"}).
						AddRow("", "", "", time.Now(), ""))
				return &Auth{
					Conn: conn,
				}
			}(),
		},
		{
			name: "success",
			auth: func() *Auth {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectPrepare(regexp.QuoteMeta(tokenQuery)).
					ExpectQuery().WithArgs(sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows(
						[]string{"id", "token", "ttl", "expiry", "last_used"}).
						AddRow("123", "token", "123", testExpiry, "2024-01-01"))
				return &Auth{
					Conn: conn,
				}
			}(),
			expectedResult: validTokenRecord(t),
		},
	}
	logger := logrus.New()
	for _, testCase := range testcases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.auth.Logger = logger
			token, err := testCase.auth.FetchLoginToken("token")
			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedResult, token)
		})
	}
}

func TestAuth_SaveLoginToken(t *testing.T) {
	testCases := []struct {
		name          string
		db            *Auth
		inputToken    *TokenRecord
		expectedError error
	}{
		{
			name:          "prepare failed",
			db:            prepareFailedAuth(t, saveTokenQuery),
			expectedError: errors.New("error"),
		},
		{
			name: "insert error",
			db: func() *Auth {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectPrepare(regexp.QuoteMeta(saveTokenQuery)).
					ExpectExec().
					WillReturnError(errors.New("error"))
				return &Auth{
					Conn:   conn,
					Logger: logrus.New(),
				}
			}(),
			inputToken:    validTokenRecord(t),
			expectedError: errors.New("error"),
		},
		{
			name: "no rows affected",
			db: func() *Auth {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectPrepare(regexp.QuoteMeta(saveTokenQuery)).
					ExpectExec().
					WillReturnResult(sqlmock.NewResult(0, 0))
				return &Auth{
					Conn:   conn,
					Logger: logrus.New(),
				}
			}(),
			inputToken:    validTokenRecord(t),
			expectedError: errors.New("failed to save token"),
		},
		{
			name: "success",
			db: func() *Auth {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectPrepare(regexp.QuoteMeta(saveTokenQuery)).
					ExpectExec().
					WillReturnResult(sqlmock.NewResult(123, 1))
				return &Auth{
					Conn:   conn,
					Logger: logrus.New(),
				}
			}(),
			inputToken: validTokenRecord(t),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := testCase.db.SaveLoginToken(context.Background(), testCase.inputToken)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func prepareFailedAuth(t *testing.T, query string) *Auth {
	t.Helper()
	conn, mock, err := sqlmock.New()
	assert.NoError(t, err)
	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		WillReturnError(errors.New("error"))
	return &Auth{
		Conn:   conn,
		Logger: logrus.New(),
	}
}

func validTokenRecord(t *testing.T) *TokenRecord {
	t.Helper()
	return &TokenRecord{
		Id:       "123",
		Token:    "token",
		Expiry:   testExpiry,
		TTL:      "123",
		LastUsed: sql.NullString{String: "2024-01-01", Valid: true},
	}
}
