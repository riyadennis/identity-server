package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/spf13/viper"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login endpoint where user enters his email
// and password to get back a token.
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
	token, err := generateToken(email)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, &CustomError{
			Code: TokenError,
			Err:  err,
		})
		return
	}

	res := newResponse(http.StatusOK,
		fmt.Sprintf("welcome  : %s", u.FirstName),
		"",
	)
	res.Token = token
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		logrus.Error(err)
	}
	return
}

func generateToken(email string) (string, error) {
	claims := &Claims{
		Email:          email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(tokenTTL).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)

	key, err := ioutil.ReadFile(viper.GetString("signing-key"))
	if err != nil{
		return "", err
	}
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(key)
	if err != nil{
		return "", err
	}
	tokenStr, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}
