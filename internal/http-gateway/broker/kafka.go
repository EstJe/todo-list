package broker

import (
	"context"
	"encoding/json"
	"github.com/EstJe/todo-list/domain/dto"
	"github.com/EstJe/todo-list/internal/lib/op"
	"github.com/segmentio/kafka-go"
	"log/slog"
)

type KafkaWriter struct {
	log    *slog.Logger
	writer *kafka.Writer
}

func NewKafkaWriter(log *slog.Logger, w *kafka.Writer) *KafkaWriter {
	return &KafkaWriter{log: log, writer: w}
}

func (kw *KafkaWriter) WriteTaskAudit(value dto.AuditEvent) error {
	oper := op.Operation()

	log := kw.log.With(
		slog.String("op", oper),
	)

	log.Info("writing audit")

	ctx := context.Background()

	jsonData, err := json.Marshal(value)
	if err != nil {
		log.Error("failed to marshal audit event", "error", err.Error())
		return op.Wrap(err)
	}

	err = kw.writer.WriteMessages(ctx,
		kafka.Message{
			Value: jsonData,
		},
	)
	if err != nil {
		log.Warn("failed to write messages", "error", err.Error())
		return op.Wrap(err)
	}

	log.Info("writing audit done")

	return nil
}
