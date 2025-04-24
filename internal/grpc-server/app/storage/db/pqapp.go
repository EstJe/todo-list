package dbapp

import (
	"database/sql"
	"github.com/EstJe/todo-list/internal/grpc-server/storage"
	"github.com/EstJe/todo-list/internal/lib/op"
	_ "github.com/lib/pq"
	"log/slog"
)

type PostgresApp struct {
	log     *slog.Logger
	conn    *sql.DB
	url     string
	storage *storage.Postgres
}

func NewPostgresApp(log *slog.Logger, url string) (*PostgresApp, *storage.Postgres) {
	storage := &storage.Postgres{}
	return &PostgresApp{
		log:     log,
		url:     url,
		storage: storage,
	}, storage
}

func (a *PostgresApp) MustRun() {
	oper := op.Operation()

	log := a.log.With(
		slog.String("op", oper),
	)

	log.Info("open postgres db", "URL", a.url)

	var err error
	if a.conn, err = sql.Open("postgres", a.url); err != nil {
		log.Error("failed to open database", "error", err.Error())
		panic(op.Wrap(err))
	}
	*a.storage = *storage.NewPostgres(a.conn)

	log.Info("successfully")
}

func (a *PostgresApp) Close() {
	oper := op.Operation()

	log := a.log.With(
		slog.String("op", oper),
	)

	log.Info("close postgres db", "URL", a.url)

	a.conn.Close()

	log.Info("successfully")
}
