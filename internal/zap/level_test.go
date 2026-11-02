package zap

import (
	"testing"

	"github.com/PavloVM7/go-logger/internal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func Test_levelToShortString(t *testing.T) {
	tests := []struct {
		name  string
		level zapcore.Level
		want  string
	}{
		{
			name:  "debug",
			level: zapcore.DebugLevel,
			want:  internal.DebugShort.String(),
		},
		{
			name:  "info",
			level: zapcore.InfoLevel,
			want:  internal.InfoShort.String(),
		},
		{
			name:  "warn",
			level: zapcore.WarnLevel,
			want:  internal.WarnShort.String(),
		},
		{
			name:  "error",
			level: zapcore.ErrorLevel,
			want:  internal.ErrorShort.String(),
		},
		{
			name:  "fatal",
			level: zapcore.FatalLevel,
			want:  internal.FatalShort.String(),
		},
		{
			name:  "unknown defaults to debug",
			level: zapcore.Level(99),
			want:  "LVL99",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := levelToShortString(tt.level)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_kindLevelToZapLevel(t *testing.T) {
	tests := []struct {
		name string
		kind internal.LevelKind
		want zapcore.Level
	}{
		{
			name: "debug",
			kind: internal.LevelKindDebug,
			want: zapcore.DebugLevel,
		},
		{
			name: "info",
			kind: internal.LevelKindInfo,
			want: zapcore.InfoLevel,
		},
		{
			name: "warn",
			kind: internal.LevelKindWarn,
			want: zapcore.WarnLevel,
		},
		{
			name: "error",
			kind: internal.LevelKindError,
			want: zapcore.ErrorLevel,
		},
		{
			name: "fatal",
			kind: internal.LevelKindFatal,
			want: zapcore.FatalLevel,
		},
		{
			name: "unknown defaults to debug",
			kind: internal.LevelKind(99),
			want: zapcore.DebugLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := kindLevelToZapLevel(tt.kind)
			assert.Equal(t, tt.want, got)
		})
	}
}
