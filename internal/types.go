// Package internal contains core types and configuration structures for the logger.
// It defines encoding formats, configuration base types, and related utilities
// used throughout the logger implementation.
package internal

// Encoding represents the output encoding format for log entries.
type Encoding string

const (
	// JSONEncoding encodes log entries as JSON.
	JSONEncoding Encoding = "json"
	// ConsoleEncoding encodes log entries in a human-readable console format.
	ConsoleEncoding Encoding = "console"
)

// String returns the string representation of the Encoding.
func (en Encoding) String() string {
	return string(en)
}

// EncodingFromString converts a string to an Encoding value.
// Returns ConsoleEncoding as the default if the string does not match any known encoding.
func EncodingFromString(str string) Encoding {
	encodings := []Encoding{
		JSONEncoding,
		ConsoleEncoding,
	}

	for _, en := range encodings {
		if str == en.String() {
			return en
		}
	}

	return ConsoleEncoding
}

// ConfigBase holds the base configuration settings for a logger.
type ConfigBase struct {
	// Name is the name of the logger instance.
	Name string
	// Level is the minimum log level to capture (e.g., "debug", "info", "warn", "error").
	Level string
	// DateFormat specifies the format string for log timestamps.
	DateFormat string
	// Encoding determines the output format for log entries.
	Encoding Encoding
}

type LevelShortString string

const (
	DebugShort LevelShortString = "DBG"
	InfoShort  LevelShortString = "INF"
	WarnShort  LevelShortString = "WRN"
	ErrorShort LevelShortString = "ERR"
	FatalShort LevelShortString = "FTL"
)

func (ls LevelShortString) String() string {
	return string(ls)
}
