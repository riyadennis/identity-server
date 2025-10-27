package rest

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"

	"github.com/riyadennis/identity-server/foundation"
)

// @Summary      Liveness probe
// @Description  Returns liveness and k8s deployment info
// @Tags         Health
// @Produce      json
// @Success      200   {object}  map[string]interface{}
// @Router       /liveness [get]
func Liveness(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "unavailable"
	}
	data := struct {
		Status    string `jsonapi:"status,omitempty"`
		Host      string `jsonapi:"host,omitempty"`
		Pod       string `jsonapi:"pod,omitempty"`
		PodIP     string `jsonapi:"podIP,omitempty"`
		Node      string `jsonapi:"node,omitempty"`
		Namespace string `jsonapi:"namespace,omitempty"`
	}{
		Status:    "up",
		Host:      hostName,
		Pod:       os.Getenv("KUBERNETES_PODNAME"),
		PodIP:     os.Getenv("KUBERNETES_NAMESPACE_POD_IP"),
		Node:      os.Getenv("KUBERNETES_NODENAME"),
		Namespace: os.Getenv("KUBERNETES_NAMESPACE"),
	}

	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(data)
}

// @Summary      Readiness probe
// @Description  Checks if API is ready for traffic (DB available)
// @Tags         Health
// @Produce      json
// @Success      200   {object}  foundation.Response
// @Failure      500   {object}  foundation.Response
// @Router       /readiness [get]
func Ready(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, request *http.Request, params httprouter.Params) {
		if err := db.Ping(); err != nil {
			foundation.ErrorResponse(w, http.StatusInternalServerError, err, foundation.DatabaseError)
			return
		}

		_ = foundation.JSONResponse(w, http.StatusOK, "OK", "")
	}
}
