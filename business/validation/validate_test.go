package validation

import (
	"errors"
	"testing"

	"github.com/riyadennis/identity-server/business/store"
)

func TestValidateUser(t *testing.T) {
	scenarios := []struct {
		name          string
		user          *store.UserRequest
		expectedError error
	}{
		{
			name:          "empty user",
			user:          nil,
			expectedError: errEmptyUser,
		},
		{
			name:          "missing first name",
			user:          &store.UserRequest{},
			expectedError: errMissingFirstName,
		},
		{
			name: "missing last name",
			user: &store.UserRequest{
				FirstName: "John",
			},
			expectedError: errMissingLastName,
		},
		{
			name: "missing email",
			user: &store.UserRequest{
				FirstName: "John",
				LastName:  "Doe",
			},
			expectedError: errMissingEmail,
		},
		{
			name: "invalid email",
			user: &store.UserRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "INVALID",
			},
			expectedError: errInvalidEmail,
		},
		{
			name: "missing terms",
			user: &store.UserRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@test.com",
			},
			expectedError: errTermsMissing,
		},
		{
			name: "missing terms",
			user: &store.UserRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@test.com",
				Terms:     true,
			},
			expectedError: nil,
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
