package logger

import (
	"errors"
)

var (
	// ErrLoggerNotFound is returned when a logger is not found in the context.
	ErrLoggerNotFound = errors.New("logger not found")
	// ErrUnknownLevel is returned when an invalid log level string is provided.
	ErrUnknownLevel = errors.New("unknown error level")
)
