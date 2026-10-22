package internal_slog

import (
	"context"
	"io"
	"log/slog"
	"time"
)

func newJsonHandler(w io.Writer, opt *slog.HandlerOptions, dateFormat string) *slog.JSONHandler {
	if len(dateFormat) == 0 {
		dateFormat = time.RFC3339Nano
	}

	originalReplace := opt.ReplaceAttr
	opt.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			if tm, ok := a.Value.Any().(time.Time); ok {
				a.Value = slog.StringValue(tm.Format(dateFormat))
				return a
			}
		}
		if originalReplace == nil {
			return a
		}
		return originalReplace(groups, a)
	}

	return slog.NewJSONHandler(w, opt)
}

type textHandler struct {
	dateFormat string
	h          slog.Handler
	w          io.Writer
}

func (h *textHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}
func (h *textHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.h.WithAttrs(attrs)
}
func (h *textHandler) WithGroup(name string) slog.Handler {
	return h.h.WithGroup(name)
}
func (h *textHandler) Handle(ctx context.Context, r slog.Record) error {
	timeStr := r.Time.Format(h.dateFormat)
	if _, err := h.w.Write([]byte(timeStr)); err != nil {
		return err
	}

	return h.h.Handle(ctx, r)
}

func newTextHandler(w io.Writer, opt *slog.HandlerOptions, dateFormat string) *textHandler {
	if len(dateFormat) == 0 {
		dateFormat = time.RFC3339Nano
	}

	originalReplace := opt.ReplaceAttr
	opt.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			return slog.Attr{}
		}

		if originalReplace == nil {
			return a
		}
		return originalReplace(groups, a)
	}

	return &textHandler{h: slog.NewTextHandler(w, opt), w: w, dateFormat: dateFormat}
}
