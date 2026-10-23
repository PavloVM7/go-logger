package slog

import (
	"fmt"
	logslog "log/slog"
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
