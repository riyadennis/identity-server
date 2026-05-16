package mocks

import (
	"context"

	"github.com/riyadennis/identity-server/business/store"
)

type Store struct {
	Error error
	*store.User
}

func (s *Store) Insert(_ context.Context, _ *store.User) (*store.User, error) {
	return s.User, s.Error
}

func (s *Store) Read(_ context.Context, _ string) (*store.User, error) {
	return s.User, s.Error
}

func (s *Store) Retrieve(_ context.Context, _ string) (*store.User, error) {
	return s.User, s.Error
}

func (s *Store) Delete(_ string) (int64, error) {
	return 0, s.Error
}

func (s *Store) Ping() error {
	return s.Error
}

func (s *Store) UpdateRole(_ context.Context, _ string, _ string) error {
	return s.Error
}

func (s *Store) ListByRole(_ context.Context, _ string) ([]*store.User, error) {
	if s.User == nil {
		return nil, s.Error
	}
	return []*store.User{s.User}, s.Error
}

func (s *Store) ListAll(_ context.Context) ([]*store.User, error) {
	if s.User == nil {
		return nil, s.Error
	}
	return []*store.User{s.User}, s.Error
}

type Authenticator struct {
	ReturnVal bool
	Error     error
	Token     *store.TokenRecord
}

func (ma *Authenticator) Authenticate(_, _ string) (bool, error) {
	return ma.ReturnVal, ma.Error
}

func (ma *Authenticator) FetchLoginToken(_ string) (*store.TokenRecord, error) {
	return ma.Token, nil
}

func (ma *Authenticator) SaveLoginToken(_ context.Context, _ *store.TokenRecord) error {
	return nil
}
