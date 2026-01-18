package logger

import (
	"log/slog"
	"os"
)

var defaultLogger *slog.Logger

// InitLogger initializes the global logger with the specified verbosity level
func InitLogger(verbose bool) {
	var level slog.Level
	if verbose {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewTextHandler(os.Stderr, opts)
	defaultLogger = slog.New(handler)
	slog.SetDefault(defaultLogger)
}

// Debug logs a debug-level message
func Debug(msg string, args ...any) {
	if defaultLogger != nil {
		defaultLogger.Debug(msg, args...)
	}
}

// Info logs an info-level message
func Info(msg string, args ...any) {
	if defaultLogger != nil {
		defaultLogger.Info(msg, args...)
	}
}

// Error logs an error-level message
func Error(msg string, args ...any) {
	if defaultLogger != nil {
		defaultLogger.Error(msg, args...)
	}
}

// With returns a logger with additional attributes
func With(args ...any) *slog.Logger {
	if defaultLogger != nil {
		return defaultLogger.With(args...)
	}
	return slog.Default()
}
