package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
	req *http.Request, _ httprouter.Params) {
	email, password, ok := req.BasicAuth()
	if !ok {
		errorResponse(w, http.StatusBadRequest, &CustomError{
			Code: InvalidRequest,
			Err:  errors.New("empty login data"),
		})
		return
	}
	source := dataSource()
	u, err := source.Read(email)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, &CustomError{
			Code: UserDoNotExist,
			Err:  err,
		})
		return
	}
	if u == nil {
		errorResponse(w, http.StatusInternalServerError, &CustomError{
			Code: UserDoNotExist,
			Err:  errors.New("email not found"),
		})
		return
	}
	valid, err := source.Authenticate(email, password)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, &CustomError{
			Code: InvalidRequest,
			Err:  err,
		})
		return
	}
	if !valid {
		errorResponse(w, http.StatusBadRequest, &CustomError{
			Code: UnAuthorised,
			Err:  err,
		})
		return
	}
	token, err := generateToken()
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, &CustomError{
			Code: TokenError,
			Err:  err,
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(token)
	if err != nil {
		logrus.Error(err)
	}
	return
}

func generateToken() (*Token, error) {
	jwtConf := viper.GetStringMapString("jwt")
	expiry := time.Now().UTC().Add(tokenTTL)
	t, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": expiry.Unix(),
		"iss": jwtConf["issuer"],
	}).SignedString([]byte(jwtConf["signing-key"]))
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
