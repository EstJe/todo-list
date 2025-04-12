package dbapp

import (
	"github.com/EstJe/todo-list/app/internal/grpc-server/app/storage/db/postgres"
	"github.com/EstJe/todo-list/app/internal/grpc-server/storage"
	"log/slog"
)

type DBApp interface {
	MustRun()
	Close()
}

func New(log *slog.Logger, url string) (DBApp, storage.Storage) {
	return pqapp.NewPostgresApp(log, url)
}
