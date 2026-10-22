package internal_slog

import (
	"log/slog"

	"github.com/PavloVM7/go-logger/logger"
)

const LevelFatal slog.Level = 12

func loggerLevelToSlogLevel(level logger.Level) slog.Level {
	switch level {
	case logger.DebugLevel:
		return slog.LevelDebug
	case logger.InfoLevel:
		return slog.LevelInfo
	case logger.WarnLevel:
		return slog.LevelWarn
	case logger.ErrorLevel:
		return slog.LevelError
	case logger.FatalLevel:
		return LevelFatal
	default:
		return slog.LevelDebug
	}
}
