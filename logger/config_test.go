package logger_test

import (
	"testing"

	"github.com/PavloVM7/go-logger/logger"
	"github.com/stretchr/testify/assert"
)

func TestConfig_GetLevel(t *testing.T) {
	tests := []struct {
		cfg  logger.Config
		want logger.Level
		err  error
	}{
		{
			cfg:  logger.Config{Name: "debug", Level: "debug"},
			want: logger.DebugLevel,
		},
		{
			cfg:  logger.Config{Name: "info", Level: "info"},
			want: logger.InfoLevel,
		},
		{
			cfg:  logger.Config{Name: "empty", Level: ""},
			want: logger.InfoLevel,
		},
		{
			cfg:  logger.Config{Name: "warn", Level: "warn"},
			want: logger.WarnLevel,
		},
		{
			cfg:  logger.Config{Name: "warning", Level: "warning"},
			want: logger.WarnLevel,
		},
		{
			cfg:  logger.Config{Name: "error", Level: "error"},
			want: logger.ErrorLevel,
		},
		{
			cfg:  logger.Config{Name: "fatal", Level: "fatal"},
			want: logger.FatalLevel,
		},
		{
			cfg:  logger.Config{Name: "debugger", Level: "debugger"},
			want: logger.InfoLevel,
			err:  logger.ErrUnknownLevel,
		},
		{
			cfg:  logger.Config{Name: "err", Level: "err"},
			want: logger.InfoLevel,
			err:  logger.ErrUnknownLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.cfg.Name, func(t *testing.T) {
			actual, err := tt.cfg.GetLevel()
			assert.Equal(t, tt.want, actual)
			if tt.err == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, logger.ErrUnknownLevel)
			}
		})
	}
}
