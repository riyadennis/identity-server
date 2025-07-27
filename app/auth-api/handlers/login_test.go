package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/foundation"
)

func TestLogin(t *testing.T) {
	scenarios := []struct {
		name     string
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
	db := setupDB(t)
	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			h := NewHandler(db, &store.TokenConfig{
				Issuer:  "TEST",
				KeyPath: os.Getenv("KEY_PATH"),
			}, testLogger)
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

func setupDB(t *testing.T) *store.DB {
	t.Helper()
	cfg := store.NewENVConfig()
	conn, err := store.Connect(cfg.DB)
	assert.NoError(t, err)
	db := &store.DB{
		Conn: conn,
	}
	return db
}
