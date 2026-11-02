package logger

import (
	"github.com/PavloVM7/go-logger/internal"
	"github.com/PavloVM7/go-logger/internal/slog"
)

// NewLoggerSlog creates a new slog-based logger with JSON formatting.
func NewLoggerSlog(cfg Config) (Logger, error) {
	c, err := newSlogConfig(cfg)
	if err != nil {
		return nil, err
	}

	return newWrapper(c)
}
func newWrapper(c slog.Config) (Logger, error) {
	return slog.NewWrapper(c)
}

func newSlogConfig(cfg Config) (slog.Config, error) {
	stringLevel, err := cfg.levelToString()
	if err != nil {
		return slog.Config{}, err
	}
	result := slog.Config{
		Level:      stringLevel,
		DateFormat: cfg.DateFormat,
		Encoding:   internal.EncodingFromString(cfg.Encoding),
		Outputs:    cfg.Outputs,
	}
	if len(result.DateFormat) == 0 {
		result.DateFormat = defaultDateFormat
	}

	return result, nil
}
