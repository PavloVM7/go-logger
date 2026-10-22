// Package slog provides custom handlers for the standard library's log/slog package.
package slog

import (
	"context"
	"fmt"
	"io"
	logslog "log/slog"
	"time"
)

func levelStringToShortString(levelStr string) string {
	switch levelStr {
	case "DEBUG", "debug":
		return "DBG"
	case "INFO", "info":
		return "INF"
	case "WARN", "warn":
		return "WRN"
	case "ERROR", "error":
		return "ERR"
	case "FATAL", "fatal", "12", "ERROR+4":
		return "FTL"
	default:
		return levelStr
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
			levelStr := a.Value.String()
			a.Value = logslog.StringValue(levelStringToShortString(levelStr))
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
	timeStr := r.Time.Format(h.dateFormat)
	if _, err := h.w.Write([]byte(timeStr)); err != nil {
		return err
	}
	space := []byte(" ")
	if _, err := h.w.Write(space); err != nil {
		return err
	}
	levelStr := levelToShortString(r.Level)
	if _, err := h.w.Write([]byte(levelStr)); err != nil {
		return err
	}
	if _, err := h.w.Write(space); err != nil {
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
