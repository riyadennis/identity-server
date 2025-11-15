package server

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/riyadennis/identity-server/app/auth-api/rest"
	"github.com/riyadennis/identity-server/business/store"
	"github.com/sirupsen/logrus"
)

const timeOut = 5 * time.Second

var (
	errEmptyPort           = errors.New("port number empty")
	errPortNotAValidNumber = errors.New("port number is not a valid number")
	errPortReserved        = errors.New("port is a reserved number")
	errPortBeyondRange     = errors.New("port is beyond the allowed range")
)

// Server have all the setup needed to run and shut down a http server
type Server struct {
	Logger      *logrus.Logger
	restServer  http.Server
	ServerError chan error
	ShutDown    chan os.Signal
}

// NewServer creates a server instance with error and shutdown channels initialized
func NewServer(logger *logrus.Logger, restPort string) (*Server, error) {
	errChan := make(chan error, 2)
	shutdown := make(chan os.Signal, 1)

	err := validatePort(restPort)
	if err != nil {
		return nil, err
	}
	return &Server{
		Logger: logger,
		restServer: http.Server{
			Addr:         ":" + restPort,
			ReadTimeout:  timeOut,
			WriteTimeout: timeOut,
			// to prevent SlowLoris attack
			ReadHeaderTimeout: timeOut,
			// close idle open connections
			IdleTimeout: 120 * time.Second,
			// set to 1MB
			MaxHeaderBytes: 1 << 20,
		},
		ServerError: errChan,
		ShutDown:    shutdown,
	}, nil
}

func (s *Server) RESTHandler(conn *sql.DB, tc *store.TokenConfig) {
	s.restServer.Handler = rest.LoadRESTEndpoints(conn, tc, s.Logger)
}

// Run registers routes and starts a webserver
// and waits to receive from shutdown and error channels
func (s *Server) Run() error {
	// Start the service
	go func() {
		s.Logger.Infof("server running on port %s", s.restServer.Addr)
		s.ServerError <- s.restServer.ListenAndServe()
	}()

	select {
	case err := <-s.ServerError:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	case sig := <-s.ShutDown:
		s.Logger.Infof("main: %v: Start shutdown", sig)
		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), timeOut)
		defer cancel()

		err := s.restServer.Shutdown(ctx)
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
