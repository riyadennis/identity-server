package handlers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Home(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if req.Header.Get("Token") == "" {
		err := jsonResponse(w, http.StatusUnauthorized,
			"missing token",
			UnAuthorised)
		if err != nil {
			logrus.Error(err)
		}
		return
	}
	t, err := jwt.Parse(req.Header.Get("Token"), tokenHandler)
	if err != nil || t == nil {
		err := jsonResponse(w, http.StatusUnauthorized,
			"invalid token",
			UnAuthorised)
		if err != nil {
			logrus.Error(err)
		}
		return
	}
	if !t.Valid {
		err := jsonResponse(w, http.StatusUnauthorized,
			"invalid token",
			UnAuthorised)
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	err = jsonResponse(w, http.StatusOK,
		"Authorised",
		"")
	if err != nil {
		logrus.Error(err)
	}

}

func tokenHandler(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("There was an error")
	}
	return mySigningKey, nil
}
