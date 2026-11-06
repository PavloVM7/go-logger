package slog

import (
	"fmt"
	"log/slog"

	"github.com/PavloVM7/go-logger/internal"
)

// LevelFatal represents the fatal log level, higher than error.
const LevelFatal slog.Level = 12

var shortStrings = map[slog.Level]string{
	slog.LevelDebug: internal.DebugShort.String(),
	slog.LevelInfo:  internal.InfoShort.String(),
	slog.LevelWarn:  internal.WarnShort.String(),
	slog.LevelError: internal.ErrorShort.String(),
	LevelFatal:      internal.FatalShort.String(),
}

func levelToShortString(level slog.Level) string {
	if res, ok := shortStrings[level]; ok {
		return res
	}
	return fmt.Sprintf("LVL%d", level)
}

func stringToLevel(level string) (slog.Level, error) {
	kind, err := internal.ParseLevelString(level)
	if err != nil {
		return slog.LevelInfo, err
	}
	return kindLevelToSlogLevel(kind)
}

var slogLevels = map[internal.LevelKind]slog.Level{
	internal.LevelKindDebug: slog.LevelDebug,
	internal.LevelKindInfo:  slog.LevelInfo,
	internal.LevelKindWarn:  slog.LevelWarn,
	internal.LevelKindError: slog.LevelError,
	internal.LevelKindFatal: LevelFatal,
}

// ToSlogLevel converts LevelKind to slog.Level.
func kindLevelToSlogLevel(lk internal.LevelKind) (slog.Level, error) {
	if sl, ok := slogLevels[lk]; ok {
		return sl, nil
	}
	return slog.LevelInfo, fmt.Errorf("unknown logging level: %v", lk)
}
