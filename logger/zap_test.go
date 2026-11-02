package logger

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/PavloVM7/go-logger/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestZapLogger_With_JSON(t *testing.T) {
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

	t.Run("JSON encoding", func(t *testing.T) {
		for _, tt := range tests {
			logFile := getFullOutputFileName("zap-logger.log")
			config := Config{
				Name:       "test-logger",
				Level:      "debug",
				DateFormat: dateFormat,
				Encoding:   "json",
				Outputs:    []string{internal.OutputStdout, logFile},
			}
			_ = os.Remove(logFile)

			l, err := NewLoggerZap(config)
			require.NoError(t, err)

			tt.logFn(l)

			closeErr := l.Close()
			if closeErr != nil &&
				!strings.Contains(closeErr.Error(), "bad file descriptor") &&
				!strings.Contains(closeErr.Error(), "stdout") {
				require.NoError(t, closeErr)
			}

			output := readLogFile(t, logFile)
			if len(output) > 0 {
				require.NoError(t, os.Remove(logFile))
			}
			require.Len(t, output, 1)

			mp := make(map[string]any)
			err = json.Unmarshal([]byte(output[0]), &mp)
			require.NoError(t, err)

			level, ok := mp["level"].(string)
			require.True(t, ok, "level key should exist in output")
			require.Equal(t, tt.abbrev, level, "level should be abbreviated")

			sdt, ok := mp["ts"].(string)
			require.True(t, ok)
			_, err = time.Parse(dateFormat, sdt)
			require.NoError(t, err)
			_, err = time.Parse(time.DateOnly, sdt)
			require.Error(t, err)

			msg := mp["msg"]
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
		}
	})
}

func TestZapLogger_With_Console(t *testing.T) {
	const (
		message    = "log message"
		key1       = "key1"
		val1       = "value 1"
		key2       = "key2"
		val2       = 123
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
			logFn:  func(l Logger) { l.Debugw(message, key1, val1, key2, val2) },
		},
		{
			name:   "infow",
			abbrev: "INF",
			logFn:  func(l Logger) { l.Infow(message, key1, val1, key2, val2) },
		},
		{
			name:   "warnw",
			abbrev: "WRN",
			logFn:  func(l Logger) { l.Warnw(message, key1, val1, key2, val2) },
		},
		{
			name:   "errorw",
			abbrev: "ERR",
			logFn:  func(l Logger) { l.Errorw(message, key1, val1, key2, val2) },
		},
	}

	t.Run("console encoding", func(t *testing.T) {
		for _, tt := range tests {
			logFile := getFullOutputFileName("zap-logger.log")
			config := Config{
				Name:       "test-logger",
				Level:      "debug",
				DateFormat: dateFormat,
				Encoding:   "console",
				Outputs:    []string{internal.OutputStdout, logFile},
			}
			_ = os.Remove(logFile)

			l, err := NewLoggerZap(config)
			require.NoError(t, err)

			tt.logFn(l)

			closeErr := l.Close()
			if closeErr != nil && !strings.Contains(closeErr.Error(), "bad file descriptor") {
				require.NoError(t, closeErr)
			}

			output := readLogFile(t, logFile)
			if len(output) > 0 {
				require.NoError(t, os.Remove(logFile))
			}
			require.Len(t, output, 1)

			logLine := output[0]
			assert.Contains(t, logLine, tt.abbrev)
			assert.Contains(t, logLine, message)
			assert.Contains(t, logLine, key1)
			assert.Contains(t, logLine, val1)
			assert.Contains(t, logLine, key2)
			assert.Contains(t, logLine, "123")
		}
	})
}

func TestZapLogger_LevelFiltering(t *testing.T) {
	const (
		message = "test message"
		key     = "key"
		val     = "value"
	)

	tests := []struct {
		name      string
		level     string
		logFn     func(Logger)
		shouldLog bool
	}{
		{
			name:      "debug level with debug log",
			level:     "debug",
			logFn:     func(l Logger) { l.Debugw(message, key, val) },
			shouldLog: true,
		},
		{
			name:      "info level with debug log",
			level:     "info",
			logFn:     func(l Logger) { l.Debugw(message, key, val) },
			shouldLog: false,
		},
		{
			name:      "info level with info log",
			level:     "info",
			logFn:     func(l Logger) { l.Infow(message, key, val) },
			shouldLog: true,
		},
		{
			name:      "warn level with info log",
			level:     "warn",
			logFn:     func(l Logger) { l.Infow(message, key, val) },
			shouldLog: false,
		},
		{
			name:      "warn level with warn log",
			level:     "warn",
			logFn:     func(l Logger) { l.Warnw(message, key, val) },
			shouldLog: true,
		},
		{
			name:      "error level with warn log",
			level:     "error",
			logFn:     func(l Logger) { l.Warnw(message, key, val) },
			shouldLog: false,
		},
		{
			name:      "error level with error log",
			level:     "error",
			logFn:     func(l Logger) { l.Errorw(message, key, val) },
			shouldLog: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logFile := getFullOutputFileName("zap-logger.log")
			config := Config{
				Name:       "test-logger",
				Level:      tt.level,
				Encoding:   "json",
				DateFormat: "2006-01-02T15:04:05.000Z-07:00",
				Outputs:    []string{internal.OutputStdout, logFile},
			}
			_ = os.Remove(logFile)

			l, err := NewLoggerZap(config)
			require.NoError(t, err)

			tt.logFn(l)

			closeErr := l.Close()
			if closeErr != nil && !strings.Contains(closeErr.Error(), "bad file descriptor") {
				require.NoError(t, closeErr)
			}

			output := readLogFile(t, logFile)
			if len(output) > 0 {
				require.NoError(t, os.Remove(logFile))
			}

			if tt.shouldLog {
				require.Len(t, output, 1, "expected exactly one log line")
				assert.Contains(t, output[0], message)
			} else {
				require.Len(t, output, 0, "expected no log lines")
			}
		})
	}
}

func TestZapLogger_DefaultDateFormat(t *testing.T) {
	logFile := getFullOutputFileName("zap-logger.log")
	config := Config{
		Name:     "test-logger",
		Level:    "debug",
		Encoding: "json",
		Outputs:  []string{internal.OutputStdout, logFile},
	}
	_ = os.Remove(logFile)

	l, err := NewLoggerZap(config)
	require.NoError(t, err)

	l.Infow("test message")

	closeErr := l.Close()
	if closeErr != nil && !strings.Contains(closeErr.Error(), "bad file descriptor") {
		require.NoError(t, closeErr)
	}

	output := readLogFile(t, logFile)
	if len(output) > 0 {
		require.NoError(t, os.Remove(logFile))
	}
	require.Len(t, output, 1)

	mp := make(map[string]any)
	err = json.Unmarshal([]byte(output[0]), &mp)
	require.NoError(t, err)

	sdt, ok := mp["ts"].(string)
	require.True(t, ok)
	_, err = time.Parse(defaultDateFormat, sdt)
	require.NoError(t, err)
}

func TestZapLogger_Name(t *testing.T) {
	const loggerName = "my-custom-logger"

	logFile := getFullOutputFileName("zap-logger.log")
	config := Config{
		Name:       loggerName,
		Level:      "debug",
		Encoding:   "json",
		DateFormat: "2006-01-02T15:04:05.000Z-07:00",
		Outputs:    []string{internal.OutputStdout, logFile},
	}
	_ = os.Remove(logFile)

	l, err := NewLoggerZap(config)
	require.NoError(t, err)

	l.Infow("test message")

	closeErr := l.Close()
	if closeErr != nil && !strings.Contains(closeErr.Error(), "bad file descriptor") {
		require.NoError(t, closeErr)
	}

	output := readLogFile(t, logFile)
	if len(output) > 0 {
		require.NoError(t, os.Remove(logFile))
	}
	require.Len(t, output, 1)

	mp := make(map[string]any)
	err = json.Unmarshal([]byte(output[0]), &mp)
	require.NoError(t, err)

	loggerNameField, ok := mp["logger"].(string)
	require.True(t, ok)
	assert.Equal(t, loggerName, loggerNameField)
}

func TestZapLogger_EmptyName(t *testing.T) {
	logFile := getFullOutputFileName("zap-logger.log")
	config := Config{
		Name:       "",
		Level:      "debug",
		Encoding:   "json",
		DateFormat: "2006-01-02T15:04:05.000Z-07:00",
		Outputs:    []string{internal.OutputStdout, logFile},
	}
	_ = os.Remove(logFile)

	l, err := NewLoggerZap(config)
	require.NoError(t, err)

	l.Infow("test message")

	closeErr := l.Close()
	if closeErr != nil && !strings.Contains(closeErr.Error(), "bad file descriptor") {
		require.NoError(t, closeErr)
	}

	output := readLogFile(t, logFile)
	if len(output) > 0 {
		require.NoError(t, os.Remove(logFile))
	}
	require.Len(t, output, 1)

	mp := make(map[string]any)
	err = json.Unmarshal([]byte(output[0]), &mp)
	require.NoError(t, err)

	loggerNameField, ok := mp["logger"].(string)
	require.True(t, ok)
	assert.Equal(t, "default", loggerNameField)
}

func TestZapLogger_MultipleOutputs(t *testing.T) {
	const message = "test message"

	logFile1 := getFullOutputFileName("zap-logger1.log")
	logFile2 := getFullOutputFileName("zap-logger2.log")
	config := Config{
		Name:       "test-logger",
		Level:      "debug",
		Encoding:   "json",
		DateFormat: "2006-01-02T15:04:05.000Z-07:00",
		Outputs:    []string{logFile1, logFile2},
	}

	for _, output := range config.Outputs {
		_ = os.Remove(output)
	}

	l, err := NewLoggerZap(config)
	require.NoError(t, err)

	l.Infow(message)

	closeErr := l.Close()
	if closeErr != nil && !strings.Contains(closeErr.Error(), "bad file descriptor") {
		require.NoError(t, closeErr)
	}

	for _, output := range config.Outputs {
		fileOutput := readLogFile(t, output)
		if len(fileOutput) > 0 {
			require.NoError(t, os.Remove(output))
		}
		require.Len(t, fileOutput, 1)
		assert.Contains(t, fileOutput[0], message)
	}
}

func TestZapLogger_Close(t *testing.T) {
	logFile := getFullOutputFileName("zap-logger.log")
	config := Config{
		Name:       "test-logger",
		Level:      "debug",
		Encoding:   "json",
		DateFormat: "2006-01-02T15:04:05.000Z-07:00",
		Outputs:    []string{internal.OutputStdout, logFile},
	}
	_ = os.Remove(logFile)

	l, err := NewLoggerZap(config)
	require.NoError(t, err)

	l.Infow("test message")

	closeErr := l.Close()
	if closeErr != nil && !strings.Contains(closeErr.Error(), "bad file descriptor") {
		require.NoError(t, closeErr)
	}

	closeErr = l.Close()
	if closeErr != nil && !strings.Contains(closeErr.Error(), "bad file descriptor") {
		require.NoError(t, closeErr, "closing already closed logger should not error")
	}

	output := readLogFile(t, logFile)
	if len(output) > 0 {
		require.NoError(t, os.Remove(logFile))
	}
	require.Len(t, output, 1)
}

func TestZapLogger_InvalidLevel(t *testing.T) {
	logFile := getFullOutputFileName("zap-logger.log")
	config := Config{
		Name:       "test-logger",
		Level:      "invalid-level",
		Encoding:   "json",
		DateFormat: "2006-01-02T15:04:05.000Z-07:00",
		Outputs:    []string{internal.OutputStdout, logFile},
	}
	_ = os.Remove(logFile)

	l, err := NewLoggerZap(config)
	require.NoError(t, err, "invalid level should default to debug")

	l.Infow("test message")

	closeErr := l.Close()
	if closeErr != nil && !strings.Contains(closeErr.Error(), "bad file descriptor") {
		require.NoError(t, closeErr)
	}

	output := readLogFile(t, logFile)
	if len(output) > 0 {
		require.NoError(t, os.Remove(logFile))
	}
	require.Len(t, output, 1)
}

func TestZapLogger_Fatalw(t *testing.T) {
	t.Skip()
	logFile := getFullOutputFileName("zap-logger.log")
	config := Config{
		Name:       "test-logger",
		Level:      "debug",
		Encoding:   "json",
		DateFormat: "2006-01-02T15:04:05.000Z-07:00",
		Outputs:    []string{internal.OutputStdout, logFile},
	}
	_ = os.Remove(logFile)

	l, err := NewLoggerZap(config)
	require.NoError(t, err)

	l.Fatalw("fatal message", "key", "value")

	output := readLogFile(t, logFile)
	require.NoError(t, os.Remove(logFile))
	require.Len(t, output, 1)
	assert.Contains(t, output[0], "fatal message")
}

func TestZapLogger_ContextIntegration(t *testing.T) {
	logFile := getFullOutputFileName("zap-logger.log")
	config := Config{
		Name:       "test-logger",
		Level:      "debug",
		Encoding:   "json",
		DateFormat: "2006-01-02T15:04:05.000Z-07:00",
		Outputs:    []string{internal.OutputStdout, logFile},
	}
	_ = os.Remove(logFile)

	l, err := NewLoggerZap(config)
	require.NoError(t, err)

	ctx := context.TODO()
	ctx = ToContext(ctx, l)

	retrievedLogger, err := FromContext(ctx)
	require.NoError(t, err)
	assert.NotNil(t, retrievedLogger)

	retrievedLogger.Infow("test message")

	closeErr := l.Close()
	if closeErr != nil && !strings.Contains(closeErr.Error(), "bad file descriptor") {
		require.NoError(t, closeErr)
	}

	output := readLogFile(t, logFile)
	if len(output) > 0 {
		require.NoError(t, os.Remove(logFile))
	}
	require.Len(t, output, 1)
	assert.Contains(t, output[0], "test message")
}

func TestZapLogger_ZapLevelMapping(t *testing.T) {
	const message = "test log message"
	tests := []struct {
		name            string
		level           string
		logFn           func(l Logger)
		unexpectedLogFn func(l Logger)
	}{
		{
			name:            "debug maps to debug",
			level:           "debug",
			logFn:           func(l Logger) { l.Debugw(message) },
			unexpectedLogFn: nil,
		},
		{
			name:            "info maps to info",
			level:           "info",
			logFn:           func(l Logger) { l.Infow(message) },
			unexpectedLogFn: func(l Logger) { l.Debugw(message) },
		},
		{
			name:            "warn maps to warn",
			level:           "warn",
			logFn:           func(l Logger) { l.Warnw(message) },
			unexpectedLogFn: func(l Logger) { l.Infow(message) },
		},
		{
			name:            "error maps to error",
			level:           "error",
			logFn:           func(l Logger) { l.Errorw(message) },
			unexpectedLogFn: func(l Logger) { l.Warnw(message) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logFile := getFullOutputFileName("zap-logger.log")
			config := Config{
				Name:       "test-logger",
				Level:      tt.level,
				Encoding:   "json",
				DateFormat: "2006-01-02T15:04:05.000Z-07:00",
				Outputs:    []string{internal.OutputStdout, logFile},
			}
			_ = os.Remove(logFile)

			l, err := NewLoggerZap(config)
			require.NoError(t, err)

			tt.logFn(l)
			if tt.unexpectedLogFn != nil {
				tt.unexpectedLogFn(l)
			}

			closeErr := l.Close()
			if closeErr != nil && !strings.Contains(closeErr.Error(), "bad file descriptor") {
				require.NoError(t, closeErr)
			}

			output := readLogFile(t, logFile)
			require.NoError(t, os.Remove(logFile))
			require.Len(t, output, 1)
			require.Contains(t, output[0], message)
		})
	}
}
