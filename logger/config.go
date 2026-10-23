// Package logger provides a structured logging interface with context support.
package logger

import "strings"

const defaultDateFormat = "2006-01-02T15:04:05.000Z-07:00"

// Config holds the configuration for a logger instance.
type Config struct {
	Name       string `json:"name" yaml:"name"`
	Level      string `json:"level" yaml:"level"`
	AddSource  bool   `json:"add_source" yaml:"add_source"`
	DateFormat string `json:"date_format" yaml:"date_format"`
}

// GetLevel parses and returns the log level from the configuration.
func (c *Config) GetLevel() (Level, error) {
	low := strings.ToLower(strings.TrimSpace(c.Level))
	return ParseLevel(low)
}
