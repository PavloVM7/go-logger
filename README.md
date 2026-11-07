# go-logger

<p align="center">
  <a href="https://github.com/PavloVM7/go-logger/actions/workflows/go.yml"><img src="https://github.com/PavloVM7/go-logger/actions/workflows/go.yml/badge.svg" alt="Test Status"></a>
  &nbsp;&nbsp;&nbsp;
  <a href="https://github.com/PavloVM7/go-logger/actions/workflows/revive.yml"><img src="https://github.com/PavloVM7/go-logger/actions/workflows/revive.yml/badge.svg" alt="Test Status"></a>
  <a href="https://golang.org/dl/" style="padding:20px"><img src="https://img.shields.io/badge/Go-1.27.3+-00ADD8?style=flat&logo=go" alt="Go Version"></a>
</p>


A structured logging library for Go that provides a unified interface over multiple logging backends (slog and zap). It supports context propagation, structured key-value logging, and multiple output destinations.

## Features

- **Unified Interface**: Single API for both slog and zap implementations
- **Context Support**: Store and retrieve loggers from context
- **Structured Logging**: Key-value pairs with both formatted and structured methods
- **Multiple Outputs**: Log to stdout, stderr, or files simultaneously
- **Flexible Configuration**: JSON or console encoding, custom date formats, log levels
- **Level Filtering**: Support for debug, info, warn, error, and fatal levels
- **Named Loggers**: Create child loggers with specific names

## Installation

```bash
go get github.com/PavloVM7/go-logger
```

## Quick Start

### Using slog (Standard Library)

```go
package main

import (
    "github.com/PavloVM7/go-logger/logger"
)

func main() {
    cfg := logger.Config{
        Name:       "my-app",
        Level:      "info",
        Encoding:   "json",
        DateFormat: "2006-01-02T15:04:05.000Z-07:00",
        Outputs:    []string{"stdout"},
    }

    log, err := logger.NewLoggerSlog(cfg)
    if err != nil {
        panic(err)
    }
    defer log.Close()

    log.Infow("Application started", "version", "1.0.0")
    log.Debugf("Debug information: %v", "some data")
}
```

### Using zap

```go
package main

import (
    "github.com/PavloVM7/go-logger/logger"
)

func main() {
    cfg := logger.Config{
        Name:       "my-app",
        Level:      "info",
        Encoding:   "json",
        DateFormat: "2006-01-02T15:04:05.000Z-07:00",
        Outputs:    []string{"stdout"},
    }

    log, err := logger.NewLoggerZap(cfg)
    if err != nil {
        panic(err)
    }
    defer log.Close()

    log.Infow("Application started", "version", "1.0.0")
    log.Debugf("Debug information: %v", "some data")
}
```

## Configuration

### Config Structure

```go
type Config struct {
    Name       string   // Logger name (defaults to "default")
    Level      string   // Log level: debug, info, warn, error, fatal
    AddSource  bool     // Add source file and line number (currently unused)
    DateFormat string   // Timestamp format (defaults to RFC3339)
    Encoding   string   // Output encoding: "json" or "console"
    Outputs    []string // Output destinations: "stdout", "stderr", or file paths
}
```

### Log Levels

- `debug` (-1): Detailed debugging information
- `info` (0): General informational messages (default)
- `warn` (1): Potentially harmful situations
- `error` (2): Error events that might allow the application to continue
- `fatal` (3): Very severe errors that will cause the application to abort

### Output Destinations

- `"stdout"`: Standard output
- `"stderr"`: Standard error
- File paths: Absolute or relative paths to log files

Example with multiple outputs:
```go
cfg := logger.Config{
    Name:    "my-app",
    Level:   "info",
    Outputs: []string{"stdout", "/var/log/my-app.log"},
}
```

## API Reference

### Logger Interface

```go
type Logger interface {
    // Formatted logging methods
    Fatalf(format string, args ...any)
    Errorf(format string, args ...any)
    Warnf(format string, args ...any)
    Infof(format string, args ...any)
    Debugf(format string, args ...any)

    // Structured logging methods
    Fatalw(msg string, keysAndValues ...any)
    Errorw(msg string, keysAndValues ...any)
    Warnw(msg string, keysAndValues ...any)
    Infow(msg string, keysAndValues ...any)
    Debugw(msg string, keysAndValues ...any)

    // Context and naming
    Named(name string) Logger
    With(keysAndValues ...any) Logger
    Name() string

    // Level checking
    IsLoggingLevelDebug() bool
    IsLoggingLevelInfo() bool
    IsLoggingLevelWarn() bool
    IsLoggingLevelError() bool

    // Cleanup
    Close() error
}
```

### Context Integration

```go
import "context"

// Store logger in context
ctx := logger.ToContext(context.Background(), log)

// Retrieve logger from context
retrievedLog, err := logger.FromContext(ctx)
if err != nil {
    // Handle error
}
```

### Named Loggers

Create child loggers with specific names for better organization:

```go
httpLogger := log.Named("http")
dbLogger := log.Named("database")

httpLogger.Infow("Request received", "path", "/api/users")
dbLogger.Infow("Query executed", "table", "users")
```

### With Context

Add structured fields to all log messages from a logger instance:

```go
requestLogger := log.With("requestID", "abc-123", "userID", 456)
requestLogger.Infow("Processing request") // Includes requestID and userID
```

## Examples

### Console Encoding

```go
cfg := logger.Config{
    Name:     "my-app",
    Level:    "debug",
    Encoding: "console",
    Outputs:  []string{"stdout"},
}
log, _ := logger.NewLoggerSlog(cfg)
log.Infow("User logged in", "username", "john")
```

### JSON Encoding with File Output

```go
cfg := logger.Config{
    Name:       "my-app",
    Level:      "info",
    Encoding:   "json",
    DateFormat: "2006-01-02T15:04:05.000Z-07:00",
    Outputs:    []string{"stdout", "./logs/app.log"},
}
log, _ := logger.NewLoggerZap(cfg)
defer log.Close()

log.Errorw("Database connection failed", "error", "connection timeout")
```

### Level Filtering

```go
cfg := logger.Config{
    Name:    "my-app",
    Level:   "warn", // Only warn, error, and fatal will be logged
    Outputs: []string{"stdout"},
}
log, _ := logger.NewLoggerSlog(cfg)

log.Debugf("This won't be logged")  // Below warn level
log.Infow("This won't be logged")   // Below warn level
log.Warnw("This will be logged")     // At warn level
log.Errorw("This will be logged")    // At error level
```

## Testing

Run tests:

```bash
make test
```

Run linting:

```bash
make lint
```

Run all verification steps:

```bash
make verify
```

## License

BSD License - see LICENSE file for details.

