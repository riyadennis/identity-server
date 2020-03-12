package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"

	_ "github.com/mattn/go-sqlite3"
	"github.com/riyadennis/identity-server/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	scenarios := []struct {
		name             string
		req              *http.Request
		expectedResponse *Response
	}{
		{
			name: "empty request",
			req:  nil,
			expectedResponse: newResponse(400,
				"empty request",
				"invalid-request"),
		},
		{
			// one validation check
			// rest of the rules are
			// testing in TestValidate
			name: "missing email",
			req: func() *http.Request {
				u := &store.User{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "",
				}
				req := httptest.NewRequest("POST", "/register",
					registerPayLoad(t, u))
				req.Header.Set("content-type", "application/json")
				return req
			}(),
			expectedResponse: newResponse(400,
				"missing email",
				"invalid-user-data"),
		},
	}
	w := httptest.NewRecorder()
	for _, sc := range scenarios {
		Register(w, sc.req, nil)
		resp := responseFromHttp(t, w.Body)
		if !cmp.Equal(resp, sc.expectedResponse) {
			t.Errorf("unexpected response,got %v, want %v", resp,
				sc.expectedResponse)
		}
	}
}

func registerPayLoad(t *testing.T, u *store.User) io.Reader {
	jB, err := json.Marshal(u)
	if err != nil {
		t.Error(err)
	}
	return bytes.NewReader(jB)
}

func responseFromHttp(t *testing.T, data io.Reader) *Response {
	resp := &Response{}
	b, err := ioutil.ReadAll(data)
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal(b, resp)
	if err != nil {
		t.Error(err)
	}
	return resp
}

func TestValidateUser(t *testing.T) {
	scenarios := []struct {
		name        string
		user        *store.User
		expectedErr error
	}{
		{
			name:        "nil user",
			user:        nil,
			expectedErr: errors.New("empty user details"),
		},
		{
			name: "missing first name",
			user: func() *store.User {
				u := user(t)
				u.FirstName = ""
				return u
			}(),
			expectedErr: errors.New("missing first name"),
		},
		{
			name: "missing last name",
			user: func() *store.User {
				u := user(t)
				u.LastName = ""
				return u
			}(),
			expectedErr: errors.New("missing last name"),
		},
		{
			name: "missing email",
			user: func() *store.User {
				u := user(t)
				u.Email = ""
				return u
			}(),
			expectedErr: errors.New("missing email"),
		},
		{
			name: "missing terms",
			user: func() *store.User {
				u := user(t)
				u.Terms = false
				return u
			}(),
			expectedErr: errors.New("missing terms"),
		},
		{
			name: "invalid email",
			user: func() *store.User {
				u := user(t)
				u.Email = "invalid"
				return u
			}(),
			expectedErr: errors.New("invalid email"),
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			err := validateUser(sc.user)
			assert.Equal(t, sc.expectedErr, err)
		})
	}

}

func TestGeneratePassword(t *testing.T) {
	pass, err := generatePassword()
	if err != nil {
		t.Error(err)
	}
	if pass == "" {
		t.Error("empty password")
	}
}

func TestUserDataFromRequest(t *testing.T) {
	scenarios := []struct {
		name          string
		request       *http.Request
		expectedUser  *store.User
		expectedError string
	}{
		{
			name:          "empty request",
			request:       nil,
			expectedUser:  nil,
			expectedError: "empty request",
		},
		{
			name: "missing content type",
			request: &http.Request{
				Header: nil,
			},
			expectedUser:  nil,
			expectedError: "invalid content type",
		},
		{
			name: "missing body",
			request: func() *http.Request {
				req := httptest.NewRequest("POST", "/register", nil)
				req.Header.Set("content-type", "application/json")
				return req
			}(),
			expectedUser:  nil,
			expectedError: "unexpected end of JSON input",
		},
		{
			name: "missing body",
			request: func() *http.Request {
				req := httptest.NewRequest("POST", "/register", nil)
				req.Header.Set("content-type", "application/json")
				return req
			}(),
			expectedUser:  nil,
			expectedError: "unexpected end of JSON input",
		},
	}
	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			u, err := userDataFromRequest(sc.request)
			if !cmp.Equal(u, sc.expectedUser) {
				t.Errorf("expected user %v, got %v", sc.expectedUser, u)
			}
			if !cmp.Equal(err.Error(), sc.expectedError) {
				t.Errorf("expected error %v, got %v", sc.expectedError, err)
			}
		})
	}
}

func user(t *testing.T) *store.User {
	t.Helper()
	u := &store.User{
		FirstName:        "John",
		LastName:         "Doe",
		Email:            "joh@doe.com",
		Password:         "testPassword",
		Company:          "testCompany",
		PostCode:         "E112QD",
		Terms:            true,
		RegistrationDate: "",
	}
	return u
}
