package handlers

import (
	"bytes"
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

	"github.com/riyadennis/identity-server/business"
	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/foundation"
)

var (
	testEmail string
)

func TestRegister(t *testing.T) {
	scenarios := []struct {
		name             string
		req              *http.Request
		store            store.Store
		expectedResponse *foundation.Response
	}{
		{
			name: "missing email",
			req: registerPayLoad(t, &store.User{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "",
			}),
			expectedResponse: foundation.NewResponse(http.StatusBadRequest,
				"missing email",
				foundation.ValidationFailed),
		},
		{
			name: "missing first name",
			req: registerPayLoad(t, &store.User{
				FirstName: "",
				LastName:  "Doe",
				Email:     testEmail,
			}),
			expectedResponse: foundation.NewResponse(http.StatusBadRequest,
				"missing first name",
				foundation.ValidationFailed),
		},
		{
			name: "missing last name",
			req: registerPayLoad(t, &store.User{
				FirstName: "John",
				LastName:  "",
				Email:     testEmail,
			}),
			expectedResponse: foundation.NewResponse(http.StatusBadRequest,
				"missing last name",
				foundation.ValidationFailed),
		},
		{
			name: "invalid email",
			req: func() *http.Request {
				u := user(t)
				u.Email = "joh@dom"
				return registerPayLoad(t, u)
			}(),
			expectedResponse: foundation.NewResponse(http.StatusBadRequest,
				"invalid email",
				foundation.ValidationFailed),
		},
		{
			name: "error reading db",
			req: func() *http.Request {
				u := user(t)
				u.Email = "joh@doe.com"
				return registerPayLoad(t, u)
			}(),
			store:            &MockStore{Error: errors.New("error")},
			expectedResponse: foundation.NewResponse(http.StatusBadRequest, "error", foundation.ValidationFailed),
		},
		{
			name: "duplicate email",
			req: func() *http.Request {
				u := user(t)
				u.Email = "joh@doe.com"
				return registerPayLoad(t, u)
			}(),
			store: &MockStore{User: &store.User{Email: "joh@doe.com"}},
			expectedResponse: foundation.NewResponse(
				http.StatusBadRequest,
				"email already exists",
				foundation.EmailAlreadyExists),
		},
	}
	logger := log.New(os.Stdout, "IDENTITY-TEST", log.LstdFlags)
	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			h := NewHandler(sc.store, &MockAuthenticator{},
				&store.TokenConfig{
					Issuer:  "TEST",
					KeyPath: os.Getenv("KEY_PATH"),
				}, logger)
			h.Register(w, sc.req, nil)
			resp := response(t, w.Body)
			// TODO assert message also
			assert.Equal(t, sc.expectedResponse, resp)
		})
	}
}

func registerPayLoad(t *testing.T, u *store.User) *http.Request {
	var buff bytes.Buffer
	err := json.NewEncoder(&buff).Encode(u)
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest("POST", "/register", strings.NewReader(buff.String()))
	req.Header.Set("content-type", "application/json")
	return req
}

func responseFromHTTP(t *testing.T, data io.Reader) *foundation.Response {
	t.Helper()

	resp := &foundation.Response{}
	b, err := io.ReadAll(data)
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal(b, resp)
	if err != nil {
		t.Error(err)
	}

	if resp.Message != "" {
		return resp
	}

	return nil
}

func TestGeneratePassword(t *testing.T) {
	t.Helper()

	pass, err := business.GeneratePassword()
	if err != nil {
		t.Error(err)
	}
	if pass == "" {
		t.Error("empty password")
	}
}

func user(t *testing.T) *store.User {
	t.Helper()

	u := &store.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     testEmail,
		Password:  "testPassword",
		Company:   "testCompany",
		PostCode:  "E112QD",
		Terms:     true,
	}
	return u
}
