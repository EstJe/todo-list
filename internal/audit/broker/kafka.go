package broker

import (
	"context"
	"encoding/json"
	"github.com/EstJe/todo-list/domain/dto"
	"github.com/EstJe/todo-list/internal/lib/op"
	"github.com/segmentio/kafka-go"
	"log/slog"
)

type KafkaReader struct {
	log    *slog.Logger
	reader *kafka.Reader
}

func NewKafkaReader(log *slog.Logger, r *kafka.Reader) *KafkaReader {
	return &KafkaReader{log: log, reader: r}
}

func (kr *KafkaReader) ReadTaskAudit() (dto.AuditEvent, error) {
	oper := op.Operation()

	log := kr.log.With(
		slog.String("op", oper),
	)

	log.Info("read audit")

	msg, err := kr.reader.ReadMessage(context.Background())
	if err != nil {
		log.Warn("read message error", "error", err.Error())
		return dto.AuditEvent{}, op.Wrap(err)
	}

	auditEvent := dto.AuditEvent{}
	if err := json.Unmarshal(msg.Value, &auditEvent); err != nil {
		log.Warn("unmarshal audit done error", "error", err.Error())
		return dto.AuditEvent{}, op.Wrap(err)
	}

	log.Info("reading audit done")

	return auditEvent, nil
}
