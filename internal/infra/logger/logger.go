package logger

import (
	"log/slog"
	"os"

	"loginhub/internal/conf"
)

func NewLogger(
	cfg *conf.Log,
) *slog.Logger {
	var level slog.Leveler
	switch cfg.Level {
	case conf.Level_DEBUG:
		level = slog.LevelDebug
	case conf.Level_INFO:
		level = slog.LevelInfo
	case conf.Level_WARN:
		level = slog.LevelWarn
	case conf.Level_ERROR:
		level = slog.LevelError
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	}))
	slog.SetDefault(logger)

	return logger
}
