package validation

import (
	"errors"
	"os"
	"regexp"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"

	"github.com/riyadennis/identity-server/business/store"
)

var (
	errEmptyUser        = errors.New("empty user details")
	errMissingFirstName = errors.New("missing first name")
	errMissingLastName  = errors.New("missing last name")
	errMissingEmail     = errors.New("missing email")
	errInvalidEmail     = errors.New("invalid email")
	errTermsMissing     = errors.New("please select terms")

	errMissingToken       = errors.New("missing token in header")
	errMissingBearerToken = errors.New("missing bearer token in header")
	errInvalidToken       = errors.New("invalid token")
	errTokenKeyNotFound   = errors.New("key to validate token not found")
	errInvalidTokenMethod = errors.New("invalid token method")
)

const (
	// BearerSchema is expected prefix for token from authorisation header
	BearerSchema = "Bearer "
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

func ValidateToken(token string, tc *store.TokenConfig) error {
	if token == "" {
		return errMissingToken
	}

	if token[len(BearerSchema):] == "" {
		return errMissingBearerToken
	}

	t, err := jwt.ParseWithClaims(
		token[len(BearerSchema):],
		jwt.MapClaims{
			"exp": time.Now().UTC().Add(tc.TokenTTL).Unix(),
			"iss": tc.Issuer,
		}, fetchKey(tc.KeyPath+tc.PublicKeyName))

	if err != nil {
		return err
	}

	if !t.Valid {
		return errInvalidToken
	}

	return nil
}

func fetchKey(keyPath string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errInvalidTokenMethod
		}

		key, err := os.ReadFile(keyPath)
		if err != nil {
			return nil, errTokenKeyNotFound
		}

		publicKey, e := jwt.ParseRSAPublicKeyFromPEM(key)
		if e != nil {
			panic(e.Error())
		}

		return publicKey, nil
	}
}
