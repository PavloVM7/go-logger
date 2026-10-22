package logger

import "context"

// Logger defines the interface for a structured logger.
type Logger interface {
	LoggerWith
}

// LoggerWith defines the logging methods with structured key-value pairs.
// revive:disable:exported
type LoggerWith interface {
	Fatalw(msg string, keysAndValues ...any)
	Errorw(msg string, keysAndValues ...any)
	Warnw(msg string, keysAndValues ...any)
	Infow(msg string, keysAndValues ...any)
	Debugw(msg string, keysAndValues ...any)
}

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
