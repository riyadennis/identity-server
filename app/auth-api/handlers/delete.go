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

// Delete is the handler for delete end point to remove a user from as per the email
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
