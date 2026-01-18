package logger

import (
	"bytes"
	"log/slog"
	"os"
	"strings"
	"testing"
)

func TestInitLogger(t *testing.T) {
	tests := []struct {
		name    string
		verbose bool
		wantMsg string
		logFunc func()
	}{
		{
			name:    "verbose mode shows debug messages",
			verbose: true,
			wantMsg: "debug message",
			logFunc: func() {
				Debug("debug message")
			},
		},
		{
			name:    "non-verbose mode hides debug messages",
			verbose: false,
			wantMsg: "",
			logFunc: func() {
				Debug("debug message")
			},
		},
		{
			name:    "info messages always shown",
			verbose: false,
			wantMsg: "info message",
			logFunc: func() {
				Info("info message")
			},
		},
		{
			name:    "error messages always shown",
			verbose: false,
			wantMsg: "error message",
			logFunc: func() {
				Error("error message")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stderr output
			var buf bytes.Buffer
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			// Initialize logger
			InitLogger(tt.verbose)

			// Execute the log function
			tt.logFunc()

			// Restore stderr and read captured output
			w.Close()
			os.Stderr = oldStderr
			buf.ReadFrom(r)
			output := buf.String()

			// Verify output
			if tt.wantMsg != "" {
				if !strings.Contains(output, tt.wantMsg) {
					t.Errorf("expected output to contain %q, got %q", tt.wantMsg, output)
				}
			} else {
				if strings.Contains(output, "debug message") {
					t.Errorf("expected no debug output in non-verbose mode, got %q", output)
				}
			}
		})
	}
}

func TestDebugWithAttributes(t *testing.T) {
	// Capture stderr output
	var buf bytes.Buffer
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Initialize logger in verbose mode
	InitLogger(true)

	// Log with attributes
	Debug("test message", "key1", "value1", "key2", 42)

	// Restore stderr and read captured output
	w.Close()
	os.Stderr = oldStderr
	buf.ReadFrom(r)
	output := buf.String()

	// Verify attributes are in output
	if !strings.Contains(output, "test message") {
		t.Errorf("expected output to contain 'test message', got %q", output)
	}
	if !strings.Contains(output, "key1=value1") {
		t.Errorf("expected output to contain 'key1=value1', got %q", output)
	}
	if !strings.Contains(output, "key2=42") {
		t.Errorf("expected output to contain 'key2=42', got %q", output)
	}
}

func TestWith(t *testing.T) {
	// Initialize logger
	InitLogger(true)

	// Create a logger with attributes
	logger := With("component", "test")

	// Verify we got a logger back
	if logger == nil {
		t.Error("expected non-nil logger from With()")
	}

	// Verify it's a slog.Logger
	if _, ok := interface{}(logger).(*slog.Logger); !ok {
		t.Error("expected With() to return *slog.Logger")
	}
}

func TestLoggerBeforeInit(t *testing.T) {
	// Reset logger to nil
	defaultLogger = nil

	// These should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("logging before init caused panic: %v", r)
		}
	}()

	Debug("test")
	Info("test")
	Error("test")
	logger := With("key", "value")
	if logger == nil {
		t.Error("With() should return default logger when not initialized")
	}
}
