package visualapp

import (
	"context"
	"github.com/EstJe/todo-list/domain/dto"
	"github.com/EstJe/todo-list/internal/lib/op"
	"log/slog"
)

type BrokerReader interface {
	ReadTaskAudit() (dto.AuditEvent, error)
}

type TempApp struct {
	log          *slog.Logger
	brokerReader BrokerReader
	cancelBR     context.CancelFunc
}

func NewTempApp(log *slog.Logger, br BrokerReader) *TempApp {
	return &TempApp{log: log, brokerReader: br}
}

func (a *TempApp) MustRun() {
	oper := op.Operation()

	log := a.log.With(
		slog.String("op", oper),
	)

	ctx, cancel := context.WithCancel(context.Background())
	a.cancelBR = cancel

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				auditEvent, err := a.brokerReader.ReadTaskAudit()
				if err != nil {
					log.Warn("audit data get error", "error", err.Error())
				}
				log.Info("audit data get success", "event", auditEvent)
			}
		}
	}()
}

func (a *TempApp) Stop() {
	oper := op.Operation()

	log := a.log.With(
		slog.String("op", oper),
	)

	log.Info("close visual app")

	a.cancelBR()

	log.Info("successfully")
}
