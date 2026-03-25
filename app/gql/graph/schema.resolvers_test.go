package graph

import (
	"context"
	"errors"
	"testing"

	"github.com/riyadennis/identity-server/app/gql/graph/model"
	"github.com/riyadennis/identity-server/app/mocks"
	"github.com/riyadennis/identity-server/business/store"
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

// insertMockStore returns empty user on Read (not found) and created on Insert
type insertMockStore struct {
	created *store.User
}

func (s *insertMockStore) Insert(ctx context.Context, u *store.User) (*store.User, error) {
	return s.created, nil
}
func (s *insertMockStore) Read(ctx context.Context, email string) (*store.User, error) {
	return &store.User{}, nil // empty — email not found
}
func (s *insertMockStore) Retrieve(ctx context.Context, id string) (*store.User, error) {
	return s.created, nil
}
func (s *insertMockStore) Delete(id string) (int64, error) { return 0, nil }
func (s *insertMockStore) Ping() error                     { return nil }