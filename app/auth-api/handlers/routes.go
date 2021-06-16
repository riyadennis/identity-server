package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/foundation/middleware"
)

func loadRoutes(conn *sql.DB, tc *store.TokenConfig, logger *log.Logger) http.Handler {
	h := NewHandler(store.NewDB(conn), tc, logger)
	router := httprouter.New()
	allowedOrigins := []string{"*"}

	// register routes here
	router.POST(RegisterEndpoint, middleware.CORS(h.Register, allowedOrigins))
	router.POST(LoginEndPoint, middleware.CORS(h.Login, allowedOrigins))
	router.DELETE(DeleteEndpoint, middleware.CORS(middleware.Auth(h.Delete, tc, logger), allowedOrigins))
	router.GET(HomeEndPoint, middleware.CORS(middleware.Auth(Home, tc, logger), allowedOrigins))

	return router
}
