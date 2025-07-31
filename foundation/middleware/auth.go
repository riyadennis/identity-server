package middleware

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/business/validation"
	"github.com/riyadennis/identity-server/foundation"
)

// Auth is the middleware that should be used for endpoints that needs jwt Token authentication.
// if Token is not present or is invalid then user is denied access to wrapped endpoint.
func Auth(next httprouter.Handle, tc *store.TokenConfig, logger *log.Logger) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		headerToken := req.Header.Get("Authorization")
		if err := validation.ValidateToken(headerToken, tc); err != nil {
			logger.Printf("invalid token: %v", err)
			foundation.ErrorResponse(w, http.StatusUnauthorized, err, foundation.UnAuthorised)
			return
		}
		next(w, req, p)
	}
}
