package zap

import (
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/PavloVM7/go-logger/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func getOutputFileName() string {
	fn := time.Now().Format("2006-01-02T15-04-05-000") + ".log"
	return path.Join("../../tests/data/tmp", fn)
}
