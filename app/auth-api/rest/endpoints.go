package rest

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/riyadennis/identity-server/business/store"
	customMiddleware "github.com/riyadennis/identity-server/foundation/middleware"
)

const (
	// RegisterEndpoint is to create a new user
	RegisterEndpoint = "/register"

	// DeleteEndpoint is to delete a user
	DeleteEndpoint = "/delete/{userID}"

	// LoginEndPoint creates a token for the  user of credentials are valid
	LoginEndPoint = "/login"

	// HomeEndPoint is the details end point that a
	// logged-in user with a valid token can access
	HomeEndPoint = "/home"

	// LivenessEndPoint is for kubernetes to check when to restart the container
	LivenessEndPoint = "/liveness"

	// ReadinessEndPoint is for kubernetes to check when the container is read to accept traffic
	ReadinessEndPoint = "/readiness"
)

// LoadRESTEndpoints adds REST endpoints to the router
func LoadRESTEndpoints(conn *sql.DB, tc *store.TokenConfig, logger *log.Logger) http.Handler {
	h := NewHandler(store.NewDB(conn), store.NewDB(conn), tc, logger)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// Set a timeout value on the request context (ctx) that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.SetHeader("Content-Type", "application/json"))
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc: func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value isn't ignored by any of the major browsers
	}))
	r.Get(LivenessEndPoint, Liveness)
	r.Get(ReadinessEndPoint, Ready(conn))
	r.Post(RegisterEndpoint, h.Register)
	r.Post(LoginEndPoint, h.Login)
	// register routes here
	r.Route("/user", func(r chi.Router) {
		ac := customMiddleware.AuthConfig{
			TokenConfig: tc,
			Logger:      logger,
		}

		r.Use(ac.Auth)
		r.Delete(DeleteEndpoint, h.Delete)
		r.Get(HomeEndPoint, Home)
	})

	return r
}
