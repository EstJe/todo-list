package config

import (
	"github.com/EstJe/todo-list/internal/lib/envconv"
	"time"
)

type Config struct {
	Env     string
	Timeout time.Duration
	GRPC    GRPCConfig
	Cache   CacheConfig
	DB      DBConfig
}

type GRPCConfig struct {
	Port int
}

type CacheConfig struct {
	URL string
	TTL time.Duration
}

type DBConfig struct {
	URL string
}

func MustLoad() Config {
	return Config{
		Env:     envconv.String("TODOAPP_ENV"),
		Timeout: envconv.TimeDuration("TODOAPP_TIMEOUT"),
		GRPC: GRPCConfig{
			Port: envconv.Int("TODOAPP_GRPC_PORT"),
		},
		Cache: CacheConfig{
			URL: envconv.String("TODOAPP_CACHE_URL"),
			TTL: envconv.TimeDuration("TODOAPP_CACHE_TTL"),
		},
		DB: DBConfig{
			URL: envconv.String("TODOAPP_DB_URL"),
		},
	}
}
