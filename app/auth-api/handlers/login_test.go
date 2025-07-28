package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/foundation"
)

func TestLogin(t *testing.T) {
	scenarios := []struct {
		name     string
		conn     *sql.DB
		request  *http.Request
		response *foundation.Response
	}{
		{
			name:    "empty request",
			request: &http.Request{},
			response: &foundation.Response{
				Status:    http.StatusBadRequest,
				Message:   "empty login data",
				ErrorCode: foundation.InvalidRequest,
			},
		},
		{
			name:    "missing email",
			request: request(t, "/login", `{}`),
			response: &foundation.Response{
				Status:    http.StatusBadRequest,
				Message:   "empty login data",
				ErrorCode: foundation.InvalidRequest,
			},
		},
		{
			name: "missing password",
			request: request(t, "/login", `{
	"email": "john4@gmail.com"
}`),
			response: &foundation.Response{
				Status:    http.StatusBadRequest,
				Message:   "empty login data",
				ErrorCode: foundation.InvalidRequest,
			},
		},
	}

	logger := log.New(os.Stdout, "IDENTITY-TEST", log.LstdFlags)
	conn, mock, err := sqlmock.New()
	assert.NoError(t, err)

	mock.ExpectPrepare(regexp.QuoteMeta("SELECT first_name, last_name, email, company, post_code, created_at, updated_at FROM identity_users email = ?")).
		ExpectQuery().
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"first_name", "last_name", "email", "company", "post_code", "created_at", "updated_at"}).
			AddRow("John", "Doe", "john.doe@test.com", "Arctura", "12345", time.Now(), time.Now()))
	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			h := NewHandler(store.NewDB(conn), &store.TokenConfig{
				Issuer:  "TEST",
				KeyPath: os.Getenv("KEY_PATH"),
			}, logger)
			h.Login(rr, sc.request, nil)
			re := response(t, rr.Body)
			assert.Equal(t, sc.response, re)
		})
	}
}

func response(t *testing.T, body io.Reader) *foundation.Response {
	t.Helper()
	var re *foundation.Response
	if body != nil {
		data, err := io.ReadAll(body)
		if err != nil {
			t.Error(err.Error())
		}
		re = &foundation.Response{}
		err = json.Unmarshal(data, re)
		if err != nil {
			t.Error(err.Error())
		}
	}
	return re
}

func request(t *testing.T, endpoint, content string) *http.Request {
	t.Helper()
	body := strings.NewReader(content)
	req, err := http.NewRequest("POST",
		endpoint, body)
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("content-type", "application/json")
	return req
}
