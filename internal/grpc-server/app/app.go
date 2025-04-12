package app

import (
	_ "fmt"
	"github.com/EstJe/todo-list/internal/grpc-server/app/grpc"
	"github.com/EstJe/todo-list/internal/grpc-server/app/storage"
	"github.com/EstJe/todo-list/internal/grpc-server/service/todo"
	"log/slog"
	"time"
)

type App struct {
	GRPC    *grpcapp.App
	Storage *storageapp.App
}

func New(
	log *slog.Logger,
	dbURL string,
	cacheURL string,
	cacheTTL time.Duration,
	grpcAddr string,
) *App {

	sapp, storageCache := storageapp.New(log, dbURL, cacheURL, cacheTTL)
	service := todosrv.New(log, storageCache)
	gapp := grpcapp.New(log, grpcAddr, service)

	return &App{
		GRPC:    gapp,
		Storage: sapp,
	}
}
