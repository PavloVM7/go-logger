package logger

import (
	"bytes"
	"encoding/json"
	logslog "log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSlogLogger_With_JSON(t *testing.T) {
	const (
		message    = "log message"
		key1       = "key1"
		val1       = "value 1"
		key2       = "key2"
		val2       = 123
		key3       = "key3"
		val3       = true
		key4       = "key4"
		val4       = 7 * time.Hour
		dateFormat = "2006-01-02T15:04:05.000Z-07:00"
	)

	tests := []struct {
		name   string
		abbrev string
		logFn  func(Logger)
	}{
		{
			name:   "debugw",
			abbrev: "DBG",
			logFn:  func(l Logger) { l.Debugw(message, key1, val1, key2, val2, key3, val3, key4, val4) },
		},
		{
			name:   "infow",
			abbrev: "INF",
			logFn:  func(l Logger) { l.Infow(message, key1, val1, key2, val2, key3, val3, key4, val4) },
		},
		{
			name:   "warnw",
			abbrev: "WRN",
			logFn:  func(l Logger) { l.Warnw(message, key1, val1, key2, val2, key3, val3, key4, val4) },
		},
		{
			name:   "errorw",
			abbrev: "ERR",
			logFn:  func(l Logger) { l.Errorw(message, key1, val1, key2, val2, key3, val3, key4, val4) },
		},
	}

	config := Config{
		Name:       "test-logger",
		Level:      "debug",
		DateFormat: dateFormat,
	}
	t.Run("JSON handler", func(t *testing.T) {
		for _, tt := range tests {
			buf := bytes.Buffer{}
			cfg, err := newSlogConfig(config, &buf)
			require.NoError(t, err)
			l := newWrapperJSON(cfg)
			tt.logFn(l)

			mp := make(map[string]any)
			err = json.Unmarshal(buf.Bytes(), &mp)
			require.NoError(t, err)

			level, ok := mp[logslog.LevelKey].(string)
			require.True(t, ok, "level key should exist in output")
			require.Equal(t, tt.abbrev, level, "level should be abbreviated")

			sdt, ok := mp[logslog.TimeKey].(string)
			require.True(t, ok)
			_, err = time.Parse(dateFormat, sdt)
			require.NoError(t, err)
			_, err = time.Parse(time.DateOnly, sdt)
			require.Error(t, err)

			msg := mp[logslog.MessageKey]
			assert.Equal(t, message, msg)

			assert.Equal(t, val1, mp[key1])

			f64, ok := mp[key2].(float64)
			require.True(t, ok)
			assert.Equal(t, val2, int(f64))

			assert.Equal(t, val3, mp[key3])

			f64, ok = mp[key4].(float64)
			require.True(t, ok)
			assert.Equal(t, val4, time.Duration(f64))
		}
	})
}
