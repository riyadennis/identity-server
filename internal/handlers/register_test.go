package handlers

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/riyadennis/identity-server/internal/store"
	"github.com/stretchr/testify/assert"
)

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

func TestUserDataFromRequest(t *testing.T){
	scenarios := []struct{
		name string
		request *http.Request
		expectedUser *store.User
		expectedError string
	}{
		{
			name: "empty request",
			request: nil,
			expectedUser: nil,
			expectedError: "empty request",
		},
		{
			name: "missing content type",
			request: &http.Request{
				Header: nil,
			},
			expectedUser: nil,
			expectedError: "invalid content type",
		},
		{
			name: "missing body",
			request: func() *http.Request{
				req := httptest.NewRequest("POST", "/register", nil )
				req.Header.Set("content-type", "application/json" )
				return req
			}(),
			expectedUser: nil,
			expectedError: "unexpected end of JSON input",
		},
		{
			name: "missing body",
			request: func() *http.Request{
				req := httptest.NewRequest("POST", "/register", nil )
				req.Header.Set("content-type", "application/json" )
				return req
			}(),
			expectedUser: nil,
			expectedError: "unexpected end of JSON input",
		},
	}
	for _, sc := range scenarios{
		t.Run(sc.name, func(t *testing.T) {
			u, err := userDataFromRequest(sc.request)
			if !cmp.Equal(u, sc.expectedUser){
				t.Errorf("expected user %v, got %v", sc.expectedUser, u)
			}
			if !cmp.Equal(err.Error(), sc.expectedError){
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
