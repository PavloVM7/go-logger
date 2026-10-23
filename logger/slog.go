package logger

import (
	"io"
	logslog "log/slog"

	"github.com/PavloVM7/go-logger/internal/slog"
)

func loggerLevelToSlogLevel(level Level) logslog.Level {
	switch level {
	case DebugLevel:
		return logslog.LevelDebug
	case InfoLevel:
		return logslog.LevelInfo
	case WarnLevel:
		return logslog.LevelWarn
	case ErrorLevel:
		return logslog.LevelError
	case FatalLevel:
		return slog.LevelFatal
	default:
		return logslog.Level(level)
	}
}

// NewLoggerSlogJSON creates a new slog-based logger with JSON formatting.
func NewLoggerSlogJSON(cfg Config) (Logger, error) {
	c, err := newSlogConfig(cfg, nil)
	if err != nil {
		return nil, err
	}

	return newWrapperJSON(c), nil
}

func newWrapperJSON(c slog.Config) Logger {
	return slog.NewWrapperJSON(c)
}

// NewLoggerSlogText creates a new slog-based logger with text formatting.
func NewLoggerSlogText(cfg Config) (Logger, error) {
	c, err := newSlogConfig(cfg, nil)
	if err != nil {
		return nil, err
	}

	return slog.NewWrapperText(c), nil
}

func newSlogConfig(cfg Config, w io.Writer) (slog.Config, error) {
	level, err := cfg.GetLevel()
	if err != nil {
		return slog.Config{}, err
	}
	result := slog.Config{
		Level:      loggerLevelToSlogLevel(level),
		DateFormat: cfg.DateFormat,
		W:          w,
	}
	if len(result.DateFormat) == 0 {
		result.DateFormat = defaultDateFormat
	}

	return result, nil
}
