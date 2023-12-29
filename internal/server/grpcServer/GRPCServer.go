package grpcServer

import (
	"authServer/internal/DTO"
	"authServer/internal/service/auth"
	"context"
	authServer "github.com/MajorShebka/protos/gen/go/authServer/majorShebka.authServer.v1"
	"google.golang.org/grpc"
	"log/slog"
)

type GRPCServer struct {
	authServer.UnimplementedAuthServer
	auth *auth.Auth
	log  *slog.Logger
}

func Register(gRPC *grpc.Server, auth *auth.Auth, log *slog.Logger) {
	authServer.RegisterAuthServer(gRPC, &GRPCServer{auth: auth, log: log})
}

func (s *GRPCServer) Login(context context.Context, request *authServer.LoginRequest) (*authServer.LoginResponse, error) {
	const op = "GRPCServer.Login"
	log := s.log.With(slog.String("op", op))

	log.Debug("request: " + request.String())

	result, err := (*s.auth).Login(DTO.CustomerDTO{Login: request.Login, Password: request.Password})

	log.Debug("result: " + result)

	return &authServer.LoginResponse{Token: result}, err
}

func (s *GRPCServer) Register(context context.Context, request *authServer.RegisterRequest) (*authServer.RegisterResponse, error) {
	result, err := (*s.auth).Register(DTO.CustomerDTO{Login: request.Login, Password: request.Password})

	return &authServer.RegisterResponse{Token: result}, err
}
