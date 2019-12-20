package handlers

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

func Auth(next httprouter.Handle) httprouter.Handle{
	return func(w http.ResponseWriter, req *http.Request,  p httprouter.Params){
		headerToken := req.Header.Get("Token")
		if headerToken == "" {
			errorResponse(w, http.StatusUnauthorized, &CustomError{
				Code: UnAuthorised,
				Err:  errors.New("missing token"),
			})
			return
		}
		t, err := jwt.ParseWithClaims(
			headerToken,
			jwt.MapClaims{
				"exp": time.Now().UTC().Add(tokenTTL).Unix(),
			},
			tokenHandler,
		)
		if err != nil || t == nil {
			errorResponse(w, http.StatusUnauthorized, &CustomError{
				Code: UnAuthorised,
				Err:  err,
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
		next(w, req, p)
	}
}

