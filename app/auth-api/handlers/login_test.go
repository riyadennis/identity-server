package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/foundation"
)

type MockStore struct {
	Error error
	*store.UserResource
}

func (m *MockStore) Insert(ctx context.Context, u *store.UserRequest) (*store.UserResource, error) {
	return m.UserResource, m.Error
}

func (m *MockStore) Read(ctx context.Context, email string) (*store.UserResource, error) {
	return m.UserResource, m.Error
}

func (m *MockStore) Delete(id string) (int64, error) {
	return 0, m.Error
}

type MockAuthenticator struct {
	ReturnVal bool
	Error     error
}

func (ma *MockAuthenticator) Authenticate(email, password string) (bool, error) {
	return ma.ReturnVal, ma.Error
}

func TestLogin(t *testing.T) {
	scenarios := []struct {
		name          string
		store         store.Store
		authenticator store.Authenticator
		request       *http.Request
		response      *foundation.Response
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
		{
			name: "login DB error",
			request: func() *http.Request {
				credentials := `{
					"email": "john4@gmail.com","password":"pass"
				}`
				req := request(t, "/login", "")

				base64 := base64.StdEncoding.EncodeToString([]byte(credentials))
				req.Header.Set("Authorization", "Basic "+base64)
				return req
			}(),
			response: &foundation.Response{
				Status:    http.StatusInternalServerError,
				Message:   "error",
				ErrorCode: foundation.UserDoNotExist,
			},
			store: &MockStore{
				Error: errors.New("error"),
			},
		},
		{
			name: "no User in DB",
			request: func() *http.Request {
				credentials := `{
					"email": "john4@gmail.com","password":"pass"
				}`
				req := request(t, "/login", "")

				base64 := base64.StdEncoding.EncodeToString([]byte(credentials))
				req.Header.Set("Authorization", "Basic "+base64)
				return req
			}(),
			response: &foundation.Response{
				Status:    http.StatusInternalServerError,
				Message:   "email not found",
				ErrorCode: foundation.UserDoNotExist,
			},
			store: &MockStore{},
		},
		{
			name: "authentication error",
			request: func() *http.Request {
				credentials := `{
					"email": "john4@gmail.com","password":"pass"
				}`
				req := request(t, "/login", "")

				base64 := base64.StdEncoding.EncodeToString([]byte(credentials))
				req.Header.Set("Authorization", "Basic "+base64)
				return req
			}(),
			response: &foundation.Response{
				Status:    http.StatusBadRequest,
				Message:   "email not found",
				ErrorCode: foundation.InvalidRequest,
			},
			store: &MockStore{
				UserResource: &store.UserResource{
					ID:        "123",
					FirstName: "Joe",
				},
			},
			authenticator: &MockAuthenticator{
				Error: errors.New("error"),
			},
		},
		{
			name: "authentication failed",
			request: func() *http.Request {
				credentials := `{
					"email": "john4@gmail.com","password":"pass"
				}`
				req := request(t, "/login", "")

				base64 := base64.StdEncoding.EncodeToString([]byte(credentials))
				req.Header.Set("Authorization", "Basic "+base64)
				return req
			}(),
			response: &foundation.Response{
				Status:    http.StatusBadRequest,
				Message:   "email not found",
				ErrorCode: foundation.UnAuthorised,
			},
			store: &MockStore{
				UserResource: &store.UserResource{
					ID:        "123",
					FirstName: "Joe",
				},
			},
			authenticator: &MockAuthenticator{},
		},
		{
			name: "authentication key not found",
			request: func() *http.Request {
				credentials := `{
					"email": "john4@gmail.com","password":"pass"
				}`
				req := request(t, "/login", "")

				base64 := base64.StdEncoding.EncodeToString([]byte(credentials))
				req.Header.Set("Authorization", "Basic "+base64)
				return req
			}(),
			response: &foundation.Response{
				Status:    http.StatusInternalServerError,
				Message:   errTokenGeneration.Error(),
				ErrorCode: foundation.KeyNotFound,
			},
			store: &MockStore{
				UserResource: &store.UserResource{
					ID:        "123",
					FirstName: "Joe",
				},
			},
			authenticator: &MockAuthenticator{
				ReturnVal: true,
			},
		},
	}

	logger := log.New(os.Stdout, "IDENTITY-TEST", log.LstdFlags)
	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			h := NewHandler(sc.store, sc.authenticator,
				&store.TokenConfig{
					Issuer: "TEST",
				}, logger)
			h.Login(rr, sc.request, nil)
			re := response(t, rr.Body)
			assert.Equal(t, sc.response, re)
		})
	}
}

func TestLoginAuthenticationKeyFound(t *testing.T) {
	logger := log.New(os.Stdout, "IDENTITY-LOGIN-TEST", log.LstdFlags)
	rr := httptest.NewRecorder()

	h := NewHandler(&MockStore{
		UserResource: &store.UserResource{
			ID:        "123",
			FirstName: "Joe",
		},
	}, &MockAuthenticator{ReturnVal: true},
		&store.TokenConfig{
			Issuer:         "TEST",
			KeyPath:        "../../../business/validation/testdata/",
			PrivateKeyName: "test_private.pem",
			PublicKeyName:  "test_public.pem",
		}, logger)

	credentials := `{
			"email": "john4@gmail.com","password":"pass"
		}`
	req := request(t, "/login", "")
	base64 := base64.StdEncoding.EncodeToString([]byte(credentials))
	req.Header.Set("Authorization", "Basic "+base64)

	h.Login(rr, req, nil)
	re := response(t, rr.Body)

	assert.Equal(t, &foundation.Response{
		Status: http.StatusOK}, re)
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
