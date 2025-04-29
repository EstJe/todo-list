package logger

import (
	"github.com/EstJe/todo-list/internal/lib/logger/prettylogger"
	"github.com/EstJe/todo-list/internal/lib/op"
	"io"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func New(env string, pathLogFile string) (*slog.Logger, error) {
	return setupLogger(env, pathLogFile)
}

func NewMock() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func setupLogger(env string, pathLogFile string) (*slog.Logger, error) {
	// TODO: сделать создание нового файла для логов при превышении 20МБ размера файла

	var log *slog.Logger

	// Открываем файл
	logFile, err := os.OpenFile(pathLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("Failed to open log file", "err", err)
		return nil, op.Wrap(err)
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	switch env {
	case envLocal:
		opts := prettylogger.PrettyHandlerOptions{
			SlogOpts: &slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		}

		handler := opts.NewPrettyHandler(multiWriter)

		log = slog.New(handler)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log, nil
}
