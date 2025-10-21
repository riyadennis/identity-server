package handlers

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"

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
// @Router       /delete/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, _ *http.Request, params httprouter.Params) {
	id := params.ByName("id")
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
