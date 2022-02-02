package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/riyadennis/identity-server/foundation/middleware"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/julienschmidt/httprouter"
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

// Token have credentials present in a token
type Token struct {
	Status      int    `json:"status"`
	AccessToken string `json:"access_token"`
	Expiry      string `json:"expiry"`
	TokenType   string `json:"token_type"`
}

// Login endpoint where user enters his email
// and password to get back a Token.
// Which can be used to authenticate other requests.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	email, password, ok := r.BasicAuth()
	if !ok {
		h.Logger.Printf("invalid request: %v", r)

		foundation.ErrorResponse(w, http.StatusBadRequest,
			errors.New("empty login data"), foundation.InvalidRequest)
		return
	}

	u, err := h.Store.Read(r.Context(), email)
	if err != nil {
		foundation.ErrorResponse(w, http.StatusInternalServerError,
			err, foundation.UserDoNotExist)
		return
	}

	if u == nil {
		h.Logger.Printf("failed to find user in DB: %v", r)
		foundation.ErrorResponse(w, http.StatusInternalServerError,
			errEmailNotFound, foundation.UserDoNotExist)
		return
	}

	valid, err := h.Authenticator.Authenticate(email, password)
	if err != nil {
		h.Logger.Printf("failed to find user in DB: %v", err)

		foundation.ErrorResponse(w, http.StatusBadRequest,
			errEmailNotFound, foundation.InvalidRequest)
		return
	}

	if !valid {
		h.Logger.Printf("failed to authenticate user: %v", err)

		foundation.ErrorResponse(w, http.StatusBadRequest,
			errEmailNotFound, foundation.UnAuthorised)
		return
	}

	key, err := fetchPrivateKey(h.TokenConfig.KeyPath+h.TokenConfig.PrivateKeyName, h.TokenConfig.KeyPath+h.TokenConfig.PublicKeyName)
	if err != nil {
		h.Logger.Printf("failed to fetch keys: %v", err)

		foundation.ErrorResponse(w, http.StatusInternalServerError,
			errTokenGeneration, foundation.KeyNotFound)

		return
	}

	token, err := generateToken(h.Logger, h.TokenConfig.Issuer, key)
	if err != nil {
		h.Logger.Printf("token generation failed: %v", err)

		foundation.ErrorResponse(w, http.StatusInternalServerError,
			errTokenGeneration, foundation.TokenError)
		return
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

	return ioutil.ReadFile(privateKeyPath)
}

func generateToken(logger *log.Logger, issuer string, key []byte) (*Token, error) {
	expiry := time.Now().UTC().Add(middleware.TokenTTL)

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(key)
	if err != nil {
		logger.Printf("failed to parser private key: %v", err)
		return nil, err
	}

	t, err := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"exp": expiry.Unix(),
		"iss": issuer,
	}).SignedString(privateKey)
	if err != nil {
		logger.Printf("failed to sign using private key: %v", err)
		return nil, err
	}

	return &Token{
		Status:      200,
		AccessToken: t,
		Expiry:      expiry.String(),
		TokenType:   "Bearer",
	}, nil
}
