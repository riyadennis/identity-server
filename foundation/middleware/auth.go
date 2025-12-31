package middleware

import (
	"context"
	"net/http"

	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/business/validation"
	"github.com/riyadennis/identity-server/foundation"
	"github.com/sirupsen/logrus"
)

const (
	UserClaimsKey  = "claims"
	AccessTokenKey = "accessToken"
)

type AuthConfig struct {
	TokenConfig *store.TokenConfig
	Logger      *logrus.Logger
}

// Auth is the middleware that should be used for endpoints that needs jwt Token authentication.
// If Token is not present or is invalid, then the user is denied access to the wrapped endpoint.
func (ac *AuthConfig) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerToken := r.Header.Get("Authorization")
		claims, err := validation.ValidateToken(headerToken, ac.TokenConfig)
		if err != nil {
			ac.Logger.Errorf("invalid token: %v", err)
			foundation.ErrorResponse(w, http.StatusUnauthorized, err, foundation.UnAuthorised)
			return
		}

		ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
		ctx = context.WithValue(ctx, AccessTokenKey, headerToken)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
