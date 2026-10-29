package slog

import (
	"bytes"
	"testing"

	"github.com/PavloVM7/go-logger/internal"
	"github.com/stretchr/testify/assert"
)

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
			W:          &buf,
		}
		l := NewWrapper(cfg)
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
			W:          &buf,
		}
		l := NewWrapper(cfg)
		l.Debugw(message)
		assert.NotEmpty(t, buf.String())
		buf.Reset()
		assert.Empty(t, buf.String())
		l.Infow(message)
		assert.NotEmpty(t, buf.String())
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
			W:          &buf,
		}
		l := NewWrapperText(cfg)
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
			W:          &buf,
		}
		l := NewWrapperText(cfg)
		l.Debugw(message)
		assert.NotEmpty(t, buf.String())
		buf.Reset()
		assert.Empty(t, buf.String())
		l.Infow(message)
		assert.NotEmpty(t, buf.String())
	})
}
