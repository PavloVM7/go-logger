package slog

import (
	"bytes"
	"encoding/json/v2"
	"os"
	"path"
	"slices"
	"strings"
	"testing"

	"github.com/PavloVM7/go-logger/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnabled(t *testing.T) {
	type loggerBoolFunc func(*slogWrapper) bool
	tests := []struct {
		level  internal.LevelKind
		trues  []loggerBoolFunc
		falses []loggerBoolFunc
	}{
		{
			level: internal.LevelKindDebug,
			trues: []loggerBoolFunc{
				func(l *slogWrapper) bool { return l.IsLoggingLevelDebug() },
				func(l *slogWrapper) bool { return l.IsLoggingLevelInfo() },
				func(l *slogWrapper) bool { return l.IsLoggingLevelWarn() },
				func(l *slogWrapper) bool { return l.IsLoggingLevelError() },
			},
			falses: nil,
		},
		{
			level: internal.LevelKindInfo,
			trues: []loggerBoolFunc{
				func(l *slogWrapper) bool { return l.IsLoggingLevelInfo() },
				func(l *slogWrapper) bool { return l.IsLoggingLevelWarn() },
				func(l *slogWrapper) bool { return l.IsLoggingLevelError() },
			},
			falses: []loggerBoolFunc{
				func(l *slogWrapper) bool { return l.IsLoggingLevelDebug() },
			},
		},
		{
			level: internal.LevelKindWarn,
			trues: []loggerBoolFunc{
				func(l *slogWrapper) bool { return l.IsLoggingLevelWarn() },
				func(l *slogWrapper) bool { return l.IsLoggingLevelError() },
			},
			falses: []loggerBoolFunc{
				func(l *slogWrapper) bool { return l.IsLoggingLevelDebug() },
				func(l *slogWrapper) bool { return l.IsLoggingLevelInfo() },
			},
		},
		{
			level: internal.LevelKindError,
			trues: []loggerBoolFunc{
				func(l *slogWrapper) bool { return l.IsLoggingLevelError() },
			},
			falses: []loggerBoolFunc{
				func(l *slogWrapper) bool { return l.IsLoggingLevelDebug() },
				func(l *slogWrapper) bool { return l.IsLoggingLevelInfo() },
				func(l *slogWrapper) bool { return l.IsLoggingLevelWarn() },
			},
		},
		{
			level: internal.LevelKindFatal,
			trues: nil,
			falses: []loggerBoolFunc{
				func(l *slogWrapper) bool { return l.IsLoggingLevelDebug() },
				func(l *slogWrapper) bool { return l.IsLoggingLevelInfo() },
				func(l *slogWrapper) bool { return l.IsLoggingLevelWarn() },
				func(l *slogWrapper) bool { return l.IsLoggingLevelError() },
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.level.String(), func(t *testing.T) {
			cfg := Config{
				Name:       "test-logger",
				Level:      tt.level.String(),
				DateFormat: "2006-01-02T15:04:05.000Z-07:00",
				Encoding:   internal.JSONEncoding,
				Outputs:    []string{internal.OutputStdout},
				w:          nil,
			}
			l, err := NewWrapper(cfg)
			require.NoError(t, err)
			for i, f := range tt.trues {
				assert.True(t, f(l), "%#q trues i=%d", tt.level.String(), i)
			}
			for i, f := range tt.falses {
				assert.False(t, f(l), "%#q falses i=%d", tt.level.String(), i)
			}
		})
	}
}

func TestWith(t *testing.T) {
	cfg := Config{
		Name:       "test-logger",
		Level:      "info",
		DateFormat: "2006-01-02T15:04:05.000Z-07:00",
		Encoding:   internal.JSONEncoding,
		Outputs:    []string{internal.OutputStdout},
		w:          nil,
	}

	t.Run("success", func(t *testing.T) {
		const (
			key1           = "key1"
			value1         = "value 1"
			key2           = "key2"
			value2         = 123
			key3           = "key3"
			value3         = "value 3"
			message1       = "test message 1"
			message2       = "test message 2"
			messageKeyName = "msg"
		)

		buf := bytes.Buffer{}
		cfg.w = &buf

		l1, err := NewWrapper(cfg)
		require.NoError(t, err)
		l2 := l1.With(key1, value1, key2, value2)

		l1.Infow(message1, key3, value3)
		l2.Errorw(message2, key3, value3)

		strs := parseBuf(t, &buf)
		require.Len(t, strs, 2)

		map1 := make(map[string]any)
		map2 := make(map[string]any)

		require.NoError(t, json.Unmarshal([]byte(strs[0]), &map1))
		require.NoError(t, json.Unmarshal([]byte(strs[1]), &map2))

		assert.Equal(t, message1, map1[messageKeyName])
		assert.EqualValues(t, value3, map1[key3])
		assert.Empty(t, map1[key1])
		assert.Empty(t, map1[key2])

		assert.Equal(t, message2, map2[messageKeyName])
		assert.EqualValues(t, value1, map2[key1])
		assert.EqualValues(t, value2, map2[key2])
		assert.EqualValues(t, value3, map2[key3])
	})

	t.Run("close file", func(t *testing.T) {
		const (
			key1           = "key1"
			value1         = "value 1"
			key2           = "key2"
			value2         = 123
			key3           = "key3"
			value3         = "value 3"
			message1       = "test message 1"
			message2       = "test message 2"
			messageKeyName = "msg"
		)

		fileName := getFullOutputFileName("slog-file.log")

		defer func() { _ = os.Remove(fileName) }()

		cfg.w = nil
		cfg.Outputs = []string{fileName}

		l1, err := NewWrapper(cfg)
		require.NoError(t, err)
		l2 := l1.With(key1, value1, key2, value2)

		l1.Infow(message1, key3, value3)
		l2.Infow(message2, key3, value3)

		require.NoError(t, l2.Close())
		require.NoError(t, l1.Close())

		strs := readLogFile(t, fileName)
		require.Len(t, strs, 2)

		map1 := make(map[string]any)
		map2 := make(map[string]any)

		require.NoError(t, json.Unmarshal([]byte(strs[0]), &map1))
		require.NoError(t, json.Unmarshal([]byte(strs[1]), &map2))

		assert.Equal(t, message1, map1[messageKeyName])
		assert.EqualValues(t, value3, map1[key3])
		assert.Empty(t, map1[key1])
		assert.Empty(t, map1[key2])

		assert.Equal(t, message2, map2[messageKeyName])
		assert.EqualValues(t, value1, map2[key1])
		assert.EqualValues(t, value2, map2[key2])
		assert.EqualValues(t, value3, map2[key3])
	})
}

func TestNewWrapperJSON(t *testing.T) {
	const (
		dateFormat = "2006-01-02T15:04:05.000Z-07:00"
		message    = "test log message"
	)

	t.Run("success", func(t *testing.T) {
		buf := bytes.Buffer{}
		cfg := Config{
			Name:       "test-logger",
			Level:      "info",
			DateFormat: dateFormat,
			Encoding:   internal.JSONEncoding,
			Outputs:    []string{internal.OutputStdout},
			w:          &buf,
		}
		l, err := NewWrapper(cfg)
		require.NoError(t, err)
		l.Debugw(message)
		assert.Empty(t, buf.String())
		l.Infow(message)
		assert.NotEmpty(t, buf.String())
	})

	t.Run("unknown log level", func(t *testing.T) {
		buf := bytes.Buffer{}
		cfg := Config{
			Name:       "test-logger",
			Level:      "unknown",
			DateFormat: dateFormat,
			Encoding:   internal.JSONEncoding,
			Outputs:    []string{internal.OutputStdout},
			w:          &buf,
		}
		_, err := NewWrapper(cfg)
		require.Error(t, err)
	})
}

func TestNewWrapperText(t *testing.T) {
	const (
		dateFormat = "2006-01-02T15:04:05.000Z-07:00"
		message    = "test log message"
	)
	t.Run("success", func(t *testing.T) {
		buf := bytes.Buffer{}
		cfg := Config{
			Name:       "test-logger",
			Level:      "info",
			DateFormat: dateFormat,
			Encoding:   internal.ConsoleEncoding,
			w:          &buf,
		}
		l, err := NewWrapper(cfg)
		require.NoError(t, err)
		l.Debugw(message)
		assert.Empty(t, buf.String())
		l.Infow(message)
		assert.NotEmpty(t, buf.String())
	})

	t.Run("unknown log level", func(t *testing.T) {
		buf := bytes.Buffer{}
		cfg := Config{
			Name:       "test-logger",
			Level:      "unknown",
			DateFormat: dateFormat,
			w:          &buf,
		}
		_, err := NewWrapper(cfg)
		require.Error(t, err)
	})
}

func parseBuf(t *testing.T, buf *bytes.Buffer) []string {
	t.Helper()

	ar := strings.Split(buf.String(), "\n")
	result := make([]string, 0, len(ar))
	for _, s := range ar {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		result = append(result, s)
	}
	return slices.Clip(result)
}

func readLogFile(t *testing.T, fileName string) []string {
	t.Helper()

	f := os.File{}
	f.Close()
	bts, err := os.ReadFile(fileName)
	require.NoError(t, err)

	res := strings.Split(string(bts), "\n")
	if len(res) > 0 && res[len(res)-1] == "" {
		res = res[:len(res)-1]
	}
	return res
}

func getFullOutputFileName(fileName string) string {
	return path.Join("../../tests/data/tmp", fileName)
}
