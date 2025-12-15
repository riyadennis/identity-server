package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/riyadennis/identity-server/app/rest"
	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/foundation"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const timeOut = 5 * time.Second

var (
	errEmptyPort           = errors.New("restPort number empty")
	errPortNotAValidNumber = errors.New("restPort number is not a valid number")
	errPortReserved        = errors.New("restPort is a reserved number")
	errPortBeyondRange     = errors.New("restPort is beyond the allowed range")
)

type ProtoServer struct {
	Server *grpc.Server
	port   string
}

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

	err := foundation.ValidatePort(restPort)
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

func (s *Server) RESTHandler(tc *store.TokenConfig, st store.Store, auth store.Authenticator) {
	s.restServer.Handler = rest.LoadRESTEndpoints(tc, s.Logger, st, auth)
}

// Run registers routes and starts a webserver
// and waits to receive from shutdown and error channels
func (s *Server) Run() error {
	// Start the rest server
	go func() {
		s.Logger.Infof("rest server running on port %s", s.restServer.Addr)
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
