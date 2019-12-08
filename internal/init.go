package internal

import (
	"net/http"

	"github.com/riyadennis/identity-server/internal/handlers"
)

func Server(port string) {
	http.HandleFunc(registerEndpoint, handlers.Register)
	http.ListenAndServe(port, nil)
}
