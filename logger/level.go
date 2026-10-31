package logger

import (
	"fmt"

	"github.com/PavloVM7/go-logger/internal"
)

// Level represents the severity level of a log message.
type Level int8

const (
	// DebugLevel is the lowest severity level for detailed debugging information.
	DebugLevel Level = -1
	// InfoLevel is the default severity level for general informational messages.
	InfoLevel Level = 0
	// WarnLevel indicates potentially harmful situations.
	WarnLevel Level = 1
	// ErrorLevel indicates error events that might still allow the application to continue running.
	ErrorLevel Level = 2
	// FatalLevel indicates very severe error events that will presumably lead the application to abort.
	FatalLevel Level = 3
)

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	default:
		return fmt.Sprintf("Level(%d)", l)
	}
}

// CapitalString returns the level name in uppercase.
func (l Level) CapitalString() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return fmt.Sprintf("LEVEL(%d)", l)
	}
}

// ShortCapitalString returns a short (3-character) uppercase representation of the level.
func (l Level) ShortCapitalString() string {
	switch l {
	case DebugLevel:
		return "DEB"
	case InfoLevel:
		return "INF"
	case WarnLevel:
		return "WAR"
	case ErrorLevel:
		return "ERR"
	case FatalLevel:
		return "FAT"
	default:
		return fmt.Sprintf("LEV(%d)", l)
	}
}

// ParseLevel converts a string to a Level, returning an error if the string is invalid.
func ParseLevel(txt string) (Level, error) {
	kind, err := internal.ParseLevelString(txt)
	if err != nil {
		return InfoLevel, fmt.Errorf("%w: '%s'", ErrUnknownLevel, txt)
	}

	levelMap := map[internal.LevelKind]Level{
		internal.LevelKindDebug: DebugLevel,
		internal.LevelKindInfo:  InfoLevel,
		internal.LevelKindWarn:  WarnLevel,
		internal.LevelKindError: ErrorLevel,
		internal.LevelKindFatal: FatalLevel,
	}
	if level, ok := levelMap[kind]; ok {
		return level, nil
	}
	return InfoLevel, nil
}
