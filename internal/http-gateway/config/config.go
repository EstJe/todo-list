package config

import (
	"github.com/EstJe/todo-list/internal/lib/envconv"
)

type Config struct {
	Env    string
	Http   HTTPConfig
	GRPC   GRPCConfig
	Broker BrokerConfig
}

type HTTPConfig struct {
	Port int
}

type GRPCConfig struct {
	Addr string
}

type BrokerConfig struct {
	Addr  string
	Topic string
}

func MustLoad() Config {
	return Config{
		Env: envconv.String("PROXY_ENV"),
		Http: HTTPConfig{
			Port: envconv.Int("PROXY_HTTP_PORT"),
		},
		GRPC: GRPCConfig{
			Addr: envconv.String("PROXY_GRPC_ADDR"),
		},
		Broker: BrokerConfig{
			Addr:  envconv.String("PROXY_BROKER_ADDR"),
			Topic: envconv.String("PROXY_BROKER_TOPIC"),
		},
	}
}
