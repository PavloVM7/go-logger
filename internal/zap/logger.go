package zap

import (
	"github.com/PavloVM7/go-logger/internal"
	"go.uber.org/zap"
)

// NewWrapperZap creates a zap logger wrapper.
// revive:disable:unexported-return
func NewWrapperZap(cfg Config) (*zapLogger, error) {
	return newZapLogger(cfg)
}

func newZapLogger(cfg Config, opt ...zap.Option) (*zapLogger, error) {
	cfg.check()

	conf := newConfig(cfg)
	zl, err := conf.Build(opt...)
	if err != nil {
		return nil, err
	}
	zl = zl.Named(cfg.Name)

	return &zapLogger{
		SugaredLogger: zl.Sugar(),
	}, nil
}

// zapLogger implements Logger interface using zap SugaredLogger.
type zapLogger struct {
	*zap.SugaredLogger
}

func (l *zapLogger) Close() error {
	return l.Sync()
}

func (l *zapLogger) With(keysAndValues ...any) internal.Logger {
	return &zapLogger{
		SugaredLogger: l.SugaredLogger.With(keysAndValues...),
	}
}
