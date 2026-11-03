package logger

import (
	"context"

	"github.com/PavloVM7/go-logger/internal"
)

// Logger is a structured logger with context support and key-value logging methods.
type Logger internal.Logger

type loggerCtxKey struct{}

var (
	logCtxKey = loggerCtxKey{}
)

// ToContext - set logger to context
func ToContext(ctx context.Context, log Logger) context.Context {
	return context.WithValue(ctx, logCtxKey, log)
}

// FromContext - get logger from context
func FromContext(ctx context.Context) (Logger, error) {
	result, ok := ctx.Value(logCtxKey).(Logger)
	if ok && result != nil {
		return result, nil
	}

	return nil, ErrLoggerNotFound
}
