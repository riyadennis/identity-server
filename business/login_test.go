package business

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/riyadennis/identity-server/app/auth-api/mocks"
	"github.com/riyadennis/identity-server/business/store"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestUserCredentialsInDB(t *testing.T) {
	testCases := []struct {
		name          string
		email         string
		password      string
		mockStore     store.Store
		mockAuth      store.Authenticator
		expectedUser  *store.User
		expectedError error
	}{
		{
			name:          "user not found - read returns error",
			email:         "test@example.com",
			password:      "password123",
			mockStore:     &mocks.Store{Error: errors.New("error")},
			expectedError: errEmailNotFound,
		},
		{
			name:          "user not found - read returns nil",
			email:         "notfound@example.com",
			password:      "password123",
			mockStore:     &mocks.Store{},
			mockAuth:      &mocks.Authenticator{},
			expectedError: errEmailNotFound,
		},
		{
			name:     "authentication fails with error",
			email:    "test@example.com",
			password: "wrongpassword",
			mockStore: &mocks.Store{
				User: &store.User{
					ID:        "user123",
					Email:     "test@example.com",
					FirstName: "John",
					LastName:  "Doe",
				},
			},
			mockAuth:      &mocks.Authenticator{Error: errors.New("err")},
			expectedUser:  nil,
			expectedError: errInvalidPassword,
		},
		{
			name:     "authentication returns false",
			email:    "test@example.com",
			password: "wrongpassword",
			mockStore: &mocks.Store{
				User: &store.User{
					ID:        "user123",
					Email:     "test@example.com",
					FirstName: "John",
					LastName:  "Doe",
				},
			},
			mockAuth:      &mocks.Authenticator{ReturnVal: false},
			expectedError: errInvalidPassword,
		},
		{
			name:     "successful authentication",
			email:    "test@example.com",
			password: "correctpassword",
			mockStore: &mocks.Store{
				User: &store.User{
					ID:        "user123",
					Email:     "test@example.com",
					FirstName: "John",
					LastName:  "Doe",
				},
			},
			mockAuth: &mocks.Authenticator{ReturnVal: true},
			expectedUser: &store.User{
				ID:        "user123",
				Email:     "test@example.com",
				FirstName: "John",
				LastName:  "Doe",
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger := logrus.New()
			logger.SetOutput(os.Stderr)

			helper := NewHelper(tc.mockStore, tc.mockAuth, logger)
			user, err := helper.UserCredentialsInDB(context.Background(), tc.email, tc.password)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedUser, user)
		})
	}
}

func TestManageToken(t *testing.T) {
	// Create a temporary directory for test keys
	tempDir, err := os.MkdirTemp("", "test-keys")
	assert.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()
	futureExpiry := time.Now().Add(2 * time.Hour)
	testCases := []struct {
		name          string
		config        *store.TokenConfig
		userID        string
		mockStore     store.Store
		mockAuth      store.Authenticator
		expectedToken *store.Token
		expectedError error
	}{
		{
			name: "fetch token returns error",
			config: &store.TokenConfig{
				Issuer:         "test-issuer",
				KeyPath:        tempDir + "/",
				PrivateKeyName: "private.pem",
				PublicKeyName:  "public.pem",
			},
			userID:        "user123",
			mockStore:     &mocks.Store{Error: errors.New("error")},
			mockAuth:      &mocks.Authenticator{},
			expectedError: errors.New("database error"),
		},
		{
			name: "token exists and not expired - reuse existing token",
			config: &store.TokenConfig{
				Issuer:         "test-issuer",
				KeyPath:        tempDir + "/",
				PrivateKeyName: "private.pem",
				PublicKeyName:  "public.pem",
			},
			userID:    "user123",
			mockStore: &mocks.Store{},
			mockAuth: &mocks.Authenticator{Token: &store.TokenRecord{
				Id:     "token123",
				Token:  "existing-jwt-token",
				Expiry: futureExpiry,
				TTL:    "3600",
			}},
		},
		{
			name: "token exists but expired - generate new token",
			config: &store.TokenConfig{
				Issuer:         "test-issuer",
				KeyPath:        tempDir + "/",
				PrivateKeyName: "private.pem",
				PublicKeyName:  "public.pem",
			},
			userID:    "user123",
			mockStore: &mocks.Store{},
			mockAuth: func() *mocks.Authenticator {
				pastExpiry := time.Now().Add(-2 * time.Hour)
				return &mocks.Authenticator{
					Token: &store.TokenRecord{
						Id:     "token123",
						Token:  "expired-jwt-token",
						Expiry: pastExpiry,
						TTL:    "3600",
					},
				}
			}(),
			expectedError: nil,
		},
		{
			name: "token does not exist - generate new token",
			config: &store.TokenConfig{
				Issuer:         "test-issuer",
				KeyPath:        tempDir + "/",
				PrivateKeyName: "private.pem",
				PublicKeyName:  "public.pem",
			},
			userID: "user123",
			mockStore: &mocks.Store{
				User: &store.User{
					ID:        "user123",
					Email:     "test@example.com",
					FirstName: "John",
					LastName:  "Doe",
				},
			},
			mockAuth: &mocks.Authenticator{Token: &store.TokenRecord{
				Id:     "token123",
				Token:  "existing-jwt-token",
				Expiry: futureExpiry,
				TTL:    "3600",
			}},
			expectedError: nil,
		},
		{
			name: "save token fails",
			config: &store.TokenConfig{
				Issuer:         "test-issuer",
				KeyPath:        tempDir + "/",
				PrivateKeyName: "private.pem",
				PublicKeyName:  "public.pem",
			},
			userID: "user123",
			mockStore: &mocks.Store{
				User: &store.User{
					ID:        "user123",
					Email:     "test@example.com",
					FirstName: "John",
					LastName:  "Doe",
				},
			},
			mockAuth: &mocks.Authenticator{Token: &store.TokenRecord{
				Id:     "token123",
				Token:  "existing-jwt-token",
				Expiry: futureExpiry,
				TTL:    "3600",
			}},
			expectedError: errors.New("failed to save token"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger := logrus.New()
			logger.SetOutput(os.Stderr)

			helper := NewHelper(tc.mockStore, tc.mockAuth, logger)

			_, err := helper.ManageToken(context.Background(), tc.config, tc.userID)
			if err != nil {
				assert.EqualError(t, tc.expectedError, err.Error())
			}
			// Clean up generated keys if any
			privateKeyPath := filepath.Join(tc.config.KeyPath, tc.config.PrivateKeyName)
			publicKeyPath := filepath.Join(tc.config.KeyPath, tc.config.PublicKeyName)
			os.Remove(privateKeyPath)
			os.Remove(publicKeyPath)
		})
	}
}
