package slog

import (
	"errors"
	"fmt"
	"io"
	logslog "log/slog"
	"os"
	"slices"
	"strings"

	"github.com/PavloVM7/go-logger/internal"
	"github.com/PavloVM7/go-logger/pkg/set"
)

// Config holds the configuration for a slog logger.
type Config struct {
	internal.ConfigBase
	w io.Writer
}

// AddStdoutOutput adds stdout to the logger outputs.
func (c *Config) AddStdoutOutput() {
	c.AddOutput(internal.OutputStdout)
}

// AddOutput adds a custom output to the logger outputs.
func (c *Config) AddOutput(output string) {
	if slices.Contains(c.Outputs, output) {
		return
	}

	c.Outputs = append(c.Outputs, output)
}

func (c *Config) checkOutputs() {
	if len(c.Outputs) == 0 {
		c.Outputs = []string{internal.OutputStdout}
		return
	}

	setOutputs := set.NewCapacity[string](len(c.Outputs))
	for _, out := range c.Outputs {
		trm := strings.TrimSpace(out)
		if len(trm) == 0 {
			continue
		}
		setOutputs.Add(trm)
	}

	c.Outputs = setOutputs.ToSlice()
	if len(c.Outputs) == 0 {
		c.Outputs = []string{internal.OutputStdout}
	}
}

func (c *Config) check() {
	c.checkOutputs()

	name := strings.TrimSpace(c.Name)
	if len(name) == 0 {
		name = "default"
	}
	c.Name = name
}

func (c *Config) createHandlerOptions() (logslog.HandlerOptions, error) {
	level, err := stringToLevel(c.Level)
	if err != nil {
		return logslog.HandlerOptions{}, fmt.Errorf("slog config error: %v", err)
	}

	return logslog.HandlerOptions{
		Level:     level,
		AddSource: false,
	}, nil
}

func (c *Config) createWriters() ([]io.Writer, error) {
	c.check()
	writers, err := createWrites(c.Outputs)
	if err != nil {
		return nil, fmt.Errorf("failed to create slog writer: %w", err)
	}

	// to tests only
	if c.w != nil {
		_ = addWriter(&writers, c.w)
	}

	return writers, nil
}

func createWriter(writers []io.Writer) (io.Writer, error) {
	if len(writers) == 0 {
		return nil, errors.New("no writers for logging")
	}

	if len(writers) == 1 {
		return writers[0], nil
	}

	return io.MultiWriter(writers...), nil
}

func createWrites(outputs []string) ([]io.Writer, error) {
	result := make([]io.Writer, 0, len(outputs))
	var err error
	for _, out := range outputs {
		switch out {
		case internal.OutputStdout:
			err = addWriter(&result, os.Stdout)
		case internal.OutputStderr:
			err = addWriter(&result, os.Stderr)
		default:
			err = addFile(&result, out)
		}

		if err != nil {
			err = fmt.Errorf("'%s' %w", out, err)
			break
		}
	}

	if err != nil {
		closeWriters(result)
		return nil, err
	}

	return slices.Clip(result), nil
}

func addFile(writers *[]io.Writer, fileName string) error {
	f, err := createLogFile(fileName)
	if err != nil {
		return err
	}
	return addWriter(writers, f)
}

func addWriter(writers *[]io.Writer, w io.Writer) error {
	if slices.Contains(*writers, w) {
		return errors.New("duplicate writer")
	}
	*writers = append(*writers, w)
	return nil
}

func createLogFile(fileName string) (*os.File, error) {
	logFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	return logFile, nil
}
