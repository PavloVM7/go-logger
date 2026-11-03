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

// Logger defines the interface for a structured logger.
type Logger interface {
	LoggerFmt
	LoggerWith
	With(keysAndValues ...any) Logger
	IsLoggingLevelDebug() bool
	IsLoggingLevelInfo() bool
	IsLoggingLevelWarn() bool
	IsLoggingLevelError() bool
	Close() error
}

// LoggerFmt defines the formatted logging methods.
type LoggerFmt interface {
	Fatalf(format string, args ...any)
	Errorf(format string, args ...any)
	Warnf(format string, args ...any)
	Infof(format string, args ...any)
	Debugf(format string, args ...any)
}

// LoggerWith defines the logging methods with structured key-value pairs.
// revive:disable:exported
type LoggerWith interface {
	Fatalw(msg string, keysAndValues ...any)
	Errorw(msg string, keysAndValues ...any)
	Warnw(msg string, keysAndValues ...any)
	Infow(msg string, keysAndValues ...any)
	Debugw(msg string, keysAndValues ...any)
}
