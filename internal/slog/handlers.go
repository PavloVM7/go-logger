// Package slog provides custom handlers for the standard library's log/slog package.
package slog

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"strings"
	"time"
)

const (
	nameKeyName = "logger"
	spaceByte   = ' '
)

type iNamedHandler interface {
	Name() string
}

func newNamedJSONHandler(w io.Writer, opt *slog.HandlerOptions, dateFormat string) *namedHandler {
	return newNamedHandler(newJSONHandler(w, opt, dateFormat))
}

// revive:disable:cognitive-complexity
func newJSONHandler(w io.Writer, opt *slog.HandlerOptions, dateFormat string) *slog.JSONHandler {
	if len(dateFormat) == 0 {
		dateFormat = time.RFC3339Nano
	}

	originalReplace := opt.ReplaceAttr
	opt.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
		switch a.Key {
		case slog.TimeKey:
			if tm, ok := a.Value.Any().(time.Time); ok {
				a.Value = slog.StringValue(tm.Format(dateFormat))
				return a
			}
		case slog.LevelKey:
			if lvl, ok := a.Value.Any().(slog.Level); ok {
				a.Value = slog.StringValue(levelToShortString(lvl))
				return a
			}
		default:
			if originalReplace != nil {
				return originalReplace(groups, a)
			}
		}
		return a
	}

	return slog.NewJSONHandler(w, opt)
}

type textHandler struct {
	dateFormat string
	h          slog.Handler
	w          io.Writer
	name       string
}

func (h *textHandler) clone() *textHandler {
	return &textHandler{
		dateFormat: h.dateFormat,
		h:          h.h,
		w:          h.w,
		name:       h.name,
	}
}

func (h *textHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}
func (h *textHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.h.WithAttrs(attrs)
}
func (h *textHandler) WithGroup(name string) slog.Handler {
	name = strings.TrimSpace(name)
	if name == "" {
		return h
	}
	result := h.clone()
	if result.name != "" {
		name = strings.Join([]string{result.name, name}, ".")
	}
	result.name = name

	return result
}
func (h *textHandler) Handle(ctx context.Context, r slog.Record) error {
	var buf bytes.Buffer
	buf.WriteString(r.Time.Format(h.dateFormat))
	buf.WriteByte(spaceByte)
	buf.WriteString(levelToShortString(r.Level))
	buf.WriteByte(spaceByte)

	if h.name != "" {
		buf.WriteString(nameKeyName)
		buf.WriteByte('=')
		buf.WriteString(h.name)
		buf.WriteByte(spaceByte)
	}

	if _, err := buf.WriteTo(h.w); err != nil {
		return err
	}

	return h.h.Handle(ctx, r)
}
func (h *textHandler) Name() string {
	return h.name
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
		if a.Key == slog.LevelKey {
			return slog.Attr{}
		}
		if a.Key == nameKeyName {
			return slog.Attr{}
		}
		if originalReplace == nil {
			return a
		}

		return originalReplace(groups, a)
	}

	return &textHandler{h: slog.NewTextHandler(w, opt), w: w, dateFormat: dateFormat}
}

// --- namedHandler ---

type namedHandler struct {
	slog.Handler
	name string
}

func (nh *namedHandler) WithGroup(name string) slog.Handler {
	name = strings.TrimSpace(name)
	if name == "" {
		return nh
	}
	if nh.name != "" {
		name = strings.Join([]string{nh.name, name}, ".")
	}
	return &namedHandler{
		Handler: nh.Handler,
		name:    name,
	}
}

func (nh *namedHandler) Handle(ctx context.Context, r slog.Record) error {
	if nh.name != "" {
		r.AddAttrs(slog.String(nameKeyName, nh.name))
	}

	return nh.Handler.Handle(ctx, r)
}

func (nh *namedHandler) Name() string {
	return nh.name
}

func newNamedHandler(handler slog.Handler) *namedHandler {
	return &namedHandler{Handler: handler}
}
