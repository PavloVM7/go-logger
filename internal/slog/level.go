package slog

import (
	logslog "log/slog"

	"github.com/PavloVM7/go-logger/logger"
)

const LevelFatal logslog.Level = 12

func loggerLevelToSlogLevel(level logger.Level) logslog.Level {
	switch level {
	case logger.DebugLevel:
		return logslog.LevelDebug
	case logger.InfoLevel:
		return logslog.LevelInfo
	case logger.WarnLevel:
		return logslog.LevelWarn
	case logger.ErrorLevel:
		return logslog.LevelError
	case logger.FatalLevel:
		return LevelFatal
	default:
		return logslog.Level(level)
	}
}
