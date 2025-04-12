package grpcapp

import (
	"github.com/EstJe/todo-list/app/internal/grpc-server/grpc/todo"
	"github.com/EstJe/todo-list/app/internal/grpc-server/service/todo"
	"github.com/EstJe/todo-list/app/internal/lib/interceptors"
	"github.com/EstJe/todo-list/app/internal/lib/op"
	todoapi "github.com/EstJe/todo-list/internal/protos/todo"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	addr       string
	service    *todosrv.TodoService
}

func New(log *slog.Logger, addr string, service *todosrv.TodoService) *App {
	gRPCServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.RequestIDInterceptor(),
		))

	return &App{log: log, addr: addr, gRPCServer: gRPCServer, service: service}
}

func (a *App) MustRun() {
	oper := op.Operation()

	log := a.log.With(
		slog.String("op", oper),
	)

	log.Info("starting tcp server", "addr", a.addr)

	l, err := net.Listen("tcp", a.addr)
	if err != nil {
		log.Error("failed to listen", "error", err.Error())
		panic(op.Wrap(err))
	}

	log.Info("starting gRPC server", "address", l.Addr().String())

	api := todogrpc.New(a.service)
	todoapi.RegisterTodoServiceServer(a.gRPCServer, api)

	go func() {
		if err := a.gRPCServer.Serve(l); err != nil {
			log.Error("failed to serve", "error", err.Error())
			panic(op.Wrap(err))
		}
	}()

	log.Info("successfully")
}

func (a *App) GracefulShutdown() {
	oper := op.Operation()

	log := a.log.With(
		slog.String("op", oper),
	)

	log.Info("shutdown grpc server", "addr", a.addr)

	a.gRPCServer.GracefulStop()

	log.Info("successfully")
}
