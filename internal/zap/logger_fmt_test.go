package zap

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/PavloVM7/go-logger/internal"
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

	cfg := Config{
		Name:       "test-logger",
		DateFormat: dateFormat,
		Level:      "debug",
	}

	tests := []struct {
		name   string
		abbrev string
		logFn  func(*zapLogger)
	}{
		{
			name:   "debugf",
			abbrev: "DBG",
			logFn:  func(sw *zapLogger) { sw.Debugf(formatString, key, value) },
		},
		{
			name:   "infof",
			abbrev: "INF",
			logFn:  func(sw *zapLogger) { sw.Infof(formatString, key, value) },
		},
		{
			name:   "warnf",
			abbrev: "WRN",
			logFn:  func(sw *zapLogger) { sw.Warnf(formatString, key, value) },
		},
		{
			name:   "errorf",
			abbrev: "ERR",
			logFn:  func(sw *zapLogger) { sw.Errorf(formatString, key, value) },
		},
		{
			name:   "fatalf",
			abbrev: "FTL",
			logFn:  func(sw *zapLogger) { sw.Fatalf(formatString, key, value) },
		},
	}

	t.Run("JSON", func(t *testing.T) {
		cfg.Encoding = internal.JSONEncoding
		conf, err := newConfig(cfg)
		require.NoError(t, err)

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				buf := bytes.Buffer{}
				l := newLogger(t, cfg, &buf)
				tt.logFn(l)
				_ = l.Close()

				mp := make(map[string]any)
				err := json.Unmarshal(buf.Bytes(), &mp)
				require.NoError(t, err)
				assert.Len(t, mp, 4)

				level, ok := mp[conf.EncoderConfig.LevelKey].(string)
				require.True(t, ok, "level key should exist in output")
				require.Equal(t, tt.abbrev, level, "level should be abbreviated")

				sdt, ok := mp[conf.EncoderConfig.TimeKey].(string)
				require.True(t, ok)
				_, err = time.Parse(dateFormat, sdt)
				require.NoError(t, err)

				msg := mp[conf.EncoderConfig.MessageKey]
				assert.Equal(t, want, msg)

				assert.EqualValues(t, cfg.Name, mp[conf.EncoderConfig.NameKey])
			})
		}
	})

	t.Run("console", func(t *testing.T) {
		cfg.Encoding = internal.ConsoleEncoding
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				buf := bytes.Buffer{}
				l := newLogger(t, cfg, &buf)
				tt.logFn(l)
				_ = l.Close()

				sub := buf.String()[0:len(dateFormat)]
				_, err := time.Parse(dateFormat, sub)
				require.NoError(t, err)
				_, err = time.Parse(time.DateOnly, sub)
				require.Error(t, err)

				output := buf.String()
				require.Contains(t, output, tt.abbrev, "output should contain abbreviated level")
				require.Contains(t, output, want, "output should contain log message")
			})
		}
	})
}
