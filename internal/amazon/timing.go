package amazon

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// TimingTransport is an http.RoundTripper that logs request timing information
// It wraps the default transport to add performance debugging capabilities
type TimingTransport struct {
	Transport http.RoundTripper
	Logger    *log.Logger
	Verbose   bool
}

// NewTimingTransport creates a new TimingTransport with default settings
func NewTimingTransport(verbose bool) *TimingTransport {
	return &TimingTransport{
		Transport: http.DefaultTransport,
		Logger:    log.New(os.Stderr, "[HTTP] ", log.LstdFlags|log.Lmicroseconds),
		Verbose:   verbose,
	}
}

// RoundTrip implements the http.RoundTripper interface
// It wraps HTTP requests with timing measurements and logging
func (t *TimingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Record start time
	start := time.Now()

	// Log request start if verbose mode is enabled
	if t.Verbose {
		t.Logger.Printf("→ %s %s", req.Method, req.URL.String())
	}

	// Execute the actual HTTP request
	resp, err := t.Transport.RoundTrip(req)

	// Calculate duration
	duration := time.Since(start)

	// Log timing information
	if err != nil {
		t.Logger.Printf("✗ %s %s - Failed: %v (Duration: %v)",
			req.Method,
			req.URL.String(),
			err,
			duration)
		return nil, err
	}

	// Log successful request with timing
	t.logRequestTiming(req, resp, duration)

	return resp, nil
}

// logRequestTiming logs detailed timing information for successful requests
func (t *TimingTransport) logRequestTiming(req *http.Request, resp *http.Response, duration time.Duration) {
	// Basic timing log (always shown)
	statusIcon := "✓"
	if resp.StatusCode >= 400 {
		statusIcon = "✗"
	}

	t.Logger.Printf("%s %s %s - Status: %d (Duration: %v)",
		statusIcon,
		req.Method,
		req.URL.Path,
		resp.StatusCode,
		duration)

	// Verbose logging includes additional details
	if t.Verbose {
		// Log response headers if available
		contentLength := resp.Header.Get("Content-Length")
		contentType := resp.Header.Get("Content-Type")

		if contentLength != "" {
			t.Logger.Printf("  ↳ Response Size: %s bytes, Type: %s", contentLength, contentType)
		}

		// Categorize request speed for debugging
		var speedCategory string
		switch {
		case duration < 100*time.Millisecond:
			speedCategory = "FAST"
		case duration < 500*time.Millisecond:
			speedCategory = "NORMAL"
		case duration < 2*time.Second:
			speedCategory = "SLOW"
		default:
			speedCategory = "VERY SLOW"
		}

		t.Logger.Printf("  ↳ Performance: %s", speedCategory)
	}
}

// RequestTimingLogger provides structured logging for non-HTTP operations
type RequestTimingLogger struct {
	logger  *log.Logger
	verbose bool
}

// NewRequestTimingLogger creates a new RequestTimingLogger
func NewRequestTimingLogger(verbose bool) *RequestTimingLogger {
	return &RequestTimingLogger{
		logger:  log.New(os.Stderr, "[OPERATION] ", log.LstdFlags|log.Lmicroseconds),
		verbose: verbose,
	}
}

// SetOutput sets the output destination for the logger
func (l *RequestTimingLogger) SetOutput(w io.Writer) {
	l.logger.SetOutput(w)
}

// LogOperation logs the timing of a generic operation
func (l *RequestTimingLogger) LogOperation(operation string, duration time.Duration, err error) {
	if err != nil {
		l.logger.Printf("✗ %s - Failed: %v (Duration: %v)", operation, err, duration)
		return
	}

	l.logger.Printf("✓ %s - Success (Duration: %v)", operation, duration)

	if l.verbose {
		var speedCategory string
		switch {
		case duration < 10*time.Millisecond:
			speedCategory = "VERY FAST"
		case duration < 50*time.Millisecond:
			speedCategory = "FAST"
		case duration < 200*time.Millisecond:
			speedCategory = "NORMAL"
		case duration < 1*time.Second:
			speedCategory = "SLOW"
		default:
			speedCategory = "VERY SLOW"
		}

		l.logger.Printf("  ↳ Performance: %s", speedCategory)
	}
}

// TimeOperation wraps a function call with timing measurement
func (l *RequestTimingLogger) TimeOperation(operation string, fn func() error) error {
	start := time.Now()
	err := fn()
	duration := time.Since(start)
	l.LogOperation(operation, duration, err)
	return err
}

// TimedResult holds the result of a timed operation with a return value
type TimedResult struct {
	Duration time.Duration
	Error    error
}

// TimeOperationWithResult wraps a function that returns a value with timing measurement
func (l *RequestTimingLogger) TimeOperationWithResult(operation string, fn func() (interface{}, error)) (interface{}, *TimedResult) {
	start := time.Now()
	result, err := fn()
	duration := time.Since(start)

	l.LogOperation(operation, duration, err)

	return result, &TimedResult{
		Duration: duration,
		Error:    err,
	}
}

// PrintTimingSummary prints a summary of request timings
func PrintTimingSummary(timings map[string][]time.Duration) {
	if len(timings) == 0 {
		return
	}

	fmt.Fprintln(os.Stderr, "\n=== Request Timing Summary ===")

	for endpoint, durations := range timings {
		if len(durations) == 0 {
			continue
		}

		var total time.Duration
		min := durations[0]
		max := durations[0]

		for _, d := range durations {
			total += d
			if d < min {
				min = d
			}
			if d > max {
				max = d
			}
		}

		avg := total / time.Duration(len(durations))

		fmt.Fprintf(os.Stderr, "%s:\n", endpoint)
		fmt.Fprintf(os.Stderr, "  Requests: %d\n", len(durations))
		fmt.Fprintf(os.Stderr, "  Average:  %v\n", avg)
		fmt.Fprintf(os.Stderr, "  Min:      %v\n", min)
		fmt.Fprintf(os.Stderr, "  Max:      %v\n", max)
		fmt.Fprintf(os.Stderr, "  Total:    %v\n", total)
	}

	fmt.Fprintln(os.Stderr, "==============================")
}
