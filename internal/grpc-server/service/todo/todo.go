package todosrv

import (
	"context"
	"errors"
	"github.com/EstJe/todo-list/domain/models"
	"github.com/EstJe/todo-list/internal/grpc-server/storage"
	"github.com/EstJe/todo-list/internal/lib/interceptors"
	"github.com/EstJe/todo-list/internal/lib/op"
	"log/slog"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type Storage interface {
	CreateTask(ctx context.Context, title string, description string) (int32, error)
	DeleteTask(ctx context.Context, id int32) error
	DoneTask(ctx context.Context, id int32) error
	Tasks(ctx context.Context) ([]models.Task, error)
}

type TodoService struct {
	log     *slog.Logger
	storage Storage
}

func New(log *slog.Logger, storage Storage) *TodoService {
	return &TodoService{log: log, storage: storage}
}

func (tds *TodoService) CreateTask(ctx context.Context, title string, description string) (int32, error) {
	oper := op.Operation()
	requestID := interceptors.GetRequestID(ctx)

	log := tds.log.With(
		slog.String("op", oper),
		slog.String("request_id", requestID))

	log.Info("start create task")

	id, err := tds.storage.CreateTask(ctx, title, description)
	if err != nil {
		log.Warn("failed to create task", "error", err.Error())
		return id, op.Wrap(err)
	}

	log.Info("task create success")

	return id, nil
}

func (tds *TodoService) DeleteTask(ctx context.Context, id int32) error {
	oper := op.Operation()
	requestID := interceptors.GetRequestID(ctx)

	log := tds.log.With(
		slog.String("op", oper),
		slog.String("request_id", requestID))

	log.Info("start delete task")

	err := tds.storage.DeleteTask(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			log.Warn("task not found")
			return op.Wrap(ErrTaskNotFound)
		}

		log.Warn("failed to delete task", "error", err.Error())
		return op.Wrap(err)
	}

	log.Info("task delete success")

	return nil
}

func (tds *TodoService) MarkTaskDone(ctx context.Context, id int32) error {
	oper := op.Operation()
	requestID := interceptors.GetRequestID(ctx)

	log := tds.log.With(
		slog.String("op", oper),
		slog.String("request_id", requestID))

	log.Info("start mark task done")

	err := tds.storage.DoneTask(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			log.Warn("task not found")
			return op.Wrap(ErrTaskNotFound)
		}

		log.Warn("failed to mark task done", "error", err.Error())
		return op.Wrap(err)
	}

	log.Info("task mark success")

	return nil
}

func (tds *TodoService) GetTasks(ctx context.Context) ([]models.Task, error) {
	oper := op.Operation()
	requestID := interceptors.GetRequestID(ctx)

	log := tds.log.With(
		slog.String("op", oper),
		slog.String("request_id", requestID))

	log.Info("start get tasks")

	tasks, err := tds.storage.Tasks(ctx)
	if err != nil {
		log.Warn("failed to get tasks", "error", err.Error())
		return nil, op.Wrap(err)
	}

	log.Info("tasks get success")

	return tasks, nil
}
