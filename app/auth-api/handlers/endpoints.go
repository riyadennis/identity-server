package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/foundation/middleware"
)

const (
	// RegisterEndpoint is to create a new user
	RegisterEndpoint = "/register"

	// DeleteEndpoint is to delete a user
	DeleteEndpoint = "/delete/:id"

	// LoginEndPoint creates a token for the  user of credentials are valid
	LoginEndPoint = "/login"

	// HomeEndPoint is the details end point that a
	// logged in user with valid token can access
	HomeEndPoint = "/home"

	// LivenessEndPoint is for kubernetes to check when to restart the container
	LivenessEndPoint = "/liveness"

	// ReadinessEndPoint is for kubernetes to check when the container is read to accept traffic
	ReadinessEndPoint = "/readiness"
)

func loadRoutes(conn *sql.DB, tc *store.TokenConfig, logger *log.Logger) http.Handler {
	h := NewHandler(store.NewDB(conn), tc, logger)
	router := httprouter.New()
	allowedOrigins := []string{"*"}

	router.GET(LivenessEndPoint, Liveness)
	router.GET(ReadinessEndPoint, Ready(conn))
	// register routes here
	router.POST(RegisterEndpoint, middleware.CORS(h.Register, allowedOrigins))
	router.POST(LoginEndPoint, middleware.CORS(h.Login, allowedOrigins))
	router.DELETE(DeleteEndpoint, middleware.CORS(middleware.Auth(h.Delete, tc, logger), allowedOrigins))
	router.GET(HomeEndPoint, middleware.CORS(middleware.Auth(Home, tc, logger), allowedOrigins))

	return router
}
