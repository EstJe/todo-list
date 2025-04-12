package config

import (
	"os"
	"strconv"
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
	Addr string
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
		Env:     envString("TODOAPP_ENV"),
		Timeout: envDuration("TODOAPP_TIMEOUT"),
		GRPC: GRPCConfig{
			Addr: envString("TODOAPP_GRPC_ADDR"),
		},
		Cache: CacheConfig{
			URL: envString("TODOAPP_CACHE_URL"),
			TTL: envDuration("TODOAPP_CACHE_TTL"),
		},
		DB: DBConfig{
			URL: envString("TODOAPP_DB_URL"),
		},
	}
}

func envString(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(key + " environment variable not set")
	}
	return value
}

func envDuration(key string) time.Duration {
	value := envString(key)

	duration, err := time.ParseDuration(value)
	if err != nil {
		panic("Invalid " + key + " value: " + err.Error())
	}
	return duration
}

func envInt(key string) int {
	value := envString(key)

	intValue, err := strconv.Atoi(value)
	if err != nil {
		panic("Invalid " + key + " value: " + err.Error())
	}
	return intValue
}
