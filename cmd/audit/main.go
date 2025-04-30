package main

import (
	"github.com/EstJe/todo-list/internal/audit/app"
	"github.com/EstJe/todo-list/internal/audit/config"
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
		cfg.Broker.Addr,
		cfg.Broker.Topic,
	)

	application.Broker.ReaderRun()
	application.Visual.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	sign := <-stop

	log.Warn("shutting down audit...", "signal", sign)

	application.Broker.ReaderStop()
	application.Visual.Stop()

	log.Info("audit stopped")
}
