package azfunc

import (
	"log/slog"
	"os"
)

// logger is the interface that wraps around methods Debug, Error, Info
// and Warn.
type logger interface {
	Debug(msg string, args ...any)
	Error(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
}

// Logger wraps around the log/slog package and satisfies the logger.
type Logger struct {
	stdout *slog.Logger
	stderr *slog.Logger
}

// NewLogger creates and returns a new Logger. Errors and warnings are
// written to stderr, while debug and info messages are written to stdout.
func NewLogger() Logger {
	return Logger{
		stderr: slog.New(slog.NewJSONHandler(os.Stderr, nil)),
		stdout: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

// Debug logs at [LevelDebug].
func (l Logger) Debug(msg string, args ...any) {
	l.stdout.Debug(msg, args...)
}

// Error logs at [LevelError].
func (l Logger) Error(msg string, args ...any) {
	l.stderr.Error(msg, args...)
}

// Info logs at [LevelInfo].
func (l Logger) Info(msg string, args ...any) {
	l.stdout.Info(msg, args...)
}

// Warn logs at [LevelWarn].
func (l Logger) Warn(msg string, args ...any) {
	l.stderr.Warn(msg, args...)
}

// noOpLogger is a placeholder for when no logger is provided to the
// function app.
type noOpLogger struct{}

// Debug together with Error, Info and Warn satisfies the logger interface.
func (l noOpLogger) Debug(msg string, args ...any) {}

// Error together with Debug, Info and Warn satisfies the logger interface.
func (l noOpLogger) Error(msg string, args ...any) {}

// Info together with Debug, Error and Warn satisfies the logger interface.
func (l noOpLogger) Info(msg string, args ...any) {}

// Warn together with Debug, Info and Error satisfies the logger interface.
func (l noOpLogger) Warn(msg string, args ...any) {}
