package handlers

import (
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// Auth is the wrapper that should be used for endpoints
// that needs jwt token authentication.
// if token is not present or is invalid then user
// is denied access to wrapped endpoint.
func Auth(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		headerToken := req.Header.Get("Token")
		if headerToken == "" {
			errorResponse(w, http.StatusUnauthorized, &CustomError{
				Code: UnAuthorised,
				Err:  errors.New("missing token"),
			})
			return
		}
		t, err := jwt.Parse(headerToken, tokenHandler)
		if err != nil || t == nil {
			logrus.Errorf("unable to parse jwt :: %v", err)
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
