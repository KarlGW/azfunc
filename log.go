package azfunc

import (
	"log/slog"
	"os"
	"strings"
	"sync"
)

// Logger is the interface that wraps around methods Debug, Error, Info
// and Warn. It is used to log to stdout and stderr.
type Logger interface {
	Debug(msg string, args ...any)
	Error(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
}

// logger wraps around the log/slog package and satisfies the logger.
type logger struct {
	stdout *slog.Logger
	stderr *slog.Logger
}

// NewLogger creates and returns a new Logger. Errors and warnings are
// written to stderr, while debug and info messages are written to stdout.
func NewLogger() Logger {
	return logger{
		stderr: slog.New(slog.NewJSONHandler(os.Stderr, nil)),
		stdout: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

// Debug logs at [LevelDebug].
func (l logger) Debug(msg string, args ...any) {
	l.stdout.Debug(msg, args...)
}

// Error logs at [LevelError].
func (l logger) Error(msg string, args ...any) {
	l.stderr.Error(msg, args...)
}

// Info logs at [LevelInfo].
func (l logger) Info(msg string, args ...any) {
	l.stdout.Info(msg, args...)
}

// Warn logs at [LevelWarn].
func (l logger) Warn(msg string, args ...any) {
	l.stderr.Warn(msg, args...)
}

// InvocationLogger is the interface that wraps around the methods Debug,
// Error, Info, Warn and Write. It is used to log to the function host.
type InvocationLogger interface {
	Logger
	Write(msg string)
	Entries() []string
}

// invocationLogger is used to log to the function host.
type invocationLogger struct {
	l  *slog.Logger
	w  *stringWriter
	mu *sync.RWMutex
}

// newInvocationLogger creates and returns a new invocationLogger.
func newInvocationLogger() invocationLogger {
	w := &stringWriter{
		entries: make([]string, 0),
	}

	return invocationLogger{
		l:  slog.New(slog.NewJSONHandler(w, nil)),
		w:  w,
		mu: &sync.RWMutex{},
	}
}

// Debug together with Error, Info and Warn satisfies the logger interface.
func (l invocationLogger) Debug(msg string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.l.Debug(msg, args...)
	index := len(l.w.entries) - 1
	l.w.entries[index] = strings.TrimSuffix(l.w.entries[index], "\n")
}

// Error together with Debug, Info and Warn satisfies the logger interface.
func (l invocationLogger) Error(msg string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.l.Error(msg, args...)
	index := len(l.w.entries) - 1
	l.w.entries[index] = strings.TrimSuffix(l.w.entries[index], "\n")
}

// Info together with Debug, Error and Warn satisfies the logger interface.
func (l invocationLogger) Info(msg string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.l.Info(msg, args...)
	index := len(l.w.entries) - 1
	l.w.entries[index] = strings.TrimSuffix(l.w.entries[index], "\n")
}

// Warn together with Debug, Info and Error satisfies the logger interface.
func (l invocationLogger) Warn(msg string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.l.Warn(msg, args...)
	index := len(l.w.entries) - 1
	l.w.entries[index] = strings.TrimSuffix(l.w.entries[index], "\n")
}

// Write a message to the log.
func (l invocationLogger) Write(msg string) {
	l.w.Write([]byte(msg))
}

// Entries returns the entries (log) written to the stringWriter.
func (l invocationLogger) Entries() []string {
	if len(l.w.entries) == 0 {
		return nil
	}
	return l.w.entries
}

// stringWriter writes to an string slice. Used by the invocationLogger.
type stringWriter struct {
	entries []string
}

// Write an entry to the stringWriter.
func (w *stringWriter) Write(p []byte) (n int, err error) {
	w.entries = append(w.entries, string(p))
	return len(p), nil
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
