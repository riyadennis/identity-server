package validation

import (
	"errors"
	"regexp"

	"github.com/riyadennis/identity-server/business/store"
)

var (
	errEmptyUser        = errors.New("empty user details")
	errMissingFirstName = errors.New("missing first name")
	errMissingLastName  = errors.New("missing last name")
	errMissingEmail     = errors.New("missing email")
	errInvalidEmail     = errors.New("invalid email")
	errTermsMissing     = errors.New("please select terms")
)

// ValidateUser checks registration request validity
func ValidateUser(u *store.UserRequest) error {
	if u == nil {
		return errEmptyUser
	}
	if u.FirstName == "" {
		return errMissingFirstName
	}
	if u.LastName == "" {
		return errMissingLastName
	}
	if u.Email == "" {
		return errMissingEmail
	}

	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !re.MatchString(u.Email) {
		return errInvalidEmail
	}

	if !u.Terms {
		return errTermsMissing
	}

	return nil
}
