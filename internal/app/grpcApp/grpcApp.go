package grpcApp

import (
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type GRPCApp struct {
	Log        *slog.Logger
	GRPCServer *grpc.Server
	Port       int
}

func (a GRPCApp) MustRun() {
	if err := a.run(); err != nil {
		panic("Cannot start gRPC server: " + err.Error())
	}
}

func (a GRPCApp) run() error {
	const op = "grpcapp.Run"
	log := a.Log.With(
		slog.String("op", op),
		slog.Int("port", a.Port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.Port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server is running", slog.String("addr", l.Addr().String()))

	if err := a.GRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a GRPCApp) Stop() {
	const op = "grcpapp.Stop"
	a.Log.With(slog.String("op", op)).Info("Stopping gRPC server", slog.Int("port", a.Port))

	a.GRPCServer.GracefulStop()
}
