package todogrpc

import (
	"context"
	"errors"
	"github.com/EstJe/todo-list/app/domain/models"
	"github.com/EstJe/todo-list/app/internal/grpc-server/service/todo"
	todoapi "github.com/EstJe/todo-list/internal/protos/todo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (s *serverAPI) CreateTask(ctx context.Context, in *todoapi.CreateTaskRequest) (*todoapi.CreateTaskResponse, error) {
	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	id, err := s.todoservice.CreateTask(ctx, in.GetTitle(), in.GetDescription())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &todoapi.CreateTaskResponse{Id: id}, nil
}

func (s *serverAPI) DeleteTask(ctx context.Context, in *todoapi.DeleteTaskRequest) (*emptypb.Empty, error) {
	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.todoservice.DeleteTask(ctx, in.GetId())
	if err != nil {
		if errors.Is(err, todosrv.ErrTaskNotFound) {
			return nil, status.Error(codes.NotFound, "task not found")
		}

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &emptypb.Empty{}, nil
}

func (s *serverAPI) MarkTaskDone(ctx context.Context, in *todoapi.MarkTaskDoneRequest) (*emptypb.Empty, error) {
	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.todoservice.MarkTaskDone(ctx, in.GetId())
	if err != nil {
		if errors.Is(err, todosrv.ErrTaskNotFound) {
			return nil, status.Error(codes.NotFound, "task not found")
		}

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &emptypb.Empty{}, nil
}

func (s *serverAPI) GetTasks(ctx context.Context, in *emptypb.Empty) (*todoapi.GetTasksResponse, error) {
	tasks, err := s.todoservice.GetTasks(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	taskspb := make([]*todoapi.Task, 0, len(tasks))
	for _, task := range tasks {
		taskspb = append(taskspb, &todoapi.Task{
			Id:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			StatusId:    task.StatusID,
			CreatedAt:   timestamppb.New(task.CreatedAt),
		})
	}

	return &todoapi.GetTasksResponse{Tasks: taskspb}, nil
}
