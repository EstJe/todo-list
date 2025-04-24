package middleware

import (
	"github.com/EstJe/todo-list/domain/dto"
	"log/slog"
	"net/http"
	"time"
)

type BrokerWriter interface {
	WriteTaskAudit(value dto.AuditEvent) error
}

func BrokerTaskAuditMiddleware(log *slog.Logger, writer BrokerWriter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auditEvent := dto.AuditEvent{
			Method:     r.Method,
			Path:       r.URL.Path,
			RemoteAddr: r.RemoteAddr,
			Timestamp:  time.Now(),
		}

		go func() {
			err := writer.WriteTaskAudit(auditEvent)
			if err != nil {
				log.Error("failed to send audit log to Kafka", "error", err.Error())
			}
		}()

		next.ServeHTTP(w, r)
	})
}
