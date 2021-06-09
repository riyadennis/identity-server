package handlers

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/riyadennis/identity-server/foundation"

	jwt "github.com/dgrijalva/jwt-go/v4"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// BearerSchema is expected prefix for token from authorisation header
const (
	BearerSchema = "Bearer "

	tokenTTL = 120 * time.Hour
)

// Auth is the wrapper that should be used for endpoints
// that needs jwt Token authentication.
// if Token is not present or is invalid then user
// is denied access to wrapped endpoint.
func Auth(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		headerToken := req.Header.Get("Authorization")
		if headerToken == "" {
			foundation.ErrorResponse(w, http.StatusUnauthorized, errors.New("missing Token"), foundation.UnAuthorised)
			return
		}
		if headerToken[len(BearerSchema):] == "" {
			foundation.ErrorResponse(w, http.StatusBadRequest, errors.New("bearer Token not present"), foundation.UnAuthorised)
			return
		}
		t, err := jwt.ParseWithClaims(
			headerToken[len(BearerSchema):],
			jwt.MapClaims{
				"exp": time.Now().UTC().Add(tokenTTL).Unix(),
				"iss": os.Getenv("ISSUER"),
			}, tokenHandler)
		if err != nil || t == nil {
			logrus.Errorf("unable to parse jwt :: %v", err)
			foundation.ErrorResponse(w, http.StatusUnauthorized, err, foundation.UnAuthorised)
			return
		}
		if !t.Valid {
			foundation.ErrorResponse(w, http.StatusUnauthorized, errors.New("invalid Token"), foundation.UnAuthorised)
			return
		}
		next(w, req, p)
	}
}
