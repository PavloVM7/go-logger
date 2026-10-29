package logger

import (
	"github.com/PavloVM7/go-logger/internal"
	"github.com/PavloVM7/go-logger/internal/zap"
)

// NewLoggerZap creates a new Logger instance backed by zap.
func NewLoggerZap(cfg Config) (Logger, error) {
	conf := newZapConfig(cfg)
	return zap.NewWrapperZap(conf)
}

func newZapConfig(cfg Config) zap.Config {
	result := zap.Config{
		Name:       cfg.Name,
		Level:      cfg.Level,
		DateFormat: cfg.DateFormat,
		Encoding:   internal.EncodingFromString(cfg.Encoding),
	}
	if len(result.DateFormat) == 0 {
		result.DateFormat = defaultDateFormat
	}
	return result
}
