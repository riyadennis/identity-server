package handlers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Home(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	err := jsonResponse(w, http.StatusOK,
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
	return []byte(mySigningKey), nil
}
