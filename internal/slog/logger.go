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
	writer, err := cfg.createWriter()
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

	return newSlogWrapper(h)
}

func newSlogWrapper(h logslog.Handler) (*slogWrapper, error) {
	if h == nil {
		return nil, errors.New("a handler couldn't be nil")
	}

	return &slogWrapper{
		ctx: context.Background(),
		lg:  logslog.New(h),
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
	return nil
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
