package azfunc

// logger is the interface that wraps around methods Debug, Error, Info
// and Warn.
type logger interface {
	Debug(msg string, args ...any)
	Error(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
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
