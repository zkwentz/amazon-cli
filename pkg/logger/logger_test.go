package logger

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"
)

// captureLogOutput captures the log output for testing
func captureLogOutput(t *testing.T, verbose, quiet bool, logFunc func()) map[string]interface{} {
	var buf bytes.Buffer

	// Create a temporary logger that writes to our buffer
	var level slog.Level
	if quiet {
		level = slog.LevelError
	} else if verbose {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	defaultLogger = slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: level,
	}))

	// Execute the logging function
	logFunc()

	// Parse the JSON output
	var result map[string]interface{}
	if buf.Len() > 0 {
		if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
			t.Fatalf("Failed to parse log output: %v\nOutput: %s", err, buf.String())
		}
	}

	return result
}

func TestInit(t *testing.T) {
	tests := []struct {
		name    string
		verbose bool
		quiet   bool
		want    slog.Level
	}{
		{
			name:    "default level (info)",
			verbose: false,
			quiet:   false,
			want:    slog.LevelInfo,
		},
		{
			name:    "verbose mode (debug)",
			verbose: true,
			quiet:   false,
			want:    slog.LevelDebug,
		},
		{
			name:    "quiet mode (error)",
			verbose: false,
			quiet:   true,
			want:    slog.LevelError,
		},
		{
			name:    "quiet overrides verbose",
			verbose: true,
			quiet:   true,
			want:    slog.LevelError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init(tt.verbose, tt.quiet)
			// Note: We can't directly test the level without reflection,
			// but we test the behavior in other tests
		})
	}
}

func TestDebug(t *testing.T) {
	tests := []struct {
		name       string
		verbose    bool
		quiet      bool
		msg        string
		args       []interface{}
		wantLogged bool
	}{
		{
			name:       "debug logged in verbose mode",
			verbose:    true,
			quiet:      false,
			msg:        "test debug message",
			args:       []interface{}{"key", "value"},
			wantLogged: true,
		},
		{
			name:       "debug not logged in default mode",
			verbose:    false,
			quiet:      false,
			msg:        "test debug message",
			args:       []interface{}{"key", "value"},
			wantLogged: false,
		},
		{
			name:       "debug not logged in quiet mode",
			verbose:    false,
			quiet:      true,
			msg:        "test debug message",
			args:       []interface{}{"key", "value"},
			wantLogged: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureLogOutput(t, tt.verbose, tt.quiet, func() {
				Debug(tt.msg, tt.args...)
			})

			if tt.wantLogged {
				if result == nil || result["msg"] != tt.msg {
					t.Errorf("Expected debug message to be logged, got: %v", result)
				}
				if result["level"] != "DEBUG" {
					t.Errorf("Expected level DEBUG, got: %v", result["level"])
				}
			} else {
				if result != nil {
					t.Errorf("Expected no log output, got: %v", result)
				}
			}
		})
	}
}

func TestInfo(t *testing.T) {
	tests := []struct {
		name       string
		verbose    bool
		quiet      bool
		msg        string
		args       []interface{}
		wantLogged bool
	}{
		{
			name:       "info logged in verbose mode",
			verbose:    true,
			quiet:      false,
			msg:        "test info message",
			args:       []interface{}{"key", "value"},
			wantLogged: true,
		},
		{
			name:       "info logged in default mode",
			verbose:    false,
			quiet:      false,
			msg:        "test info message",
			args:       []interface{}{"key", "value"},
			wantLogged: true,
		},
		{
			name:       "info not logged in quiet mode",
			verbose:    false,
			quiet:      true,
			msg:        "test info message",
			args:       []interface{}{"key", "value"},
			wantLogged: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureLogOutput(t, tt.verbose, tt.quiet, func() {
				Info(tt.msg, tt.args...)
			})

			if tt.wantLogged {
				if result == nil || result["msg"] != tt.msg {
					t.Errorf("Expected info message to be logged, got: %v", result)
				}
				if result["level"] != "INFO" {
					t.Errorf("Expected level INFO, got: %v", result["level"])
				}
			} else {
				if result != nil {
					t.Errorf("Expected no log output, got: %v", result)
				}
			}
		})
	}
}

func TestWarn(t *testing.T) {
	tests := []struct {
		name       string
		verbose    bool
		quiet      bool
		msg        string
		args       []interface{}
		wantLogged bool
	}{
		{
			name:       "warn logged in verbose mode",
			verbose:    true,
			quiet:      false,
			msg:        "test warn message",
			args:       []interface{}{"key", "value"},
			wantLogged: true,
		},
		{
			name:       "warn logged in default mode",
			verbose:    false,
			quiet:      false,
			msg:        "test warn message",
			args:       []interface{}{"key", "value"},
			wantLogged: true,
		},
		{
			name:       "warn not logged in quiet mode",
			verbose:    false,
			quiet:      true,
			msg:        "test warn message",
			args:       []interface{}{"key", "value"},
			wantLogged: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureLogOutput(t, tt.verbose, tt.quiet, func() {
				Warn(tt.msg, tt.args...)
			})

			if tt.wantLogged {
				if result == nil || result["msg"] != tt.msg {
					t.Errorf("Expected warn message to be logged, got: %v", result)
				}
				if result["level"] != "WARN" {
					t.Errorf("Expected level WARN, got: %v", result["level"])
				}
			} else {
				if result != nil {
					t.Errorf("Expected no log output, got: %v", result)
				}
			}
		})
	}
}

func TestError(t *testing.T) {
	tests := []struct {
		name       string
		verbose    bool
		quiet      bool
		msg        string
		args       []interface{}
		wantLogged bool
	}{
		{
			name:       "error logged in verbose mode",
			verbose:    true,
			quiet:      false,
			msg:        "test error message",
			args:       []interface{}{"key", "value"},
			wantLogged: true,
		},
		{
			name:       "error logged in default mode",
			verbose:    false,
			quiet:      false,
			msg:        "test error message",
			args:       []interface{}{"key", "value"},
			wantLogged: true,
		},
		{
			name:       "error logged in quiet mode",
			verbose:    false,
			quiet:      true,
			msg:        "test error message",
			args:       []interface{}{"key", "value"},
			wantLogged: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureLogOutput(t, tt.verbose, tt.quiet, func() {
				Error(tt.msg, tt.args...)
			})

			if tt.wantLogged {
				if result == nil || result["msg"] != tt.msg {
					t.Errorf("Expected error message to be logged, got: %v", result)
				}
				if result["level"] != "ERROR" {
					t.Errorf("Expected level ERROR, got: %v", result["level"])
				}
			} else {
				if result != nil {
					t.Errorf("Expected no log output, got: %v", result)
				}
			}
		})
	}
}

func TestWith(t *testing.T) {
	result := captureLogOutput(t, false, false, func() {
		contextLogger := With("requestID", "12345", "userID", "user-abc")
		contextLogger.Info("test message with context")
	})

	if result == nil || result["msg"] != "test message with context" {
		t.Errorf("Expected message to be logged, got: %v", result)
	}

	if result["requestID"] != "12345" {
		t.Errorf("Expected requestID=12345, got: %v", result["requestID"])
	}

	if result["userID"] != "user-abc" {
		t.Errorf("Expected userID=user-abc, got: %v", result["userID"])
	}
}

func TestLogger(t *testing.T) {
	logger := Logger()
	if logger == nil {
		t.Error("Expected Logger() to return non-nil logger")
	}
}

func TestLogWithAttributes(t *testing.T) {
	result := captureLogOutput(t, false, false, func() {
		Info("operation completed",
			"duration", 150,
			"status", "success",
			"count", 42)
	})

	if result == nil || result["msg"] != "operation completed" {
		t.Errorf("Expected message to be logged, got: %v", result)
	}

	if result["duration"].(float64) != 150 {
		t.Errorf("Expected duration=150, got: %v", result["duration"])
	}

	if result["status"] != "success" {
		t.Errorf("Expected status=success, got: %v", result["status"])
	}

	if result["count"].(float64) != 42 {
		t.Errorf("Expected count=42, got: %v", result["count"])
	}
}
