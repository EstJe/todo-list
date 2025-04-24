package config

import (
	"github.com/EstJe/todo-list/internal/lib/envconv"
)

type Config struct {
	Env    string
	Broker BrokerConfig
}

type BrokerConfig struct {
	Addr  string
	Topic string
}

func MustLoad() Config {
	return Config{
		Env: envconv.String("AUDIT_ENV"),
		Broker: BrokerConfig{
			Addr:  envconv.String("AUDIT_BROKER_ADDR"),
			Topic: envconv.String("AUDIT_BROKER_TOPIC"),
		},
	}
}
