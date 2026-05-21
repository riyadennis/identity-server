package graph

import (
	"context"
	"errors"
	"testing"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/riyadennis/identity-server/app/gql/graph/model"
	"github.com/riyadennis/identity-server/app/mocks"
	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/foundation/middleware"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testEmail    = "john@example.com"
	testPassword = "secret123"
)

func newResolver(st store.Store, au store.Authenticator, tc *store.TokenConfig) *Resolver {
	return NewResolver(logrus.New(), tc, st, au)
}

func tokenConfig() *store.TokenConfig {
	return &store.TokenConfig{
		KeyPath:        "../../../business/validation/testdata/",
		PrivateKeyName: "test_private.pem",
		PublicKeyName:  "test_public.pem",
		Issuer:         "test",
	}
}

// --- Login ---

func TestLogin_InvalidEmail(t *testing.T) {
	r := &mutationResolver{newResolver(nil, nil, tokenConfig())}
	email := "not-an-email"
	password := testPassword
	_, err := r.Login(context.Background(), model.LoginInput{Email: &email, Password: &password})
	require.Error(t, err)
}

func TestLogin_StoreError(t *testing.T) {
	r := &mutationResolver{newResolver(
		&mocks.Store{Error: errors.New("db down")},
		&mocks.Authenticator{},
		tokenConfig(),
	)}
	email := testEmail
	password := testPassword
	_, err := r.Login(context.Background(), model.LoginInput{Email: &email, Password: &password})
	require.Error(t, err)
}

func TestLogin_AuthError(t *testing.T) {
	r := &mutationResolver{newResolver(
		&mocks.Store{User: &store.User{ID: "1", Email: testEmail}},
		&mocks.Authenticator{Error: errors.New("bad password")},
		tokenConfig(),
	)}
	email := testEmail
	password := testPassword
	_, err := r.Login(context.Background(), model.LoginInput{Email: &email, Password: &password})
	require.Error(t, err)
}

func TestLogin_Success(t *testing.T) {
	r := &mutationResolver{newResolver(
		&mocks.Store{User: &store.User{ID: "1", Email: testEmail}},
		&mocks.Authenticator{ReturnVal: true},
		tokenConfig(),
	)}
	email := testEmail
	password := testPassword
	resp, err := r.Login(context.Background(), model.LoginInput{Email: &email, Password: &password})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
}

// --- Register ---

func TestRegister_ValidationError(t *testing.T) {
	r := &mutationResolver{newResolver(&mocks.Store{}, &mocks.Authenticator{}, tokenConfig())}
	// missing last name — validation should fail
	_, err := r.Register(context.Background(), model.RegisterInput{
		FirstName: "John",
		Email:     testEmail,
		Password:  testPassword,
		Terms:     true,
	})
	require.Error(t, err)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	existing := &store.User{Email: testEmail}
	r := &mutationResolver{newResolver(
		&mocks.Store{User: existing},
		&mocks.Authenticator{},
		tokenConfig(),
	)}
	_, err := r.Register(context.Background(), model.RegisterInput{
		FirstName: "John",
		LastName:  "Doe",
		Email:     testEmail,
		Password:  testPassword,
		Terms:     true,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "email already exists")
}

func TestRegister_StoreReadError(t *testing.T) {
	r := &mutationResolver{newResolver(
		&mocks.Store{Error: errors.New("db error")},
		&mocks.Authenticator{},
		tokenConfig(),
	)}
	_, err := r.Register(context.Background(), model.RegisterInput{
		FirstName: "John",
		LastName:  "Doe",
		Email:     testEmail,
		Password:  testPassword,
		Terms:     true,
	})
	require.Error(t, err)
}

func TestRegister_Success(t *testing.T) {
	company := "Acme"
	postCode := "SW1A"
	created := &store.User{
		ID:        "abc-123",
		FirstName: "John",
		LastName:  "Doe",
		Email:     testEmail,
		Company:   company,
		PostCode:  postCode,
		Terms:     true,
		CreatedAt: "2026-01-01 00:00:00",
	}
	// Read returns empty user (email doesn't exist yet), Insert returns created user
	st := &insertMockStore{created: created}
	r := &mutationResolver{newResolver(st, &mocks.Authenticator{}, tokenConfig())}

	resp, err := r.Register(context.Background(), model.RegisterInput{
		FirstName: "John",
		LastName:  "Doe",
		Email:     testEmail,
		Password:  testPassword,
		Company:   &company,
		PostCode:  &postCode,
		Terms:     true,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "abc-123", *resp.ID)
	assert.Equal(t, testEmail, *resp.Email)
}

// insertMockStore returns empty user on Read (not found) and created on Insert.
type insertMockStore struct {
	created *store.User
}

func (s *insertMockStore) Insert(_ context.Context, _ *store.User) (*store.User, error) {
	return s.created, nil
}

func (s *insertMockStore) Read(_ context.Context, _ string) (*store.User, error) {
	return &store.User{}, nil // empty — email not found
}

func (s *insertMockStore) Retrieve(_ context.Context, _ string) (*store.User, error) {
	return s.created, nil
}
func (s *insertMockStore) Delete(_ string) (int64, error) { return 0, nil }
func (s *insertMockStore) Ping() error                    { return nil }
func (s *insertMockStore) UpdateRole(_ context.Context, _ string, _ string) error {
	return nil
}

func (s *insertMockStore) ListByRole(_ context.Context, _ string) ([]*store.User, error) {
	return nil, nil
}

func (s *insertMockStore) ListAll(_ context.Context) ([]*store.User, error) {
	return nil, nil
}

func (s *insertMockStore) UpdateUser(_ context.Context, _ string, _ *store.User) (*store.User, error) {
	return s.created, nil
}

func (s *insertMockStore) ToggleActive(_ context.Context, _ string) (bool, error) {
	return false, nil
}

// --- UpdateUser ---

func adminCtx(adminID string) context.Context {
	claims := &jwt.RegisteredClaims{Subject: adminID}
	ctx := context.WithValue(context.Background(), middleware.UserClaimsKey, claims)
	return ctx
}

func TestUpdateUser_NoAuth(t *testing.T) {
	r := &mutationResolver{newResolver(&mocks.Store{}, &mocks.Authenticator{}, tokenConfig())}
	firstName := "Jane"
	_, err := r.UpdateUser(context.Background(), "user-1", model.UpdateUserInput{FirstName: &firstName})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
}

func TestUpdateUser_NotAdmin(t *testing.T) {
	// Caller exists but has USER role
	st := &mocks.Store{User: &store.User{ID: "caller-1", Role: "USER"}}
	r := &mutationResolver{newResolver(st, &mocks.Authenticator{}, tokenConfig())}
	firstName := "Jane"
	_, err := r.UpdateUser(adminCtx("caller-1"), "user-1", model.UpdateUserInput{FirstName: &firstName})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "admin role required")
}

func TestUpdateUser_NotCreator(t *testing.T) {
	adminID := "admin-1"
	targetUser := &store.User{
		ID:        "user-1",
		FirstName: "John",
		CreatedBy: "other-admin",
		Role:      "ADMIN",
	}
	// Retrieve returns admin first (for callerIsAdmin), then target user
	st := &updateMockStore{
		admin:  &store.User{ID: adminID, Role: "ADMIN"},
		target: targetUser,
	}
	r := &mutationResolver{newResolver(st, &mocks.Authenticator{}, tokenConfig())}
	firstName := "Jane"
	_, err := r.UpdateUser(adminCtx(adminID), "user-1", model.UpdateUserInput{FirstName: &firstName})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "only the admin who created this user")
}

func TestUpdateUser_Success(t *testing.T) {
	adminID := "admin-1"
	targetUser := &store.User{
		ID:        "user-1",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Company:   "OldCo",
		PostCode:  "12345",
		CreatedBy: adminID,
		Role:      "USER",
	}
	updatedUser := &store.User{
		ID:        "user-1",
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "john@example.com",
		Company:   "NewCo",
		PostCode:  "12345",
		CreatedBy: adminID,
		UpdatedAt: "2026-05-20 10:00:00",
	}
	st := &updateMockStore{
		admin:   &store.User{ID: adminID, Role: "ADMIN"},
		target:  targetUser,
		updated: updatedUser,
	}
	r := &mutationResolver{newResolver(st, &mocks.Authenticator{}, tokenConfig())}
	firstName := "Jane"
	company := "NewCo"
	resp, err := r.UpdateUser(adminCtx(adminID), "user-1", model.UpdateUserInput{
		FirstName: &firstName,
		Company:   &company,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "user-1", resp.ID)
	assert.Equal(t, "Jane", *resp.FirstName)
	assert.Equal(t, "NewCo", *resp.Company)
}

// updateMockStore returns different users for admin (callerIsAdmin) vs target (by ID).
type updateMockStore struct {
	admin   *store.User
	target  *store.User
	updated *store.User
	calls   int
}

func (s *updateMockStore) Insert(_ context.Context, _ *store.User) (*store.User, error) {
	return nil, nil
}

func (s *updateMockStore) Read(_ context.Context, _ string) (*store.User, error) {
	return &store.User{}, nil
}

func (s *updateMockStore) Retrieve(_ context.Context, _ string) (*store.User, error) {
	s.calls++
	// First call is callerIsAdmin checking the admin, second is fetching the target user
	if s.calls == 1 {
		return s.admin, nil
	}
	return s.target, nil
}
func (s *updateMockStore) Delete(_ string) (int64, error) { return 0, nil }
func (s *updateMockStore) Ping() error                    { return nil }
func (s *updateMockStore) UpdateRole(_ context.Context, _ string, _ string) error {
	return nil
}

func (s *updateMockStore) ListByRole(_ context.Context, _ string) ([]*store.User, error) {
	return nil, nil
}

func (s *updateMockStore) ListAll(_ context.Context) ([]*store.User, error) { return nil, nil }
func (s *updateMockStore) UpdateUser(_ context.Context, _ string, _ *store.User) (*store.User, error) {
	return s.updated, nil
}

func (s *updateMockStore) ToggleActive(_ context.Context, _ string) (bool, error) {
	return false, nil
}
