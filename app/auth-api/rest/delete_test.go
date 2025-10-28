package rest

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/foundation"
)

func TestHandlerDelete(t *testing.T) {
	scenarios := []struct {
		name             string
		userID           string
		conn             *sql.DB
		expectedStatus   int
		expectedResponse *foundation.Response
	}{
		{
			name:           "no id",
			expectedStatus: http.StatusBadRequest,
			expectedResponse: &foundation.Response{
				Status:    http.StatusBadRequest,
				Message:   errInvalidID.Error(),
				ErrorCode: foundation.InvalidRequest,
			},
		},
		{
			name:   "database error",
			userID: "INVALID",
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
			name:   "user do not exist",
			userID: "INVALID",
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
			name:   "success",
			userID: "123",
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
	
	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			dbCOnn := store.NewDB(sc.conn)
			handler := NewHandler(dbCOnn, &store.Auth{
				Conn: sc.conn,
			}, &store.TokenConfig{}, logrus.New())

			r := httptest.NewRequest("GET", "/admin/delete/{userID}", nil)
			routeContext := chi.NewRouteContext()
			routeContext.URLParams.Add("userID", sc.userID)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, routeContext))

			handler.Delete(w, r)

			assert.Equal(t, w.Code, sc.expectedStatus)

			resp := responseFromHTTP(t, w.Body)
			if resp != nil {
				assert.Equal(t, sc.expectedResponse.ErrorCode, resp.ErrorCode)
			}
			assert.Equal(t, sc.expectedResponse.Status, w.Code)
		})
	}
}
