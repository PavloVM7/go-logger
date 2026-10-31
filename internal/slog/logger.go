package slog

import (
	"context"
	"errors"
	"fmt"
	"io"
	logslog "log/slog"
	"os"
	"slices"

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
	ctx     context.Context
	lg      *logslog.Logger
	writers []io.Writer
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

func (l *slogWrapper) Close() error {
	errs := make([]error, 0, len(l.writers))
	for i := 0; i < len(l.writers); i++ {
		w := l.writers[i]
		if w == os.Stdout || w == os.Stderr || w == os.Stdin {
			continue
		}
		if cl, ok := w.(closable); ok {
			errs = append(errs, cl.Close())
		}
	}

	return errors.Join(errs...)
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

func closeClosable(cl closable) {
	if cl == os.Stdout || cl == os.Stderr || cl == os.Stdin {
		return
	}
	if err := cl.Close(); err != nil {
		fmt.Printf("failed to close closable %T, err: %v\n", cl, err)
	}
}
