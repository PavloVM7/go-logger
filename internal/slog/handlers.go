// Package slog provides custom handlers for the standard library's log/slog package.
package slog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	logslog "log/slog"
	"time"
)

func levelStringToShortString(levelStr string) logslog.Level {
	switch levelStr {
	case "DEBUG", "debug":
		return logslog.LevelDebug
	case "INFO", "info":
		return logslog.LevelInfo
	case "WARN", "warn":
		return logslog.LevelWarn
	case "ERROR", "error":
		return logslog.LevelError
	case "FATAL", "fatal", "12", "ERROR+4":
		return LevelFatal
	default:
		return logslog.LevelDebug - 4
	}
}

func levelToShortString(level logslog.Level) string {
	switch level {
	case logslog.LevelDebug:
		return "DBG"
	case logslog.LevelInfo:
		return "INF"
	case logslog.LevelWarn:
		return "WRN"
	case logslog.LevelError:
		return "ERR"
	case LevelFatal:
		return "FTL"
	default:
		return fmt.Sprintf("LVL%d", level)
	}
}

// revive:disable:cognitive-complexity
func newJSONHandler(w io.Writer, opt *logslog.HandlerOptions, dateFormat string) *logslog.JSONHandler {
	if len(dateFormat) == 0 {
		dateFormat = time.RFC3339Nano
	}

	originalReplace := opt.ReplaceAttr
	opt.ReplaceAttr = func(groups []string, a logslog.Attr) logslog.Attr {
		if a.Key == logslog.TimeKey {
			if tm, ok := a.Value.Any().(time.Time); ok {
				a.Value = logslog.StringValue(tm.Format(dateFormat))
				return a
			}
		}
		if a.Key == logslog.LevelKey {
			lvl := levelStringToShortString(a.Value.String())
			a.Value = logslog.StringValue(levelToShortString(lvl))
			return a
		}
		if originalReplace == nil {
			return a
		}
		return originalReplace(groups, a)
	}

	return logslog.NewJSONHandler(w, opt)
}

type textHandler struct {
	dateFormat string
	h          logslog.Handler
	w          io.Writer
}

func (h *textHandler) Enabled(ctx context.Context, level logslog.Level) bool {
	return h.h.Enabled(ctx, level)
}
func (h *textHandler) WithAttrs(attrs []logslog.Attr) logslog.Handler {
	return h.h.WithAttrs(attrs)
}
func (h *textHandler) WithGroup(name string) logslog.Handler {
	return h.h.WithGroup(name)
}
func (h *textHandler) Handle(ctx context.Context, r logslog.Record) error {
	var buf bytes.Buffer
	buf.WriteString(r.Time.Format(h.dateFormat))
	buf.WriteByte(' ')
	buf.WriteString(levelToShortString(r.Level))
	buf.WriteByte(' ')

	if _, err := buf.WriteTo(h.w); err != nil {
		return err
	}

	return h.h.Handle(ctx, r)
}

func newTextHandler(w io.Writer, opt *logslog.HandlerOptions, dateFormat string) *textHandler {
	if len(dateFormat) == 0 {
		dateFormat = time.RFC3339Nano
	}

	originalReplace := opt.ReplaceAttr
	opt.ReplaceAttr = func(groups []string, a logslog.Attr) logslog.Attr {
		if a.Key == logslog.TimeKey {
			return logslog.Attr{}
		}
		if a.Key == logslog.LevelKey {
			return logslog.Attr{}
		}

		if originalReplace == nil {
			return a
		}

		return originalReplace(groups, a)
	}

	return &textHandler{h: logslog.NewTextHandler(w, opt), w: w, dateFormat: dateFormat}
}
