package middleware

import (
	"log"
	"net/http"

	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/business/validation"
	"github.com/riyadennis/identity-server/foundation"
)

type AuthConfig struct {
	TokenConfig *store.TokenConfig
	Logger      *log.Logger
}

// Auth is the middleware that should be used for endpoints that needs jwt Token authentication.
// If Token is not present or is invalid, then the user is denied access to the wrapped endpoint.
func (ac *AuthConfig) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerToken := r.Header.Get("Authorization")
		if err := validation.ValidateToken(headerToken, ac.TokenConfig); err != nil {
			ac.Logger.Printf("invalid token: %v", err)
			foundation.ErrorResponse(w, http.StatusUnauthorized, err, foundation.UnAuthorised)
			return
		}
		next.ServeHTTP(w, r)
	})
}
