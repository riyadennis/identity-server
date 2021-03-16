package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

// Server registers routes and starts web server
func Server(port string) {
	router := httprouter.New()
	// register routes here
	router.POST(RegisterEndpoint, Register)
	router.POST(LoginEndPoint, Login)
	router.POST(DeleteEndpoint, Auth(Delete))
	router.GET(HomeEndPoint, Auth(Home))

	handler := cors.Default().Handler(router)
	logrus.Fatal(http.ListenAndServe(":"+port, handler))
}
