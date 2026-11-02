package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	logslog "log/slog"
	"os"
	"path"
	"strings"
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
		Outputs:    []string{"stdout", getFullOutputFileName("slog-logger.log")},
	}

	t.Run("JSON handler", func(t *testing.T) {
		config.Encoding = "json"
		_ = os.Remove(config.Outputs[1])

		defer func() {
			_ = os.Remove(config.Outputs[1])
		}()

		for _, tt := range tests {
			cfg, err := newSlogConfig(config)
			require.NoError(t, err)
			l, err := newWrapper(cfg)
			require.NoError(t, err)

			tt.logFn(l)

			require.NoError(t, l.Close())

			output := readLogFile(t, config.Outputs[1])
			require.NoError(t, os.Remove(config.Outputs[1]))
			require.Len(t, output, 1)

			mp := make(map[string]any)
			err = json.Unmarshal([]byte(output[0]), &mp)
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

func readLogFile(t *testing.T, fileName string) []string {
	t.Helper()

	bts, err := os.ReadFile(fileName)
	require.NoError(t, err)

	res := strings.Split(string(bts), "\n")
	if len(res) > 0 && res[len(res)-1] == "" {
		res = res[:len(res)-1]
	}
	return res
}

func getFullOutputFileName(fileName string) string {
	return path.Join("../tests/data/tmp", fileName)
}

func captureOutput(f func()) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, r); err != nil {
			fmt.Print("copy err:", err)
		}
		outC <- buf.String()
	}()

	f()

	_ = w.Close()
	os.Stdout = oldStdout
	return <-outC
}
