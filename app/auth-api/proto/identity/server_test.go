package identity

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/riyadennis/identity-server/app/auth-api/mocks"
	"github.com/riyadennis/identity-server/business/store"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const (
	testEmail    = "john.doe@gmail.com"
	testPassword = "password123"
	testUserID   = "test-user-id"
	testToken    = "test-access-token"
)

func TestNewServer(t *testing.T) {
	// Create a mock database connection
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	logger := logrus.New()
	tokenConfig := &store.TokenConfig{
		Issuer:         "test-issuer",
		KeyPath:        "/tmp/keys/",
		PrivateKeyName: "private.pem",
		PublicKeyName:  "public.pem",
	}

	// Test NewServer creates a valid server instance
	server := NewServer(logger, tokenConfig, &mocks.Store{}, &mocks.Authenticator{})

	assert.NotNil(t, server)
	assert.NotNil(t, server.Store)
	assert.NotNil(t, server.Authenticator)
	assert.NotNil(t, server.Logger)
	assert.NotNil(t, server.TokenConfig)
	assert.Equal(t, tokenConfig, server.TokenConfig)
	assert.Equal(t, logger, server.Logger)
}

func TestLogin(t *testing.T) {
	scenarios := []struct {
		name             string
		request          *LoginRequest
		mockStore        store.Store
		mockAuth         store.Authenticator
		expectedError    error
		expectedResponse *LoginResponse
	}{
		{
			name: "invalid email",
			request: func() *LoginRequest {
				// Test with invalid email format
				invalidEmail := "invalid-email"
				password := testPassword
				return &LoginRequest{
					Email:    &invalidEmail,
					Password: &password,
				}
			}(),
			expectedError: errors.New("invalid email"),
		},
		{
			name: "user not found",
			request: func() *LoginRequest {
				// Test with invalid email format
				email := testEmail
				password := testPassword
				return &LoginRequest{
					Email:    &email,
					Password: &password,
				}
			}(),
			mockStore: &mocks.Store{
				Error: errors.New("user not found"),
			},
			expectedError: errors.New("email not found"),
		},
		{
			name: "authentication failed",
			request: func() *LoginRequest {
				// Test with invalid email format
				email := testEmail
				password := testPassword
				return &LoginRequest{
					Email:    &email,
					Password: &password,
				}
			}(),
			mockStore: &mocks.Store{},
			mockAuth: &mocks.Authenticator{
				ReturnVal: false,
				Error:     errors.New("authentication failed"),
			},
			expectedError: errors.New("email not found"),
		},
		{
			name: "success",
			request: func() *LoginRequest {
				// Test with invalid email format
				email := testEmail
				password := testPassword
				return &LoginRequest{
					Email:    &email,
					Password: &password,
				}
			}(),
			mockStore: &mocks.Store{
				User: &store.User{Email: testEmail, Password: testPassword},
			},
			mockAuth: &mocks.Authenticator{
				ReturnVal: true,
				Token: &store.TokenRecord{
					Expiry: time.Now().Add(2 * time.Hour),
					TTL:    "123",
				},
			},
			expectedResponse: func() *LoginResponse {
				status := int32(http.StatusOK)
				exp := time.Now().Add(2 * time.Hour).String()
				tType := "Bearer"
				return &LoginResponse{
					Status:    &status,
					Expiry:    &exp,
					TokenType: &tType,
				}
			}(),
		},
	}
	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			server := &Server{
				Logger:        logrus.New(),
				Store:         sc.mockStore,
				Authenticator: sc.mockAuth,
				TokenConfig: &store.TokenConfig{
					Issuer:         "test-issuer",
					KeyPath:        "/tmp/",
					PrivateKeyName: "private.pem",
					PublicKeyName:  "public.pem",
				},
			}
			resp, err := server.Login(context.Background(), sc.request)
			assert.Equal(t, sc.expectedError, err)
			checkResponse(t, sc.expectedResponse, resp)
		})
	}
}

func checkResponse(t *testing.T, expected, actual *LoginResponse) {
	if expected == nil && actual != nil {
		t.Error("unexpected response")
	} else if actual == nil && expected != nil {
		t.Error("expected response but not found")
	} else if expected != nil {
		assert.Equal(t, expected.Status, actual.Status)
		assert.Equal(t, expected.Status, actual.Status)
	}
}
func TestLogin_InvalidTokenTTL(t *testing.T) {
	mockStore := &mocks.Store{
		User: &store.User{
			ID:    testUserID,
			Email: testEmail,
		},
	}

	expiryTime := time.Now().Add(120 * time.Hour)
	mockAuth := &mocks.Authenticator{
		ReturnVal: true,
		Token: &store.TokenRecord{
			Id:     "token-id",
			UserId: testUserID,
			Token:  testToken,
			TTL:    "invalid-ttl", // This will cause strconv.ParseInt to fail
			Expiry: expiryTime,
		},
	}

	server := &Server{
		Store:         mockStore,
		Authenticator: mockAuth,
		Logger:        logrus.New(),
		TokenConfig: &store.TokenConfig{
			Issuer:         "test-issuer",
			KeyPath:        "/tmp/",
			PrivateKeyName: "private.pem",
			PublicKeyName:  "public.pem",
		},
	}

	email := testEmail
	password := testPassword
	req := &LoginRequest{
		Email:    &email,
		Password: &password,
	}

	resp, err := server.Login(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestMustEmbedUnimplementedIdentityServer(t *testing.T) {
	server := &Server{}

	// Test that this method exists and doesn't panic
	assert.NotPanics(t, func() {
		server.mustEmbedUnimplementedIdentityServer()
	})
}
