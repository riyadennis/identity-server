package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"

	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/foundation"
)

func TestHandlerDelete(t *testing.T) {
	scenarios := []struct {
		name             string
		params           httprouter.Params
		conn             *sql.DB
		expectedStatus   int
		expectedResponse *foundation.Response
	}{
		{
			name:           "no id",
			params:         nil,
			expectedStatus: http.StatusBadRequest,
			expectedResponse: &foundation.Response{
				Status:    http.StatusBadRequest,
				Message:   errInvalidID.Error(),
				ErrorCode: foundation.InvalidRequest,
			},
		},
		{
			name: "database error",
			params: httprouter.Params{
				{
					Key:   "id",
					Value: "INVALID",
				},
			},
			conn: func() *sql.DB {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectPrepare("DELETE FROM identity_users WHERE id = ?").
					ExpectExec().
					WithArgs("INVALID").
					WillReturnError(errors.New("error"))
				return conn
			}(),
			expectedStatus: http.StatusBadRequest,
			expectedResponse: &foundation.Response{
				Status:    http.StatusBadRequest,
				Message:   errDeleteFailed.Error(),
				ErrorCode: foundation.DatabaseError,
			},
		},
		{
			name: "user do not exist",
			params: httprouter.Params{
				{
					Key:   "id",
					Value: "INVALID",
				},
			},
			conn: func() *sql.DB {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectPrepare("DELETE FROM identity_users WHERE id = ?").
					ExpectExec().
					WithArgs("INVALID").
					WillReturnResult(sqlmock.NewResult(0, 0))
				return conn
			}(),
			expectedStatus: http.StatusBadRequest,
			expectedResponse: &foundation.Response{
				Status:    http.StatusBadRequest,
				Message:   errDeleteFailed.Error(),
				ErrorCode: foundation.UserDoNotExist,
			},
		},
		{
			name: "success",
			params: httprouter.Params{
				{
					Key:   "id",
					Value: "123",
				},
			},
			conn: func() *sql.DB {
				conn, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectPrepare("DELETE FROM identity_users WHERE id = ?").
					ExpectExec().
					WithArgs("123").
					WillReturnResult(sqlmock.NewResult(1, 1))
				return conn
			}(),
			expectedStatus: http.StatusNoContent,
			expectedResponse: &foundation.Response{
				Status:    http.StatusNoContent,
				Message:   "",
				ErrorCode: "",
			},
		},
	}

	logger := log.New(os.Stdout, "IDENTITY-TEST", log.LstdFlags)
	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			dbCOnn := store.NewDB(sc.conn)
			handler := NewHandler(dbCOnn, &store.TokenConfig{}, logger)

			handler.Delete(w, nil, sc.params)

			assert.Equal(t, w.Code, sc.expectedStatus)

			resp := responseFromHTTP(t, w.Body)
			if resp != nil {
				assert.Equal(t, sc.expectedResponse.ErrorCode, resp.ErrorCode)
			}
			assert.Equal(t, sc.expectedResponse.Status, w.Code)
		})
	}
}
