package storage

import (
	"app/domain/models"
	"app/internal/lib/op"
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"time"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(redisURL string) Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "",
		DB:       0,
	})

	return Redis{client: client}
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
	val, err := r.client.Get(ctx, "tasks").Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}

		return nil, op.Wrap(err)
	}

	var tasks []models.Task
	err = json.Unmarshal([]byte(val), &tasks)
	if err != nil {
		return nil, op.Wrap(err)
	}

	return tasks, err
}
