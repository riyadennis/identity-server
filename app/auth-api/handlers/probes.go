package handlers

import (
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/riyadennis/identity-server/foundation"
)

// Liveness returns liveness status of the service
func Liveness(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	_ = foundation.JSONResponse(w, http.StatusOK, "OK", "")
}

// Ready returns readiness status of the service
func Ready(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, request *http.Request, params httprouter.Params) {
		if err := db.Ping(); err != nil {
			foundation.ErrorResponse(w, http.StatusInternalServerError, err, foundation.DatabaseError)
			return
		}

		_ = foundation.JSONResponse(w, http.StatusOK, "OK", "")
	}
}
