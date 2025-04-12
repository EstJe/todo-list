package storage

import (
	"context"
	"encoding/json"
	"github.com/EstJe/todo-list/app/domain/models"
	"github.com/EstJe/todo-list/app/internal/lib/op"
	"github.com/go-redis/redis/v8"
	"time"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(client *redis.Client) *Redis {
	return &Redis{client: client}
}

func (r *Redis) SetTasks(ctx context.Context, tasks []models.Task, ttl time.Duration) error {
	jsonData, err := json.Marshal(tasks)
	if err != nil {
		return op.Wrap(err)
	}

	err = r.client.Set(ctx, "tasks", jsonData, ttl).Err()
	if err != nil {
		return op.Wrap(err)
	}

	return nil
}

func (r *Redis) GetTasks(ctx context.Context) ([]models.Task, error) {
	val, err := r.client.Get(ctx, "tasks").Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}

		return nil, op.Wrap(err)
	}

	var tasks []models.Task
	err = json.Unmarshal(val, &tasks)
	if err != nil {
		return nil, op.Wrap(err)
	}

	return tasks, err
}
