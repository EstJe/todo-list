package todogrpc

import (
	"context"
	"github.com/EstJe/todo-list/app/domain/models"
	todoapi "github.com/EstJe/todo-list/internal/protos/todo"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TodoService interface {
	CreateTask(ctx context.Context, title string, description string) (int32, error)
	DeleteTask(ctx context.Context, id int32) error
	MarkTaskDone(ctx context.Context, id int32) error
	GetTasks(ctx context.Context) ([]models.Task, error)
}

type serverAPI struct {
	todoapi.UnimplementedTodoServiceServer
	todoservice TodoService
}

func New(service TodoService) *serverAPI {
	return &serverAPI{todoservice: service}
}

func Register(gRPC *grpc.Server, api *serverAPI) {
	todoapi.RegisterTodoServiceServer(gRPC, api)
}

func (s *serverAPI) CreateTask(ctx context.Context, req *todoapi.CreateTaskRequest) (*todoapi.CreateTaskResponse, error) {

}

func (s *serverAPI) DeleteTask(ctx context.Context, req *todoapi.DeleteTaskRequest) (*emptypb.Empty, error) {

}

func (s *serverAPI) MarkTaskDone(ctx context.Context, req *todoapi.MarkTaskDoneRequest) (*emptypb.Empty, error) {

}

func (s *serverAPI) GetTasks(ctx context.Context, empty *emptypb.Empty) (*todoapi.ListTasksResponse, error) {

}
