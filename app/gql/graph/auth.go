package graph

import (
	"context"
	"fmt"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/riyadennis/identity-server/app/gql/graph/model"
	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/foundation/middleware"
	"github.com/sirupsen/logrus"
)

// callerID extracts the authenticated caller's user ID from the request context.
func callerID(ctx context.Context) (string, error) {
	claims, ok := ctx.Value(middleware.UserClaimsKey).(*jwt.RegisteredClaims)
	if !ok || claims == nil || claims.Subject == "" {
		return "", fmt.Errorf("unauthorized: missing or invalid token claims")
	}
	return claims.Subject, nil
}

// callerIsAdmin looks up the caller's role from the store and returns an error
// if they are not an ADMIN. Use this at the top of any resolver that requires
// elevated permissions.
func callerIsAdmin(ctx context.Context, s store.Store, logger *logrus.Logger) error {
	userID, err := callerID(ctx)
	if err != nil {
		return err
	}

	user, err := s.Retrieve(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to verify caller permissions: %w", err)
	}
	if user == nil {
		return fmt.Errorf("unauthorized: caller not found")
	}
	logger.Infof("got user %v", model.Role(user.Role))
	if model.Role(user.Role) != model.RoleAdmin {
		return fmt.Errorf("forbidden: admin role required")
	}

	return nil
}
