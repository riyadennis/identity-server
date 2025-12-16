package graph

import (
	"github.com/riyadennis/identity-server/business/store"
	"github.com/sirupsen/logrus"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	Logger        *logrus.Logger
	tokenConfig   *store.TokenConfig
	Store         store.Store
	Authenticator store.Authenticator
}

func NewResolver(l *logrus.Logger, tc *store.TokenConfig, st store.Store, au store.Authenticator) *Resolver {
	return &Resolver{
		Logger:        l,
		tokenConfig:   tc,
		Store:         st,
		Authenticator: au,
	}
}
