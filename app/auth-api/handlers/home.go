package handlers

import (
	"fmt"
	"net/http"

	"github.com/riyadennis/identity-server/foundation"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Home is the rest endpoint a logged in user with valid token can access
func Home(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	err := foundation.JSONResponse(w, http.StatusOK,
		"Authorised",
		"")
	if err != nil {
		logrus.Error(err)
	}
}

func tokenHandler(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unable to handle Token")
	}
	return []byte(viper.GetStringMapString("jwt")["signing-key"]), nil
}
