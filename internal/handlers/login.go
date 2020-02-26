package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type LoginDetails struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	user, password, ok:= req.BasicAuth()
	if !ok{
		errorResponse(w, http.StatusBadRequest, &CustomError{
			Code: InvalidRequest,
			Err:  errors.New("empty login data"),
		})
		return
	}
	ld := &LoginDetails{
		Email: user,
		Password: password,
	}
	if ld.Email == "" {
		errorResponse(w, http.StatusBadRequest, &CustomError{
			Code: EmailMissing,
			Err:  errors.New("email missing"),
		})
		return
	}
	if ld.Password == "" {
		logrus.Error("no password")
		errorResponse(w, http.StatusBadRequest, &CustomError{
			Code: PassWordError,
			Err:  errors.New("password missing"),
		})
		return
	}
	fName, err := Idb.Authenticate(ld.Email, ld.Password)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, &CustomError{
			Code: InvalidRequest,
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
	res := newResponse(http.StatusOK,
		fmt.Sprintf("welcome  : %s", fName),
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

func generateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().UTC().Add(tokenTTL).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(mySigningKey))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}
