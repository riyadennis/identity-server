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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/runtime/protoimpl"
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
	helper := business.NewHelper(s.Store, s.Authenticator, s.Logger)
	token, err := helper.Login(ctx, s.TokenConfig, *request.Email, *request.Password)
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

func (s *Server) Me(ctx context.Context, ur *UserRequest) (*UserResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}
	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing authorization header")
	}
	claims, err := validation.ValidateToken(authHeaders[0], s.TokenConfig)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "token failed validation")
	}
	user, err := s.Store.Retrieve(ctx, claims.Subject)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if user == nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to find user for ID %s", claims.Subject))
	}
	fullName := user.FirstName + " " + user.LastName
	return &UserResponse{
		state:     protoimpl.MessageState{},
		ID:        &user.ID,
		Email:     &user.Email,
		Name:      &fullName,
		sizeCache: 0,
	}, nil
}

func (s *Server) mustEmbedUnimplementedIdentityServer() {}
