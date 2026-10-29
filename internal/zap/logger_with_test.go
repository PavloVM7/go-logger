package zap

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/PavloVM7/go-logger/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

	tests := []struct {
		name   string
		abbrev string
		logFn  func(*zapLogger)
	}{
		{
			name:   "debugw",
			abbrev: "DBG",
			logFn:  func(sw *zapLogger) { sw.Debugw(message, key1, val1, key2, val2, key3, val3, key4, val4) },
		},
		{
			name:   "infow",
			abbrev: "INF",
			logFn:  func(sw *zapLogger) { sw.Infow(message, key1, val1, key2, val2, key3, val3, key4, val4) },
		},
		{
			name:   "warnw",
			abbrev: "WRN",
			logFn:  func(sw *zapLogger) { sw.Warnw(message, key1, val1, key2, val2, key3, val3, key4, val4) },
		},
		{
			name:   "errorw",
			abbrev: "ERR",
			logFn:  func(sw *zapLogger) { sw.Errorw(message, key1, val1, key2, val2, key3, val3, key4, val4) },
		},
		{
			name:   "fatalw",
			abbrev: "FTL",
			logFn:  func(sw *zapLogger) { sw.Fatalw(message, key1, val1, key2, val2, key3, val3, key4, val4) },
		},
	}

	cfg := Config{
		Name:       "test-logger",
		DateFormat: dateFormat,
		Level:      "debug",
	}

	t.Run("JSON", func(t *testing.T) {
		cfg.Encoding = internal.JSONEncoding
		conf := newConfig(cfg)

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				buf := bytes.Buffer{}
				l := newLogger(t, cfg, &buf)
				tt.logFn(l)
				_ = l.Close()

				mp := make(map[string]any)
				err := json.Unmarshal(buf.Bytes(), &mp)
				require.NoError(t, err)

				level, ok := mp[conf.EncoderConfig.LevelKey].(string)
				require.True(t, ok, "level key should exist in output")
				require.Equal(t, tt.abbrev, level, "level should be abbreviated")

				sdt, ok := mp[conf.EncoderConfig.TimeKey].(string)
				require.True(t, ok)
				_, err = time.Parse(dateFormat, sdt)
				require.NoError(t, err)
				_, err = time.Parse(time.DateOnly, sdt)
				require.Error(t, err)

				msg := mp[conf.EncoderConfig.MessageKey]
				assert.Equal(t, message, msg)

				assert.Equal(t, val1, mp[key1])

				f64, ok := mp[key2].(float64)
				require.True(t, ok)
				assert.Equal(t, val2, int(f64))

				assert.Equal(t, val3, mp[key3])

				ds, ok := mp[key4].(string)
				require.True(t, ok)
				dur, err := time.ParseDuration(ds)
				require.NoError(t, err)
				assert.Equal(t, val4, dur)
			})
		}
	})

	t.Run("console", func(t *testing.T) {
		cfg.Encoding = internal.ConsoleEncoding
		// conf := newConfig(cfg)

		for _, tt := range tests {
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
			require.Contains(t, output, message, "output should contain log message")

			sub = getJson(t, output)
			mp := jsonToMap(t, []byte(sub))

			assert.Equal(t, val1, mp[key1])

			f64, ok := mp[key2].(float64)
			require.True(t, ok)
			assert.Equal(t, val2, int(f64))

			assert.Equal(t, val3, mp[key3])

			ds, ok := mp[key4].(string)
			require.True(t, ok)
			dur, err := time.ParseDuration(ds)
			require.NoError(t, err)
			assert.Equal(t, val4, dur)
		}
	})
}

func jsonToMap(t *testing.T, bts []byte) map[string]any {
	t.Helper()

	result := make(map[string]any)
	err := json.Unmarshal(bts, &result)
	require.NoError(t, err, "src: %q", string(bts))
	return result
}

func getJson(t *testing.T, src string) string {
	t.Helper()

	b := strings.IndexByte(src, '{')
	e := strings.IndexByte(src, '}')
	return src[b : e+1]
}

func newLogger(t *testing.T, cfg Config, buf *bytes.Buffer) *zapLogger {
	cfg.check()
	conf := newConfig(cfg)

	writeSyncer := zapcore.AddSync(buf)

	encoder := newEncoder(t, cfg.Encoding, conf.EncoderConfig)
	zc := zapcore.NewCore(encoder, writeSyncer, conf.Level)
	opt := zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(c, zc)
	})
	hook := zap.WithFatalHook(&fatalHook{})

	result, err := newZapLogger(cfg, opt, hook)
	if err != nil {
		t.Fatalf("failed to create zapLogger: %q", err)
	}
	return result
}

func newEncoder(t *testing.T, encoding internal.Encoding, encConf zapcore.EncoderConfig) zapcore.Encoder {
	t.Helper()

	switch encoding {
	case internal.JSONEncoding:
		return zapcore.NewJSONEncoder(encConf)
	case internal.ConsoleEncoding:
		return zapcore.NewConsoleEncoder(encConf)
	default:
		t.Fatalf("unknown encoding %q", encoding)
		return nil
	}
}

type fatalHook struct {
}

func (a fatalHook) OnWrite(ce *zapcore.CheckedEntry, _ []zapcore.Field) {}
