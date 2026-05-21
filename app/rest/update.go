package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/foundation"
	"github.com/riyadennis/identity-server/foundation/middleware"
)

var (
	errUpdateFailed = errors.New("failed to update user")
	errNotCreator   = errors.New("only the admin who created this user can update them")
)

// UpdateUser @Summary      Update user data
//
//	@Description	Update editable fields for a user (admin only, must be creator)
//	@Tags			Admin
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		string		true	"User ID"
//	@Param			user	body		store.User	true	"Updated user data"
//	@Success		200		{object}	store.User
//	@Failure		400		{object}	foundation.Response
//	@Failure		403		{object}	foundation.Response
//	@Failure		500		{object}	foundation.Response
//	@Router			/admin/update/{userID} [put]
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		foundation.ErrorResponse(w, http.StatusBadRequest,
			errInvalidID, foundation.InvalidRequest)
		return
	}

	// Get the authenticated admin's ID from JWT claims.
	claims, ok := r.Context().Value(middleware.UserClaimsKey).(*jwt.RegisteredClaims)
	if !ok || claims == nil {
		foundation.ErrorResponse(w, http.StatusUnauthorized,
			errors.New("invalid token claims"), foundation.UnAuthorised)
		return
	}
	adminID := claims.Subject

	// Fetch the target user to verify the admin is the creator.
	existing, err := h.Store.Retrieve(r.Context(), userID)
	if err != nil {
		h.Logger.Errorf("failed to retrieve user: %v", err)
		foundation.ErrorResponse(w, http.StatusInternalServerError,
			errUpdateFailed, foundation.DatabaseError)
		return
	}

	if existing == nil {
		foundation.ErrorResponse(w, http.StatusBadRequest,
			errUpdateFailed, foundation.UserDoNotExist)
		return
	}

	if existing.CreatedBy != adminID {
		foundation.ErrorResponse(w, http.StatusForbidden,
			errNotCreator, foundation.UnAuthorised)
		return
	}

	// Decode the update payload.
	u := &store.User{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		h.Logger.Errorf("invalid data in request: %v", err)
		foundation.ErrorResponse(w, http.StatusBadRequest,
			err, foundation.InvalidRequest)
		return
	}

	updated, err := h.Store.UpdateUser(r.Context(), userID, u)
	if err != nil {
		h.Logger.Errorf("failed to update user: %v", err)
		foundation.ErrorResponse(w, http.StatusInternalServerError,
			errUpdateFailed, foundation.DatabaseError)
		return
	}

	updated.Password = "********"
	_ = foundation.Resource(w, http.StatusOK, updated)
}
