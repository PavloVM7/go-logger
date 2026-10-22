package logger

import "fmt"

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
	switch string(txt) {
	case "debug":
		return DebugLevel, nil
	case "info", "": // make the zero value useful
		return InfoLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "error":
		return ErrorLevel, nil
	case "fatal":
		return FatalLevel, nil
	default:
		return InfoLevel, fmt.Errorf("%w: '%s'", ErrUnknownLevel, txt)
	}
}
