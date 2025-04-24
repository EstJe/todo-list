package grpcapp

import (
	"context"
	todoapi "github.com/EstJe/todo-list/api/gen"
	"github.com/EstJe/todo-list/domain/models"
	"github.com/EstJe/todo-list/internal/grpc-server/grpc/todo"
	"github.com/EstJe/todo-list/internal/lib/interceptors"
	"github.com/EstJe/todo-list/internal/lib/op"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"strconv"
)

type TodoService interface {
	CreateTask(ctx context.Context, title string, description string) (int32, error)
	DeleteTask(ctx context.Context, id int32) error
	MarkTaskDone(ctx context.Context, id int32) error
	GetTasks(ctx context.Context) ([]models.Task, error)
}

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
	service    TodoService
}

func New(log *slog.Logger, port int, service TodoService) *App {
	gRPCServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.RequestIDInterceptor(),
		))

	return &App{log: log, port: port, gRPCServer: gRPCServer, service: service}
}

func (a *App) MustRun() {
	oper := op.Operation()

	log := a.log.With(
		slog.String("op", oper),
	)

	log.Info("starting tcp server", "port", a.port)

	l, err := net.Listen("tcp", ":"+strconv.Itoa(a.port))
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

	log.Info("shutdown grpc server", "port", a.port)

	a.gRPCServer.GracefulStop()

	log.Info("successfully")
}
