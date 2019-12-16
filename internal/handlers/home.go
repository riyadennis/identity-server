package handlers

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Home(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if req.Header.Get("Token") == "" {
		errorResponse(w, http.StatusUnauthorized, &CustomError{
			Code: UnAuthorised,
			Err:  errors.New("missing token"),
		})
		return
	}
	t, err := jwt.Parse(req.Header.Get("Token"), tokenHandler)
	if err != nil || t == nil {
		errorResponse(w, http.StatusUnauthorized, &CustomError{
			Code: UnAuthorised,
			Err:  errors.New("invalid token"),
		})
		return
	}
	if !t.Valid {
		errorResponse(w, http.StatusUnauthorized, &CustomError{
			Code: UnAuthorised,
			Err:  errors.New("invalid token"),
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = jsonResponse(w, http.StatusOK,
		"Authorised",
		"")
	if err != nil {
		logrus.Error(err)
	}

}

func tokenHandler(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unable to handle token")
	}
	return mySigningKey, nil
}
