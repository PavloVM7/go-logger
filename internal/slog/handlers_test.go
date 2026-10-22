package internal_slog

import (
	"bytes"
	"encoding/json/v2"
	"log/slog"
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
		opt := slog.HandlerOptions{}
		jh := newJsonHandler(&buf, &opt, dateFormat)
		l := slog.New(jh)
		l.Info("debug message", "key1", 123)
		mp := make(map[string]any)
		err := json.Unmarshal(buf.Bytes(), &mp)
		require.NoError(t, err)
		sdt, ok := mp[slog.TimeKey].(string)
		require.True(t, ok)
		_, err = time.Parse(dateFormat, sdt)
		require.NoError(t, err)
		_, err = time.Parse(time.DateOnly, sdt)
		require.Error(t, err)
	})

	t.Run("text handler", func(t *testing.T) {
		buf := bytes.Buffer{}
		opt := slog.HandlerOptions{}
		jh := newTextHandler(&buf, &opt, dateFormat)
		l := slog.New(jh)
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
		logFn  func(*slog.Logger)
	}{
		{"debug", "DBG", func(l *slog.Logger) { l.Debug("test message") }},
		{"info", "INF", func(l *slog.Logger) { l.Info("test message") }},
		{"warn", "WRN", func(l *slog.Logger) { l.Warn("test message") }},
		{"error", "ERR", func(l *slog.Logger) { l.Error("test message") }},
	}

	t.Run("JSON handler", func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				buf := bytes.Buffer{}
				opt := slog.HandlerOptions{Level: slog.LevelDebug}
				jh := newJsonHandler(&buf, &opt, time.RFC3339Nano)
				l := slog.New(jh)
				tt.logFn(l)

				mp := make(map[string]any)
				err := json.Unmarshal(buf.Bytes(), &mp)
				require.NoError(t, err)
				level, ok := mp[slog.LevelKey].(string)
				require.True(t, ok, "level key should exist in output")
				require.Equal(t, tt.abbrev, level, "level should be abbreviated")
			})
		}
	})

	t.Run("text handler", func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				buf := bytes.Buffer{}
				opt := slog.HandlerOptions{Level: slog.LevelDebug}
				th := newTextHandler(&buf, &opt, time.RFC3339Nano)
				l := slog.New(th)
				tt.logFn(l)

				output := buf.String()
				require.Contains(t, output, tt.abbrev, "output should contain abbreviated level")
			})
		}
	})

	t.Run("fatal level with slogWrapper", func(t *testing.T) {
		t.Run("JSON handler", func(t *testing.T) {
			buf := bytes.Buffer{}
			opt := slog.HandlerOptions{Level: slog.LevelDebug}
			jh := newJsonHandler(&buf, &opt, time.RFC3339Nano)
			wrapper := newSlogWrapper(jh)
			wrapper.Fatalw("fatal message", "key", "value")

			mp := make(map[string]any)
			err := json.Unmarshal(buf.Bytes(), &mp)
			require.NoError(t, err)
			level, ok := mp[slog.LevelKey].(string)
			require.True(t, ok, "level key should exist in output")
			require.Equal(t, "FTL", level, "fatal level should be abbreviated")
		})

		t.Run("text handler", func(t *testing.T) {
			buf := bytes.Buffer{}
			opt := slog.HandlerOptions{Level: slog.LevelDebug}
			th := newTextHandler(&buf, &opt, time.RFC3339Nano)
			wrapper := newSlogWrapper(th)
			wrapper.Fatalw("fatal message", "key", "value")

			output := buf.String()
			require.Contains(t, output, "FTL", "output should contain abbreviated fatal level")
		})
	})
}
