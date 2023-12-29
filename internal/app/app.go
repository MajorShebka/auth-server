package app

import (
	"authServer/internal/app/grpcApp"
	"authServer/internal/repository/customerRepository"
	"authServer/internal/server/grpcServer"
	"authServer/internal/service/auth"
	"google.golang.org/grpc"
	"log/slog"
	"time"
)

type App interface {
	MustRun()
	Stop()
}

func NewGRPC(log *slog.Logger,
	grpcPort int,
	storageUrl string,
	tokenTTL time.Duration) App {
	gRPCServer := grpc.NewServer()

	repo := customerRepository.InitCustomerRepository(storageUrl, log)
	authService := auth.InitAuth(repo, log, tokenTTL)
	grpcServer.Register(gRPCServer, authService, log)

	var app App
	app = grpcApp.GRPCApp{
		Log:        log,
		GRPCServer: gRPCServer,
		Port:       grpcPort,
	}

	return app
}
