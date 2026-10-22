package slog

import (
	"context"
	logslog "log/slog"
	"os"
)

var osExit = os.Exit

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
