package amazon

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewTimingTransport(t *testing.T) {
	tests := []struct {
		name    string
		verbose bool
	}{
		{"verbose mode enabled", true},
		{"verbose mode disabled", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transport := NewTimingTransport(tt.verbose)

			if transport == nil {
				t.Fatal("NewTimingTransport returned nil")
			}

			if transport.Transport == nil {
				t.Error("Transport should not be nil")
			}

			if transport.Logger == nil {
				t.Error("Logger should not be nil")
			}

			if transport.Verbose != tt.verbose {
				t.Errorf("Verbose = %v, want %v", transport.Verbose, tt.verbose)
			}
		})
	}
}

func TestTimingTransport_RoundTrip(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		delay          time.Duration
		verbose        bool
		wantLogContent []string
	}{
		{
			name:       "successful request with 200",
			statusCode: 200,
			delay:      10 * time.Millisecond,
			verbose:    false,
			wantLogContent: []string{
				"✓",
				"GET",
				"Status: 200",
				"Duration:",
			},
		},
		{
			name:       "failed request with 404",
			statusCode: 404,
			delay:      5 * time.Millisecond,
			verbose:    false,
			wantLogContent: []string{
				"✗",
				"GET",
				"Status: 404",
			},
		},
		{
			name:       "verbose mode includes extra details",
			statusCode: 200,
			delay:      50 * time.Millisecond,
			verbose:    true,
			wantLogContent: []string{
				"→",
				"✓",
				"GET",
				"Performance:",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server with configurable delay
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(tt.delay)
				w.WriteHeader(tt.statusCode)
				w.Write([]byte("test response"))
			}))
			defer server.Close()

			// Capture log output
			var logBuf bytes.Buffer
			transport := NewTimingTransport(tt.verbose)
			transport.Logger.SetOutput(&logBuf)

			// Create HTTP client with timing transport
			client := &http.Client{Transport: transport}

			// Make request
			resp, err := client.Get(server.URL)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			// Verify response
			if resp.StatusCode != tt.statusCode {
				t.Errorf("StatusCode = %d, want %d", resp.StatusCode, tt.statusCode)
			}

			// Verify log output
			logOutput := logBuf.String()
			for _, want := range tt.wantLogContent {
				if !strings.Contains(logOutput, want) {
					t.Errorf("Log output missing %q\nGot: %s", want, logOutput)
				}
			}
		})
	}
}

func TestTimingTransport_RoundTrip_NetworkError(t *testing.T) {
	// Capture log output
	var logBuf bytes.Buffer
	transport := NewTimingTransport(false)
	transport.Logger.SetOutput(&logBuf)

	// Create HTTP client with timing transport
	client := &http.Client{Transport: transport}

	// Make request to invalid URL
	_, err := client.Get("http://invalid-host-that-does-not-exist-12345.com")
	if err == nil {
		t.Fatal("Expected error for invalid host")
	}

	// Verify error logging
	logOutput := logBuf.String()
	expectedContent := []string{"✗", "GET", "Failed:", "Duration:"}
	for _, want := range expectedContent {
		if !strings.Contains(logOutput, want) {
			t.Errorf("Log output missing %q\nGot: %s", want, logOutput)
		}
	}
}

func TestNewRequestTimingLogger(t *testing.T) {
	tests := []struct {
		name    string
		verbose bool
	}{
		{"verbose enabled", true},
		{"verbose disabled", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewRequestTimingLogger(tt.verbose)

			if logger == nil {
				t.Fatal("NewRequestTimingLogger returned nil")
			}

			if logger.logger == nil {
				t.Error("logger should not be nil")
			}

			if logger.verbose != tt.verbose {
				t.Errorf("verbose = %v, want %v", logger.verbose, tt.verbose)
			}
		})
	}
}

func TestRequestTimingLogger_LogOperation(t *testing.T) {
	tests := []struct {
		name           string
		operation      string
		duration       time.Duration
		err            error
		verbose        bool
		wantLogContent []string
	}{
		{
			name:      "successful operation",
			operation: "TestOperation",
			duration:  50 * time.Millisecond,
			err:       nil,
			verbose:   false,
			wantLogContent: []string{
				"✓",
				"TestOperation",
				"Success",
				"Duration:",
			},
		},
		{
			name:      "failed operation",
			operation: "FailedOperation",
			duration:  10 * time.Millisecond,
			err:       errors.New("test error"),
			verbose:   false,
			wantLogContent: []string{
				"✗",
				"FailedOperation",
				"Failed:",
				"test error",
			},
		},
		{
			name:      "verbose mode includes performance category",
			operation: "VerboseOperation",
			duration:  100 * time.Millisecond,
			err:       nil,
			verbose:   true,
			wantLogContent: []string{
				"✓",
				"VerboseOperation",
				"Performance:",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture log output
			var logBuf bytes.Buffer
			logger := NewRequestTimingLogger(tt.verbose)
			logger.SetOutput(&logBuf)

			// Log operation
			logger.LogOperation(tt.operation, tt.duration, tt.err)

			// Verify log output
			logOutput := logBuf.String()
			for _, want := range tt.wantLogContent {
				if !strings.Contains(logOutput, want) {
					t.Errorf("Log output missing %q\nGot: %s", want, logOutput)
				}
			}
		})
	}
}

func TestRequestTimingLogger_TimeOperation(t *testing.T) {
	tests := []struct {
		name         string
		operation    string
		operationFn  func() error
		expectError  bool
		wantDuration time.Duration
	}{
		{
			name:      "successful operation",
			operation: "TestOp",
			operationFn: func() error {
				time.Sleep(10 * time.Millisecond)
				return nil
			},
			expectError:  false,
			wantDuration: 10 * time.Millisecond,
		},
		{
			name:      "failed operation",
			operation: "FailOp",
			operationFn: func() error {
				time.Sleep(5 * time.Millisecond)
				return errors.New("operation failed")
			},
			expectError:  true,
			wantDuration: 5 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture log output
			var logBuf bytes.Buffer
			logger := NewRequestTimingLogger(false)
			logger.SetOutput(&logBuf)

			// Time operation
			start := time.Now()
			err := logger.TimeOperation(tt.operation, tt.operationFn)
			actualDuration := time.Since(start)

			// Verify error
			if tt.expectError && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Verify duration (allow some margin for execution overhead)
			if actualDuration < tt.wantDuration {
				t.Errorf("Duration = %v, want at least %v", actualDuration, tt.wantDuration)
			}

			// Verify logging occurred
			logOutput := logBuf.String()
			if !strings.Contains(logOutput, tt.operation) {
				t.Errorf("Log output missing operation name %q", tt.operation)
			}
		})
	}
}

func TestRequestTimingLogger_TimeOperationWithResult(t *testing.T) {
	tests := []struct {
		name        string
		operation   string
		resultValue interface{}
		resultError error
	}{
		{
			name:        "operation with string result",
			operation:   "GetString",
			resultValue: "test result",
			resultError: nil,
		},
		{
			name:        "operation with error",
			operation:   "GetError",
			resultValue: nil,
			resultError: errors.New("test error"),
		},
		{
			name:        "operation with int result",
			operation:   "GetInt",
			resultValue: 42,
			resultError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture log output
			var logBuf bytes.Buffer
			logger := NewRequestTimingLogger(false)
			logger.SetOutput(&logBuf)

			// Time operation with result
			result, timedResult := logger.TimeOperationWithResult(tt.operation, func() (interface{}, error) {
				time.Sleep(5 * time.Millisecond)
				return tt.resultValue, tt.resultError
			})

			// Verify result
			if result != tt.resultValue {
				t.Errorf("Result = %v, want %v", result, tt.resultValue)
			}

			// Verify timed result
			if timedResult == nil {
				t.Fatal("TimedResult is nil")
			}

			if timedResult.Error != tt.resultError {
				t.Errorf("TimedResult.Error = %v, want %v", timedResult.Error, tt.resultError)
			}

			if timedResult.Duration < 5*time.Millisecond {
				t.Errorf("TimedResult.Duration = %v, want at least 5ms", timedResult.Duration)
			}

			// Verify logging
			logOutput := logBuf.String()
			if !strings.Contains(logOutput, tt.operation) {
				t.Errorf("Log output missing operation name %q", tt.operation)
			}
		})
	}
}

func TestPrintTimingSummary(t *testing.T) {
	tests := []struct {
		name           string
		timings        map[string][]time.Duration
		wantOutput     []string
		wantEmptyLog   bool
	}{
		{
			name: "single endpoint with multiple requests",
			timings: map[string][]time.Duration{
				"/api/cart": {
					100 * time.Millisecond,
					200 * time.Millisecond,
					150 * time.Millisecond,
				},
			},
			wantOutput: []string{
				"Request Timing Summary",
				"/api/cart",
				"Requests: 3",
				"Average:",
				"Min:",
				"Max:",
				"Total:",
			},
			wantEmptyLog: false,
		},
		{
			name: "multiple endpoints",
			timings: map[string][]time.Duration{
				"/api/cart":     {100 * time.Millisecond, 200 * time.Millisecond},
				"/api/checkout": {300 * time.Millisecond},
			},
			wantOutput: []string{
				"Request Timing Summary",
				"/api/cart",
				"/api/checkout",
				"Requests:",
			},
			wantEmptyLog: false,
		},
		{
			name:         "empty timings map",
			timings:      map[string][]time.Duration{},
			wantOutput:   []string{},
			wantEmptyLog: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stderr output
			r, w, _ := os.Pipe()
			oldStderr := os.Stderr
			os.Stderr = w

			// Print summary
			PrintTimingSummary(tt.timings)

			// Restore stderr
			w.Close()
			os.Stderr = oldStderr

			// Read captured output
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			// Verify output
			if tt.wantEmptyLog {
				if output != "" {
					t.Errorf("Expected empty output for empty timings, got: %s", output)
				}
				return
			}

			for _, want := range tt.wantOutput {
				if !strings.Contains(output, want) {
					t.Errorf("Output missing %q\nGot: %s", want, output)
				}
			}
		})
	}
}

func TestClientWithTimingIntegration(t *testing.T) {
	// Test that NewClientWithOptions creates client with timing enabled
	client := NewClientWithOptions(true)

	if client == nil {
		t.Fatal("NewClientWithOptions returned nil")
	}

	if client.httpClient == nil {
		t.Error("httpClient should not be nil")
	}

	if client.timingLogger == nil {
		t.Error("timingLogger should not be nil")
	}

	if !client.enableTiming {
		t.Error("enableTiming should be true")
	}

	// Test SetTimingEnabled
	client.SetTimingEnabled(false)
	if client.enableTiming {
		t.Error("enableTiming should be false after SetTimingEnabled(false)")
	}
}

func TestTimingTransport_PerformanceCategories(t *testing.T) {
	// This test verifies that different response times are correctly categorized
	tests := []struct {
		name           string
		delay          time.Duration
		verbose        bool
		wantCategory   string
	}{
		{
			name:         "fast request",
			delay:        50 * time.Millisecond,
			verbose:      true,
			wantCategory: "FAST",
		},
		{
			name:         "normal request",
			delay:        300 * time.Millisecond,
			verbose:      true,
			wantCategory: "NORMAL",
		},
		{
			name:         "slow request",
			delay:        1 * time.Second,
			verbose:      true,
			wantCategory: "SLOW",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(tt.delay)
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			// Capture log output
			var logBuf bytes.Buffer
			transport := NewTimingTransport(tt.verbose)
			transport.Logger.SetOutput(&logBuf)

			// Make request
			client := &http.Client{Transport: transport}
			resp, err := client.Get(server.URL)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			// Verify performance category appears in logs
			logOutput := logBuf.String()
			if tt.verbose && !strings.Contains(logOutput, tt.wantCategory) {
				t.Errorf("Expected category %q in logs\nGot: %s", tt.wantCategory, logOutput)
			}
		})
	}
}
