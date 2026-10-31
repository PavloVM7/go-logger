package slog

import (
	logslog "log/slog"

	"github.com/PavloVM7/go-logger/internal"
)

// LevelFatal represents the fatal log level, higher than error.
const LevelFatal logslog.Level = 12

func stringToLevel(level string) (logslog.Level, error) {
	kind, err := internal.ParseLevelString(level)
	if err != nil {
		return logslog.LevelInfo, err
	}
	return kind.ToSlogLevel(), nil
}

func levelToShortString(level logslog.Level) string {
	kind := internal.FromSlogLevel(level)
	return kind.ToShortString()
}
