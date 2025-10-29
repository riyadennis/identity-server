package rest

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/foundation"
)

const (
	testEmail    = "john.doe@gmail.com"
	testPassword = "pass"
)

type MockStore struct {
	Error error
	*store.User
}

func (m *MockStore) Insert(ctx context.Context, u *store.User) (*store.User, error) {
	return m.User, m.Error
}

func (m *MockStore) Read(ctx context.Context, email string) (*store.User, error) {
	return m.User, m.Error
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
func (ma *MockAuthenticator) FetchLoginToken(userID string) (*store.TokenRecord, error) {
	return nil, nil
}
func (ma *MockAuthenticator) SaveLoginToken(ctx context.Context, t *store.TokenRecord) error {
	return nil
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
			name:     "invalid email",
			request:  loginRequest(t, "invalid", "pass"),
			response: expectedResponse(t, "invalid email"),
		},
		{
			name:     "login DB error",
			request:  loginRequest(t, testEmail, testPassword),
			response: expectedResponse(t, "email not found"),
			store: &MockStore{
				Error: errors.New("error"),
			},
		},
		{
			name:     "authentication error",
			request:  loginRequest(t, testEmail, testPassword),
			response: expectedResponse(t, "invalid password"),
			store: &MockStore{
				User: &store.User{
					ID:        "123",
					FirstName: "Joe",
				},
			},
			authenticator: &MockAuthenticator{
				Error: errors.New("error"),
			},
		},
		{
			name:    "authentication key not found",
			request: loginRequest(t, testEmail, testPassword),
			response: &foundation.Response{
				Status:    http.StatusInternalServerError,
				Message:   errTokenGeneration.Error(),
				ErrorCode: foundation.TokenError,
			},
			store: &MockStore{
				User: &store.User{
					ID:        "123",
					FirstName: "Joe",
				},
			},
			authenticator: &MockAuthenticator{
				ReturnVal: true,
			},
		},
	}

	logger := logrus.New()
	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			h := NewHandler(sc.store, sc.authenticator,
				&store.TokenConfig{
					Issuer: "TEST",
				}, logger)
			h.Login(rr, sc.request)
			re := response(t, rr.Body)
			assert.Equal(t, sc.response, re)
		})
	}
}

func TestLoginAuthenticationKeyFound(t *testing.T) {
	logger := logrus.New()
	rr := httptest.NewRecorder()

	h := NewHandler(&MockStore{
		User: &store.User{
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

	h.Login(rr, loginRequest(t, "john@gmail.com", "pass"))
	re := response(t, rr.Body)

	assert.Equal(t, &foundation.Response{
		Status: http.StatusOK}, re)
}

func loginRequest(t *testing.T, email, password string) *http.Request {
	t.Helper()
	credentials := email + ":" + password
	req := request(t, "/login", "")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(credentials)))
	return req
}
func expectedResponse(t *testing.T, message string) *foundation.Response {
	t.Helper()
	return &foundation.Response{
		Status:    http.StatusBadRequest,
		Message:   message,
		ErrorCode: foundation.InvalidRequest,
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
