package mocks

import (
	"context"

	"github.com/riyadennis/identity-server/business/store"
)

type Store struct {
	Error error
	*store.User
}

func (m *Store) Insert(ctx context.Context, u *store.User) (*store.User, error) {
	return m.User, m.Error
}

func (m *Store) Read(ctx context.Context, email string) (*store.User, error) {
	return m.User, m.Error
}

func (m *Store) Delete(id string) (int64, error) {
	return 0, m.Error
}

type Authenticator struct {
	ReturnVal bool
	Error     error
	Token     *store.TokenRecord
}

func (ma *Authenticator) Authenticate(email, password string) (bool, error) {
	return ma.ReturnVal, ma.Error
}
func (ma *Authenticator) FetchLoginToken(userID string) (*store.TokenRecord, error) {
	return ma.Token, nil
}
func (ma *Authenticator) SaveLoginToken(ctx context.Context, t *store.TokenRecord) error {
	return nil
}
