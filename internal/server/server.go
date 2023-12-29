package server

import (
	"context"
	authServer "github.com/MajorShebka/protos/gen/go/authServer/majorShebka.authServer.v1"
)

type Server interface {
	Login(ctx context.Context, request *authServer.LoginRequest) (*authServer.LoginResponse, error)
	Register(ctx context.Context, request *authServer.RegisterRequest) (*authServer.RegisterResponse, error)
}
