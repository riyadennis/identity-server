package validation

import (
	"errors"
	"regexp"

	"github.com/riyadennis/identity-server/business/store"
)

// ValidateUser checks registration request validity
func ValidateUser(u *store.User) error {
	if u == nil {
		return errors.New("empty user details")
	}
	if u.FirstName == "" {
		return errors.New("missing first name")
	}
	if u.LastName == "" {
		return errors.New("missing last name")
	}
	if u.Email == "" {
		return errors.New("missing email")
	}

	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !re.MatchString(u.Email) {
		return errors.New("invalid email")
	}

	if !u.Terms {
		return errors.New("please select terms")
	}
	return nil
}