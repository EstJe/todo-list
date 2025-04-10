package todosrv

import (
	"context"
	"github.com/EstJe/todo-list/app/domain/models"
)

type Storage interface {
	CreateTask(ctx context.Context, title string, description string) (int32, error)
	DeleteTask(ctx context.Context, id int32) error
	DoneTask(ctx context.Context, id int32) error
	Tasks(ctx context.Context) ([]models.Task, error)
}

type TodoList struct {
	storage Storage
}

func (tdl *TodoList) Create(ctx context.Context, title string, description string) (int32, error) {
	re
}
