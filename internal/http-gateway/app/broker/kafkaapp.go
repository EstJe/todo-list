package brokerapp

import (
	"github.com/EstJe/todo-list/internal/http-gateway/broker"
	"github.com/EstJe/todo-list/internal/lib/op"
	"github.com/segmentio/kafka-go"
	"log/slog"
)

type KafkaApp struct {
	log   *slog.Logger
	conn  *kafka.Writer
	addr  string
	topic string
	kafka *broker.KafkaWriter
}

func NewKafkaApp(log *slog.Logger, addr string, topic string) (*KafkaApp, *broker.KafkaWriter) {
	kafka := &broker.KafkaWriter{}

	return &KafkaApp{log: log, addr: addr, topic: topic, kafka: kafka}, kafka
}

func (a *KafkaApp) WriterRun() {
	oper := op.Operation()

	log := a.log.With(
		slog.String("op", oper),
	)

	log.Info("create kafka writer conn", "address", a.addr, "topic", a.topic)

	a.conn = &kafka.Writer{
		Addr:     kafka.TCP(a.addr),
		Topic:    a.topic,
		Balancer: &kafka.LeastBytes{},
		Async:    false,
	}

	*a.kafka = *broker.NewKafkaWriter(log, a.conn)

	log.Info("successfully")

}

func (a *KafkaApp) WriterStop() {
	oper := op.Operation()

	log := a.log.With(
		slog.String("op", oper),
	)

	log.Info("close kafka writer conn", "address", a.addr)

	a.conn.Close()

	log.Info("successfully")
}
