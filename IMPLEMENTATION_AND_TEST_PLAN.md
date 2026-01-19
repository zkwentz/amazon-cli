# Amazon CLI - Comprehensive Implementation and Test Plan

## Executive Summary

This document outlines a complete implementation and testing strategy for the amazon-cli project. It provides structured phases for development, comprehensive test coverage requirements, quality gates, and validation criteria to ensure the CLI works reliably for AI agent integration.

**Current State**: Foundation is in place with Cobra CLI framework, basic command structure, and placeholder implementations. Most Amazon API integration is mocked.

**Target State**: Fully functional CLI with real Amazon integration, comprehensive test coverage, CI/CD automation, and production-ready distribution.

---

## Table of Contents

1. [Implementation Plan](#implementation-plan)
2. [Test Strategy](#test-strategy)
3. [Quality Gates](#quality-gates)
4. [Validation Checklist](#validation-checklist)
5. [Risk Mitigation](#risk-mitigation)

---

## Implementation Plan

### Phase 1: Core Infrastructure Completion (Week 1)

#### 1.1 Amazon API Integration Research
**Goal**: Understand Amazon's actual API/web structure before implementing

**Tasks**:
- [ ] Research Amazon's Shopping API availability and access
- [ ] Document Amazon's authentication mechanisms (OAuth, cookie-based, etc.)
- [ ] Map Amazon web pages for scraping (orders, cart, search, etc.)
- [ ] Identify required HTTP headers, cookies, and CSRF tokens
- [ ] Test basic HTTP requests against Amazon to verify approach
- [ ] Document rate limiting thresholds through experimentation

**Deliverables**:
- `docs/amazon-api-research.md` - Complete documentation of Amazon's API/web structure
- `docs/authentication-strategy.md` - Chosen auth approach with rationale

**Testing**:
- Manual testing of HTTP requests against Amazon
- Verify authentication flow works with real Amazon account
- Confirm rate limiting strategy prevents blocks

---

#### 1.2 Authentication System Implementation
**Goal**: Replace mock auth with real Amazon authentication

**Tasks**:
- [ ] Implement chosen authentication method (OAuth or browser session)
- [ ] Create secure token/cookie storage in `~/.amazon-cli/config.json`
- [ ] Implement automatic token refresh logic
- [ ] Add session validation before API calls
- [ ] Handle 2FA scenarios gracefully
- [ ] Implement auth error recovery (redirect to login)

**Files to Modify**:
- `internal/amazon/auth.go` - Real auth implementation
- `cmd/auth.go` - Wire up real auth client
- `internal/config/config.go` - Secure credential storage

**Testing Requirements**:
- Unit tests for token refresh logic
- Integration tests for full OAuth flow
- Manual testing with real Amazon account
- Test 2FA handling
- Test expired token recovery
- Security audit of credential storage (file permissions)

**Test Coverage Target**: 85%+ for auth code

---

#### 1.3 HTTP Client & Rate Limiting Enhancement
**Goal**: Bulletproof HTTP client with proper error handling and rate limiting

**Tasks**:
- [ ] Implement retry logic with exponential backoff
- [ ] Add circuit breaker pattern for repeated failures
- [ ] Implement request/response logging (redact credentials)
- [ ] Add HTTP timeout configurations
- [ ] Implement User-Agent rotation
- [ ] Add CAPTCHA detection and handling
- [ ] Create mock HTTP server for testing

**Files to Create/Modify**:
- `internal/ratelimit/limiter.go` - Enhanced rate limiter
- `internal/amazon/client.go` - Production-ready HTTP client
- `internal/amazon/client_test.go` - HTTP client tests
- `internal/testutil/mock_server.go` - Mock Amazon server for tests

**Testing Requirements**:
- Unit tests for rate limiter (jitter, backoff calculations)
- Integration tests with mock server
- Test timeout scenarios
- Test retry logic with 429/503 responses
- Test circuit breaker triggering
- Load testing to verify rate limits

**Test Coverage Target**: 90%+ for HTTP client and rate limiter

---

### Phase 2: Orders & Returns Features (Week 2)

#### 2.1 Orders Implementation
**Goal**: Real Amazon orders retrieval and parsing

**Tasks**:
- [ ] Implement HTML parsing for Amazon order history page
- [ ] Create order detail page parser
- [ ] Implement tracking information extraction
- [ ] Handle pagination for order history
- [ ] Add filtering by status and date range
- [ ] Implement error handling for missing/malformed data

**Files to Modify**:
- `internal/amazon/orders.go` - Real implementation
- `pkg/models/order.go` - Ensure models match Amazon data
- `cmd/orders.go` - Wire up real client

**Testing Requirements**:
```go
// Test categories needed:
1. Parser Tests
   - Test with sample Amazon HTML (save fixtures)
   - Test with missing optional fields
   - Test with various order statuses
   - Test with multiple items per order

2. Integration Tests
   - Test full order list flow
   - Test order detail retrieval
   - Test tracking information
   - Test pagination

3. Error Cases
   - Test with network errors
   - Test with auth failures
   - Test with malformed HTML responses
   - Test with empty order history
```

**Test Coverage Target**: 80%+ for orders code

**Test Fixtures**:
- Save real Amazon HTML responses (anonymized) in `testdata/orders/`
- Create edge case fixtures (empty orders, canceled items, etc.)

---

#### 2.2 Returns Implementation
**Goal**: Real Amazon returns initiation and tracking

**Tasks**:
- [ ] Implement returnable items fetching
- [ ] Create return options parser
- [ ] Implement return initiation flow
- [ ] Add return label retrieval
- [ ] Implement return status tracking
- [ ] Add return reason validation

**Files to Modify**:
- `internal/amazon/returns.go` - Real implementation
- Create `pkg/models/return.go` - Return data models
- Create `cmd/returns.go` - Returns commands

**Testing Requirements**:
- Unit tests for return reason validation
- Parser tests with Amazon HTML fixtures
- Integration tests for return flow
- Test return window expiration
- Test items not eligible for return
- Test label generation

**Test Coverage Target**: 80%+ for returns code

---

### Phase 3: Search & Product Features (Week 3)

#### 3.1 Search Implementation
**Goal**: Amazon product search with filtering

**Tasks**:
- [ ] Implement search query builder
- [ ] Create search results parser
- [ ] Add price range filtering
- [ ] Implement category filtering
- [ ] Add Prime-only filtering
- [ ] Handle pagination for search results

**Files to Modify**:
- `internal/amazon/search.go` - Real implementation
- `pkg/models/product.go` - Product data models

**Testing Requirements**:
- Parser tests with search result HTML
- Test various search queries
- Test all filter combinations
- Test pagination
- Test zero results
- Test product with missing data (price, rating, etc.)

**Test Coverage Target**: 80%+

---

#### 3.2 Product Details Implementation
**Goal**: Detailed product information retrieval

**Tasks**:
- [ ] Implement product detail page parser
- [ ] Extract product images
- [ ] Parse product features and description
- [ ] Implement reviews retrieval
- [ ] Handle products with variants (size, color, etc.)
- [ ] Add availability checking

**Files to Modify**:
- `internal/amazon/product.go` - Real implementation

**Testing Requirements**:
- Parser tests with product detail HTML
- Test products with all optional fields
- Test products missing optional data
- Test variant products
- Test out-of-stock products
- Test review parsing with various ratings

**Test Coverage Target**: 80%+

---

### Phase 4: Cart & Checkout (Week 4)

#### 4.1 Cart Operations Implementation
**Goal**: Full cart management functionality

**Tasks**:
- [ ] Implement add to cart (handle CSRF tokens)
- [ ] Create cart retrieval and parsing
- [ ] Implement remove from cart
- [ ] Add update quantity functionality
- [ ] Implement clear cart
- [ ] Handle cart errors (item unavailable, quantity limits)

**Files to Modify**:
- `internal/amazon/cart.go` - Replace mock with real implementation
- `internal/amazon/cart_test.go` - Expand tests for real implementation

**Testing Requirements**:
```go
// Expand existing cart_test.go with:
1. Real API Integration Tests (optional, can use test account)
   - Test add valid item
   - Test add invalid ASIN
   - Test add out-of-stock item
   - Test quantity limits
   - Test cart retrieval
   - Test remove item
   - Test clear cart

2. Parser Tests
   - Parse cart HTML with fixtures
   - Handle missing prices
   - Handle promotional discounts
   - Calculate totals correctly

3. Error Handling
   - Session expiration during cart operation
   - Network failures mid-operation
   - CSRF token refresh
```

**Test Coverage Target**: 85%+ (cart is critical)

---

#### 4.2 Checkout Implementation
**Goal**: Secure checkout flow with safety guards

**Tasks**:
- [ ] Implement checkout preview (dry run)
- [ ] Create address/payment selection
- [ ] Implement order submission
- [ ] Add order confirmation parsing
- [ ] Implement --confirm flag validation
- [ ] Add pre-checkout validation (stock, prices)

**Files to Modify**:
- `internal/amazon/cart.go` - Add checkout methods
- `cmd/cart.go` - Enhance checkout command

**Critical Testing Requirements**:
```go
// CRITICAL: Checkout must be thoroughly tested
// Use test Amazon account or mock server only!

1. Preview Mode Tests (Safe - No Real Purchases)
   - Test checkout preview without --confirm
   - Verify dry_run flag in output
   - Test address selection
   - Test payment method selection
   - Test delivery option parsing

2. Validation Tests
   - Verify --confirm is required
   - Test empty cart rejection
   - Test out-of-stock item detection
   - Test price change detection

3. Mock Checkout Tests
   - Use mock server for actual purchase flow
   - Test order confirmation parsing
   - Test order ID extraction
   - Test success/failure scenarios

4. Safety Tests
   - Verify commands fail without --confirm
   - Test that preview doesn't create order
   - Verify error messages are clear
```

**Test Coverage Target**: 95%+ (critical functionality)

**IMPORTANT**: Never test real checkout against production Amazon without explicit safeguards!

---

### Phase 5: Subscriptions (Week 5)

#### 5.1 Subscribe & Save Implementation
**Goal**: Subscription management features

**Tasks**:
- [ ] Create `pkg/models/subscription.go` - Data models
- [ ] Create `internal/amazon/subscriptions.go` - API client
- [ ] Implement subscription listing
- [ ] Add skip delivery functionality
- [ ] Implement frequency changes
- [ ] Add subscription cancellation
- [ ] Create `cmd/subscriptions.go` - Commands

**Testing Requirements**:
- Parser tests with subscription page HTML
- Test frequency validation (1-26 weeks)
- Test skip delivery with --confirm
- Test cancellation flow
- Test upcoming deliveries sorting

**Test Coverage Target**: 75%+

---

### Phase 6: Error Handling & Polish (Week 6)

#### 6.1 Comprehensive Error Handling
**Goal**: Consistent error handling across all features

**Tasks**:
- [ ] Audit all error paths for proper error codes
- [ ] Implement error context wrapping
- [ ] Add user-friendly error messages
- [ ] Create error recovery suggestions
- [ ] Implement verbose error logging
- [ ] Add CAPTCHA/login redirect detection

**Files to Modify**:
- `pkg/models/errors.go` - Enhanced error types
- All `internal/amazon/*.go` - Consistent error handling
- `internal/output/output.go` - Error formatting

**Testing Requirements**:
```go
// Error scenario tests for each feature:
1. Network Errors
   - Timeout scenarios
   - Connection refused
   - DNS failures

2. Amazon Errors
   - 401 (auth required)
   - 403 (forbidden/CAPTCHA)
   - 404 (not found)
   - 429 (rate limited)
   - 500/503 (server errors)

3. Validation Errors
   - Invalid ASIN format
   - Invalid order ID
   - Empty required fields
   - Out of range values

4. State Errors
   - Empty cart checkout
   - Expired return window
   - Canceled order tracking
```

**Test Coverage Target**: 85%+ for error handling paths

---

#### 6.2 Input Validation
**Goal**: Validate all user inputs before API calls

**Tasks**:
- [ ] Create validation utilities in `internal/validation/`
- [ ] Implement ASIN format validation
- [ ] Add order ID format validation
- [ ] Implement price range validation
- [ ] Add quantity validation
- [ ] Create subscription ID validation

**Files to Create**:
- `internal/validation/validators.go`
- `internal/validation/validators_test.go`

**Testing Requirements**:
- Unit tests for each validator
- Test valid and invalid inputs
- Test edge cases (empty, too long, special chars)
- Benchmark validators for performance

**Test Coverage Target**: 100% (validators must be bulletproof)

---

### Phase 7: Integration Testing (Week 7)

#### 7.1 End-to-End Test Suite
**Goal**: Full workflow testing from command to output

**Tasks**:
- [ ] Create `test/e2e/` directory for integration tests
- [ ] Build mock Amazon server with realistic responses
- [ ] Create test scenarios for each workflow
- [ ] Implement test fixtures for all response types
- [ ] Add performance benchmarks

**Test Scenarios**:
```bash
# E2E test workflows to implement:

1. Authentication Flow
   - Login -> Verify Status -> Logout

2. Order Management Flow
   - Login -> List Orders -> Get Order Detail -> Track Order

3. Search to Purchase Flow
   - Search Product -> Get Product Details -> Add to Cart ->
   -> View Cart -> Preview Checkout (no confirm)

4. Return Flow
   - List Returnable Items -> Get Return Options ->
   -> Create Return (with confirm) -> Check Status

5. Subscription Management Flow
   - List Subscriptions -> Get Subscription -> Skip Delivery (with confirm)

6. Error Recovery Flow
   - Expired Token -> Retry -> Auto Refresh -> Success

7. Rate Limiting Flow
   - Rapid requests -> Rate limited -> Backoff -> Success
```

**Implementation**:
```go
// test/e2e/e2e_test.go structure:

package e2e_test

import (
    "testing"
    "github.com/zkwentz/amazon-cli/internal/testutil"
)

func TestE2E_CompleteOrderFlow(t *testing.T) {
    // Setup mock Amazon server
    server := testutil.NewMockAmazonServer()
    defer server.Close()

    // Configure CLI to use mock server
    config := testutil.NewTestConfig(server.URL)

    // Run complete order workflow
    // ... test implementation
}
```

**Test Coverage**: All major workflows covered

---

#### 7.2 Performance Testing
**Goal**: Ensure CLI meets performance requirements

**Tasks**:
- [ ] Create performance benchmarks
- [ ] Test response times for all commands
- [ ] Measure memory usage
- [ ] Test with large datasets (1000+ orders)
- [ ] Benchmark JSON parsing performance

**Performance Targets**:
- Command execution: < 5 seconds (typical)
- JSON parsing: < 100ms for typical responses
- Memory usage: < 50MB for typical operations
- Rate limiter overhead: < 10ms per request

**Benchmarks to Create**:
```go
// internal/amazon/orders_test.go

func BenchmarkParseOrders_10(b *testing.B)    { benchParseOrders(b, 10) }
func BenchmarkParseOrders_100(b *testing.B)   { benchParseOrders(b, 100) }
func BenchmarkParseOrders_1000(b *testing.B)  { benchParseOrders(b, 1000) }

func BenchmarkSearchResults(b *testing.B) { /* ... */ }
func BenchmarkCartOperations(b *testing.B) { /* ... */ }
```

---

### Phase 8: CI/CD & Automation (Week 8)

#### 8.1 GitHub Actions Setup
**Goal**: Automated testing and releases

**Tasks**:
- [ ] Create `.github/workflows/ci.yml` - Run tests on every push/PR
- [ ] Create `.github/workflows/release.yml` - Automated releases on tags
- [ ] Create `.github/workflows/lint.yml` - Code quality checks
- [ ] Add test coverage reporting (codecov.io)
- [ ] Create PR templates with testing checklist

**CI Workflow**:
```yaml
# .github/workflows/ci.yml

name: CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go-version: [1.25.x]

    runs-on: ${{ matrix.os }}

    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go-version }}

      - name: Run tests
        run: |
          go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          file: ./coverage.txt

      - name: Run benchmarks
        run: go test -bench=. -benchmem ./...

  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: golangci/golangci-lint-action@v4
        with:
          version: latest
```

**Release Workflow**:
```yaml
# .github/workflows/release.yml

name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: 1.25.x

      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v5
        with:
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          HOMEBREW_TAP_GITHUB_TOKEN: ${{ secrets.HOMEBREW_TAP_GITHUB_TOKEN }}
```

---

#### 8.2 Quality Gates
**Goal**: Prevent regressions and maintain quality

**Tasks**:
- [ ] Require 80%+ test coverage for PRs
- [ ] Require all tests passing before merge
- [ ] Require linter passing (zero warnings)
- [ ] Add CODEOWNERS for review requirements
- [ ] Create pre-commit hooks for local testing

**Pre-commit Hook**:
```bash
#!/bin/bash
# .git/hooks/pre-commit

echo "Running pre-commit checks..."

# Run tests
go test ./... || exit 1

# Run linter
golangci-lint run || exit 1

# Check formatting
test -z "$(gofmt -l .)" || {
    echo "Go files must be formatted with gofmt"
    exit 1
}

echo "All pre-commit checks passed!"
```

---

### Phase 9: Documentation & Release (Week 9)

#### 9.1 Documentation Completion
**Goal**: Comprehensive user and developer documentation

**Tasks**:
- [ ] Complete README.md with all features
- [ ] Create CONTRIBUTING.md with development guide
- [ ] Create CHANGELOG.md for release notes
- [ ] Write API documentation (if exposing Go packages)
- [ ] Create troubleshooting guide
- [ ] Document common error scenarios and fixes

**Documentation Files to Create**:
```
docs/
├── DEVELOPMENT.md          # How to build and test locally
├── TROUBLESHOOTING.md      # Common issues and solutions
├── API.md                  # Go package documentation
├── SECURITY.md             # Security considerations
└── examples/
    ├── search-workflow.md
    ├── order-workflow.md
    └── return-workflow.md
```

---

#### 9.2 skills.md for ClawdHub
**Goal**: AI agent integration documentation

**Tasks**:
- [ ] Create `skills.md` with proper metadata
- [ ] Document all commands with input/output schemas
- [ ] Add AI agent usage examples
- [ ] Document error handling for agents
- [ ] Add safety guidelines for purchase operations
- [ ] Include JSON schema definitions

**Skills.md Structure**:
```markdown
---
name: amazon-cli
description: CLI tool for Amazon shopping automation
version: 1.0.0
author: zkwentz
repository: https://github.com/zkwentz/amazon-cli
tags: [shopping, e-commerce, amazon, automation]
---

# Amazon CLI Skill

## Overview
[Description for AI agents]

## Installation
[Installation instructions]

## Actions

### search
Search for products on Amazon

**Inputs:**
- query (required): Search query string
- category (optional): Product category
- min_price (optional): Minimum price filter
- max_price (optional): Maximum price filter
- prime_only (optional): Only show Prime items

**Output Schema:**
\```json
{
  "query": "string",
  "results": [
    {
      "asin": "string",
      "title": "string",
      "price": 0.00,
      ...
    }
  ]
}
\```

**Example:**
\```bash
amazon-cli search "wireless headphones" --prime-only --max-price 100
\```

[... continue for all commands ...]
```

---

#### 9.3 Release Preparation
**Goal**: Production-ready release

**Tasks**:
- [ ] Version all packages to v1.0.0
- [ ] Create release notes in CHANGELOG.md
- [ ] Build and test binaries for all platforms
- [ ] Test Homebrew formula installation
- [ ] Create demo video/GIF for README
- [ ] Update all documentation with v1.0.0 references
- [ ] Create GitHub release with binaries

**Release Checklist**:
```markdown
## v1.0.0 Release Checklist

- [ ] All tests passing (100% of test suite)
- [ ] Test coverage >= 80%
- [ ] All documentation complete
- [ ] CHANGELOG.md updated
- [ ] Version bumped in all files
- [ ] Binaries built for: darwin/amd64, darwin/arm64, linux/amd64, linux/arm64, windows/amd64
- [ ] Homebrew formula tested
- [ ] skills.md validated
- [ ] Security audit completed
- [ ] Performance benchmarks run
- [ ] Example commands tested
- [ ] README.md screenshots/demos added
```

---

## Test Strategy

### Test Pyramid

```
                    /\
                   /  \
                  / E2E \ (10% - Full workflows)
                 /______\
                /        \
               / Integr.  \ (30% - Component integration)
              /____________\
             /              \
            /  Unit Tests    \ (60% - Individual functions)
           /___________________\
```

### Test Coverage Targets by Component

| Component | Target Coverage | Priority | Notes |
|-----------|----------------|----------|-------|
| `internal/amazon/auth.go` | 85% | Critical | Security-sensitive |
| `internal/amazon/cart.go` | 90% | Critical | Purchase operations |
| `internal/amazon/orders.go` | 80% | High | Core feature |
| `internal/amazon/search.go` | 80% | High | Core feature |
| `internal/amazon/product.go` | 80% | High | Core feature |
| `internal/amazon/returns.go` | 80% | Medium | Important feature |
| `internal/amazon/subscriptions.go` | 75% | Medium | Nice to have |
| `internal/ratelimit/` | 90% | Critical | Prevents bans |
| `internal/validation/` | 100% | Critical | Input safety |
| `internal/output/` | 85% | High | User-facing |
| `internal/config/` | 85% | High | Data integrity |
| `cmd/*.go` | 70% | Medium | Integration tested |
| `pkg/models/*.go` | 60% | Low | Mostly structs |

**Overall Target**: 80% code coverage

---

### Test Categories

#### 1. Unit Tests
**Purpose**: Test individual functions in isolation

**Guidelines**:
- Test happy path and error cases
- Mock external dependencies (HTTP, file I/O)
- Fast execution (< 1s for entire suite)
- No network calls
- Deterministic results

**Example**:
```go
// internal/validation/validators_test.go

func TestValidateASIN(t *testing.T) {
    tests := []struct {
        name    string
        asin    string
        wantErr bool
    }{
        {"valid ASIN", "B08N5WRWNW", false},
        {"too short", "B08N5", true},
        {"too long", "B08N5WRWNWEXTRA", true},
        {"empty", "", true},
        {"special chars", "B08N5@RWNW", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateASIN(tt.asin)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateASIN() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

---

#### 2. Integration Tests
**Purpose**: Test component interactions

**Guidelines**:
- Use mock HTTP server for Amazon responses
- Test with real config files (temporary test dirs)
- Test error propagation between layers
- Slower execution acceptable (< 10s total)

**Example**:
```go
// internal/amazon/orders_test.go

func TestGetOrders_Integration(t *testing.T) {
    // Setup mock server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Serve fixture HTML
        html, _ := os.ReadFile("testdata/orders/order_list.html")
        w.Write(html)
    }))
    defer server.Close()

    // Create client pointed at mock server
    client := NewClient()
    client.baseURL = server.URL

    // Test
    orders, err := client.GetOrders(10, "")
    if err != nil {
        t.Fatalf("GetOrders() error = %v", err)
    }

    if len(orders.Orders) == 0 {
        t.Error("Expected orders, got none")
    }
}
```

---

#### 3. End-to-End Tests
**Purpose**: Test complete user workflows

**Guidelines**:
- Test CLI binary directly (not Go code)
- Use mock server or test Amazon account
- Test JSON output parsing
- Test exit codes
- Slow execution acceptable (< 60s total)

**Example**:
```go
// test/e2e/search_test.go

func TestE2E_SearchWorkflow(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping E2E test in short mode")
    }

    // Build CLI binary
    cmd := exec.Command("go", "build", "-o", "amazon-cli-test", ".")
    if err := cmd.Run(); err != nil {
        t.Fatalf("Failed to build: %v", err)
    }
    defer os.Remove("amazon-cli-test")

    // Run search command
    cmd = exec.Command("./amazon-cli-test", "search", "headphones", "--output", "json")
    output, err := cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("Command failed: %v\nOutput: %s", err, output)
    }

    // Parse JSON output
    var result map[string]interface{}
    if err := json.Unmarshal(output, &result); err != nil {
        t.Fatalf("Invalid JSON output: %v", err)
    }

    // Validate structure
    if _, ok := result["query"]; !ok {
        t.Error("Missing 'query' field in output")
    }
}
```

---

#### 4. Table-Driven Tests
**Best Practice**: Use for comprehensive scenario coverage

**Example**:
```go
func TestAddToCart(t *testing.T) {
    tests := []struct {
        name        string
        asin        string
        quantity    int
        setupMock   func(*httptest.Server)
        wantErr     bool
        errContains string
        validate    func(*testing.T, *models.Cart)
    }{
        {
            name:     "valid item",
            asin:     "B08N5WRWNW",
            quantity: 1,
            setupMock: func(s *httptest.Server) {
                // Configure mock response
            },
            wantErr: false,
            validate: func(t *testing.T, cart *models.Cart) {
                if cart.ItemCount != 1 {
                    t.Errorf("Expected 1 item, got %d", cart.ItemCount)
                }
            },
        },
        // ... more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // ... test implementation
        })
    }
}
```

---

### Test Fixtures

**Directory Structure**:
```
testdata/
├── auth/
│   ├── login_success.html
│   ├── login_failed.html
│   └── oauth_response.json
├── orders/
│   ├── order_list.html
│   ├── order_detail.html
│   ├── empty_orders.html
│   └── order_tracking.html
├── cart/
│   ├── cart_empty.html
│   ├── cart_with_items.html
│   └── checkout_preview.html
├── search/
│   ├── search_results.html
│   ├── no_results.html
│   └── filtered_results.html
└── products/
    ├── product_detail.html
    ├── out_of_stock.html
    └── reviews.html
```

**Fixture Guidelines**:
- Save real Amazon HTML (anonymized - remove personal data)
- Create edge case fixtures (empty, malformed, etc.)
- Version fixtures (Amazon changes layout)
- Document what each fixture tests

---

### Mock Server Implementation

**Create reusable mock server**:
```go
// internal/testutil/mock_amazon.go

package testutil

import (
    "net/http"
    "net/http/httptest"
    "os"
    "path/filepath"
)

type MockAmazonServer struct {
    *httptest.Server
    Fixtures map[string]string
}

func NewMockAmazonServer() *MockAmazonServer {
    m := &MockAmazonServer{
        Fixtures: make(map[string]string),
    }

    m.Server = httptest.NewServer(http.HandlerFunc(m.handler))
    return m
}

func (m *MockAmazonServer) handler(w http.ResponseWriter, r *http.Request) {
    // Route based on path
    switch r.URL.Path {
    case "/orders":
        m.serveFixture(w, "testdata/orders/order_list.html")
    case "/cart":
        m.serveFixture(w, "testdata/cart/cart_with_items.html")
    default:
        http.NotFound(w, r)
    }
}

func (m *MockAmazonServer) serveFixture(w http.ResponseWriter, path string) {
    data, err := os.ReadFile(path)
    if err != nil {
        http.Error(w, "Fixture not found", 404)
        return
    }
    w.Write(data)
}

func (m *MockAmazonServer) WithFixture(path, fixture string) *MockAmazonServer {
    m.Fixtures[path] = fixture
    return m
}
```

**Usage in tests**:
```go
func TestOrderRetrieval(t *testing.T) {
    server := testutil.NewMockAmazonServer()
    defer server.Close()

    client := amazon.NewClient()
    client.BaseURL = server.URL

    orders, err := client.GetOrders(10, "")
    // ... assertions
}
```

---

## Quality Gates

### Pre-Commit Checks
- [ ] All tests pass (`go test ./...`)
- [ ] Code formatted (`gofmt -w .`)
- [ ] No linter errors (`golangci-lint run`)
- [ ] Test coverage maintained (>= 80%)

### Pull Request Requirements
- [ ] All CI checks pass
- [ ] Code review approved
- [ ] Test coverage >= 80%
- [ ] No decrease in coverage from main
- [ ] Documentation updated for new features
- [ ] CHANGELOG.md updated

### Release Requirements
- [ ] All tests passing on main branch
- [ ] Coverage >= 80% overall
- [ ] All documentation complete
- [ ] CHANGELOG.md updated
- [ ] Version bumped appropriately
- [ ] No known critical bugs
- [ ] Security audit passed
- [ ] Performance benchmarks acceptable

---

## Validation Checklist

### Functional Validation

#### Authentication
- [ ] OAuth login flow works with real Amazon account
- [ ] Tokens stored securely (file permissions 0600)
- [ ] Automatic token refresh works
- [ ] Auth status command shows correct information
- [ ] Logout clears credentials
- [ ] 2FA scenarios handled gracefully

#### Orders
- [ ] Can list recent orders
- [ ] Can get order details
- [ ] Can track shipments
- [ ] Order history retrieval works
- [ ] Filtering by status works
- [ ] Pagination works for large order histories
- [ ] JSON output matches documented schema

#### Returns
- [ ] Can list returnable items
- [ ] Can get return options
- [ ] Can initiate return with --confirm
- [ ] Return without --confirm shows preview
- [ ] Return label retrieval works
- [ ] Return status tracking works
- [ ] Invalid reason codes rejected

#### Search & Products
- [ ] Search returns relevant results
- [ ] Price filtering works
- [ ] Category filtering works
- [ ] Prime-only filter works
- [ ] Product details retrieval works
- [ ] Reviews retrieval works
- [ ] ASIN validation works

#### Cart & Checkout
- [ ] Can add items to cart
- [ ] Can view cart contents
- [ ] Can remove items from cart
- [ ] Can clear cart with --confirm
- [ ] Checkout preview works (without --confirm)
- [ ] Checkout requires --confirm flag
- [ ] **CRITICAL**: Checkout without --confirm never places order
- [ ] Cart totals calculated correctly

#### Subscriptions
- [ ] Can list subscriptions
- [ ] Can get subscription details
- [ ] Can skip delivery with --confirm
- [ ] Can change frequency with --confirm
- [ ] Can cancel subscription with --confirm
- [ ] Upcoming deliveries sorted correctly

---

### Non-Functional Validation

#### Performance
- [ ] Commands respond in < 5 seconds (typical case)
- [ ] Large datasets (1000+ orders) handled efficiently
- [ ] Memory usage < 50MB for typical operations
- [ ] JSON parsing < 100ms

#### Security
- [ ] Credentials never logged (even in verbose mode)
- [ ] Config file has restrictive permissions (0600)
- [ ] All HTTP requests use HTTPS
- [ ] No hardcoded secrets in code
- [ ] Input validation prevents injection attacks
- [ ] Purchase operations require explicit confirmation

#### Reliability
- [ ] Rate limiting prevents Amazon blocks
- [ ] Retry logic handles transient failures
- [ ] Circuit breaker prevents cascading failures
- [ ] Error messages are clear and actionable
- [ ] Graceful handling of CAPTCHA/login redirects

#### Usability
- [ ] Help text clear for all commands
- [ ] Error messages suggest fixes
- [ ] JSON output consistently formatted
- [ ] Exit codes match documentation
- [ ] --verbose provides useful debug info

#### Compatibility
- [ ] Works on macOS (arm64 and amd64)
- [ ] Works on Linux (arm64 and amd64)
- [ ] Works on Windows (amd64)
- [ ] Homebrew installation works
- [ ] Binary installation works

---

## Risk Mitigation

### High-Risk Areas

#### 1. Amazon API Changes
**Risk**: Amazon changes their website structure, breaking parsers

**Mitigation**:
- Version test fixtures with timestamps
- Monitor for parsing errors in production
- Implement parser fallbacks for common changes
- Alert on high failure rates
- Maintain parser abstraction layer

**Recovery Plan**:
- Update parsers based on new HTML structure
- Update test fixtures
- Release patch version
- Communicate to users via release notes

---

#### 2. Account Bans
**Risk**: Aggressive scraping leads to Amazon blocking accounts

**Mitigation**:
- Conservative rate limiting (1-2 sec between requests)
- Randomized jitter
- User-Agent rotation
- Respect robots.txt
- Circuit breaker on repeated 429s
- Documentation warning users about risks

**Recovery Plan**:
- Provide instructions for appealing bans
- Adjust rate limits if needed
- Add opt-in "slow mode" with longer delays

---

#### 3. Accidental Purchases
**Risk**: Bugs in checkout code cause unintended purchases

**Mitigation**:
- **NEVER test real checkout against production**
- --confirm flag strictly required
- Comprehensive preview mode
- Clear dry_run indicators in output
- Integration tests with mock server only
- Code review required for checkout changes
- Documentation warnings about --confirm

**Recovery Plan**:
- Provide clear instructions for order cancellation
- Document how to contact Amazon support
- Implement checkout rollback if possible

---

#### 4. Credential Theft
**Risk**: Config file with credentials compromised

**Mitigation**:
- Restrictive file permissions (0600)
- Documentation on credential security
- Credential rotation instructions
- Consider encrypted storage in future
- Never commit config files to git

**Recovery Plan**:
- Provide instructions for revoking tokens
- Force re-authentication
- Alert users to change Amazon password

---

#### 5. CAPTCHA Blocking
**Risk**: Amazon presents CAPTCHAs, blocking automation

**Mitigation**:
- Detect CAPTCHA responses
- Clear error message to user
- Suggest manual browser login
- Consider headless browser for auth only
- Document workarounds

**Recovery Plan**:
- Provide alternative auth methods
- Manual cookie extraction guide
- Consider browser extension for auth

---

## Testing Tools & Libraries

### Required Tools
```bash
# Install testing tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install gotest.tools/gotestsum@latest

# Install test coverage tools
go install github.com/axw/gocov/gocov@latest
go install github.com/matm/gocov-html/cmd/gocov-html@latest
```

### Recommended Libraries
```go
// go.mod additions for testing

require (
    github.com/stretchr/testify v1.8.4      // Assertions
    github.com/PuerkitoBio/goquery v1.8.1   // HTML parsing for tests
    github.com/jarcoal/httpmock v1.3.1      // HTTP mocking
    github.com/google/go-cmp v0.6.0         // Deep comparison
)
```

---

## Test Execution Commands

### Run All Tests
```bash
# Run all tests with coverage
go test -v -race -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out

# Run with verbose output
go test -v ./...

# Run specific package
go test -v ./internal/amazon

# Run specific test
go test -v -run TestAddToCart ./internal/amazon
```

### Run Tests by Category
```bash
# Unit tests only (fast)
go test -short ./...

# Integration tests (exclude E2E)
go test -run Integration ./...

# E2E tests only
go test ./test/e2e/...

# Run benchmarks
go test -bench=. -benchmem ./...
```

### Coverage Analysis
```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total

# Find untested code
go tool cover -html=coverage.out

# Coverage by package
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | sort -k3 -n
```

---

## Success Criteria

### Must Have (v1.0.0)
- [ ] 80%+ test coverage overall
- [ ] All core features working (auth, orders, cart, search)
- [ ] Zero critical bugs
- [ ] CI/CD pipeline operational
- [ ] Documentation complete
- [ ] Works on macOS, Linux, Windows
- [ ] Homebrew distribution working
- [ ] skills.md published to ClawdHub

### Nice to Have (v1.1.0+)
- [ ] Returns management
- [ ] Subscriptions management
- [ ] 90%+ test coverage
- [ ] Performance optimizations
- [ ] Table output format
- [ ] Bash completion scripts

---

## Appendix: Test Examples

### Example 1: Parser Test with Fixtures
```go
// internal/amazon/orders_test.go

func TestParseOrders(t *testing.T) {
    // Load fixture
    html, err := os.ReadFile("testdata/orders/order_list.html")
    if err != nil {
        t.Fatalf("Failed to load fixture: %v", err)
    }

    // Parse
    doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(html))
    orders := parseOrders(doc)

    // Assertions
    if len(orders) != 3 {
        t.Errorf("Expected 3 orders, got %d", len(orders))
    }

    // Validate first order structure
    order := orders[0]
    if order.OrderID == "" {
        t.Error("Order ID should not be empty")
    }
    if order.Total <= 0 {
        t.Error("Order total should be positive")
    }
}
```

### Example 2: Mock Server Integration Test
```go
// internal/amazon/cart_test.go

func TestAddToCart_Integration(t *testing.T) {
    // Setup mock server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/cart/add" {
            // Simulate successful add
            w.Header().Set("Content-Type", "text/html")
            fixture, _ := os.ReadFile("testdata/cart/add_success.html")
            w.Write(fixture)
        }
    }))
    defer server.Close()

    // Create client
    client := NewClient()
    client.baseURL = server.URL

    // Test
    cart, err := client.AddToCart("B08N5WRWNW", 1)

    // Assert
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
    if cart.ItemCount != 1 {
        t.Errorf("Expected 1 item, got %d", cart.ItemCount)
    }
}
```

### Example 3: E2E CLI Test
```go
// test/e2e/cli_test.go

func TestCLI_SearchCommand(t *testing.T) {
    // Build binary
    buildCmd := exec.Command("go", "build", "-o", "amazon-cli-test", "./...")
    if err := buildCmd.Run(); err != nil {
        t.Fatal(err)
    }
    defer os.Remove("amazon-cli-test")

    // Run command
    cmd := exec.Command("./amazon-cli-test", "search", "headphones", "-o", "json")
    output, err := cmd.Output()

    // Check exit code
    if err != nil {
        if exitErr, ok := err.(*exec.ExitError); ok {
            t.Fatalf("Command exited with code %d", exitErr.ExitCode())
        }
        t.Fatal(err)
    }

    // Validate JSON
    var result map[string]interface{}
    if err := json.Unmarshal(output, &result); err != nil {
        t.Fatalf("Invalid JSON: %v\nOutput: %s", err, output)
    }

    // Check required fields
    if _, ok := result["query"]; !ok {
        t.Error("Missing 'query' field")
    }
    if _, ok := result["results"]; !ok {
        t.Error("Missing 'results' field")
    }
}
```

---

## Timeline Summary

| Week | Phase | Deliverables | Test Coverage Target |
|------|-------|--------------|---------------------|
| 1 | Core Infrastructure | Auth, HTTP client, Rate limiter | 85% |
| 2 | Orders & Returns | Orders API, Returns API | 80% |
| 3 | Search & Products | Search, Product details, Reviews | 80% |
| 4 | Cart & Checkout | Cart ops, Checkout flow | 90% |
| 5 | Subscriptions | Subscription management | 75% |
| 6 | Error Handling | Validation, Error handling | 85% |
| 7 | Integration Tests | E2E tests, Performance tests | N/A |
| 8 | CI/CD | GitHub Actions, Quality gates | N/A |
| 9 | Documentation | README, skills.md, Release | N/A |

**Total**: 9 weeks to production-ready v1.0.0

---

## Final Notes

This plan provides a structured approach to building a production-ready Amazon CLI with comprehensive testing. The key principles are:

1. **Test Early**: Write tests alongside implementation
2. **Test Thoroughly**: Aim for 80%+ coverage with meaningful tests
3. **Test Safely**: Never test purchases against production
4. **Automate**: CI/CD catches regressions immediately
5. **Document**: Every feature needs docs and examples

Success is measured not just by feature completion, but by reliability, maintainability, and user trust. The --confirm flag and comprehensive testing ensure users can trust the CLI won't make accidental purchases.
