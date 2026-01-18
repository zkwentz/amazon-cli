# Integration Tests for Orders Commands

This directory contains integration tests for the `amazon-cli` orders commands.

## Overview

The integration tests cover all orders-related commands:
- `amazon-cli orders list` - List recent orders with optional filters
- `amazon-cli orders get <order-id>` - Get detailed order information
- `amazon-cli orders track <order-id>` - Track shipment status
- `amazon-cli orders history` - Get order history by year

## Test Coverage

### Functional Tests
- **OrdersListCommand**: Tests listing orders with various parameters (limit, status filters)
- **OrdersGetCommand**: Tests retrieving specific order details by order ID
- **OrdersTrackCommand**: Tests tracking order shipments
- **OrdersHistoryCommand**: Tests retrieving order history by year

### Error Handling Tests
- Authentication errors (AUTH_REQUIRED, AUTH_EXPIRED)
- Not found errors (invalid order IDs)
- Network errors (connection failures)
- Invalid input errors (missing required parameters)

### Non-Functional Tests
- **Rate Limiting**: Verifies rate limiting is enforced between requests
- **JSON Output**: Ensures all commands output valid JSON
- **Exit Codes**: Validates correct exit codes per PRD specification

## Running the Tests

### Prerequisites
- Go 1.21 or higher
- Amazon CLI project structure set up

### Run All Integration Tests
```bash
cd tests/integration
go test -v
```

### Run Specific Test
```bash
go test -v -run TestOrdersListCommand
```

### Run with Coverage
```bash
go test -v -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Test Structure

Each test uses a mock HTTP server to simulate Amazon's API responses. This allows testing without:
- Making real API calls to Amazon
- Requiring valid authentication
- Being subject to rate limits

### Mock Server
The `createMockAmazonServer()` helper creates an HTTP test server that returns predefined JSON responses, simulating various scenarios:
- Successful responses with valid data
- Error responses (404, 401, 429, etc.)
- Edge cases (empty results, missing fields)

### Command Execution
The `executeCommand()` helper:
1. Creates a temporary config file with test credentials
2. Executes the CLI command with the mock server URL
3. Captures output and exit code
4. Cleans up temporary files

## Exit Codes Tested

Per PRD specifications:
- `0` - Success
- `2` - Invalid arguments
- `3` - Authentication error
- `4` - Network error
- `6` - Not found

## Expected JSON Output Schemas

### Orders List Response
```json
{
  "orders": [
    {
      "order_id": "string",
      "date": "string",
      "total": number,
      "status": "string",
      "items": []
    }
  ],
  "total_count": number
}
```

### Order Details Response
```json
{
  "order_id": "string",
  "date": "string",
  "total": number,
  "status": "string",
  "items": [...],
  "tracking": {
    "carrier": "string",
    "tracking_number": "string",
    "status": "string",
    "delivery_date": "string"
  }
}
```

### Error Response
```json
{
  "error": {
    "code": "string",
    "message": "string",
    "details": {}
  }
}
```

## CI/CD Integration

These tests are designed to run in CI/CD pipelines:
- No external dependencies required
- Fast execution (mock servers)
- Deterministic results
- Clear pass/fail criteria

### GitHub Actions Example
```yaml
- name: Run Integration Tests
  run: |
    cd tests/integration
    go test -v -race -coverprofile=coverage.out
    go tool cover -func=coverage.out
```

## Future Enhancements

- [ ] Add tests for pagination
- [ ] Add tests for concurrent requests
- [ ] Add tests for different output formats (table, raw)
- [ ] Add performance benchmarks
- [ ] Add tests for config file edge cases
- [ ] Add tests for verbose logging output

## Troubleshooting

### Tests Fail with "command not found"
Ensure you're running tests from the correct directory and `main.go` is present.

### Tests Timeout
Check that the mock server is properly starting and shutting down. Increase timeout if needed.

### JSON Parsing Errors
Verify the mock responses match the expected schema. Check for trailing commas or invalid JSON syntax.
