package logger

import (
	"log/slog"
	"os"
)

var defaultLogger *slog.Logger

func init() {
	// Initialize with Info level by default
	defaultLogger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

// Init initializes the logger with the specified configuration
func Init(verbose, quiet bool) {
	var level slog.Level
	if quiet {
		level = slog.LevelError
	} else if verbose {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	defaultLogger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	}))
}

// Debug logs a debug message with optional attributes
func Debug(msg string, args ...any) {
	defaultLogger.Debug(msg, args...)
}

// Info logs an info message with optional attributes
func Info(msg string, args ...any) {
	defaultLogger.Info(msg, args...)
}

// Warn logs a warning message with optional attributes
func Warn(msg string, args ...any) {
	defaultLogger.Warn(msg, args...)
}

// Error logs an error message with optional attributes
func Error(msg string, args ...any) {
	defaultLogger.Error(msg, args...)
}

// With returns a new logger with the given attributes
func With(args ...any) *slog.Logger {
	return defaultLogger.With(args...)
}

// Logger returns the default logger instance
func Logger() *slog.Logger {
	return defaultLogger
}
