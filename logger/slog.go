package logger

import (
	"io"

	"github.com/PavloVM7/go-logger/internal"
	"github.com/PavloVM7/go-logger/internal/slog"
)

// NewLoggerSlogJSON creates a new slog-based logger with JSON formatting.
func NewLoggerSlogJSON(cfg Config) (Logger, error) {
	c, err := newSlogConfig(cfg, nil)
	if err != nil {
		return nil, err
	}

	return newWrapperJSON(c), nil
}

func newWrapperJSON(c slog.Config) Logger {
	return slog.NewWrapper(c)
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
	result := slog.Config{
		Level:      cfg.Level,
		DateFormat: cfg.DateFormat,
		Encoding:   internal.EncodingFromString(cfg.Encoding),
		W:          w,
	}
	if len(result.DateFormat) == 0 {
		result.DateFormat = defaultDateFormat
	}

	return result, nil
}
