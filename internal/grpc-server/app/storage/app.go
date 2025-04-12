package storageapp

import (
	"github.com/EstJe/todo-list/app/internal/grpc-server/app/storage/cache"
	"github.com/EstJe/todo-list/app/internal/grpc-server/app/storage/db"
	"github.com/EstJe/todo-list/app/internal/grpc-server/storage"
	"log/slog"
	"time"
)

type App struct {
	log          *slog.Logger
	DB           dbapp.DBApp
	Cache        cacheapp.CacheApp
	storageCache *storage.StorageCache
}

func New(log *slog.Logger, dbURL string, cacheURL string, cacheTTL time.Duration) (*App, *storage.StorageCache) {
	storageCache := &storage.StorageCache{}

	dbApp, db := dbapp.New(log, dbURL)
	cacheApp, cache := cacheapp.New(log, cacheURL, cacheTTL)
	storageCache = storage.NewStorageCache(log, db, cache, cacheTTL)

	return &App{
		log:          log,
		DB:           dbApp,
		Cache:        cacheApp,
		storageCache: storageCache,
	}, storageCache
}
