package slog

import (
	"context"
	"errors"
	"fmt"
	"io"
	logslog "log/slog"
	"os"
	"slices"
	"sync"

	"github.com/PavloVM7/go-logger/internal"
)

var osExit = os.Exit

// NewWrapper returns a new slog wrapper with JSON formatting.
// revive:disable:unexported-return
func NewWrapper(cfg Config) (*slogWrapper, error) {
	cfg.check()
	opt := cfg.createHandlerOptions()
	writers, err := cfg.createWriters()
	if err != nil {
		return nil, fmt.Errorf("failed to create slog wrapper: %w", err)
	}
	writer, err := createWriter(writers)
	if err != nil {
		return nil, fmt.Errorf("failed to create slog wrapper: %w", err)
	}
	var h logslog.Handler
	encoding := internal.EncodingFromString(cfg.Encoding.String())
	if encoding == internal.JSONEncoding {
		h = newJSONHandler(writer, &opt, cfg.DateFormat)
	} else {
		h = newTextHandler(writer, &opt, cfg.DateFormat)
	}

	return newSlogWrapper(h, writers)
}

func newSlogWrapper(h logslog.Handler, writers []io.Writer) (*slogWrapper, error) {
	if h == nil {
		return nil, errors.New("a handler couldn't be nil")
	}

	return &slogWrapper{
		ctx:     context.Background(),
		lg:      logslog.New(h),
		writers: writers,
	}, nil
}

type slogWrapper struct {
	ctx       context.Context
	lg        *logslog.Logger
	writers   []io.Writer
	closeOnce sync.Once
}

func (l *slogWrapper) addWriter(w io.Writer) bool {
	if w == nil {
		return false
	}
	for i := 0; i < len(l.writers); i++ {
		if w == l.writers[i] {
			return false
		}
	}
	l.writers = append(l.writers, w)
	l.writers = slices.Clip(l.writers)
	return true
}

func (l *slogWrapper) Fatalw(msg string, keysAndValues ...any) {
	l.lg.Log(l.ctx, LevelFatal, msg, keysAndValues...)
	_ = l.Close()
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

// Close closes all closable writers (excluding stdout, stderr, stdin).
// It is idempotent and safe to call multiple times.
func (l *slogWrapper) Close() error {
	var closeErr error
	l.closeOnce.Do(func() {
		closeErr = l.closeWriters()
	})
	return closeErr
}

// revive:disable:cognitive-complexity
func (l *slogWrapper) closeWriters() error {
	errs := make([]error, 0, len(l.writers))
	for i := 0; i < len(l.writers); i++ {
		w := l.writers[i]
		if isStdWriter(w) {
			continue
		}
		if cl, ok := w.(closable); ok {
			if err := cl.Close(); err != nil {
				errs = append(errs, fmt.Errorf("failed to close writer %T: %w", w, err))
			}
		}
	}
	return errors.Join(errs...)
}

func isStdWriter(w io.Writer) bool {
	return w == os.Stdout || w == os.Stderr || w == os.Stdin
}

func closeWriters(writers []io.Writer) {
	for _, w := range writers {
		cl, ok := w.(closable)
		if !ok {
			continue
		}
		closeClosable(cl)
	}
}

type closable interface {
	Close() error
}

type flusher interface {
	Flush() error
}

func closeClosable(cl closable) {
	if cl == os.Stdout || cl == os.Stderr || cl == os.Stdin {
		return
	}
	if err := cl.Close(); err != nil {
		fmt.Printf("failed to close closable %T, err: %v\n", cl, err)
	}
}
