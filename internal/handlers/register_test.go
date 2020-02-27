package handlers

import (
	"errors"
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
