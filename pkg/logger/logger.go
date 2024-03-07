package logger

import (
	"fmt"
	"log/slog"
	"os"
	"tn/internal/config"
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func SetupLogger(cfg *config.Config) (*slog.Logger, error) {
	var log *slog.Logger
	file, err := os.OpenFile(fmt.Sprintf("%s.log", cfg.ServiceName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	switch cfg.Env {
	case "local":
		log = slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		log = slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log.With("env", cfg.Env), nil

}
