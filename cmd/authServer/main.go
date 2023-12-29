package main

import (
	"authServer/internal/app"
	"authServer/internal/config"
	"authServer/internal/config/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetUpLogger(cfg.Env)

	log.Debug("Config and logger was initialized")
	log.Debug("Starting application!")

	var application app.App
	application = app.NewGRPC(log, cfg.GRPC.Port, cfg.StorageUrl, cfg.TokenTTL)

	go application.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	log.Info("Stopping application")
	application.Stop()
	log.Info("Application stopped")
}
