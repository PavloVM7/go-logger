package zap

import (
	"fmt"

	"github.com/PavloVM7/go-logger/internal"
	"go.uber.org/zap/zapcore"
)

func levelShortStringEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(levelToShortString(l))
}

func levelToShortString(l zapcore.Level) string {
	switch l {
	case zapcore.DebugLevel:
		return internal.DebugShort.String()
	case zapcore.InfoLevel:
		return internal.InfoShort.String()
	case zapcore.WarnLevel:
		return internal.WarnShort.String()
	case zapcore.ErrorLevel:
		return internal.ErrorShort.String()
	case zapcore.FatalLevel:
		return internal.FatalShort.String()
	default:
		return fmt.Sprintf("LVL%d", l)
	}
}
