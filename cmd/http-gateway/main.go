package main

import (
	"github.com/EstJe/todo-list/internal/http-gateway/app"
	"github.com/EstJe/todo-list/internal/http-gateway/config"
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

	// init app
	application := app.New(
		log,
		cfg.GRPC.Addr,
		cfg.Http.Port,
		cfg.Broker.Addr,
		cfg.Broker.Topic,
	)

	application.Broker.WriterRun()
	application.Proxy.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	sign := <-stop

	log.Warn("shutting down http-gateway...", "signal", sign)

	application.Proxy.Shutdown()
	application.Broker.WriterStop()

	log.Info("http-gateway stopped")
}
