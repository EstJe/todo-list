package cacheapp

import (
	"github.com/EstJe/todo-list/internal/grpc-server/storage"
	"github.com/EstJe/todo-list/internal/lib/op"
	"github.com/go-redis/redis/v8"
	"log/slog"
	"time"
)

type RedisApp struct {
	log     *slog.Logger
	conn    *redis.Client
	url     string
	ttl     time.Duration
	storage *storage.Redis
}

func NewRedisApp(log *slog.Logger, url string, ttl time.Duration) (*RedisApp, *storage.Redis) {
	storage := &storage.Redis{}
	return &RedisApp{
		log:     log,
		url:     url,
		ttl:     ttl,
		storage: storage,
	}, storage
}

func (a *RedisApp) MustRun() {
	oper := op.Operation()

	log := a.log.With(
		slog.String("op", oper),
	)

	log.Info("open redis cache", "URL", a.url)

	opts, err := redis.ParseURL(a.url)
	if err != nil {
		panic(op.Wrap(err))
	}
	a.conn = redis.NewClient(opts)
	*a.storage = *storage.NewRedis(a.conn)

	log.Info("successfully")
}

func (a *RedisApp) Close() {
	oper := op.Operation()

	log := a.log.With(
		slog.String("op", oper),
	)

	log.Info("close redis cache", "URL", a.url)

	a.conn.Close()

	log.Info("successfully")
}
