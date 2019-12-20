package internal

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/riyadennis/identity-server/internal/handlers"
	"github.com/sirupsen/logrus"
)

func Server(port string) {
	router := httprouter.New()
	// register routes here
	router.POST(RegisterEndpoint, handlers.Register)
	router.POST(LoginEndPoint, handlers.Login)
	router.GET(HomeEndPoint, handlers.Auth(handlers.Home))
	logrus.Fatal(http.ListenAndServe(port, router))
}
