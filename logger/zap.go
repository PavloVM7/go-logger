package logger

import (
	"fmt"

	"github.com/PavloVM7/go-logger/internal"
	"github.com/PavloVM7/go-logger/internal/zap"
)

// NewLoggerZap creates a new Logger instance backed by zap.
func NewLoggerZap(cfg Config) (Logger, error) {
	conf, err := newZapConfig(cfg)
	if err != nil {
		return nil, err
	}
	return zap.NewWrapperZap(conf)
}

func newZapConfig(cfg Config) (zap.Config, error) {
	strLevel, err := cfg.levelToString()
	if err != nil {
		return zap.Config{}, fmt.Errorf("failed to parse string level: %w", err)
	}
	result := zap.Config{
		Name:       cfg.Name,
		Level:      strLevel,
		DateFormat: cfg.DateFormat,
		Encoding:   internal.EncodingFromString(cfg.Encoding),
		Outputs:    cfg.Outputs,
	}
	if len(result.DateFormat) == 0 {
		result.DateFormat = defaultDateFormat
	}
	return result, nil
}
