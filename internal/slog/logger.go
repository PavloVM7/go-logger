package internal_slog

import (
	"context"
	"log/slog"
)

type slogWrapper struct {
	ctx        context.Context
	lg         *slog.Logger
	dateFormat string
}

func newSlogWrapper(h slog.Handler) *slogWrapper {
	if h == nil {
		return nil
	}
	return &slogWrapper{
		ctx: context.Background(),
		lg:  slog.New(h),
	}
}

func (l *slogWrapper) Fatalw(msg string, keysAndValues ...any) {
	l.lg.Log(l.ctx, LevelFatal, msg, keysAndValues...)
}
func (l *slogWrapper) Errorw(msg string, keysAndValues ...any) {
	l.lg.Log(l.ctx, slog.LevelError, msg, keysAndValues...)
}
func (l *slogWrapper) Warnw(msg string, keysAndValues ...any) {
	l.lg.Log(l.ctx, slog.LevelWarn, msg, keysAndValues...)
}
func (l *slogWrapper) Infow(msg string, keysAndValues ...any) {
	l.lg.Log(l.ctx, slog.LevelInfo, msg, keysAndValues...)
}
func (l *slogWrapper) Debugw(msg string, keysAndValues ...any) {
	l.lg.Log(l.ctx, slog.LevelDebug, msg, keysAndValues...)
}
