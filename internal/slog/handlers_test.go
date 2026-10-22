package internal_slog

import (
	"bytes"
	"encoding/json/v2"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

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
