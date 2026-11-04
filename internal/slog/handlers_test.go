package slog

import (
	"bytes"
	"encoding/json/v2"
	"io"
	logslog "log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDateFormat(t *testing.T) {
	const dateFormat = "2006-01-02 15:04:05.000Z-07:00"
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
		l.Info("debug message", "key1", 123, "key2", "value 2")
		t.Log("buf:", buf.String())
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

	logOpt := option{exitFunc: func(code int) {}}

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
			wrapper, err := newSlogWrapper(jh, nil, &logOpt)
			require.NoError(t, err)

			wrapper.Fatalw("fatal message", "key", "value")

			mp := make(map[string]any)
			err = json.Unmarshal(buf.Bytes(), &mp)
			require.NoError(t, err)
			level, ok := mp[logslog.LevelKey].(string)
			require.True(t, ok, "level key should exist in output")
			require.Equal(t, "FTL", level, "fatal level should be abbreviated")
		})

		t.Run("text handler", func(t *testing.T) {
			buf := bytes.Buffer{}
			opt := logslog.HandlerOptions{Level: logslog.LevelDebug}
			th := newTextHandler(&buf, &opt, time.RFC3339Nano)
			wrapper, err := newSlogWrapper(th, nil, &logOpt)
			require.NoError(t, err)

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

func Test_namedHandler(t *testing.T) {
	const message = "test log message"
	t.Run("JSON handler", func(t *testing.T) {
		buf := bytes.Buffer{}
		opt := logslog.HandlerOptions{Level: logslog.LevelDebug}
		h := newNamedJSONHandler(&buf, &opt, time.RFC3339Nano)
		l1, err := newSlogWrapper(h, []io.Writer{&buf}, nil)
		require.NoError(t, err)
		l1.Debugw(message)

		l2 := l1.Named("name2")
		l2.Infow(message)

		l3 := l2.Named("name3")
		l3.Warnw(message)

		l4 := l1.Named("name4")
		l4.Errorw(message)

		require.NoError(t, l1.Close())
		require.NoError(t, l2.Close())
		require.NoError(t, l3.Close())
		require.NoError(t, l4.Close())

		strs := parseBuf(t, &buf)
		require.Len(t, strs, 4)

		mp1 := make(map[string]any)
		mp2 := make(map[string]any)
		mp3 := make(map[string]any)
		mp4 := make(map[string]any)
		require.NoError(t, json.Unmarshal([]byte(strs[0]), &mp1))
		require.NoError(t, json.Unmarshal([]byte(strs[1]), &mp2))
		require.NoError(t, json.Unmarshal([]byte(strs[2]), &mp3))
		require.NoError(t, json.Unmarshal([]byte(strs[3]), &mp4))
		actual1, _ := mp1[nameKeyName].(string)
		actual2, _ := mp2[nameKeyName].(string)
		actual3, _ := mp3[nameKeyName].(string)
		actual4, _ := mp4[nameKeyName].(string)
		assert.Equal(t, "", actual1)
		assert.Equal(t, "name2", actual2)
		assert.Equal(t, "name2.name3", actual3)
		assert.Equal(t, "name4", actual4)
	})

	t.Run("text handler", func(t *testing.T) {
		buf := bytes.Buffer{}
		opt := logslog.HandlerOptions{Level: logslog.LevelDebug}
		h := newTextHandler(&buf, &opt, time.RFC3339Nano)
		l1, err := newSlogWrapper(h, []io.Writer{&buf}, nil)
		require.NoError(t, err)
		l1.Debugw(message)

		l2 := l1.Named("name2")
		l2.Infow(message)

		l3 := l2.Named("name3")
		l3.Warnw(message)

		l4 := l1.Named("name4")
		l4.Errorw(message)

		require.NoError(t, l1.Close())
		require.NoError(t, l2.Close())
		require.NoError(t, l3.Close())
		require.NoError(t, l4.Close())

		t.Log("buf:", buf.String())

		strs := parseBuf(t, &buf)
		require.Len(t, strs, 4)

		assert.NotContains(t, strs[0], nameKeyName+"=")
		assert.Contains(t, strs[1], nameKeyName+"=name2")
		assert.Contains(t, strs[2], nameKeyName+"=name2.name3")
		assert.Contains(t, strs[3], nameKeyName+"=name4")
	})
}
