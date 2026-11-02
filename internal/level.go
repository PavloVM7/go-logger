// Package internal contains core types and configuration structures for the logger.
package internal

import (
	"fmt"
	"strings"
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

var shortStrings = map[LevelKind]string{
	LevelKindDebug: DebugShort.String(),
	LevelKindInfo:  InfoShort.String(),
	LevelKindWarn:  WarnShort.String(),
	LevelKindError: ErrorShort.String(),
	LevelKindFatal: FatalShort.String(),
}

// ToShortString converts LevelKind to its short string representation.
func (lk LevelKind) ToShortString() string {
	if s, ok := shortStrings[lk]; ok {
		return s
	}
	return fmt.Sprintf("LVL%d", lk)
}
