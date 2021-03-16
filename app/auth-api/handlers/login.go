package handlers

import (
	"encoding/json"
	"errors"
	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/foundation"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Status      int    `json:"status"`
	AccessToken string `json:"access_token"`
	Expiry      string `json:"expiry"`
	TokenType   string `json:"token_type"`
}

// Login endpoint where user enters his email
// and password to get back a Token.
// Which can be used to authenticate other requests.
func Login(w http.ResponseWriter,
	r *http.Request, _ httprouter.Params) {
	email, password, ok := r.BasicAuth()
	if !ok {
		foundation.ErrorResponse(w, http.StatusBadRequest,
			errors.New("empty login data"), foundation.InvalidRequest)
		return
	}
	ctx := r.Context()
	db, err := store.Connect()
	if err != nil {
		foundation.ErrorResponse(w, http.StatusInternalServerError, err, foundation.DatabaseError)
	}

	source := store.NewDB(db)
	if source == nil {
		foundation.ErrorResponse(w, http.StatusInternalServerError,
			errors.New("empty database connection"), foundation.DatabaseError)
		return
	}

	u, err := source.Read(ctx, email)
	if err != nil {
		foundation.ErrorResponse(w, http.StatusInternalServerError,
			err, foundation.UserDoNotExist)
		return
	}
	if u == nil {
		foundation.ErrorResponse(w, http.StatusInternalServerError,
			errors.New("email not found"), foundation.UserDoNotExist)
		return
	}
	valid, err := source.Authenticate(email, password)
	if err != nil {
		foundation.ErrorResponse(w, http.StatusBadRequest,
			errors.New("email not found"), foundation.InvalidRequest)
		return
	}
	if !valid {
		foundation.ErrorResponse(w, http.StatusBadRequest,
			err, foundation.UnAuthorised)
		return
	}
	token, err := generateToken()
	if err != nil {
		foundation.ErrorResponse(w, http.StatusInternalServerError,
			err, foundation.TokenError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(token)
	if err != nil {
		logrus.Error(err)
	}
}

func generateToken() (*Token, error) {
	expiry := time.Now().UTC().Add(tokenTTL)
	t, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": expiry.Unix(),
		"iss": os.Getenv("ISSUER"),
	}).SignedString([]byte(os.Getenv("KEY")))
	if err != nil {
		return nil, err
	}
	return &Token{
		Status:      200,
		AccessToken: t,
		Expiry:      expiry.String(),
		TokenType:   "Bearer",
	}, nil
}
