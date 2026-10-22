package logger_test

import (
	"testing"

	"github.com/PavloVM7/go-logger/logger"
	"github.com/stretchr/testify/assert"
)

func TestLevel_String(t *testing.T) {
	tests := []struct {
		name  string
		level logger.Level
		want  string
	}{
		{
			name:  "debug",
			level: logger.DebugLevel,
			want:  "debug",
		},
		{
			name:  "info",
			level: logger.InfoLevel,
			want:  "info",
		},
		{
			name:  "warn",
			level: logger.WarnLevel,
			want:  "warn",
		},
		{
			name:  "error",
			level: logger.ErrorLevel,
			want:  "error",
		},
		{
			name:  "fatal",
			level: logger.FatalLevel,
			want:  "fatal",
		},
		{
			name:  "unknown",
			level: 123,
			want:  "Level(123)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.level.String()
			assert.Equal(t, tt.want, actual)
		})
	}
}

func TestLevel_CapitalString(t *testing.T) {
	tests := []struct {
		name  string
		level logger.Level
		want  string
	}{
		{
			name:  "debug",
			level: logger.DebugLevel,
			want:  "DEBUG",
		},
		{
			name:  "info",
			level: logger.InfoLevel,
			want:  "INFO",
		},
		{
			name:  "warn",
			level: logger.WarnLevel,
			want:  "WARN",
		},
		{
			name:  "error",
			level: logger.ErrorLevel,
			want:  "ERROR",
		},
		{
			name:  "fatal",
			level: logger.FatalLevel,
			want:  "FATAL",
		},
		{
			name:  "unknown",
			level: 123,
			want:  "LEVEL(123)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.level.CapitalString()
			assert.Equal(t, tt.want, actual)
		})
	}
}

func TestLevel_ShortCapitalString(t *testing.T) {
	tests := []struct {
		name  string
		level logger.Level
		want  string
	}{
		{
			name:  "debug",
			level: logger.DebugLevel,
			want:  "DEB",
		},
		{
			name:  "info",
			level: logger.InfoLevel,
			want:  "INF",
		},
		{
			name:  "warn",
			level: logger.WarnLevel,
			want:  "WAR",
		},
		{
			name:  "error",
			level: logger.ErrorLevel,
			want:  "ERR",
		},
		{
			name:  "fatal",
			level: logger.FatalLevel,
			want:  "FAT",
		},
		{
			name:  "unknown",
			level: 123,
			want:  "LEV(123)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.level.ShortCapitalString()
			assert.Equal(t, tt.want, actual)
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name string
		text string
		want logger.Level
		err  error
	}{
		{
			name: logger.DebugLevel.String(),
			text: "debug",
			want: logger.DebugLevel,
		},
		{
			name: logger.InfoLevel.String(),
			text: "info",
			want: logger.InfoLevel,
		},
		{
			name: "empty",
			text: "",
			want: logger.InfoLevel,
		},
		{
			name: logger.WarnLevel.String(),
			text: "warn",
			want: logger.WarnLevel,
		},
		{
			name: "warning",
			text: "warning",
			want: logger.WarnLevel,
		},
		{
			name: logger.ErrorLevel.String(),
			text: "error",
			want: logger.ErrorLevel,
		},
		{
			name: logger.FatalLevel.String(),
			text: "fatal",
			want: logger.FatalLevel,
		},
		{
			name: "debugger",
			text: "debugger",
			want: logger.InfoLevel,
			err:  logger.ErrUnknownLevel,
		},
		{
			name: "err",
			text: "err",
			want: logger.InfoLevel,
			err:  logger.ErrUnknownLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := logger.ParseLevel(tt.text)
			assert.Equal(t, tt.want, actual)
			if tt.err == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.ErrorIs(t, err, logger.ErrUnknownLevel)
			}
		})
	}
}
