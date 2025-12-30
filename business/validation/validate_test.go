package validation

import (
	"errors"
	"os"
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/riyadennis/identity-server/business/store"
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
		{
			name: "missing terms",
			user: &store.User{
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

func TestValidateToken(t *testing.T) {
	// privateKey, _ := loadTestKeys(t)
	scenarios := []struct {
		name          string
		token         string
		tokenConfig   *store.TokenConfig
		expectedError error
	}{
		{
			name:          "empty token",
			token:         "",
			tokenConfig:   &store.TokenConfig{},
			expectedError: errMissingToken,
		},
		{
			name:          "missing bearer token",
			token:         "Bearer ",
			tokenConfig:   &store.TokenConfig{},
			expectedError: errMissingBearerToken,
		},
		{
			name:  "invalid token format",
			token: "Bearer invalid.token.format",
			tokenConfig: &store.TokenConfig{
				TokenTTL:      time.Hour,
				Issuer:        "test-issuer",
				KeyPath:       "./testdata/",
				PublicKeyName: "test_public.pem",
			},
			expectedError: jwt.ErrTokenMalformed,
		},
		{
			name: "missing key file",
			token: func() string {
				validToken := generateTestToken(t, "test-issuer", time.Now().UTC().Add(1*time.Hour))
				return validToken
			}(),
			tokenConfig: &store.TokenConfig{
				TokenTTL:      time.Hour,
				Issuer:        "test-issuer",
				KeyPath:       "./testdata/",
				PublicKeyName: "nonexistent.pem",
			},
			expectedError: errTokenKeyNotFound,
		},
		{
			name: "expired token",
			token: func() string {
				expiredToken := generateTestToken(t, "test-issuer", time.Now().UTC().Add(-1*time.Hour))
				return expiredToken
			}(),
			tokenConfig: &store.TokenConfig{
				TokenTTL:      time.Hour,
				Issuer:        "test-issuer",
				KeyPath:       "./testdata/",
				PublicKeyName: "test_public.pem",
			},
			expectedError: jwt.ErrTokenExpired,
		},
		{
			name: "wrong issuer",
			token: func() string {
				wrongIssuerToken := generateTestToken(t, "wrong-issuer", time.Now().UTC().Add(1*time.Hour))
				return wrongIssuerToken
			}(),
			tokenConfig: &store.TokenConfig{
				TokenTTL:      time.Hour,
				Issuer:        "test-issuer",
				KeyPath:       "./testdata/",
				PublicKeyName: "test_public.pem",
			},
			expectedError: errInvalidToken,
		},
		{
			name: "valid token",
			token: func() string {
				validToken := generateTestToken(t, "test-issuer", time.Now().UTC().Add(1*time.Hour))
				return validToken
			}(),
			tokenConfig: &store.TokenConfig{
				TokenTTL:      time.Hour,
				Issuer:        "test-issuer",
				KeyPath:       "./testdata/",
				PublicKeyName: "test_public.pem",
			},
			expectedError: nil,
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			_, err := ValidateToken(sc.token, sc.tokenConfig)
			if !errors.Is(err, sc.expectedError) {
				t.Fatalf("expected err %v, got %v", sc.expectedError, err)
			}
		})
	}
}

func generateTestToken(t *testing.T, issuer string, ttl time.Time) string {
	t.Helper()

	privateKeyData, err := os.ReadFile("testdata/test_private.pem")
	assert.NoError(t, err)

	signedToken, err := store.GenerateToken(logrus.New(), issuer, privateKeyData, ttl)
	assert.NoError(t, err)

	return "Bearer " + signedToken.AccessToken
}
