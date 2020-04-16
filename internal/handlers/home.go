package handlers

import (
	"context"
	"fmt"

	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/riyadennis/aes-encryption/ex"
	"github.com/riyadennis/aes-encryption/ex/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Home(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	cl := ex.NewClient()
	ctx := context.Background()
	in := &api.DataRequest{
		Data: &api.Data{
			Message: "hello",
		},
	}
	resp, err := cl.Store(ctx, in)
	if err != nil {
		logrus.Error(err)
		jsonResponse(w, http.StatusInternalServerError,
			fmt.Sprintf("unable to save message :: %s", err.Error()),
			"cannot talk to aes")
	}
	logrus.Infof("got response from aes :: %s", resp)
	err = jsonResponse(w, http.StatusOK,
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
