package storage

import (
	"context"
	"errors"
	"github.com/EstJe/todo-list/domain/models"
	"github.com/EstJe/todo-list/internal/lib/op"
	"log/slog"
	"time"
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

type Cache interface {
	SetTasks(ctx context.Context, tasks []models.Task, ttl time.Duration) error
	GetTasks(ctx context.Context) ([]models.Task, error)
}

type StorageCache struct {
	log         *slog.Logger
	storage     Storage
	cache       Cache
	cacheTTL    time.Duration
	cacheActual bool
}

func NewStorageCache(log *slog.Logger, s Storage, c Cache, cacheTTL time.Duration) *StorageCache {
	return &StorageCache{
		log:         log,
		storage:     s,
		cache:       c,
		cacheTTL:    cacheTTL,
		cacheActual: false,
	}
}

func (sc *StorageCache) CreateTask(ctx context.Context, title string, description string) (int32, error) {
	id, err := sc.storage.CreateTask(ctx, title, description)
	if err != nil {
		return 0, op.Wrap(err)
	}

	sc.cacheActual = false

	return id, nil
}

func (sc *StorageCache) DeleteTask(ctx context.Context, id int32) error {
	err := sc.storage.DeleteTask(ctx, id)
	if err != nil {
		return op.Wrap(err)
	}

	sc.cacheActual = false

	return nil
}

func (sc *StorageCache) DoneTask(ctx context.Context, id int32) error {
	err := sc.storage.DoneTask(ctx, id)
	if err != nil {
		return op.Wrap(err)
	}

	sc.cacheActual = false

	return nil
}

func (sc *StorageCache) Tasks(ctx context.Context) ([]models.Task, error) {
	oper := op.Operation()

	log := sc.log.With(
		slog.String("op", oper),
	)

	tasks, err := sc.cache.GetTasks(ctx)
	if err != nil {
		return nil, op.Wrap(err)
	}

	if len(tasks) == 0 || !sc.cacheActual {
		log.Debug("cache not active yet")

		tasks, err = sc.storage.Tasks(ctx)
		if err != nil {
			return nil, op.Wrap(err)
		}

		err = sc.cache.SetTasks(ctx, tasks, sc.cacheTTL)
		if err != nil {
			return nil, op.Wrap(err)
		}

		sc.cacheActual = true
	} else {
		log.Debug("cache active")
	}

	return tasks, nil

}
