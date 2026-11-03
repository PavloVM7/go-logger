package slog

import (
	"bytes"
	"encoding/json/v2"
	logslog "log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Logger_fmt(t *testing.T) {
	const (
		formatString = "test message key=%s value=%v"
		key          = "super"
		value        = 123
		dateFormat   = "2006-01-02T15:04:05.000Z-07:00"
		want         = "test message key=super value=123"
	)

	optFn := func() logslog.HandlerOptions {
		return logslog.HandlerOptions{Level: logslog.LevelDebug, AddSource: false}
	}

	tests := []struct {
		name   string
		abbrev string
		logFn  func(*slogWrapper)
	}{
		{
			name:   "debugf",
			abbrev: "DBG",
			logFn:  func(sw *slogWrapper) { sw.Debugf(formatString, key, value) },
		},
		{
			name:   "infof",
			abbrev: "INF",
			logFn:  func(sw *slogWrapper) { sw.Infof(formatString, key, value) },
		},
		{
			name:   "warnf",
			abbrev: "WRN",
			logFn:  func(sw *slogWrapper) { sw.Warnf(formatString, key, value) },
		},
		{
			name:   "errorf",
			abbrev: "ERR",
			logFn:  func(sw *slogWrapper) { sw.Errorf(formatString, key, value) },
		},
		{
			name:   "fatalf",
			abbrev: "FTL",
			logFn:  func(sw *slogWrapper) { sw.Fatalf(formatString, key, value) },
		},
	}

	t.Run("JSON", func(t *testing.T) {
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
				assert.Len(t, mp, 3)

				level, ok := mp[logslog.LevelKey].(string)
				require.True(t, ok, "level key should exist in output")
				require.Equal(t, tt.abbrev, level, "level should be abbreviated")

				sdt, ok := mp[logslog.TimeKey].(string)
				require.True(t, ok)
				_, err = time.Parse(dateFormat, sdt)
				require.NoError(t, err)

				msg := mp[logslog.MessageKey]
				assert.Equal(t, want, msg)
			})
		}
	})

	t.Run("console", func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
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

				require.Contains(t, output, " "+tt.abbrev+" ", "output should contain abbreviated level")
				require.Contains(t, output, " msg=\""+want+"\"", "output should contain log message")
			})
		}
	})
}
