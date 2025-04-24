package storageapp

import (
	"github.com/EstJe/todo-list/internal/grpc-server/app/storage/cache"
	"github.com/EstJe/todo-list/internal/grpc-server/app/storage/db"
	"github.com/EstJe/todo-list/internal/grpc-server/storage"
	"log/slog"
	"time"
)

type DBApp interface {
	MustRun()
	Close()
}

type CacheApp interface {
	MustRun()
	Close()
}

type App struct {
	log          *slog.Logger
	DB           DBApp
	Cache        CacheApp
	storageCache *storage.StorageCache
}

func New(log *slog.Logger, dbURL string, cacheURL string, cacheTTL time.Duration) (*App, *storage.StorageCache) {
	storageCache := &storage.StorageCache{}

	dbApp, db := dbapp.NewPostgresApp(log, dbURL)
	cacheApp, cache := cacheapp.NewRedisApp(log, cacheURL, cacheTTL)
	storageCache = storage.NewStorageCache(log, db, cache, cacheTTL)

	return &App{
		log:          log,
		DB:           dbApp,
		Cache:        cacheApp,
		storageCache: storageCache,
	}, storageCache
}
