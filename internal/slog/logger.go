package slog

import (
	"context"
	"io"
	logslog "log/slog"
	"os"
)

var osExit = os.Exit

// Config holds the configuration for a slog logger.
type Config struct {
	Name       string
	Level      logslog.Level
	DateFormat string
	W          io.Writer
}

func (c *Config) check() {
	if c.W == nil {
		c.W = os.Stdout
	}
}

func (c *Config) createHandlerOptions() logslog.HandlerOptions {
	return logslog.HandlerOptions{
		Level:     c.Level,
		AddSource: false,
	}
}

// NewWrapperJSON returns a new slog wrapper with JSON formatting.
// revive:disable-next-line:unexported-return
func NewWrapperJSON(cfg Config) *slogWrapper {
	cfg.check()
	opt := cfg.createHandlerOptions()
	h := newJSONHandler(cfg.W, &opt, cfg.DateFormat)

	return newSlogWrapper(h)
}

// NewWrapperText returns a new slog wrapper with text formatting.
// revive:disable-next-line:unexported-return
func NewWrapperText(cfg Config) *slogWrapper {
	cfg.check()
	opt := cfg.createHandlerOptions()
	h := newTextHandler(cfg.W, &opt, cfg.DateFormat)

	return newSlogWrapper(h)
}

type slogWrapper struct {
	ctx context.Context
	lg  *logslog.Logger
}

func newSlogWrapper(h logslog.Handler) *slogWrapper {
	if h == nil {
		return nil
	}
	return &slogWrapper{
		ctx: context.Background(),
		lg:  logslog.New(h),
	}
}

func (l *slogWrapper) Fatalw(msg string, keysAndValues ...any) {
	l.lg.Log(l.ctx, LevelFatal, msg, keysAndValues...)
	osExit(1)
}
func (l *slogWrapper) Errorw(msg string, keysAndValues ...any) {
	l.lg.Log(l.ctx, logslog.LevelError, msg, keysAndValues...)
}
func (l *slogWrapper) Warnw(msg string, keysAndValues ...any) {
	l.lg.Log(l.ctx, logslog.LevelWarn, msg, keysAndValues...)
}
func (l *slogWrapper) Infow(msg string, keysAndValues ...any) {
	l.lg.Log(l.ctx, logslog.LevelInfo, msg, keysAndValues...)
}
func (l *slogWrapper) Debugw(msg string, keysAndValues ...any) {
	l.lg.Log(l.ctx, logslog.LevelDebug, msg, keysAndValues...)
}
