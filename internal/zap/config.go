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

func newConfig(conf Config) zap.Config {
	levelKind, err := internal.ParseLevelString(conf.Level)
	if err != nil {
		fmt.Printf("unknown level string '%s'", conf.Level)
		levelKind = internal.LevelKindDebug
	}

	result := zap.Config{
		Level:             zap.NewAtomicLevelAt(levelKind.ToZapLevel()),
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
	return result
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
