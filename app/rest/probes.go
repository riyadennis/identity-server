package rest

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/foundation"
)

// @Router			/liveness [get].
func Liveness(w http.ResponseWriter, _ *http.Request) {
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

// @Router			/readiness [get].
func Ready(store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		if err := store.Ping(); err != nil {
			foundation.ErrorResponse(w, http.StatusInternalServerError, err, foundation.DatabaseError)
			return
		}

		_ = foundation.JSONResponse(w, http.StatusOK, "OK", "")
	}
}
