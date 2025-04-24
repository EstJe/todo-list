package app

import (
	brokerapp "github.com/EstJe/todo-list/internal/audit/app/broker"
	"github.com/EstJe/todo-list/internal/audit/app/visual"
	"log/slog"
)

type BrokerApp interface {
	ReaderRun()
	ReaderStop()
}

type VisualApp interface {
	MustRun()
	Stop()
}

type App struct {
	Broker BrokerApp
	Visual VisualApp
}

func New(log *slog.Logger, addrBroker, topicBroker string) *App {
	bapp, breader := brokerapp.NewKafkaApp(log, addrBroker, topicBroker)

	vapp := visualapp.NewTempApp(log, breader)

	return &App{Broker: bapp, Visual: vapp}
}
