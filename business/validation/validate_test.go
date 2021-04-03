package validation

import (
	"errors"
	"github.com/riyadennis/identity-server/business/store"
	"testing"
)

func TestValidateUser(t *testing.T) {
	scenarios := []struct {
		name          string
		user          *store.User
		expectedError error
	}{
		{
			name:          "empty user",
			user:          nil,
			expectedError: errEmptyUser,
		},
		{
			name:          "missing first name",
			user:          &store.User{},
			expectedError: errMissingFirstName,
		},
		{
			name: "missing last name",
			user: &store.User{
				FirstName: "John",
			},
			expectedError: errMissingLastName,
		},
		{
			name: "missing email",
			user: &store.User{
				FirstName: "John",
				LastName:  "Doe",
			},
			expectedError: errMissingEmail,
		},
		{
			name: "invalid email",
			user: &store.User{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "INVALID",
			},
			expectedError: errInvalidEmail,
		},
		{
			name: "missing terms",
			user: &store.User{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@test.com",
			},
			expectedError: errTermsMissing,
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			err := ValidateUser(sc.user)
			if !errors.Is(err, sc.expectedError) {
				t.Fatalf("expected err %v, got %v", sc.expectedError, err)
			}
		})
	}
}
