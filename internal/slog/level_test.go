package slog

import (
	"log/slog"
	logslog "log/slog"
	"testing"

	"github.com/PavloVM7/go-logger/internal"
	"github.com/stretchr/testify/assert"
)

func TestFromSlogLevel(t *testing.T) {
	tests := []struct {
		name  string
		level slog.Level
		want  string
	}{
		{
			name:  "debug",
			level: slog.LevelDebug,
			want:  internal.DebugShort.String(),
		},
		{
			name:  "info",
			level: slog.LevelInfo,
			want:  internal.InfoShort.String(),
		},
		{
			name:  "warn",
			level: slog.LevelWarn,
			want:  internal.WarnShort.String(),
		},
		{
			name:  "error",
			level: slog.LevelError,
			want:  internal.ErrorShort.String(),
		},
		{
			name:  "fatal custom",
			level: LevelFatal,
			want:  internal.FatalShort.String(),
		},
		{
			name:  "unknown defaults to debug",
			level: slog.Level(99),
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

func TestLevelKind_kindLevelToSlogLevel(t *testing.T) {
	tests := []struct {
		name string
		kind internal.LevelKind
		want logslog.Level
	}{
		{
			name: "debug",
			kind: internal.LevelKindDebug,
			want: logslog.LevelDebug,
		},
		{
			name: "info",
			kind: internal.LevelKindInfo,
			want: logslog.LevelInfo,
		},
		{
			name: "warn",
			kind: internal.LevelKindWarn,
			want: logslog.LevelWarn,
		},
		{
			name: "error",
			kind: internal.LevelKindError,
			want: logslog.LevelError,
		},
		{
			name: "fatal",
			kind: internal.LevelKindFatal,
			want: logslog.Level(12),
		},
		{
			name: "unknown defaults to info",
			kind: internal.LevelKind(99),
			want: logslog.LevelInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := kindLevelToSlogLevel(tt.kind)
			assert.Equal(t, tt.want, got)
		})
	}
}
