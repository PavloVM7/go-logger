package internal

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestParseLevelString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    LevelKind
		wantErr bool
	}{
		{
			name:    "debug",
			input:   "debug",
			want:    LevelKindDebug,
			wantErr: false,
		},
		{
			name:    "debug uppercase",
			input:   "DEBUG",
			want:    LevelKindDebug,
			wantErr: false,
		},
		{
			name:    "debug mixed case",
			input:   "DeBuG",
			want:    LevelKindDebug,
			wantErr: false,
		},
		{
			name:    "debug with spaces",
			input:   "  debug  ",
			want:    LevelKindDebug,
			wantErr: false,
		},
		{
			name:    "info",
			input:   "info",
			want:    LevelKindInfo,
			wantErr: false,
		},
		{
			name:    "empty string defaults to info",
			input:   "",
			want:    LevelKindInfo,
			wantErr: false,
		},
		{
			name:    "warn",
			input:   "warn",
			want:    LevelKindWarn,
			wantErr: false,
		},
		{
			name:    "warning",
			input:   "warning",
			want:    LevelKindWarn,
			wantErr: false,
		},
		{
			name:    "error",
			input:   "error",
			want:    LevelKindError,
			wantErr: false,
		},
		{
			name:    "fatal",
			input:   "fatal",
			want:    LevelKindFatal,
			wantErr: false,
		},
		{
			name:    "unknown level",
			input:   "unknown",
			want:    LevelKindInfo,
			wantErr: true,
		},
		{
			name:    "invalid level",
			input:   "invalid",
			want:    LevelKindInfo,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLevelString(tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLevelKind_ToSlogLevel(t *testing.T) {
	tests := []struct {
		name string
		kind LevelKind
		want slog.Level
	}{
		{
			name: "debug",
			kind: LevelKindDebug,
			want: slog.LevelDebug,
		},
		{
			name: "info",
			kind: LevelKindInfo,
			want: slog.LevelInfo,
		},
		{
			name: "warn",
			kind: LevelKindWarn,
			want: slog.LevelWarn,
		},
		{
			name: "error",
			kind: LevelKindError,
			want: slog.LevelError,
		},
		{
			name: "fatal",
			kind: LevelKindFatal,
			want: slog.Level(12),
		},
		{
			name: "unknown defaults to info",
			kind: LevelKind(99),
			want: slog.LevelInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.kind.ToSlogLevel()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLevelKind_ToZapLevel(t *testing.T) {
	tests := []struct {
		name string
		kind LevelKind
		want zapcore.Level
	}{
		{
			name: "debug",
			kind: LevelKindDebug,
			want: zapcore.DebugLevel,
		},
		{
			name: "info",
			kind: LevelKindInfo,
			want: zapcore.InfoLevel,
		},
		{
			name: "warn",
			kind: LevelKindWarn,
			want: zapcore.WarnLevel,
		},
		{
			name: "error",
			kind: LevelKindError,
			want: zapcore.ErrorLevel,
		},
		{
			name: "fatal",
			kind: LevelKindFatal,
			want: zapcore.FatalLevel,
		},
		{
			name: "unknown defaults to debug",
			kind: LevelKind(99),
			want: zapcore.DebugLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.kind.ToZapLevel()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLevelKind_ToShortString(t *testing.T) {
	tests := []struct {
		name string
		kind LevelKind
		want string
	}{
		{
			name: "debug",
			kind: LevelKindDebug,
			want: "DBG",
		},
		{
			name: "info",
			kind: LevelKindInfo,
			want: "INF",
		},
		{
			name: "warn",
			kind: LevelKindWarn,
			want: "WRN",
		},
		{
			name: "error",
			kind: LevelKindError,
			want: "ERR",
		},
		{
			name: "fatal",
			kind: LevelKindFatal,
			want: "FTL",
		},
		{
			name: "unknown",
			kind: LevelKind(99),
			want: "LVL99",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.kind.ToShortString()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFromSlogLevel(t *testing.T) {
	tests := []struct {
		name  string
		level slog.Level
		want  LevelKind
	}{
		{
			name:  "debug",
			level: slog.LevelDebug,
			want:  LevelKindDebug,
		},
		{
			name:  "info",
			level: slog.LevelInfo,
			want:  LevelKindInfo,
		},
		{
			name:  "warn",
			level: slog.LevelWarn,
			want:  LevelKindWarn,
		},
		{
			name:  "error",
			level: slog.LevelError,
			want:  LevelKindError,
		},
		{
			name:  "fatal custom",
			level: SlogFatalLevel,
			want:  LevelKindFatal,
		},
		{
			name:  "unknown defaults to debug",
			level: slog.Level(99),
			want:  LevelKindDebug,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromSlogLevel(tt.level)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFromZapLevel(t *testing.T) {
	tests := []struct {
		name  string
		level zapcore.Level
		want  LevelKind
	}{
		{
			name:  "debug",
			level: zapcore.DebugLevel,
			want:  LevelKindDebug,
		},
		{
			name:  "info",
			level: zapcore.InfoLevel,
			want:  LevelKindInfo,
		},
		{
			name:  "warn",
			level: zapcore.WarnLevel,
			want:  LevelKindWarn,
		},
		{
			name:  "error",
			level: zapcore.ErrorLevel,
			want:  LevelKindError,
		},
		{
			name:  "fatal",
			level: zapcore.FatalLevel,
			want:  LevelKindFatal,
		},
		{
			name:  "unknown defaults to debug",
			level: zapcore.Level(99),
			want:  LevelKindDebug,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromZapLevel(tt.level)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLevelConversionRoundTrip(t *testing.T) {
	tests := []LevelKind{
		LevelKindDebug,
		LevelKindInfo,
		LevelKindWarn,
		LevelKindError,
		LevelKindFatal,
	}

	for _, kind := range tests {
		t.Run(kind.ToShortString(), func(t *testing.T) {
			// Slog round trip
			slogLevel := kind.ToSlogLevel()
			backFromSlog := FromSlogLevel(slogLevel)
			assert.Equal(t, kind, backFromSlog)

			// Zap round trip
			zapLevel := kind.ToZapLevel()
			backFromZap := FromZapLevel(zapLevel)
			assert.Equal(t, kind, backFromZap)
		})
	}
}
