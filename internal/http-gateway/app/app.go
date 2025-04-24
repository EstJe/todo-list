package app

import (
	brokerapp "github.com/EstJe/todo-list/internal/http-gateway/app/broker"
	proxyapp "github.com/EstJe/todo-list/internal/http-gateway/app/proxy"
	"log/slog"
)

type BrokerApp interface {
	WriterRun()
	WriterStop()
}

type App struct {
	Proxy  *proxyapp.App
	Broker BrokerApp
}

func New(log *slog.Logger, grpcAddrProxy string, httpPortProxy int, addrBroker, topicBroker string) *App {
	bapp, bwriter := brokerapp.NewKafkaApp(log, addrBroker, topicBroker)
	papp := proxyapp.New(log, grpcAddrProxy, httpPortProxy, bwriter)

	return &App{Proxy: papp, Broker: bapp}
}
