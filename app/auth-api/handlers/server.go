package handlers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"

	"github.com/riyadennis/identity-server/business/store"
)

const timeOut = 5 * time.Second

var (
	errEmptyPort           = errors.New("port number empty")
	errPortNotAValidNumber = errors.New("port number is not a valid number")
	errPortReserved        = errors.New("port is a reserved number")
	errPortBeyondRange     = errors.New("port is beyond the allowed range")
)

// Server have all the set up needed to run and shut down a http server
type Server struct {
	httpServer  http.Server
	listenAddr  string
	serverError chan error
	shutDown    chan os.Signal
}

// NewServer creates a server instance with error and shutdown channels initialised
func NewServer(addr string) *Server {
	errChan := make(chan error, 1)
	shutdown := make(chan os.Signal, 1)

	err := validatePort(addr)
	if err != nil {
		errChan <- err
	}

	return &Server{
		httpServer: http.Server{
			Addr:         ":" + addr,
			ReadTimeout:  timeOut,
			WriteTimeout: timeOut,
		},
		listenAddr:  addr,
		serverError: errChan,
		shutDown:    shutdown,
	}
}

// Run registers routes and starts web server
// and waits to receive from shutdown and error channels
func (s *Server) Run(conn *sql.DB, tc *store.TokenConfig, logger *log.Logger) error {
	h := NewHandler(store.NewDB(conn), tc, logger)

	router := httprouter.New()
	// register routes here
	router.POST(RegisterEndpoint, h.Register)
	router.POST(LoginEndPoint, h.Login)
	router.POST(DeleteEndpoint, Auth(h.Delete, tc, logger))
	router.GET(HomeEndPoint, Auth(Home, tc, logger))

	s.httpServer.Handler = cors.Default().Handler(router)

	// Start the service
	go func() {
		logger.Printf("server running on port %s", s.httpServer.Addr)
		s.serverError <- s.httpServer.ListenAndServe()
	}()

	select {
	case err := <-s.serverError:
		return err
	case sig := <-s.shutDown:
		logger.Printf("main: %v: Start shutdown", sig)
		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), timeOut)
		defer cancel()

		err := s.httpServer.Shutdown(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func validatePort(port string) error {
	if port == "" {
		return errEmptyPort
	}

	addr, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		return errPortNotAValidNumber
	}

	if addr < 1024 {
		return errPortReserved
	}

	if addr > 65535 {
		return errPortBeyondRange
	}

	return nil
}
