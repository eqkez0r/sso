package authgrpc

import (
	"context"
	ssov1 "github.com/eqkez0r/ssoproto/gen/go/sso"
	"google.golang.org/grpc"
)

type Server struct {
	ssov1.UnimplementedAuthServer
}

func RegisterServerAPI(s *grpc.Server) {
	ssov1.RegisterAuthServer(s, &Server{})
}

func (s *Server) Login(
	ctx context.Context,
	in *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	return &ssov1.LoginResponse{}, nil
}

func (s *Server) Register(
	ctx context.Context,
	in *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	return &ssov1.RegisterResponse{}, nil
}

func (s *Server) IsAdmin(
	ctx context.Context,
	in *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
	return &ssov1.IsAdminResponse{}, nil
}
