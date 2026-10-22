package slog

import (
	"bytes"
	"encoding/json/v2"
	logslog "log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func init() {
	osExit = func(code int) {}
}

func TestDateFormat(t *testing.T) {
	const dateFormat = "2006-01-02 15:04:05.000Z07:00"
	t.Run("JSON handler", func(t *testing.T) {
		buf := bytes.Buffer{}
		opt := logslog.HandlerOptions{}
		jh := newJSONHandler(&buf, &opt, dateFormat)
		l := logslog.New(jh)
		l.Info("debug message", "key1", 123)
		mp := make(map[string]any)
		err := json.Unmarshal(buf.Bytes(), &mp)
		require.NoError(t, err)
		sdt, ok := mp[logslog.TimeKey].(string)
		require.True(t, ok)
		_, err = time.Parse(dateFormat, sdt)
		require.NoError(t, err)
		_, err = time.Parse(time.DateOnly, sdt)
		require.Error(t, err)
	})

	t.Run("text handler", func(t *testing.T) {
		buf := bytes.Buffer{}
		opt := logslog.HandlerOptions{}
		jh := newTextHandler(&buf, &opt, dateFormat)
		l := logslog.New(jh)
		l.Info("debug message", "key1", 123)
		sub := buf.String()[0:len(dateFormat)]
		_, err := time.Parse(dateFormat, sub)
		require.NoError(t, err)
		_, err = time.Parse(time.DateOnly, sub)
		require.Error(t, err)
	})
}

func TestLevelAbbreviations(t *testing.T) {
	tests := []struct {
		name   string
		abbrev string
		logFn  func(*logslog.Logger)
	}{
		{"debug", "DBG", func(l *logslog.Logger) { l.Debug("test message") }},
		{"info", "INF", func(l *logslog.Logger) { l.Info("test message") }},
		{"warn", "WRN", func(l *logslog.Logger) { l.Warn("test message") }},
		{"error", "ERR", func(l *logslog.Logger) { l.Error("test message") }},
	}

	t.Run("JSON handler", func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				buf := bytes.Buffer{}
				opt := logslog.HandlerOptions{Level: logslog.LevelDebug}
				jh := newJSONHandler(&buf, &opt, time.RFC3339Nano)
				l := logslog.New(jh)
				tt.logFn(l)

				mp := make(map[string]any)
				err := json.Unmarshal(buf.Bytes(), &mp)
				require.NoError(t, err)
				level, ok := mp[logslog.LevelKey].(string)
				require.True(t, ok, "level key should exist in output")
				require.Equal(t, tt.abbrev, level, "level should be abbreviated")
			})
		}
	})

	t.Run("text handler", func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				buf := bytes.Buffer{}
				opt := logslog.HandlerOptions{Level: logslog.LevelDebug}
				th := newTextHandler(&buf, &opt, time.RFC3339Nano)
				l := logslog.New(th)
				tt.logFn(l)

				output := buf.String()
				require.Contains(t, output, tt.abbrev, "output should contain abbreviated level")
			})
		}
	})

	t.Run("fatal level with slogWrapper", func(t *testing.T) {
		t.Run("JSON handler", func(t *testing.T) {
			buf := bytes.Buffer{}
			opt := logslog.HandlerOptions{Level: logslog.LevelDebug}
			jh := newJSONHandler(&buf, &opt, time.RFC3339Nano)
			wrapper := newSlogWrapper(jh)
			wrapper.Fatalw("fatal message", "key", "value")

			mp := make(map[string]any)
			err := json.Unmarshal(buf.Bytes(), &mp)
			require.NoError(t, err)
			level, ok := mp[logslog.LevelKey].(string)
			require.True(t, ok, "level key should exist in output")
			require.Equal(t, "FTL", level, "fatal level should be abbreviated")
		})

		t.Run("text handler", func(t *testing.T) {
			buf := bytes.Buffer{}
			opt := logslog.HandlerOptions{Level: logslog.LevelDebug}
			th := newTextHandler(&buf, &opt, time.RFC3339Nano)
			wrapper := newSlogWrapper(th)
			wrapper.Fatalw("fatal message", "key", "value")

			output := buf.String()
			require.Contains(t, output, "FTL", "output should contain abbreviated fatal level")
		})
	})
}

func TestTextLogger_empty_message(t *testing.T) {
	tests := []struct {
		name   string
		abbrev string
		logFn  func(*logslog.Logger)
	}{
		{"debug", "DBG", func(l *logslog.Logger) { l.Debug("") }},
		{"info", "INF", func(l *logslog.Logger) { l.Info("") }},
		{"warn", "WRN", func(l *logslog.Logger) { l.Warn("") }},
		{"error", "ERR", func(l *logslog.Logger) { l.Error("") }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.Buffer{}
			opt := logslog.HandlerOptions{Level: logslog.LevelDebug}
			th := newTextHandler(&buf, &opt, time.RFC3339Nano)
			l := logslog.New(th)
			tt.logFn(l)

			output := buf.String()
			require.Contains(t, output, tt.abbrev, "output should contain abbreviated level")
		})
	}
}
