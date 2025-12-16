package business

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/riyadennis/identity-server/business/validation"
	"github.com/riyadennis/identity-server/foundation"
	"github.com/sirupsen/logrus"

	"github.com/riyadennis/identity-server/business/store"
)

type Helper struct {
	Store         store.Store
	Authenticator store.Authenticator
	Logger        *logrus.Logger
}

var (
	errEmailNotFound   = errors.New("email not found")
	errInvalidPassword = errors.New("invalid password")
)

func NewHelper(s store.Store, a store.Authenticator, l *logrus.Logger) *Helper {
	return &Helper{
		Store:         s,
		Authenticator: a,
		Logger:        l,
	}
}
func (h *Helper) Login(ctx context.Context, tc *store.TokenConfig, email, password string) (*store.Token, error) {
	err := validation.ValidateEmail(email)
	if err != nil {
		return nil, err
	}
	user, err := h.UserCredentialsInDB(ctx, email, password)
	if err != nil {
		//already logged
		return nil, err
	}
	token, err := h.ManageToken(ctx, tc, user.ID)
	if err != nil {
		//already logged
		return nil, err
	}

	return token, nil
}

func (h *Helper) UserCredentialsInDB(ctx context.Context, email, password string) (*store.User, error) {
	user, err := h.Store.Read(ctx, email)
	if err != nil {
		h.Logger.Errorf("failed to find user in DB")
		return nil, errEmailNotFound
	}
	if user == nil {
		h.Logger.Printf("user not found in DB")
		return nil, errEmailNotFound
	}
	valid, err := h.Authenticator.Authenticate(email, password)
	if err != nil {
		h.Logger.Errorf("failed to authenticate provided password %v", err)
		return nil, errInvalidPassword
	}
	if !valid {
		h.Logger.Errorf("failed to authenticate user: %v", err)
		return nil, errInvalidPassword
	}
	return user, nil
}

func (h *Helper) ManageToken(ctx context.Context, config *store.TokenConfig, userID string) (*store.Token, error) {
	tr, err := h.Authenticator.FetchLoginToken(userID)
	if err != nil {
		h.Logger.Errorf("failed to fetch token from DB: %v", err)
		return nil, err
	}

	// token is present and not expired
	if tr != nil && tr.Expiry.After(time.Now()) {
		h.Logger.Printf("token already exists with id: %s", tr.Id)
		return &store.Token{
			Status:      http.StatusOK,
			AccessToken: tr.Token,
			TokenType:   "Bearer",
			Expiry:      tr.Expiry.String(),
			TokenTTL:    tr.TTL,
		}, nil
	}
	key, err := fetchPrivateKey(config.KeyPath+config.PrivateKeyName, config.KeyPath+config.PublicKeyName)
	expiryTime := time.Now().UTC().Add(120 * time.Hour)
	token, err := store.GenerateToken(h.Logger, config.Issuer, key, expiryTime)
	if err != nil {
		h.Logger.Errorf("failed to generate token: %v", err)
		return nil, err
	}
	err = h.Authenticator.SaveLoginToken(ctx, &store.TokenRecord{
		UserId: userID,
		Token:  token.AccessToken,
		Expiry: expiryTime,
		TTL:    fmt.Sprintf("%d", expiryTime.Unix()),
	})
	if err != nil {
		h.Logger.Printf("token saving failed: %v", err)
		return nil, err
	}

	return token, nil
}

func fetchPrivateKey(privateKeyPath, publicKeyPath string) ([]byte, error) {
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		err := foundation.GenerateKeys(privateKeyPath, publicKeyPath)
		if err != nil {
			return nil, err
		}
	}

	return os.ReadFile(privateKeyPath)
}
