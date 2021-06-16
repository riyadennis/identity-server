package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/riyadennis/identity-server/foundation"
)

// Home is the rest endpoint a logged in user with valid token can access
func Home(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	_ = foundation.JSONResponse(w, http.StatusOK, "Authorised", "")
}
