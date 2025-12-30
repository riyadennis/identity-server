package graph

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/riyadennis/identity-server/app/gql/graph/generated"
	"github.com/riyadennis/identity-server/business/store"
	customMiddleware "github.com/riyadennis/identity-server/foundation/middleware"
	"github.com/sirupsen/logrus"
	"github.com/vektah/gqlparser/v2/ast"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var (
	// ErrFailedToStartListener means that the listener couldn't be started
	ErrFailedToStartListener = errors.New("failed to start listener")

	// ErrFailedToStartServer means that the server couldn't be started
	ErrFailedToStartServer = errors.New("failed to start server")
)

// HTTPServer encapsulates two http server operations  that we need to execute in the service
// it is mainly helpful for testing, by creating mocks for http calls.
type HTTPServer interface {
	Shutdown(ctx context.Context) error
	Serve(l net.Listener) error
}

type Server struct {
	Server        HTTPServer
	Store         store.Store
	Authenticator store.Authenticator
	Logger        *logrus.Logger
	TokenConfig   *store.TokenConfig
	ShutDown      chan os.Signal
}

func NewServer(logger *logrus.Logger, port string, store store.Store,
	auth store.Authenticator, tc *store.TokenConfig) *Server {
	resolver := NewResolver(logger, tc, store, auth)
	srv := handler.New(generated.NewExecutableSchema(
		generated.Config{
			Resolvers: resolver,
		},
	))
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})
	addr := fmt.Sprintf(":%s", port)
	return &Server{
		Server: &http.Server{
			Addr:    addr,
			Handler: newRouter(logger, tc, srv),
		},
		Store:         store,
		Authenticator: auth,
		Logger:        logger,
		TokenConfig:   tc,
		ShutDown:      make(chan os.Signal, 1),
	}
}

func (s *Server) Start(port string) error {
	s.Logger.Info("starting service", "port", port)
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		s.Logger.Errorf("failed to start http listener: %v", err)
		return ErrFailedToStartListener
	}
	var sErr error
	go func() {
		s.Logger.Info("service finished starting and is now ready to accept requests")

		// start http listener
		sErr = s.Server.Serve(listener)
		if sErr != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				s.Logger.Errorf("failed to start http server: %v", sErr)
				return
			}
		}
	}()

	return sErr
}

type GraphQLRequest struct {
	OperationName string `json:"operationName"`
	Query         string `json:"query"`
}

// Public operations that don't require authentication
var publicOperations = map[string]bool{
	"Login":        true,
	"register":     true,
	"refreshToken": true,
}

func NeedsAuthMiddleWare(ac customMiddleware.AuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to read request", http.StatusBadRequest)
				return
			}

			// Parse GraphQL request
			var gqlReq GraphQLRequest
			if err := json.Unmarshal(body, &gqlReq); err != nil {
				http.Error(w, "Invalid GraphQL request", http.StatusBadRequest)
				return
			}

			// Restore body for next handler
			r.Body = io.NopCloser(strings.NewReader(string(body)))

			// Check if operation is public
			if publicOperations[gqlReq.OperationName] {
				next.ServeHTTP(w, r)
			} else {
				ac.Auth(next).ServeHTTP(w, r)
			}
		})
	}
}

func newRouter(logger *logrus.Logger, tc *store.TokenConfig, srv *handler.Server) http.Handler {
	chiRouter := chi.NewRouter()

	chiRouter.Use(middleware.RequestID)
	chiRouter.Use(middleware.Recoverer)
	chiRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
	}))
	ac := customMiddleware.AuthConfig{
		TokenConfig: tc,
		Logger:      logger,
	}
	chiRouter.Use(NeedsAuthMiddleWare(ac))
	chiRouter.Handle("/", otelhttp.NewHandler(
		playground.Handler("GraphQL playground", "/graphql"),
		"graphql"))

	chiRouter.Handle("/graphql", srv)
	return chiRouter
}
