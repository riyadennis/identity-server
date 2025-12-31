package validation

import (
	"errors"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"

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
func ValidateUser(u *store.User) error {
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

	if err := ValidateEmail(u.Email); err != nil {
		return err
	}
	if !u.Terms {
		return errTermsMissing
	}

	return nil
}

func ValidateEmail(email string) error {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !re.MatchString(email) {
		return errInvalidEmail
	}

	return nil
}

func ValidateToken(token string, tc *store.TokenConfig) (*jwt.RegisteredClaims, error) {
	if token == "" {
		return nil, errMissingToken
	}

	if token[len(BearerSchema):] == "" {
		return nil, errMissingBearerToken
	}
	claims := &jwt.RegisteredClaims{
		Issuer:    tc.Issuer,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(tc.TokenTTL)),
	}
	t, err := jwt.ParseWithClaims(
		token[len(BearerSchema):], claims, fetchKey(tc.KeyPath+tc.PublicKeyName),
	)
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, errInvalidToken
	}
	claims, ok := t.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, errInvalidToken
	}
	if claims.Issuer != tc.Issuer {
		return nil, errInvalidToken
	}
	if claims.ExpiresAt.Time.Before(time.Now().UTC()) {
		return nil, jwt.ErrTokenExpired
	}
	return claims, nil
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
			return nil, e
		}

		return publicKey, nil
	}
}
