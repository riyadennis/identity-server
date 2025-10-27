package rest

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/riyadennis/identity-server/foundation"
)

// Home @Summary      Get user dashboard
// @Description  Returns dashboard info for authenticated user
// @Tags         User
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200   {object}  foundation.Response
// @Failure      401   {object}  foundation.Response
// @Router       /home [get]
func Home(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	_ = foundation.JSONResponse(w, http.StatusOK, "Authorised", "")
}
