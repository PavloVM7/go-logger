package slog

import (
	"fmt"
	logslog "log/slog"

	"github.com/PavloVM7/go-logger/internal"
)

// LevelFatal represents the fatal log level, higher than error.
const LevelFatal logslog.Level = 12

func stringToLevel(level string) (logslog.Level, error) {
	var levelMap = map[string]logslog.Level{
		"debug": logslog.LevelDebug,
		"info":  logslog.LevelInfo,
		"warn":  logslog.LevelWarn,
		"error": logslog.LevelError,
		"fatal": LevelFatal,
	}
	result, ok := levelMap[level]
	if !ok {
		return logslog.LevelInfo, fmt.Errorf("unknown level string '%s'", level)
	}

	return result, nil
}

func levelStringToShortString(levelStr string) logslog.Level {
	switch levelStr {
	case "DEBUG", "debug":
		return logslog.LevelDebug
	case "INFO", "info":
		return logslog.LevelInfo
	case "WARN", "warn":
		return logslog.LevelWarn
	case "ERROR", "error":
		return logslog.LevelError
	case "FATAL", "fatal", "12", "ERROR+4":
		return LevelFatal
	default:
		return logslog.LevelDebug - 4
	}
}

func levelToShortString(level logslog.Level) string {
	switch level {
	case logslog.LevelDebug:
		return internal.DebugShort.String()
	case logslog.LevelInfo:
		return internal.InfoShort.String()
	case logslog.LevelWarn:
		return internal.WarnShort.String()
	case logslog.LevelError:
		return internal.ErrorShort.String()
	case LevelFatal:
		return internal.FatalShort.String()
	default:
		return fmt.Sprintf("LVL%d", level)
	}
}
