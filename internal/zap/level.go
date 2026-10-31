package zap

import (
	"github.com/PavloVM7/go-logger/internal"
	"go.uber.org/zap/zapcore"
)

func levelShortStringEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(levelToShortString(l))
}

func levelToShortString(l zapcore.Level) string {
	kind := internal.FromZapLevel(l)
	return kind.ToShortString()
}
