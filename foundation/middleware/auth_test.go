package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/riyadennis/identity-server/business/store"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testKeyPath = "../../business/validation/testdata/"

func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func newAuthConfig() *AuthConfig {
	return &AuthConfig{
		TokenConfig: &store.TokenConfig{
			KeyPath:        testKeyPath,
			PrivateKeyName: "test_private.pem",
			PublicKeyName:  "test_public.pem",
			Issuer:         "test",
		},
		Logger: logrus.New(),
	}
}

func validToken(t *testing.T) string {
	t.Helper()
	pemBytes, err := os.ReadFile(testKeyPath + "test_private.pem")
	require.NoError(t, err)
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(pemBytes)
	require.NoError(t, err)

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		Issuer:    "test",
		Subject:   "user-123",
		Audience:  jwt.ClaimStrings{"local"},
	}
	tok, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(privateKey)
	require.NoError(t, err)
	return tok
}

func TestAuth_MissingToken(t *testing.T) {
	ac := newAuthConfig()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	ac.Auth(okHandler()).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuth_InvalidToken(t *testing.T) {
	ac := newAuthConfig()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer not-a-real-token")
	rr := httptest.NewRecorder()
	ac.Auth(okHandler()).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuth_ValidToken(t *testing.T) {
	ac := newAuthConfig()
	tok := "Bearer " + validToken(t)
	reached := false
	handler := ac.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reached = true
		assert.NotNil(t, r.Context().Value(UserClaimsKey))
		assert.Equal(t, tok, r.Context().Value(AccessTokenKey))
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", tok)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.True(t, reached)
}