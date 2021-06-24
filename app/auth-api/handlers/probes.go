package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/riyadennis/identity-server/foundation"
)

// Liveness returns liveness status of the service,
// if app is deployed in kubernetes cluster it will return pod, node and namespace.
// the env vars are to be set from the manifest.
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
