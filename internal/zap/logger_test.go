package zap

import (
	"encoding/json"
	"os"
	"path"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/PavloVM7/go-logger/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWith(t *testing.T) {
	const (
		dateFormat = "2006-01-02T15:04:05.000Z-07:00"
	)

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

		logFile := getOutputFileName()
		defer func() {
			_ = os.Remove(logFile)
		}()

		cfg := Config{
			Name:       "test-logger",
			DateFormat: dateFormat,
			Level:      "debug",
			Encoding:   internal.JSONEncoding,
			Outputs:    []string{logFile},
		}

		l1, err := NewWrapperZap(cfg)
		require.NoError(t, err)
		l2 := l1.With(key1, value1, key2, value2)

		l1.Infow(message1, key3, value3)
		l2.Errorw(message2, key3, value3)

		strs := readLogFile(t, logFile)
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

		logFile := getOutputFileName()
		defer func() {
			_ = os.Remove(logFile)
		}()

		cfg := Config{
			Name:       "test-logger",
			DateFormat: dateFormat,
			Level:      "debug",
			Encoding:   internal.JSONEncoding,
			Outputs:    []string{logFile},
		}

		l1, err := NewWrapperZap(cfg)
		require.NoError(t, err)
		l2 := l1.With(key1, value1, key2, value2)

		l1.Infow(message1, key3, value3)
		l2.Infow(message2, key3, value3)

		require.NoError(t, l2.Close())
		require.NoError(t, l1.Close())

		strs := readLogFile(t, logFile)
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

func TestNewWrapperZapJSON(t *testing.T) {
	const (
		dateFormat = "2006-01-02T15:04:05.000Z-07:00"
		message    = "test log message"
	)
	logFile := getOutputFileName()
	t.Run("success", func(t *testing.T) {
		cfg := Config{
			Name:       "test-logger",
			DateFormat: dateFormat,
			Level:      "debug",
			Encoding:   internal.JSONEncoding,
			Outputs:    []string{internal.OutputStdout, logFile},
		}
		defer func() {
			_ = os.Remove(logFile)
		}()
		l, err := NewWrapperZap(cfg)
		require.NoError(t, err)
		l.Debugw(message)
		l.Infow(message)
		l.Warnw(message)
		l.Errorw(message)
		readFileAndCheck(t, logFile, 4, message)
	})
}

func readFileAndCheck(t *testing.T, filePath string, ln int, msg string) {
	t.Helper()

	bts, err := os.ReadFile(filePath)
	require.NoError(t, err)
	ar := strings.Split(string(bts), "\n")
	if ar[len(ar)-1] == "" {
		ar = ar[:len(ar)-1]
	}
	assert.Len(t, ar, ln)
	for _, s := range ar {
		assert.Contains(t, s, msg)
	}
}

func readLogFile(t *testing.T, fileName string) []string {
	t.Helper()

	bts, err := os.ReadFile(fileName)
	require.NoError(t, err)

	res := strings.Split(string(bts), "\n")
	if len(res) > 0 && res[len(res)-1] == "" {
		res = res[:len(res)-1]
	}
	return slices.Clip(res)
}

func getOutputFileName() string {
	fn := time.Now().Format("2006-01-02T15-04-05-000") + ".log"
	return path.Join("../../tests/data/tmp", fn)
}
