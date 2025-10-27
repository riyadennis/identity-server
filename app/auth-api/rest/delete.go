package rest

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/riyadennis/identity-server/foundation"
)

var (
	errInvalidID    = errors.New("invalid userID in request")
	errDeleteFailed = errors.New("failed to remove user")
)

// Delete @Summary      Delete Endpoint
// @Description  Permanently remove a user by ID
// @Tags         User
// @Security     ApiKeyAuth
// @Produce      json
// @Param        id path string true "User ID"
// @Success      204   {string}  string  "No Content"
// @Failure      400   {object}  foundation.Response
// @Failure      404   {object}  foundation.Response
// @Router       /user/delete/{userID} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "userID")
	if id == "" {
		foundation.ErrorResponse(w, http.StatusBadRequest,
			errInvalidID, foundation.InvalidRequest)
		return
	}

	deleted, err := h.Store.Delete(id)
	if err != nil {
		h.Logger.Printf("user deletion failed: %v", err)

		foundation.ErrorResponse(w, http.StatusBadRequest,
			errDeleteFailed, foundation.DatabaseError)
		return
	}

	if deleted == 0 {
		foundation.ErrorResponse(w, http.StatusBadRequest,
			errDeleteFailed, foundation.UserDoNotExist)
		return
	}
	_ = foundation.JSONResponse(w, http.StatusNoContent, "", "")
}
