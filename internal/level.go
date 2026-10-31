// Package internal contains core types and configuration structures for the logger.
package internal

import (
	"fmt"
	"log/slog"
	"strings"

	"go.uber.org/zap/zapcore"
)

const (
	// SlogFatalLevel represents the custom fatal log level for slog.
	SlogFatalLevel slog.Level = 12
)

// LevelShortString represents a short string representation of log levels.
type LevelShortString string

const (
	// DebugShort is the short string representation for debug level.
	DebugShort LevelShortString = "DBG"
	// InfoShort is the short string representation for info level.
	InfoShort LevelShortString = "INF"
	// WarnShort is the short string representation for warn level.
	WarnShort LevelShortString = "WRN"
	// ErrorShort is the short string representation for error level.
	ErrorShort LevelShortString = "ERR"
	// FatalShort is the short string representation for fatal level.
	FatalShort LevelShortString = "FTL"
)

func (ls LevelShortString) String() string {
	return string(ls)
}

// LevelKind represents the different level types across logging libraries.
type LevelKind int

const (
	// LevelKindDebug represents debug level.
	LevelKindDebug LevelKind = iota
	// LevelKindInfo represents info level.
	LevelKindInfo
	// LevelKindWarn represents warn level.
	LevelKindWarn
	// LevelKindError represents error level.
	LevelKindError
	// LevelKindFatal represents fatal level.
	LevelKindFatal
)

// ParseLevelString parses a level string and returns a LevelKind.
// It handles case-insensitive input and common aliases.
// Returns LevelKindInfo as default for empty or unknown strings.
func ParseLevelString(level string) (LevelKind, error) {
	normalized := strings.ToLower(strings.TrimSpace(level))

	switch normalized {
	case "debug":
		return LevelKindDebug, nil
	case "info", "":
		return LevelKindInfo, nil
	case "warn", "warning":
		return LevelKindWarn, nil
	case "error":
		return LevelKindError, nil
	case "fatal":
		return LevelKindFatal, nil
	default:
		return LevelKindInfo, fmt.Errorf("unknown level string '%s'", level)
	}
}

// ToSlogLevel converts LevelKind to slog.Level.
func (lk LevelKind) ToSlogLevel() slog.Level {
	slogLevels := map[LevelKind]slog.Level{
		LevelKindDebug: slog.LevelDebug,
		LevelKindInfo:  slog.LevelInfo,
		LevelKindWarn:  slog.LevelWarn,
		LevelKindError: slog.LevelError,
		LevelKindFatal: SlogFatalLevel,
	}
	if sl, ok := slogLevels[lk]; ok {
		return sl
	}
	return slog.LevelInfo
}

// ToZapLevel converts LevelKind to zapcore.Level.
func (lk LevelKind) ToZapLevel() zapcore.Level {
	zapLevels := map[LevelKind]zapcore.Level{
		LevelKindDebug: zapcore.DebugLevel,
		LevelKindInfo:  zapcore.InfoLevel,
		LevelKindWarn:  zapcore.WarnLevel,
		LevelKindError: zapcore.ErrorLevel,
		LevelKindFatal: zapcore.FatalLevel,
	}
	if zl, ok := zapLevels[lk]; ok {
		return zl
	}
	return zapcore.DebugLevel
}

// ToShortString converts LevelKind to its short string representation.
func (lk LevelKind) ToShortString() string {
	shortStrings := map[LevelKind]string{
		LevelKindDebug: DebugShort.String(),
		LevelKindInfo:  InfoShort.String(),
		LevelKindWarn:  WarnShort.String(),
		LevelKindError: ErrorShort.String(),
		LevelKindFatal: FatalShort.String(),
	}
	if s, ok := shortStrings[lk]; ok {
		return s
	}
	return fmt.Sprintf("LVL%d", lk)
}

// FromSlogLevel converts slog.Level to LevelKind.
func FromSlogLevel(level slog.Level) LevelKind {
	levelMap := map[slog.Level]LevelKind{
		slog.LevelDebug: LevelKindDebug,
		slog.LevelInfo:  LevelKindInfo,
		slog.LevelWarn:  LevelKindWarn,
		slog.LevelError: LevelKindError,
		SlogFatalLevel:  LevelKindFatal,
	}
	if kind, ok := levelMap[level]; ok {
		return kind
	}
	return LevelKindDebug
}

// FromZapLevel converts zapcore.Level to LevelKind.
func FromZapLevel(level zapcore.Level) LevelKind {
	levelMap := map[zapcore.Level]LevelKind{
		zapcore.DebugLevel: LevelKindDebug,
		zapcore.InfoLevel:  LevelKindInfo,
		zapcore.WarnLevel:  LevelKindWarn,
		zapcore.ErrorLevel: LevelKindError,
		zapcore.FatalLevel: LevelKindFatal,
	}
	if kind, ok := levelMap[level]; ok {
		return kind
	}
	return LevelKindDebug
}
