package logx

import (
	"context"
	"log/slog"

	"github.com/go-kratos/kratos/v2/log"
)

var _ log.Logger = (*KratosToSlog)(nil)

type KratosToSlog struct {
	logger *slog.Logger
}

func NewKratosToSlog(logger *slog.Logger) *KratosToSlog {
	return &KratosToSlog{
		logger: logger,
	}
}

func KratosLevelToSlog(level log.Level) slog.Level {
	switch level {
	case log.LevelDebug:
		return slog.LevelDebug
	case log.LevelError:
		return slog.LevelError
	case log.LevelFatal:
		return slog.LevelError
	case log.LevelWarn:
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}

func (l *KratosToSlog) Log(level log.Level, keyvals ...interface{}) error {
	l.logger.Log(context.Background(), KratosLevelToSlog(level), "", keyvals...)
	return nil
}
