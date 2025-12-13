package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/riyadennis/identity-server/app/auth-api/proto"
	"github.com/riyadennis/identity-server/app/auth-api/rest"
	"github.com/riyadennis/identity-server/business/store"
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
	GRPCServer  *ProtoServer
	ServerError chan error
	ShutDown    chan os.Signal
}

// NewServer creates a server instance with error and shutdown channels initialized
func NewServer(logger *logrus.Logger, restPort, gRPCPort string) (*Server, error) {
	errChan := make(chan error, 2)
	shutdown := make(chan os.Signal, 1)

	err := validatePorts(restPort, gRPCPort)
	if err != nil {
		return nil, err
	}

	return &Server{
		Logger: logger,
		GRPCServer: &ProtoServer{
			Server: grpc.NewServer(),
			port:   gRPCPort,
		},
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

func (s *Server) GRPCHandler(logger *logrus.Logger, tc *store.TokenConfig, st store.Store, auth store.Authenticator) {
	identityServer := proto.NewServer(logger, tc, st, auth)
	proto.RegisterIdentityServer(s.GRPCServer.Server, identityServer)
}

// Run registers routes and starts a webserver
// and waits to receive from shutdown and error channels
func (s *Server) Run() error {
	// Start the rest server
	go func() {
		s.Logger.Infof("rest server running on port %s", s.restServer.Addr)
		s.ServerError <- s.restServer.ListenAndServe()
	}()
	// Start the gRPC server
	go func() {
		s.Logger.Infof("gRPC server running on port :%s", s.GRPCServer.port)
		listener, err := net.Listen("tcp", fmt.Sprintf(":%s", s.GRPCServer.port))
		if err != nil {
			s.ServerError <- err
		}
		err = s.GRPCServer.Server.Serve(listener)
		if err != nil {
			s.ServerError <- err
		}
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

func validatePorts(ports ...string) error {
	for _, port := range ports {
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
	}

	return nil
}
