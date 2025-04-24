package proxyapp

import (
	"github.com/EstJe/todo-list/domain/dto"
	"github.com/EstJe/todo-list/internal/http-gateway/http/middleware"
	proxy "github.com/EstJe/todo-list/internal/http-gateway/http/proxy"
	"github.com/EstJe/todo-list/internal/lib/op"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"log/slog"
	"net/http"
	"strconv"
)

type App struct {
	log        *slog.Logger
	httpServer *http.Server
}

type BrokerWriter interface {
	WriteTaskAudit(value dto.AuditEvent) error
}

func New(log *slog.Logger, grpcAddr string, httpPort int, brokerWriter BrokerWriter) *App {
	grpcMux := runtime.NewServeMux()

	proxy := proxy.New(log, grpcMux, grpcAddr)
	proxy.RegisterHandlers()

	wrappedMux := middleware.BrokerTaskAuditMiddleware(log, brokerWriter, grpcMux)

	httpServer := &http.Server{
		Addr:    ":" + strconv.Itoa(httpPort),
		Handler: wrappedMux,
	}

	return &App{log: log, httpServer: httpServer}
}

func (a *App) MustRun() {
	oper := op.Operation()

	log := a.log.With(
		slog.String("op", oper),
	)

	log.Info("start listening HTTP", "address", a.httpServer.Addr)

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Error("failed to start HTTP server", "error", err.Error())
			panic(op.Wrap(err))
		}
	}()

	log.Info("successfully")
}

func (a *App) Shutdown() {
	oper := op.Operation()

	log := a.log.With(
		slog.String("op", oper),
	)

	log.Info("shutdown HTTP server", "address", a.httpServer.Addr)
}
