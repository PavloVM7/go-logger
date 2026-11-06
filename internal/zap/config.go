// Package zap provides configuration and utilities for the zap logger implementation.
package zap

import (
	"fmt"
	"strings"
	"time"

	"github.com/PavloVM7/go-logger/internal"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config holds the configuration settings for the zap logger.
type Config struct {
	internal.ConfigBase
	IsDevMode bool
}

func (c *Config) check() {
	if len(c.Encoding) == 0 {
		c.Encoding = internal.JSONEncoding
	}
	name := strings.TrimSpace(c.Name)
	if len(name) == 0 {
		name = "default"
	}
	c.Name = name
}

func newConfig(conf Config) (zap.Config, error) {
	levelKind, err := internal.ParseLevelString(conf.Level)
	if err != nil {
		return zap.Config{}, fmt.Errorf("unknown logging level string '%s'", conf.Level)
	}

	result := zap.Config{
		Level:             zap.NewAtomicLevelAt(kindLevelToZapLevel(levelKind)),
		Development:       conf.IsDevMode,
		DisableCaller:     !conf.IsDevMode,
		DisableStacktrace: true,
		Sampling:          nil,
		Encoding:          conf.Encoding.String(),
		EncoderConfig:     newEncoderConfig(),
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
	}
	if len(conf.DateFormat) > 0 {
		result.EncoderConfig.EncodeTime = timeEncoder(conf.DateFormat)
	}
	if len(conf.Outputs) > 0 {
		result.OutputPaths = conf.Outputs
	}
	return result, nil
}

func newEncoderConfig() zapcore.EncoderConfig {
	result := zap.NewProductionEncoderConfig()
	result.EncodeLevel = levelShortStringEncoder
	result.EncodeDuration = zapcore.StringDurationEncoder

	return result
}

func timeEncoder(dateFormat string) func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(dateFormat))
	}
}
