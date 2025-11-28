package proto

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/riyadennis/identity-server/business"
	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/business/validation"
	"github.com/sirupsen/logrus"
)

type Server struct {
	unImplementedServer UnimplementedIdentityServer
	Store               store.Store
	Authenticator       store.Authenticator
	Logger              *logrus.Logger
	TokenConfig         *store.TokenConfig
}

func NewServer(conn *sql.DB, logger *logrus.Logger, tc *store.TokenConfig) *Server {
	auth := &store.Auth{
		Conn:   conn,
		Logger: logger,
	}
	return &Server{
		unImplementedServer: UnimplementedIdentityServer{},
		Store:               store.NewDB(conn),
		Authenticator:       auth,
		Logger:              logger,
		TokenConfig:         tc,
	}
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
