package main

import (
	"app/internal/config"
	"app/internal/lib/logger"
)

func main() {
	// TODO: init config
	cfg := config.MustLoad()

	// TODO: init logger
	log := logger.New(cfg.Env)

	log.Error("test")

	// TODO: init PostgreSQL, Redis

	// TODO: init server
}
