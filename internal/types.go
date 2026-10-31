// Package internal contains core types and configuration structures for the logger.
// It defines encoding formats, configuration base types, and related utilities
// used throughout the logger implementation.
package internal

const (
	// OutputStdout represents the standard output stream.
	OutputStdout = "stdout"
	// OutputStderr represents the standard error stream.
	OutputStderr = "stderr"
)

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
	Outputs  []string
}

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
