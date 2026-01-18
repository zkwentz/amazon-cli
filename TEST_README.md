# Integration Tests with Mock Server

This directory contains comprehensive integration tests for the Amazon CLI HTTP client foundation.

## Test Coverage

### HTTP Client Tests (`internal/amazon/client_test.go`)

The client integration tests use Go's `httptest` package to create mock HTTP servers that simulate various Amazon API behaviors:

1. **TestClientGet** - Verifies basic GET requests work correctly
   - Tests that headers are properly set (User-Agent, Accept, Accept-Language)
   - Validates response parsing

2. **TestClientRetryOn429** - Tests automatic retry on rate limiting
   - Mock server returns 429 (Too Many Requests) twice, then succeeds
   - Verifies client automatically retries with exponential backoff
   - Confirms 3 total attempts are made

3. **TestClientRetryOn503** - Tests automatic retry on service unavailable
   - Mock server returns 503 (Service Unavailable) once, then succeeds
   - Validates retry logic for transient failures

4. **TestClientMaxRetries** - Tests retry limit enforcement
   - Mock server always returns 429
   - Verifies client stops after MaxRetries attempts (3 by default)
   - Ensures client doesn't retry indefinitely

5. **TestClientNoRetryOn404** - Tests that non-retryable errors aren't retried
   - Mock server returns 404 (Not Found)
   - Confirms only one attempt is made
   - Validates that only 429/503 trigger retries

6. **TestClientRateLimiting** - Tests rate limiting delays
   - Makes multiple sequential requests
   - Verifies minimum delay (100ms in test) between requests
   - Confirms jitter is added (delay varies between requests)

7. **TestClientUserAgentRotation** - Tests user agent rotation
   - Makes 15 requests
   - Verifies that user agents rotate through all 10 available options
   - Confirms rotation cycles back to the first after 10 requests

8. **TestClientPostForm** - Tests POST requests with form data
   - Validates POST method is used
   - Checks Content-Type header is set correctly
   - Verifies form data is encoded in query parameters

9. **TestClientCookieJar** - Tests cookie persistence
   - First request receives a session cookie
   - Subsequent requests automatically send the cookie
   - Validates cookie jar functionality for maintaining sessions

10. **TestClientTimeout** - Tests request timeout handling
    - Mock server delays response beyond timeout
    - Verifies timeout error is returned
    - Ensures client doesn't hang indefinitely

11. **TestClientConcurrentRequests** - Tests thread safety
    - Launches 5 concurrent goroutines making requests
    - Verifies all requests complete successfully
    - Validates rate limiter works correctly with concurrent access

### Rate Limiter Tests (`internal/ratelimit/limiter_test.go`)

Tests for the rate limiting logic:

1. **TestNewRateLimiter** - Tests rate limiter initialization
2. **TestWaitFirstRequest** - Verifies first request doesn't wait (only adds jitter)
3. **TestWaitSubsequentRequests** - Tests MinDelayMs enforcement
4. **TestWaitMultipleRequests** - Tests cumulative delay over multiple requests
5. **TestWaitWithBackoff** - Tests exponential backoff (2^n seconds, capped at 60s)
6. **TestShouldRetry** - Tests retry decision logic for various status codes
7. **TestShouldRetryMaxRetries** - Tests max retry limit enforcement

### Config Tests (`internal/config/config_test.go`)

Tests for configuration management:

1. **TestDefaultConfig** - Validates default configuration values
2. **TestSaveAndLoadConfig** - Tests saving and loading config files
3. **TestLoadConfigNonExistent** - Tests handling of missing config files
4. **TestSaveConfigCreatesDirectory** - Tests directory creation
5. **TestConfigFilePermissions** - Validates file permissions (0600 for security)
6. **TestConfigDirectoryPermissions** - Validates directory permissions (0700)
7. **TestConfigJSONFormat** - Tests JSON formatting with indentation
8. **TestLoadConfigWithInvalidJSON** - Tests error handling for invalid JSON
9. **TestConfigRoundTrip** - Tests data integrity through save/load cycles

## Mock Server Architecture

The tests use Go's standard `httptest` package to create lightweight HTTP servers:

```go
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Mock server logic here
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status": "ok"}`))
}))
defer server.Close()
```

### Benefits of this approach:

1. **No external dependencies** - Uses only Go standard library
2. **Fast execution** - In-memory HTTP servers start instantly
3. **Deterministic** - Complete control over server responses
4. **Thread-safe** - Each test gets its own isolated server
5. **No network I/O** - Tests run entirely in-process

## Running the Tests

Run all tests:
```bash
go test ./...
```

Run tests with verbose output:
```bash
go test -v ./...
```

Run only client tests:
```bash
go test -v ./internal/amazon/
```

Run a specific test:
```bash
go test -v ./internal/amazon/ -run TestClientRetryOn429
```

Run tests with coverage:
```bash
go test -cover ./...
```

Generate coverage report:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

Run benchmarks:
```bash
go test -bench=. ./internal/ratelimit/
```

## Test Duration

Most tests run quickly (< 1 second), but some tests intentionally include delays:

- **TestClientRetryOn429** - ~3.4s (tests exponential backoff)
- **TestClientMaxRetries** - ~3.5s (tests retry limits with backoff)
- **TestWaitWithBackoff** - ~75s (tests all backoff levels including 60s max)

You can run faster tests by excluding the backoff test:
```bash
go test -v ./... -skip TestWaitWithBackoff
```

## Key Features Tested

### Retry Logic
- Automatic retry on 429 (rate limited) and 503 (service unavailable)
- Exponential backoff: 2^n seconds (1s, 2s, 4s, 8s, ...)
- Max backoff capped at 60 seconds
- Configurable max retries (default: 3)

### Rate Limiting
- Minimum delay between requests (default: 1000ms)
- Random jitter (0-500ms) to avoid thundering herd
- First request executes immediately
- Thread-safe for concurrent use

### HTTP Headers
- User-Agent rotation (10 different browser user agents)
- Accept, Accept-Language, Accept-Encoding headers
- Content-Type for POST requests

### Session Management
- Cookie jar for maintaining sessions
- Automatic cookie persistence across requests

### Security
- Config files saved with 0600 permissions (owner read/write only)
- Config directories created with 0700 permissions
- Timeout enforcement to prevent hanging requests

## Test Data

All tests use synthetic test data:
- Mock tokens: "test-access-token", "test-refresh-token"
- Mock IDs: "addr_123", "pay_456"
- Mock responses: `{"status": "ok"}`

No real Amazon credentials or data are used in tests.

## Continuous Integration

These tests are designed to run in CI/CD environments:
- No external dependencies required
- No network access needed
- Deterministic results
- Fast execution (total runtime ~90s including long backoff tests)

For faster CI, use:
```bash
go test -short ./...
```

Then add this to slow tests:
```go
if testing.Short() {
    t.Skip("Skipping slow test in short mode")
}
```

## Future Enhancements

Potential additions to test suite:
- Table-driven tests for more HTTP status codes
- Mock tests for Amazon-specific responses (HTML parsing)
- Integration tests with actual Amazon sandbox (if available)
- Load testing for high concurrency scenarios
- Memory profiling to detect leaks
