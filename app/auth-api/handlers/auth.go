package handlers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go/v4"
	"github.com/julienschmidt/httprouter"

	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/foundation"
)

const (
	// BearerSchema is expected prefix for token from authorisation header
	BearerSchema = "Bearer "

	// tokenTTL is the expiry time for a token
	tokenTTL = 120 * time.Hour
)

var (
	errInvalidToken = errors.New("invalid token")
)

// Auth is the middleware that should be used for endpoints that needs jwt Token authentication.
// if Token is not present or is invalid then user is denied access to wrapped endpoint.
func Auth(next httprouter.Handle, tc *store.TokenConfig, logger *log.Logger) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		headerToken := req.Header.Get("Authorization")
		if headerToken == "" {
			logger.Printf("missing Authorisation in header")

			foundation.ErrorResponse(w, http.StatusUnauthorized, errors.New("missing Token"), foundation.UnAuthorised)
			return
		}

		if headerToken[len(BearerSchema):] == "" {
			logger.Printf("missing Bearer token in header")

			foundation.ErrorResponse(w, http.StatusBadRequest, errors.New("bearer Token not present"), foundation.UnAuthorised)
			return
		}

		t, err := jwt.ParseWithClaims(
			headerToken[len(BearerSchema):],
			jwt.MapClaims{
				"exp": time.Now().UTC().Add(tokenTTL).Unix(),
				"iss": tc.Issuer,
			}, fetchKey(tc.KeyPath+foundation.PublicKeyFileName))

		if err != nil || t == nil {
			logger.Printf("failed to parse the token: %v", err)

			foundation.ErrorResponse(w, http.StatusUnauthorized, errInvalidToken, foundation.TokenError)
			return
		}

		if !t.Valid {
			logger.Printf("invalid token: %v", t)

			foundation.ErrorResponse(w, http.StatusUnauthorized, errInvalidToken, foundation.UnAuthorised)
			return
		}

		next(w, req, p)
	}
}

func fetchKey(keyPath string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unable to handle Token")
		}

		key, err := ioutil.ReadFile(keyPath)
		if err != nil {
			return nil, fmt.Errorf("unable to handle Token")
		}

		publicKey, e := jwt.ParseRSAPublicKeyFromPEM(key)
		if e != nil {
			panic(e.Error())
		}

		return publicKey, nil
	}
}
