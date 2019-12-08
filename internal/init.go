package internal

import (
	"github.com/riyadennis/identity-server/internal/handlers"
	"net/http"
)

func Server(port string){
	http.HandleFunc(registerEndpoint, handlers.Register)
	http.ListenAndServe(port, nil)
}
