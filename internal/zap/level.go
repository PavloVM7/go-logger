package zap

import (
	"fmt"

	"github.com/PavloVM7/go-logger/internal"
	"go.uber.org/zap/zapcore"
)

func levelShortStringEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(levelToShortString(l))
}

var shortStrings = map[zapcore.Level]string{
	zapcore.DebugLevel: internal.DebugShort.String(),
	zapcore.InfoLevel:  internal.InfoShort.String(),
	zapcore.WarnLevel:  internal.WarnShort.String(),
	zapcore.ErrorLevel: internal.ErrorShort.String(),
	zapcore.FatalLevel: internal.FatalShort.String(),
}

func levelToShortString(l zapcore.Level) string {
	if str, ok := shortStrings[l]; ok {
		return str
	}
	return fmt.Sprintf("LVL%d", l)
}

var zapLevels = map[internal.LevelKind]zapcore.Level{
	internal.LevelKindDebug: zapcore.DebugLevel,
	internal.LevelKindInfo:  zapcore.InfoLevel,
	internal.LevelKindWarn:  zapcore.WarnLevel,
	internal.LevelKindError: zapcore.ErrorLevel,
	internal.LevelKindFatal: zapcore.FatalLevel,
}

// ToZapLevel converts LevelKind to zapcore.Level.
func kindLevelToZapLevel(lk internal.LevelKind) zapcore.Level {
	if zl, ok := zapLevels[lk]; ok {
		return zl
	}
	return zapcore.DebugLevel
}
