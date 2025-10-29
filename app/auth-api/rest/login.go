package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/business/validation"
	"github.com/riyadennis/identity-server/foundation"
)

var (
	errEmailNotFound   = errors.New("email not found")
	errTokenGeneration = errors.New("key not found")
)

// UserLogin have data needed for a user to login
type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login @Summary      Login Endpoint
//
//	@Description	Authenticate a user and return a JWT token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Basic base64(email:password)"
//	@Success		200				{object}	store.Token
//	@Failure		400				{object}	foundation.Response
//	@Failure		401				{object}	foundation.Response
//	@Failure		500				{object}	foundation.Response
//	@Router			/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	email, password, ok := r.BasicAuth()
	if !ok {
		h.Logger.Printf("invalid request, username or password is empty")

		foundation.ErrorResponse(w, http.StatusBadRequest,
			errors.New("empty login data"), foundation.InvalidRequest)
		return
	}
	err := validation.ValidateEmail(email)
	if err != nil {
		h.Logger.Printf("invalid request, username is invalid: %v", err)
		foundation.ErrorResponse(w, http.StatusBadRequest,
			err, foundation.InvalidRequest)
		return
	}
	u, err := h.Store.Read(r.Context(), email)
	if err != nil {
		foundation.ErrorResponse(w, http.StatusInternalServerError,
			err, foundation.UserDoNotExist)
		return
	}

	if u == nil {
		h.Logger.Printf("failed to find user in DB")
		foundation.ErrorResponse(w, http.StatusInternalServerError,
			errEmailNotFound, foundation.UserDoNotExist)
		return
	}

	valid, err := h.Authenticator.Authenticate(email, password)
	if err != nil {
		h.Logger.Errorf("failed to authenticate provided password %v", err)

		foundation.ErrorResponse(w, http.StatusBadRequest,
			errEmailNotFound, foundation.InvalidRequest)
		return
	}

	if !valid {
		h.Logger.Errorf("failed to authenticate user: %v", err)

		foundation.ErrorResponse(w, http.StatusBadRequest,
			errEmailNotFound, foundation.UnAuthorised)
		return
	}
	tr, err := h.Authenticator.FetchLoginToken(u.ID)
	if err != nil {
		h.Logger.Errorf("failed to fetch token from DB: %v", err)

		foundation.ErrorResponse(w, http.StatusInternalServerError,
			errTokenGeneration, foundation.KeyNotFound)

		return
	}
	var token *store.Token
	// token is present and not expired
	if tr != nil && tr.Expiry.After(time.Now()) {
		h.Logger.Printf("token already exists with id: %s", tr.Id)
		token = &store.Token{
			Status:      http.StatusOK,
			AccessToken: tr.Token,
			TokenType:   "Bearer",
			Expiry:      tr.Expiry.String(),
			TokenTTL:    tr.TTL,
		}
	} else {
		key, err := fetchPrivateKey(h.TokenConfig.KeyPath+h.TokenConfig.PrivateKeyName, h.TokenConfig.KeyPath+h.TokenConfig.PublicKeyName)
		if err != nil {
			h.Logger.Errorf("failed to fetch keys: %v", err)

			foundation.ErrorResponse(w, http.StatusInternalServerError,
				errTokenGeneration, foundation.KeyNotFound)

			return
		}
		expiryTime := time.Now().UTC().Add(120 * time.Hour)
		token, err = store.GenerateToken(h.Logger, h.TokenConfig.Issuer, key, expiryTime)
		if err != nil {
			h.Logger.Printf("token generation failed: %v", err)

			foundation.ErrorResponse(w, http.StatusInternalServerError,
				errTokenGeneration, foundation.TokenError)
			return
		}
		err = h.Authenticator.SaveLoginToken(r.Context(), &store.TokenRecord{
			UserId: u.ID,
			Token:  token.AccessToken,
			Expiry: expiryTime,
			TTL:    fmt.Sprintf("%d", expiryTime.Unix()),
		})
		if err != nil {
			h.Logger.Printf("token saving failed: %v", err)

			foundation.ErrorResponse(w, http.StatusInternalServerError,
				errTokenGeneration, foundation.TokenError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(token)
	if err != nil {
		h.Logger.Printf("json encoding failed: %v", err)
	}
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
