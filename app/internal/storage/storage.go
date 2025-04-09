package storage

import (
	"app/domain/models"
	"app/internal/lib/op"
	"context"
	"errors"
	"time"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type Storage struct {
	db          Postgres
	redis       Redis
	redisTTL    time.Duration
	cacheActual bool
}

func New(dbURL string, redisURL string, redisTTL time.Duration) (*Storage, error) {
	db, err := NewPostgres(dbURL)
	if err != nil {
		return nil, op.Wrap(err)
	}

	rd := NewRedis(redisURL)

	return &Storage{db: db, redis: rd, redisTTL: redisTTL}, nil
}

func (s *Storage) CreateTask(ctx context.Context, title string, description string) (int32, error) {
	id, err := s.db.CreateTask(ctx, title, description)
	if err != nil {
		return 0, op.Wrap(err)
	}

	s.cacheActual = false

	return id, nil
}

func (s *Storage) DeleteTask(ctx context.Context, id int32) error {
	err := s.db.DeleteTask(ctx, id)
	if err != nil {
		return op.Wrap(err)
	}

	s.cacheActual = false

	return nil
}

func (s *Storage) DoneTask(ctx context.Context, id int32) error {
	err := s.db.DoneTask(ctx, id)
	if err != nil {
		return op.Wrap(err)
	}

	s.cacheActual = false

	return nil
}

func (s *Storage) Tasks(ctx context.Context) ([]models.Task, error) {
	if !s.cacheActual {
		tasks, err := s.db.Tasks(ctx)
		if err != nil {
			return nil, op.Wrap(err)
		}

		err = s.redis.SetTasks(ctx, tasks, s.redisTTL)
		if err != nil {
			return nil, op.Wrap(err)
		}

		s.cacheActual = true

		return tasks, nil
	} else {
		tasks, err := s.redis.GetTasks(ctx)
		if err != nil {
			return nil, op.Wrap(err)
		}
		return tasks, nil
	}
}
