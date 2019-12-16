package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Home(w http.ResponseWriter, req *http.Request, _ httprouter.Params){
	if req.Header.Get("Token")== ""{
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(newResponse(
			http.StatusUnauthorized,
			"missing token",
			UnAuthorised,
		))
		return
	}
	t, err := jwt.Parse(req.Header.Get("Token"), tokenHandler)
	if err != nil || t == nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(newResponse(
			http.StatusUnauthorized,
			"invalid token",
			UnAuthorised,
		))
		return
	}
	if !t.Valid{
		json.NewEncoder(w).Encode(newResponse(
			http.StatusUnauthorized,
			"invalid token",
			UnAuthorised,
		))
		return
	}
	json.NewEncoder(w).Encode(newResponse(
		http.StatusOK,
		"Authorised",
		"",
	))
}

func tokenHandler(token *jwt.Token) (interface{}, error){
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("There was an error")
	}
	return mySigningKey, nil
}