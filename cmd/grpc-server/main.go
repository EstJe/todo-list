package main

import (
	"github.com/EstJe/todo-list/internal/grpc-server/app"
	"github.com/EstJe/todo-list/internal/grpc-server/config"
	"github.com/EstJe/todo-list/internal/lib/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// init config
	cfg := config.MustLoad()

	// init logger
	log := logger.New(cfg.Env)

	// init server
	application := app.New(
		log,
		cfg.DB.URL,
		cfg.Cache.URL,
		cfg.Cache.TTL,
		cfg.GRPC.Addr,
	)

	application.Storage.Cache.MustRun()
	application.Storage.DB.MustRun()
	application.GRPC.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	sign := <-stop

	log.Warn("shutting down application...", "signal", sign)

	application.GRPC.GracefulShutdown()
	application.Storage.DB.Close()
	application.Storage.Cache.Close()

	log.Info("application stopped")
}
