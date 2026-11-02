package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
