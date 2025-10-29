package business

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/riyadennis/identity-server/business/store"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStore is a mock implementation of store.Store interface
type MockStore struct {
	mock.Mock
}

func (m *MockStore) Insert(ctx context.Context, u *store.User) (*store.User, error) {
	args := m.Called(ctx, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*store.User), args.Error(1)
}

func (m *MockStore) Read(ctx context.Context, email string) (*store.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*store.User), args.Error(1)
}

func (m *MockStore) Delete(id string) (int64, error) {
	args := m.Called(id)
	return args.Get(0).(int64), args.Error(1)
}

// MockAuthenticator is a mock implementation of store.Authenticator interface
type MockAuthenticator struct {
	mock.Mock
}

func (m *MockAuthenticator) Authenticate(email, password string) (bool, error) {
	args := m.Called(email, password)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthenticator) FetchLoginToken(userID string) (*store.TokenRecord, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*store.TokenRecord), args.Error(1)
}

func (m *MockAuthenticator) SaveLoginToken(ctx context.Context, t *store.TokenRecord) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}

func TestUserCredentialsInDB(t *testing.T) {
	testCases := []struct {
		name          string
		email         string
		password      string
		setupMocks    func(*MockStore, *MockAuthenticator)
		expectedUser  *store.User
		expectedError error
	}{
		{
			name:     "user not found - read returns error",
			email:    "test@example.com",
			password: "password123",
			setupMocks: func(ms *MockStore, ma *MockAuthenticator) {
				ms.On("Read", mock.Anything, "test@example.com").
					Return(nil, errors.New("database error"))
			},
			expectedError: errEmailNotFound,
		},
		{
			name:     "user not found - read returns nil",
			email:    "notfound@example.com",
			password: "password123",
			setupMocks: func(ms *MockStore, ma *MockAuthenticator) {
				ms.On("Read", mock.Anything, "notfound@example.com").
					Return(nil, nil)
			},
			expectedError: errEmailNotFound,
		},
		{
			name:     "authentication fails with error",
			email:    "test@example.com",
			password: "wrongpassword",
			setupMocks: func(ms *MockStore, ma *MockAuthenticator) {
				user := &store.User{
					ID:        "user123",
					Email:     "test@example.com",
					FirstName: "John",
					LastName:  "Doe",
				}
				ms.On("Read", mock.Anything, "test@example.com").
					Return(user, nil)
				ma.On("Authenticate", "test@example.com", "wrongpassword").
					Return(false, errors.New("authentication error"))
			},
			expectedUser:  nil,
			expectedError: errInvalidPassword,
		},
		{
			name:     "authentication returns false",
			email:    "test@example.com",
			password: "wrongpassword",
			setupMocks: func(ms *MockStore, ma *MockAuthenticator) {
				user := &store.User{
					ID:        "user123",
					Email:     "test@example.com",
					FirstName: "John",
					LastName:  "Doe",
				}
				ms.On("Read", mock.Anything, "test@example.com").
					Return(user, nil)
				ma.On("Authenticate", "test@example.com", "wrongpassword").
					Return(false, nil)
			},
			expectedError: errInvalidPassword,
		},
		{
			name:     "successful authentication",
			email:    "test@example.com",
			password: "correctpassword",
			setupMocks: func(ms *MockStore, ma *MockAuthenticator) {
				user := &store.User{
					ID:        "user123",
					Email:     "test@example.com",
					FirstName: "John",
					LastName:  "Doe",
				}
				ms.On("Read", mock.Anything, "test@example.com").
					Return(user, nil)
				ma.On("Authenticate", "test@example.com", "correctpassword").
					Return(true, nil)
			},
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
			mockStore := new(MockStore)
			mockAuth := new(MockAuthenticator)
			logger := logrus.New()
			logger.SetOutput(os.Stderr)

			tc.setupMocks(mockStore, mockAuth)

			helper := NewHelper(mockStore, mockAuth, logger)
			user, err := helper.UserCredentialsInDB(context.Background(), tc.email, tc.password)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedUser, user)

			mockStore.AssertExpectations(t)
			mockAuth.AssertExpectations(t)
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
		setupMocks    func(*MockStore, *MockAuthenticator)
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
			userID: "user123",
			setupMocks: func(ms *MockStore, ma *MockAuthenticator) {
				ma.On("FetchLoginToken", "user123").
					Return(nil, errors.New("database error"))
			},
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
			userID: "user123",
			setupMocks: func(ms *MockStore, ma *MockAuthenticator) {
				existingToken := &store.TokenRecord{
					Id:     "token123",
					Token:  "existing-jwt-token",
					Expiry: futureExpiry,
					TTL:    "3600",
				}
				ma.On("FetchLoginToken", "user123").
					Return(existingToken, nil)
			},
		},
		{
			name: "token exists but expired - generate new token",
			config: &store.TokenConfig{
				Issuer:         "test-issuer",
				KeyPath:        tempDir + "/",
				PrivateKeyName: "private.pem",
				PublicKeyName:  "public.pem",
			},
			userID: "user123",
			setupMocks: func(ms *MockStore, ma *MockAuthenticator) {
				pastExpiry := time.Now().Add(-2 * time.Hour)
				expiredToken := &store.TokenRecord{
					Id:     "token123",
					Token:  "expired-jwt-token",
					Expiry: pastExpiry,
					TTL:    "3600",
				}
				ma.On("FetchLoginToken", "user123").
					Return(expiredToken, nil)
				ma.On("SaveLoginToken", mock.Anything, mock.AnythingOfType("*store.TokenRecord")).
					Return(nil)
			},
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
			setupMocks: func(ms *MockStore, ma *MockAuthenticator) {
				ma.On("FetchLoginToken", "user123").
					Return(nil, nil)
				ma.On("SaveLoginToken", mock.Anything, mock.AnythingOfType("*store.TokenRecord")).
					Return(nil)
			},
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
			setupMocks: func(ms *MockStore, ma *MockAuthenticator) {
				ma.On("FetchLoginToken", "user123").
					Return(nil, nil)
				ma.On("SaveLoginToken", mock.Anything, mock.AnythingOfType("*store.TokenRecord")).
					Return(errors.New("failed to save token"))
			},
			expectedError: errors.New("failed to save token"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockStore := new(MockStore)
			mockAuth := new(MockAuthenticator)
			logger := logrus.New()
			logger.SetOutput(os.Stderr)

			tc.setupMocks(mockStore, mockAuth)

			helper := NewHelper(mockStore, mockAuth, logger)

			_, err := helper.ManageToken(context.Background(), tc.config, tc.userID)
			if err != nil {
				assert.EqualError(t, tc.expectedError, err.Error())
			}
			mockStore.AssertExpectations(t)
			mockAuth.AssertExpectations(t)

			// Clean up generated keys if any
			privateKeyPath := filepath.Join(tc.config.KeyPath, tc.config.PrivateKeyName)
			publicKeyPath := filepath.Join(tc.config.KeyPath, tc.config.PublicKeyName)
			os.Remove(privateKeyPath)
			os.Remove(publicKeyPath)
		})
	}
}
