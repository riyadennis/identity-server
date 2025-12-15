package identity

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/riyadennis/identity-server/business"
	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/business/validation"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server struct {
	unImplementedServer UnimplementedIdentityServer
	Server              *grpc.Server
	Store               store.Store
	Authenticator       store.Authenticator
	Logger              *logrus.Logger
	TokenConfig         *store.TokenConfig
	ServerError         chan error
	ShutDown            chan os.Signal
}

func NewServer(logger *logrus.Logger, tc *store.TokenConfig, st store.Store, auth store.Authenticator) *Server {
	gs := grpc.NewServer()
	s := &Server{
		unImplementedServer: UnimplementedIdentityServer{},
		Server:              gs,
		Store:               st,
		Authenticator:       auth,
		Logger:              logger,
		TokenConfig:         tc,
		ShutDown:            make(chan os.Signal, 1),
	}
	RegisterIdentityServer(gs, s)
	return s
}

func (s *Server) Run(port string) error {
	s.Logger.Infof("gRPC server running on port :%s", port)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return err
	}
	err = s.Server.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Login(ctx context.Context, request *LoginRequest) (*LoginResponse, error) {
	s.Logger.Info("processing gRPC request to login")
	err := validation.ValidateEmail(*request.Email)
	if err != nil {
		return nil, err
	}
	helper := business.NewHelper(s.Store, s.Authenticator, s.Logger)
	user, err := helper.UserCredentialsInDB(ctx, *request.Email, *request.Password)
	if err != nil {
		return nil, err
	}
	token, err := helper.ManageToken(ctx, s.TokenConfig, user.ID)
	if err != nil {
		return nil, err
	}
	status := int32(http.StatusOK)
	ttl, err := strconv.ParseInt(token.TokenTTL, 10, 32)
	if err != nil {
		return nil, err
	}
	int32ttl := int32(ttl)
	return &LoginResponse{
		Status:      &status,
		AccessToken: &token.AccessToken,
		Expiry:      &token.Expiry,
		TokenType:   &token.TokenType,
		LastRefresh: &token.LastRefresh,
		TokenTtl:    &int32ttl,
	}, nil
}

func (s *Server) mustEmbedUnimplementedIdentityServer() {}
