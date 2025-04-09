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
	Redis   RedisConfig
	DB      DBConfig
}

type GRPCConfig struct {
	Port int
}

type RedisConfig struct {
	URL string
	TTL time.Duration
}

type DBConfig struct {
	URL string
}

func MustLoad() Config {
	return Config{
		Env:     envString("ENV"),
		Timeout: envDuration("TIMEOUT"),
		GRPC: GRPCConfig{
			Port: envInt("GRPC_PORT"),
		},
		Redis: RedisConfig{
			URL: envString("REDIS_URL"),
			TTL: envDuration("REDIS_TTL"),
		},
		DB: DBConfig{
			URL: envString("DB_URL"),
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
