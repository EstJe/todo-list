package brokerapp

import (
	"github.com/EstJe/todo-list/internal/audit/broker"
	"github.com/EstJe/todo-list/internal/lib/op"
	"github.com/segmentio/kafka-go"
	"log/slog"
)

type KafkaApp struct {
	log   *slog.Logger
	conn  *kafka.Reader
	addr  string
	topic string
	kafka *broker.KafkaReader
}

func NewKafkaApp(log *slog.Logger, addr string, topic string) (*KafkaApp, *broker.KafkaReader) {
	kafka := &broker.KafkaReader{}

	return &KafkaApp{log: log, addr: addr, topic: topic, kafka: kafka}, kafka
}

func (a *KafkaApp) ReaderRun() {
	oper := op.Operation()

	log := a.log.With(
		slog.String("op", oper),
	)

	log.Info("create kafka reader conn", "address", a.addr, "topic", a.topic)

	a.conn = kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{a.addr},
		Topic:       a.topic,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		StartOffset: kafka.FirstOffset,
	})

	*a.kafka = *broker.NewKafkaReader(log, a.conn)

	log.Info("successfully")
}

func (a *KafkaApp) ReaderStop() {
	oper := op.Operation()

	log := a.log.With(
		slog.String("op", oper),
	)

	log.Info("close kafka writer conn", "address", a.addr)

	a.conn.Close()

	log.Info("successfully")
}
