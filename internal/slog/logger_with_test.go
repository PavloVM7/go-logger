package slog

import (
	"bytes"
	"encoding/json/v2"
	logslog "log/slog"
	"testing"
	"time"

	"github.com/PavloVM7/go-logger/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Logger_With(t *testing.T) {
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

	optFn := func() logslog.HandlerOptions {
		return logslog.HandlerOptions{Level: logslog.LevelDebug, AddSource: false}
	}

	tests := []struct {
		name   string
		abbrev string
		logFn  func(internal.Logger)
	}{
		{
			name:   "debugw",
			abbrev: "DBG",
			logFn:  func(sw internal.Logger) { sw.Debugw(message, key1, val1, key2, val2, key3, val3, key4, val4) },
		},
		{
			name:   "infow",
			abbrev: "INF",
			logFn:  func(sw internal.Logger) { sw.Infow(message, key1, val1, key2, val2, key3, val3, key4, val4) },
		},
		{
			name:   "warnw",
			abbrev: "WRN",
			logFn:  func(sw internal.Logger) { sw.Warnw(message, key1, val1, key2, val2, key3, val3, key4, val4) },
		},
		{
			name:   "errorw",
			abbrev: "ERR",
			logFn:  func(sw internal.Logger) { sw.Errorw(message, key1, val1, key2, val2, key3, val3, key4, val4) },
		},
		{
			name:   "fatalw",
			abbrev: "FTL",
			logFn:  func(sw internal.Logger) { sw.Fatalw(message, key1, val1, key2, val2, key3, val3, key4, val4) },
		},
	}

	t.Run("JSON handler", func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				oldExit := osExit
				osExit = func(code int) {}
				defer func() {
					osExit = oldExit
				}()
				buf := bytes.Buffer{}
				opt := optFn()
				h := newJSONHandler(&buf, &opt, dateFormat)
				l, err := newSlogWrapper(h, nil)
				require.NoError(t, err)
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
			})
		}
	})

	t.Run("text handler", func(t *testing.T) {
		for _, tt := range tests {
			oldExit := osExit
			osExit = func(code int) {}
			defer func() {
				osExit = oldExit
			}()

			buf := bytes.Buffer{}
			opt := optFn()
			h := newTextHandler(&buf, &opt, dateFormat)
			l, err := newSlogWrapper(h, nil)
			require.NoError(t, err)
			tt.logFn(l)

			output := buf.String()

			sub := output[0:len(dateFormat)]
			_, err = time.Parse(dateFormat, sub)
			require.NoError(t, err)
			_, err = time.Parse(time.DateOnly, sub)
			require.Error(t, err)

			require.Contains(t, output, tt.abbrev, "output should contain abbreviated level")
			require.Contains(t, output, message, "output should contain log message")

			require.Contains(t, output, "key1=\"value 1\"")
			require.Contains(t, output, "key2=123")
			require.Contains(t, output, "key3=true")
			require.Contains(t, output, "key4=7h0m0s")
		}

	})
}
