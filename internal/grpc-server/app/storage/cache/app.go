package cacheapp

import (
	"github.com/EstJe/todo-list/internal/grpc-server/app/storage/cache/redis"
	"github.com/EstJe/todo-list/internal/grpc-server/storage"
	"log/slog"
	"time"
)

type CacheApp interface {
	MustRun()
	Close()
}

func New(log *slog.Logger, url string, ttl time.Duration) (CacheApp, storage.Cache) {
	return rdapp.NewRedisApp(log, url, ttl)
}
