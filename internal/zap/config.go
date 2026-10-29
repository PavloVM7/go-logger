// Package zap provides configuration and utilities for the zap logger implementation.
package zap

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/PavloVM7/go-logger/internal"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config holds the configuration settings for the zap logger.
type Config struct {
	Name       string
	DateFormat string
	Level      string
	Encoding   internal.Encoding
	IsDevMode  bool
}

func (c *Config) check() {
	if len(c.Encoding) == 0 {
		c.Encoding = "json"
	}
	name := strings.TrimSpace(c.Name)
	if len(name) == 0 {
		name = "default"
	}
	c.Name = name
}

func newConfig(conf Config) zap.Config {
	result := zap.Config{
		Level:             zap.NewAtomicLevelAt(stringToZapLevel(conf.Level)),
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
	return result
}

func newEncoderConfig() zapcore.EncoderConfig {
	result := zap.NewProductionEncoderConfig()
	result.EncodeLevel = levelShortStringEncoder
	result.EncodeDuration = zapcore.StringDurationEncoder

	return result
}

func stringToZapLevel(level string) zapcore.Level {
	var levelsMap = map[string]zapcore.Level{
		"debug": zap.DebugLevel,
		"info":  zap.InfoLevel,
		"warn":  zap.WarnLevel,
		"error": zap.ErrorLevel,
		"fatal": zap.FatalLevel,
	}
	result, ok := levelsMap[level]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown level string '%s'", level)
		return zap.DebugLevel
	}
	return result
}

func timeEncoder(dateFormat string) func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(dateFormat))
	}
}
