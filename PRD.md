# Amazon CLI - Product Requirements Document

## Overview

**amazon-cli** is a command-line interface that replaces the Amazon web interface, enabling programmatic access to core Amazon shopping functionality. The primary target users are AI agents, with the tool designed for publication on ClawdHub alongside a companion `skills.md` file.

## Goals

1. Provide full CLI access to Amazon shopping features (orders, returns, purchases, subscriptions)
2. Output structured JSON for seamless AI agent integration
3. Implement safety rails to prevent accidental purchases
4. Handle Amazon's rate limiting gracefully with built-in retry logic
5. Distribute via Homebrew and standalone binaries for easy installation

## Technical Stack

| Component | Choice | Rationale |
|-----------|--------|-----------|
| Language | Go | Single binary distribution, excellent CLI ecosystem (Cobra), cross-platform |
| CLI Framework | Cobra | Industry standard for Go CLIs, subcommand support, auto-generated help |
| HTTP Client | net/http + colly | Native Go HTTP with colly for scraping where needed |
| Auth | Browser-based OAuth | Opens browser for Amazon login, stores refresh token locally |
| Output | JSON | Structured output for AI agent consumption |
| Config Storage | Plain config file | `~/.amazon-cli/config.json` for credentials and settings |

## Target Marketplace

- **US only** (amazon.com) for initial release
- Future expansion to other marketplaces can be considered post-launch

## Authentication

### Flow
1. User runs `amazon-cli auth login`
2. CLI opens default browser to Amazon OAuth consent page
3. User authenticates with Amazon credentials
4. OAuth callback captures tokens
5. Tokens stored in `~/.amazon-cli/config.json`

### Token Management
- Automatic token refresh when expired
- `amazon-cli auth status` - Check current auth status
- `amazon-cli auth logout` - Clear stored credentials

## Core Features

### 1. Orders Management

```bash
# List recent orders
amazon-cli orders list [--limit N] [--status pending|delivered|returned]

# Get order details
amazon-cli orders get <order-id>

# Track shipment
amazon-cli orders track <order-id>

# Get order history (extended)
amazon-cli orders history [--year YYYY] [--format json]
```

**Output Schema (orders list):**
```json
{
  "orders": [
    {
      "order_id": "123-4567890-1234567",
      "date": "2024-01-15",
      "total": 29.99,
      "status": "delivered",
      "items": [
        {
          "asin": "B08N5WRWNW",
          "title": "Product Name",
          "quantity": 1,
          "price": 29.99
        }
      ],
      "tracking": {
        "carrier": "UPS",
        "tracking_number": "1Z999AA10123456784",
        "status": "delivered",
        "delivery_date": "2024-01-17"
      }
    }
  ],
  "total_count": 1
}
```

### 2. Returns Management

```bash
# List returnable items
amazon-cli returns list

# Get return options for an item
amazon-cli returns options <order-id> <item-id>

# Initiate a return
amazon-cli returns create <order-id> <item-id> --reason <reason-code> --confirm

# Get return label
amazon-cli returns label <return-id>

# Check return status
amazon-cli returns status <return-id>
```

**Return Reason Codes:**
- `defective` - Item is defective or doesn't work
- `wrong_item` - Received wrong item
- `not_as_described` - Item not as described
- `no_longer_needed` - No longer needed
- `better_price` - Found better price elsewhere
- `other` - Other reason

### 3. Search & Purchase

```bash
# Search products
amazon-cli search "<query>" [--category <cat>] [--min-price N] [--max-price N] [--prime-only]

# Get product details
amazon-cli product get <asin>

# Get product reviews summary
amazon-cli product reviews <asin> [--limit N]

# Add to cart
amazon-cli cart add <asin> [--quantity N]

# View cart
amazon-cli cart list

# Remove from cart
amazon-cli cart remove <asin>

# Clear cart
amazon-cli cart clear --confirm

# Checkout (REQUIRES --confirm)
amazon-cli cart checkout --confirm [--address-id <id>] [--payment-id <id>]

# Quick buy (REQUIRES --confirm)
amazon-cli buy <asin> --confirm [--quantity N] [--address-id <id>]
```

**Safety: Purchase commands REQUIRE `--confirm` flag to execute. Without it, the command will display what would happen but not execute.**

**Search Output Schema:**
```json
{
  "query": "wireless headphones",
  "results": [
    {
      "asin": "B08N5WRWNW",
      "title": "Sony WH-1000XM4 Wireless Headphones",
      "price": 278.00,
      "original_price": 349.99,
      "rating": 4.7,
      "review_count": 52431,
      "prime": true,
      "in_stock": true,
      "delivery_estimate": "Tomorrow"
    }
  ],
  "total_results": 1000,
  "page": 1
}
```

### 4. Subscribe & Save Management

```bash
# List all subscriptions
amazon-cli subscriptions list

# Get subscription details
amazon-cli subscriptions get <subscription-id>

# Skip next delivery
amazon-cli subscriptions skip <subscription-id> --confirm

# Change frequency
amazon-cli subscriptions frequency <subscription-id> --interval <weeks> --confirm

# Cancel subscription
amazon-cli subscriptions cancel <subscription-id> --confirm

# View upcoming deliveries
amazon-cli subscriptions upcoming
```

**Subscription Output Schema:**
```json
{
  "subscriptions": [
    {
      "subscription_id": "S01-1234567-8901234",
      "asin": "B00EXAMPLE",
      "title": "Coffee Pods 100 Count",
      "price": 45.99,
      "discount_percent": 15,
      "frequency_weeks": 4,
      "next_delivery": "2024-02-01",
      "status": "active",
      "quantity": 1
    }
  ]
}
```

## Global Flags

```bash
--output, -o      Output format: json (default), table, raw
--quiet, -q       Suppress non-essential output
--verbose, -v     Enable verbose logging
--config          Path to config file (default: ~/.amazon-cli/config.json)
--no-color        Disable colored output
```

## Configuration File

Location: `~/.amazon-cli/config.json`

```json
{
  "auth": {
    "access_token": "...",
    "refresh_token": "...",
    "expires_at": "2024-01-20T12:00:00Z"
  },
  "defaults": {
    "address_id": "addr_default",
    "payment_id": "pay_default",
    "output_format": "json"
  },
  "rate_limiting": {
    "min_delay_ms": 1000,
    "max_delay_ms": 5000,
    "max_retries": 3
  }
}
```

## Rate Limiting & Retry Strategy

To avoid triggering Amazon's anti-automation measures:

1. **Minimum delay**: 1 second between requests (configurable)
2. **Jitter**: Random 0-500ms added to each delay
3. **Exponential backoff**: On 429/503 responses, wait 2^n seconds (max 60s)
4. **Max retries**: 3 attempts before failing (configurable)
5. **User-Agent rotation**: Rotate through common browser user agents

## Error Handling

All errors return JSON with consistent schema:

```json
{
  "error": {
    "code": "AUTH_EXPIRED",
    "message": "Authentication token has expired. Run 'amazon-cli auth login' to re-authenticate.",
    "details": {}
  }
}
```

**Error Codes:**
| Code | Description |
|------|-------------|
| `AUTH_REQUIRED` | Not logged in |
| `AUTH_EXPIRED` | Token expired |
| `NOT_FOUND` | Resource not found |
| `RATE_LIMITED` | Too many requests |
| `INVALID_INPUT` | Invalid command input |
| `PURCHASE_FAILED` | Purchase could not be completed |
| `NETWORK_ERROR` | Network connectivity issue |
| `AMAZON_ERROR` | Amazon returned an error |

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Invalid arguments |
| 3 | Authentication error |
| 4 | Network error |
| 5 | Rate limited |
| 6 | Not found |

## Distribution

### Homebrew (macOS/Linux)

```bash
brew tap amazon-cli/tap
brew install amazon-cli
```

### Binary Releases

Pre-compiled binaries for:
- macOS (arm64, amd64)
- Linux (arm64, amd64)
- Windows (amd64)

Available on GitHub Releases page.

### Build from Source

```bash
go install github.com/zkwentz/amazon-cli@latest
```

## Project Structure

```
amazon-cli/
├── cmd/
│   ├── root.go           # Root command setup
│   ├── auth.go           # Auth commands
│   ├── orders.go         # Orders commands
│   ├── returns.go        # Returns commands
│   ├── search.go         # Search command
│   ├── product.go        # Product commands
│   ├── cart.go           # Cart commands
│   ├── buy.go            # Buy command
│   └── subscriptions.go  # Subscription commands
├── internal/
│   ├── amazon/           # Amazon API/scraping client
│   │   ├── client.go
│   │   ├── auth.go
│   │   ├── orders.go
│   │   ├── returns.go
│   │   ├── search.go
│   │   ├── cart.go
│   │   └── subscriptions.go
│   ├── config/           # Configuration management
│   │   └── config.go
│   ├── output/           # Output formatting
│   │   └── json.go
│   └── ratelimit/        # Rate limiting logic
│       └── limiter.go
├── pkg/
│   └── models/           # Shared data models
│       ├── order.go
│       ├── product.go
│       ├── subscription.go
│       └── errors.go
├── skills.md             # ClawdHub skills definition
├── main.go
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Development Phases

---

### Phase 1: Project Setup & Foundation

#### 1.1 Repository Initialization
- [x] Create new GitHub repository `amazon-cli`
- [x] Initialize Go module: `go mod init github.com/zkwentz/amazon-cli`
- [x] Create `.gitignore` file with Go defaults (binaries, vendor/, .env, etc.)
- [x] Create initial `main.go` entry point
- [x] Set up directory structure:
  ```
  mkdir -p cmd internal/amazon internal/config internal/output internal/ratelimit pkg/models
  ```

#### 1.2 Cobra CLI Framework Setup
- [x] Install Cobra: `go get -u github.com/spf13/cobra@latest`
- [x] Install Viper for config: `go get -u github.com/spf13/viper@latest`
- [x] Create `cmd/root.go` with root command:
  - [ ] Set `Use: "amazon-cli"`
  - [ ] Set `Short` and `Long` descriptions
  - [ ] Add persistent flags: `--output`, `--quiet`, `--verbose`, `--config`, `--no-color`
  - [ ] Initialize Viper config binding in `init()`
- [x] Create `cmd/version.go` - simple version command that prints version string
- [x] Wire up `main.go` to execute root command
- [x] Verify CLI runs: `go run main.go --help` should show help text

#### 1.3 Configuration Management
- [x] Create `internal/config/config.go`:
  - [ ] Define `Config` struct matching the JSON schema in PRD (auth, defaults, rate_limiting)
  - [ ] Define `AuthConfig` struct with `AccessToken`, `RefreshToken`, `ExpiresAt` fields
  - [ ] Define `DefaultsConfig` struct with `AddressID`, `PaymentID`, `OutputFormat` fields
  - [ ] Define `RateLimitConfig` struct with `MinDelayMs`, `MaxDelayMs`, `MaxRetries` fields
- [x] Implement `LoadConfig(path string) (*Config, error)`:
  - [ ] Default path: `~/.amazon-cli/config.json`
  - [ ] Create directory if not exists with `0700` permissions
  - [ ] Return empty config if file doesn't exist (first run)
  - [ ] Parse JSON and return Config struct
- [x] Implement `SaveConfig(config *Config, path string) error`:
  - [ ] Marshal to JSON with indentation
  - [ ] Write to file with `0600` permissions
- [x] Implement `GetConfigPath() string` helper that respects `--config` flag
- [x] Write unit tests for config load/save operations

#### 1.4 JSON Output System
- [x] Create `internal/output/output.go`:
  - [ ] Define `OutputFormat` type (JSON, Table, Raw constants)
  - [ ] Create `Printer` struct that holds format preference and quiet mode
- [x] Implement `NewPrinter(format string, quiet bool) *Printer`
- [x] Implement `Print(data interface{}) error`:
  - [ ] For JSON: marshal with `json.MarshalIndent` and print to stdout
  - [ ] For Table: use `tablewriter` package for human-readable output
  - [ ] For Raw: print raw string representation
- [x] Implement `PrintError(err error) error`:
  - [ ] Format errors as JSON: `{"error": {"code": "...", "message": "...", "details": {}}}`
  - [ ] Use error codes from PRD (AUTH_REQUIRED, AUTH_EXPIRED, etc.)
- [x] Create `pkg/models/errors.go`:
  - [ ] Define `CLIError` struct with `Code`, `Message`, `Details` fields
  - [ ] Define error code constants matching PRD table
  - [ ] Implement `Error()` method on CLIError for error interface
- [x] Write unit tests for JSON output formatting

#### 1.5 Rate Limiting Infrastructure
- [x] Create `internal/ratelimit/limiter.go`:
  - [ ] Define `RateLimiter` struct with config, last request time, retry count
- [x] Implement `NewRateLimiter(config RateLimitConfig) *RateLimiter`
- [x] Implement `Wait() error`:
  - [ ] Calculate time since last request
  - [ ] If less than MinDelayMs, sleep for the difference
  - [ ] Add random jitter (0-500ms) using `crypto/rand`
  - [ ] Update last request timestamp
- [x] Implement `WaitWithBackoff(attempt int) error`:
  - [ ] Calculate exponential backoff: `min(2^attempt * 1000ms, 60000ms)`
  - [ ] Sleep for calculated duration
  - [ ] Log backoff duration if verbose mode
- [x] Implement `ShouldRetry(statusCode int, attempt int) bool`:
  - [ ] Return true for 429 (rate limited) or 503 (service unavailable)
  - [ ] Return false if attempt >= MaxRetries
  - [ ] Return false for other status codes
- [x] Write unit tests for rate limiter logic

#### 1.6 HTTP Client Foundation
- [x] Create `internal/amazon/client.go`:
  - [ ] Define `Client` struct with http.Client, RateLimiter, Config, user agents list
  - [ ] Define list of 10+ common browser User-Agent strings for rotation
- [x] Implement `NewClient(config *Config) *Client`:
  - [ ] Create http.Client with 30 second timeout
  - [ ] Initialize cookie jar for session management
  - [ ] Create rate limiter from config
- [x] Implement `Do(req *http.Request) (*http.Response, error)`:
  - [ ] Call rate limiter `Wait()` before request
  - [ ] Set random User-Agent from rotation list
  - [ ] Set common headers (Accept, Accept-Language, etc.)
  - [ ] Execute request with retry logic:
    - [ ] If response is 429/503 and ShouldRetry is true, call WaitWithBackoff and retry
    - [ ] Return final response or error after max retries
- [x] Implement `Get(url string) (*http.Response, error)` convenience method
- [x] Implement `PostForm(url string, data url.Values) (*http.Response, error)` convenience method
- [x] Write integration tests with mock server

---

### Phase 2: Authentication System

#### 2.1 OAuth Flow Research & Setup
- [x] Research Amazon's OAuth/Login with Amazon (LWA) API:
  - [ ] Register app at https://developer.amazon.com/ to get Client ID/Secret
  - [ ] Document required OAuth scopes for order/profile access
  - [ ] Note: If official API insufficient, plan for browser session approach
- [x] Create `internal/amazon/auth.go`:
  - [ ] Define OAuth constants (auth URL, token URL, redirect URI)
  - [ ] Define `AuthTokens` struct with access token, refresh token, expiry

#### 2.2 Login Command Implementation
- [x] Create `cmd/auth.go`:
  - [ ] Add `auth` parent command with subcommands
- [x] Implement `auth login` command:
  - [ ] Generate random state parameter for CSRF protection
  - [ ] Build OAuth authorization URL with scopes and state
  - [ ] Start local HTTP server on random available port (e.g., 8085-8095)
  - [ ] Open browser to authorization URL using `github.com/pkg/browser`
  - [ ] Handle OAuth callback on local server:
    - [ ] Verify state parameter matches
    - [ ] Extract authorization code from query params
    - [ ] Exchange code for tokens via POST to token endpoint
    - [ ] Store tokens in config file
    - [ ] Display success message and close browser tab
  - [ ] Handle timeout (2 minutes) if user doesn't complete login
  - [ ] Print JSON output: `{"status": "authenticated", "expires_at": "..."}`

#### 2.3 Token Management
- [x] Implement `auth status` command:
  - [ ] Load config and check for tokens
  - [ ] If no tokens, output: `{"authenticated": false}`
  - [ ] If tokens exist, check expiry time
  - [ ] Output: `{"authenticated": true, "expires_at": "...", "expires_in_seconds": N}`
- [x] Implement `auth logout` command:
  - [ ] Load config
  - [ ] Clear auth section (set tokens to empty)
  - [ ] Save config
  - [ ] Output: `{"status": "logged_out"}`
- [x] Implement `RefreshTokenIfNeeded(config *Config) error` in `internal/amazon/auth.go`:
  - [ ] Check if access token expires within 5 minutes
  - [ ] If so, use refresh token to get new access token
  - [ ] Update config with new tokens
  - [ ] Save config to disk
- [x] Add auth check middleware to client:
  - [ ] Before any authenticated request, call RefreshTokenIfNeeded
  - [ ] If no tokens exist, return AUTH_REQUIRED error

#### 2.4 Alternative: Browser Session Auth (if OAuth insufficient)
- [x] If Amazon OAuth doesn't provide needed access, implement cookie-based auth:
  - [ ] Implement `auth login --browser` flag that:
    - [ ] Opens Amazon login page in browser
    - [ ] Instructs user to complete login
    - [ ] Uses browser automation (Rod/Chromedp) to capture session cookies
    - [ ] Stores cookies in config file
  - [ ] Implement cookie refresh detection and re-auth prompts
- [x] Document which auth method is being used in README

---

### Phase 3: Orders Management

#### 3.1 Data Models
- [x] Create `pkg/models/order.go`:
  - [ ] Define `Order` struct:
    ```go
    type Order struct {
        OrderID     string      `json:"order_id"`
        Date        string      `json:"date"`
        Total       float64     `json:"total"`
        Status      string      `json:"status"`
        Items       []OrderItem `json:"items"`
        Tracking    *Tracking   `json:"tracking,omitempty"`
    }
    ```
  - [ ] Define `OrderItem` struct with ASIN, Title, Quantity, Price
  - [ ] Define `Tracking` struct with Carrier, TrackingNumber, Status, DeliveryDate
  - [ ] Define `OrdersResponse` struct with Orders slice and TotalCount

#### 3.2 Orders API Client
- [x] Create `internal/amazon/orders.go`:
  - [ ] Research Amazon order history page structure (HTML selectors, API endpoints)
  - [ ] Document the URLs and request format needed
- [x] Implement `GetOrders(limit int, status string) (*OrdersResponse, error)`:
  - [ ] Build request to Amazon order history page/API
  - [ ] Parse HTML response using `goquery` or parse JSON if API available
  - [ ] Extract order data into Order structs
  - [ ] Filter by status if provided
  - [ ] Limit results to requested count
  - [ ] Return OrdersResponse
- [x] Implement `GetOrder(orderID string) (*Order, error)`:
  - [ ] Fetch individual order details page
  - [ ] Parse complete order information including all items
  - [ ] Return Order struct with full details
- [x] Implement `GetOrderTracking(orderID string) (*Tracking, error)`:
  - [ ] Fetch tracking information for order
  - [ ] Parse carrier, tracking number, status, delivery date
  - [ ] Return Tracking struct
- [x] Implement `GetOrderHistory(year int) (*OrdersResponse, error)`:
  - [ ] Fetch orders from specific year
  - [ ] Handle pagination if Amazon paginates results
  - [ ] Return all orders for that year

#### 3.3 Orders Commands
- [x] Create `cmd/orders.go`:
  - [ ] Add `orders` parent command
- [x] Implement `orders list` command:
  - [ ] Add `--limit` flag (default 10)
  - [ ] Add `--status` flag (pending, delivered, returned, or empty for all)
  - [ ] Call client.GetOrders with parameters
  - [ ] Output JSON response via Printer
- [x] Implement `orders get <order-id>` command:
  - [ ] Validate order-id argument is provided
  - [ ] Call client.GetOrder
  - [ ] Output JSON response
- [x] Implement `orders track <order-id>` command:
  - [ ] Validate order-id argument
  - [ ] Call client.GetOrderTracking
  - [ ] Output tracking JSON
- [x] Implement `orders history` command:
  - [ ] Add `--year` flag (default current year)
  - [ ] Call client.GetOrderHistory
  - [ ] Output JSON response

#### 3.4 Orders Testing
- [x] Create mock Amazon responses for testing
- [x] Write unit tests for order parsing logic
- [x] Write integration tests for orders commands
- [x] Test error cases: invalid order ID, auth expired, network errors

---

### Phase 4: Returns Management

#### 4.1 Data Models
- [x] Create `pkg/models/return.go`:
  - [ ] Define `ReturnableItem` struct:
    ```go
    type ReturnableItem struct {
        OrderID      string    `json:"order_id"`
        ItemID       string    `json:"item_id"`
        ASIN         string    `json:"asin"`
        Title        string    `json:"title"`
        Price        float64   `json:"price"`
        PurchaseDate string    `json:"purchase_date"`
        ReturnWindow string    `json:"return_window"`
    }
    ```
  - [ ] Define `ReturnOption` struct with Method, Label, DropoffLocation, Fee
  - [ ] Define `Return` struct with ReturnID, OrderID, ItemID, Status, Reason, CreatedAt
  - [ ] Define `ReturnLabel` struct with URL, Carrier, Instructions

#### 4.2 Returns API Client
- [x] Create `internal/amazon/returns.go`:
  - [ ] Research Amazon returns flow (URLs, form submissions, API calls)
- [x] Implement `GetReturnableItems() ([]ReturnableItem, error)`:
  - [ ] Fetch returnable items from Amazon returns center
  - [ ] Parse items with their return eligibility
  - [ ] Return slice of ReturnableItem
- [x] Implement `GetReturnOptions(orderID, itemID string) ([]ReturnOption, error)`:
  - [ ] Fetch return options for specific item
  - [ ] Parse available return methods (UPS, Amazon Locker, Whole Foods, etc.)
  - [ ] Return slice of ReturnOption
- [x] Implement `CreateReturn(orderID, itemID, reason string) (*Return, error)`:
  - [ ] Validate reason code against allowed values
  - [ ] Submit return request to Amazon
  - [ ] Parse confirmation response
  - [ ] Return Return struct with return ID
- [x] Implement `GetReturnLabel(returnID string) (*ReturnLabel, error)`:
  - [ ] Fetch return label for initiated return
  - [ ] Extract label URL/PDF link
  - [ ] Return ReturnLabel struct
- [x] Implement `GetReturnStatus(returnID string) (*Return, error)`:
  - [ ] Fetch current status of return
  - [ ] Parse status (initiated, shipped, received, refunded)
  - [ ] Return updated Return struct

#### 4.3 Returns Commands
- [x] Create `cmd/returns.go`:
  - [ ] Add `returns` parent command
- [x] Implement `returns list` command:
  - [ ] Call client.GetReturnableItems
  - [ ] Output JSON array of returnable items
- [x] Implement `returns options <order-id> <item-id>` command:
  - [ ] Validate both arguments provided
  - [ ] Call client.GetReturnOptions
  - [ ] Output JSON array of return options
- [x] Implement `returns create <order-id> <item-id>` command:
  - [ ] Add `--reason` required flag
  - [ ] Add `--confirm` required flag for safety
  - [ ] Validate reason is in allowed list (defective, wrong_item, etc.)
  - [ ] If --confirm not provided:
    - [ ] Show what would be returned (dry run)
    - [ ] Output: `{"dry_run": true, "would_return": {...}, "message": "Add --confirm to execute"}`
  - [ ] If --confirm provided:
    - [ ] Call client.CreateReturn
    - [ ] Output return confirmation JSON
- [x] Implement `returns label <return-id>` command:
  - [ ] Call client.GetReturnLabel
  - [ ] Output JSON with label URL and instructions
- [x] Implement `returns status <return-id>` command:
  - [ ] Call client.GetReturnStatus
  - [ ] Output return status JSON

#### 4.4 Returns Testing
- [x] Write unit tests for return reason validation
- [x] Write tests for dry run vs confirmed behavior
- [x] Test error cases: item not returnable, return window expired

---

### Phase 5: Search & Product Features

#### 5.1 Data Models
- [x] Create `pkg/models/product.go`:
  - [ ] Define `Product` struct:
    ```go
    type Product struct {
        ASIN            string   `json:"asin"`
        Title           string   `json:"title"`
        Price           float64  `json:"price"`
        OriginalPrice   *float64 `json:"original_price,omitempty"`
        Rating          float64  `json:"rating"`
        ReviewCount     int      `json:"review_count"`
        Prime           bool     `json:"prime"`
        InStock         bool     `json:"in_stock"`
        DeliveryEstimate string  `json:"delivery_estimate"`
        Description     string   `json:"description,omitempty"`
        Features        []string `json:"features,omitempty"`
        Images          []string `json:"images,omitempty"`
    }
    ```
  - [ ] Define `SearchResponse` struct with Query, Results, TotalResults, Page
  - [ ] Define `Review` struct with Rating, Title, Body, Author, Date, Verified
  - [ ] Define `ReviewsResponse` struct with ASIN, Reviews slice, AverageRating, TotalReviews

#### 5.2 Search API Client
- [x] Create `internal/amazon/search.go`:
  - [ ] Research Amazon search page structure and parameters
- [x] Implement `Search(query string, opts SearchOptions) (*SearchResponse, error)`:
  - [ ] Define SearchOptions struct with Category, MinPrice, MaxPrice, PrimeOnly, Page
  - [ ] Build search URL with query parameters
  - [ ] Fetch search results page
  - [ ] Parse product listings from HTML/JSON
  - [ ] Extract: ASIN, title, price, rating, review count, Prime badge, stock status
  - [ ] Return SearchResponse

#### 5.3 Product API Client
- [x] Create `internal/amazon/product.go`:
- [x] Implement `GetProduct(asin string) (*Product, error)`:
  - [ ] Fetch product detail page
  - [ ] Parse full product information:
    - [ ] Title, price, original price (for discounts)
    - [ ] Rating, review count
    - [ ] Prime eligibility
    - [ ] Stock status
    - [ ] Delivery estimate
    - [ ] Description
    - [ ] Feature bullets
    - [ ] Image URLs
  - [ ] Return Product struct
- [x] Implement `GetProductReviews(asin string, limit int) (*ReviewsResponse, error)`:
  - [ ] Fetch product reviews page
  - [ ] Parse individual reviews
  - [ ] Include: rating, title, body, author, date, verified purchase badge
  - [ ] Limit to requested count
  - [ ] Return ReviewsResponse

#### 5.4 Search & Product Commands
- [x] Create `cmd/search.go`:
  - [ ] Implement `search "<query>"` command:
    - [ ] Add `--category` flag
    - [ ] Add `--min-price` flag
    - [ ] Add `--max-price` flag
    - [ ] Add `--prime-only` flag
    - [ ] Add `--page` flag (default 1)
    - [ ] Call client.Search with options
    - [ ] Output SearchResponse JSON
- [x] Create `cmd/product.go`:
  - [ ] Add `product` parent command
- [x] Implement `product get <asin>` command:
  - [ ] Validate ASIN format (10 alphanumeric characters)
  - [ ] Call client.GetProduct
  - [ ] Output Product JSON
- [x] Implement `product reviews <asin>` command:
  - [ ] Add `--limit` flag (default 10)
  - [ ] Call client.GetProductReviews
  - [ ] Output ReviewsResponse JSON

#### 5.5 Search Testing
- [x] Test search with various query types
- [x] Test price range filtering
- [x] Test Prime-only filtering
- [x] Test ASIN validation
- [x] Test handling of out-of-stock products

---

### Phase 6: Cart & Checkout

#### 6.1 Data Models
- [x] Create `pkg/models/cart.go`:
  - [ ] Define `CartItem` struct:
    ```go
    type CartItem struct {
        ASIN      string  `json:"asin"`
        Title     string  `json:"title"`
        Price     float64 `json:"price"`
        Quantity  int     `json:"quantity"`
        Subtotal  float64 `json:"subtotal"`
        Prime     bool    `json:"prime"`
        InStock   bool    `json:"in_stock"`
    }
    ```
  - [ ] Define `Cart` struct with Items, Subtotal, EstimatedTax, Total, ItemCount
  - [ ] Define `Address` struct with ID, Name, Street, City, State, Zip, Country, Default
  - [ ] Define `PaymentMethod` struct with ID, Type, Last4, Default
  - [ ] Define `CheckoutPreview` struct with Cart, Address, PaymentMethod, DeliveryOptions
  - [ ] Define `OrderConfirmation` struct with OrderID, Total, EstimatedDelivery

#### 6.2 Cart API Client
- [x] Create `internal/amazon/cart.go`:
  - [ ] Research Amazon cart operations (add, remove, update quantity)
- [x] Implement `AddToCart(asin string, quantity int) (*Cart, error)`:
  - [ ] Submit add-to-cart request
  - [ ] Handle quantity limits
  - [ ] Return updated cart
- [x] Implement `GetCart() (*Cart, error)`:
  - [ ] Fetch current cart contents
  - [ ] Parse all cart items with prices
  - [ ] Calculate totals
  - [ ] Return Cart struct
- [x] Implement `RemoveFromCart(asin string) (*Cart, error)`:
  - [ ] Submit remove item request
  - [ ] Return updated cart
- [x] Implement `ClearCart() error`:
  - [ ] Remove all items from cart
- [x] Implement `GetAddresses() ([]Address, error)`:
  - [ ] Fetch saved addresses
  - [ ] Return slice of Address
- [x] Implement `GetPaymentMethods() ([]PaymentMethod, error)`:
  - [ ] Fetch saved payment methods
  - [ ] Return slice of PaymentMethod

#### 6.3 Checkout API Client
- [x] Implement `PreviewCheckout(addressID, paymentID string) (*CheckoutPreview, error)`:
  - [ ] Initiate checkout flow without completing
  - [ ] Fetch order preview with totals, delivery estimates
  - [ ] Return CheckoutPreview struct
- [x] Implement `CompleteCheckout(addressID, paymentID string) (*OrderConfirmation, error)`:
  - [ ] Submit final checkout
  - [ ] Handle payment authorization
  - [ ] Parse order confirmation
  - [ ] Return OrderConfirmation with order ID

#### 6.4 Cart Commands
- [x] Create `cmd/cart.go`:
  - [ ] Add `cart` parent command
- [x] Implement `cart add <asin>` command:
  - [ ] Add `--quantity` flag (default 1)
  - [ ] Validate ASIN format
  - [ ] Call client.AddToCart
  - [ ] Output updated cart JSON
- [x] Implement `cart list` command:
  - [ ] Call client.GetCart
  - [ ] Output cart JSON with all items and totals
- [x] Implement `cart remove <asin>` command:
  - [ ] Validate ASIN format
  - [ ] Call client.RemoveFromCart
  - [ ] Output updated cart JSON
- [x] Implement `cart clear` command:
  - [ ] Require `--confirm` flag
  - [ ] Without --confirm: output dry run message
  - [ ] With --confirm: call client.ClearCart, output success
- [x] Implement `cart checkout` command:
  - [ ] Require `--confirm` flag
  - [ ] Add `--address-id` optional flag
  - [ ] Add `--payment-id` optional flag
  - [ ] Without --confirm:
    - [ ] Call client.PreviewCheckout
    - [ ] Output preview JSON with `"dry_run": true`
  - [ ] With --confirm:
    - [ ] Call client.CompleteCheckout
    - [ ] Output OrderConfirmation JSON

#### 6.5 Buy Command (Quick Purchase)
- [x] Create `cmd/buy.go`:
- [x] Implement `buy <asin>` command:
  - [ ] Require `--confirm` flag
  - [ ] Add `--quantity` flag (default 1)
  - [ ] Add `--address-id` optional flag
  - [ ] Add `--payment-id` optional flag
  - [ ] Without --confirm:
    - [ ] Fetch product details
    - [ ] Output what would be purchased: `{"dry_run": true, "product": {...}, "quantity": N, "total": X}`
  - [ ] With --confirm:
    - [ ] Add to cart
    - [ ] Complete checkout
    - [ ] Output OrderConfirmation JSON

#### 6.6 Cart & Checkout Testing
- [x] Test add/remove/clear cart operations
- [x] Test checkout preview without --confirm
- [x] Test that checkout fails without --confirm
- [x] Test checkout with explicit address/payment
- [x] Test quick buy flow

---

### Phase 7: Subscriptions Management

#### 7.1 Data Models
- [x] Create `pkg/models/subscription.go`:
  - [ ] Define `Subscription` struct:
    ```go
    type Subscription struct {
        SubscriptionID  string  `json:"subscription_id"`
        ASIN            string  `json:"asin"`
        Title           string  `json:"title"`
        Price           float64 `json:"price"`
        DiscountPercent int     `json:"discount_percent"`
        FrequencyWeeks  int     `json:"frequency_weeks"`
        NextDelivery    string  `json:"next_delivery"`
        Status          string  `json:"status"`
        Quantity        int     `json:"quantity"`
    }
    ```
  - [ ] Define `SubscriptionsResponse` struct with Subscriptions slice
  - [ ] Define `UpcomingDelivery` struct with SubscriptionID, ASIN, Title, DeliveryDate, Quantity

#### 7.2 Subscriptions API Client
- [x] Create `internal/amazon/subscriptions.go`:
  - [ ] Research Amazon Subscribe & Save page structure
- [x] Implement `GetSubscriptions() (*SubscriptionsResponse, error)`:
  - [ ] Fetch Subscribe & Save dashboard
  - [ ] Parse all active and paused subscriptions
  - [ ] Return SubscriptionsResponse
- [x] Implement `GetSubscription(subscriptionID string) (*Subscription, error)`:
  - [ ] Fetch specific subscription details
  - [ ] Return full Subscription struct
- [x] Implement `SkipDelivery(subscriptionID string) (*Subscription, error)`:
  - [ ] Submit skip next delivery request
  - [ ] Return updated subscription with new next delivery date
- [x] Implement `UpdateFrequency(subscriptionID string, weeks int) (*Subscription, error)`:
  - [ ] Validate weeks is valid frequency (1, 2, 3, 4, 5, 6 months converted to weeks)
  - [ ] Submit frequency change request
  - [ ] Return updated subscription
- [x] Implement `CancelSubscription(subscriptionID string) (*Subscription, error)`:
  - [ ] Submit cancellation request
  - [ ] Return subscription with status "cancelled"
- [x] Implement `GetUpcomingDeliveries() ([]UpcomingDelivery, error)`:
  - [ ] Fetch upcoming deliveries across all subscriptions
  - [ ] Return slice of UpcomingDelivery sorted by date

#### 7.3 Subscription Commands
- [x] Create `cmd/subscriptions.go`:
  - [ ] Add `subscriptions` parent command
- [x] Implement `subscriptions list` command:
  - [ ] Call client.GetSubscriptions
  - [ ] Output JSON array of subscriptions
- [x] Implement `subscriptions get <subscription-id>` command:
  - [ ] Validate subscription-id argument
  - [ ] Call client.GetSubscription
  - [ ] Output subscription JSON
- [x] Implement `subscriptions skip <subscription-id>` command:
  - [ ] Require `--confirm` flag
  - [ ] Without --confirm: show what would be skipped
  - [ ] With --confirm: call client.SkipDelivery, output updated subscription
- [x] Implement `subscriptions frequency <subscription-id>` command:
  - [ ] Require `--interval` flag (weeks)
  - [ ] Require `--confirm` flag
  - [ ] Validate interval is reasonable (1-26 weeks)
  - [ ] Without --confirm: show what would change
  - [ ] With --confirm: call client.UpdateFrequency, output updated subscription
- [x] Implement `subscriptions cancel <subscription-id>` command:
  - [ ] Require `--confirm` flag
  - [ ] Without --confirm: show cancellation preview
  - [ ] With --confirm: call client.CancelSubscription, output confirmation
- [x] Implement `subscriptions upcoming` command:
  - [ ] Call client.GetUpcomingDeliveries
  - [ ] Output JSON array sorted by delivery date

#### 7.4 Subscription Testing
- [x] Test listing all subscriptions
- [x] Test skip delivery with/without confirm
- [x] Test frequency change validation
- [x] Test cancellation flow
- [x] Test upcoming deliveries sorting

---

### Phase 8: Error Handling & Polish

#### 8.1 Comprehensive Error Handling
- [x] Audit all API client methods for error handling:
  - [ ] Network errors → NETWORK_ERROR code
  - [ ] 401 responses → AUTH_EXPIRED code
  - [ ] 404 responses → NOT_FOUND code
  - [ ] 429 responses → RATE_LIMITED code
  - [ ] 5xx responses → AMAZON_ERROR code
  - [ ] Parse errors → AMAZON_ERROR with details
- [x] Implement error wrapping with context using `fmt.Errorf` with `%w`
- [x] Add `--verbose` flag handling to print debug info on errors
- [x] Ensure all errors output valid JSON (never panic or print stack traces)
- [x] Implement graceful handling of unexpected HTML responses (CAPTCHA, login redirects)

#### 8.2 Exit Codes
- [x] Implement exit code system in `cmd/root.go`:
  - [ ] 0 for success
  - [ ] 1 for general error
  - [ ] 2 for invalid arguments (use Cobra's built-in)
  - [ ] 3 for authentication error
  - [ ] 4 for network error
  - [ ] 5 for rate limited
  - [ ] 6 for not found
- [x] Map CLIError codes to exit codes
- [x] Ensure all commands use consistent exit codes

#### 8.3 Input Validation
- [x] Add ASIN format validation helper (10 alphanumeric characters)
- [x] Add order ID format validation
- [x] Add subscription ID format validation
- [x] Add price range validation (min < max, both positive)
- [x] Add quantity validation (positive integer, reasonable max)
- [x] Validate all user inputs before API calls

#### 8.4 Logging
- [x] Implement structured logging using `log/slog`:
  - [ ] Info level for normal operations
  - [ ] Debug level for request/response details
  - [ ] Error level for failures
- [x] Add `--verbose` flag to enable debug logging
- [x] Ensure tokens/cookies never appear in logs
- [x] Add request timing logging for performance debugging

---

### Phase 9: Distribution & Release

#### 9.1 Build System
- [x] Create `Makefile` with targets:
  - [ ] `build` - build for current platform
  - [ ] `build-all` - build for all platforms (darwin/amd64, darwin/arm64, linux/amd64, linux/arm64, windows/amd64)
  - [ ] `test` - run all tests
  - [ ] `lint` - run golangci-lint
  - [ ] `clean` - remove build artifacts
  - [ ] `install` - install to $GOPATH/bin
- [x] Add `ldflags` to embed version info at build time:
  ```
  -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)
  ```
- [x] Create `.goreleaser.yml` for automated releases

#### 9.2 GitHub Actions CI/CD
- [x] Create `.github/workflows/ci.yml`:
  - [ ] Trigger on push to main and PRs
  - [ ] Run `go test ./...`
  - [ ] Run `golangci-lint`
  - [ ] Build for all platforms
  - [ ] Upload artifacts
- [x] Create `.github/workflows/release.yml`:
  - [ ] Trigger on tag push (v*)
  - [ ] Use goreleaser to build and create release
  - [ ] Upload binaries to GitHub Release
  - [ ] Generate changelog from commits

#### 9.3 Homebrew Distribution
- [x] Create separate repo: `homebrew-amazon-cli` (or `homebrew-tap`)
- [x] Create Homebrew formula `amazon-cli.rb`:
  ```ruby
  class AmazonCli < Formula
    desc "CLI for Amazon shopping - orders, returns, purchases, subscriptions"
    homepage "https://github.com/michaelshimeles/amazon-cli"
    url "https://github.com/michaelshimeles/amazon-cli/releases/download/v#{version}/amazon-cli_#{version}_darwin_amd64.tar.gz"
    sha256 "..."
    license "MIT"

    def install
      bin.install "amazon-cli"
    end

    test do
      system "#{bin}/amazon-cli", "--version"
    end
  end
  ```
- [x] Add GitHub Action to update formula on release
- [x] Test installation: `brew tap michaelshimeles/tap - [ ] Test installation: `brew tap michaelshimeles/tap && brew install amazon-cli`- [ ] Test installation: `brew tap michaelshimeles/tap && brew install amazon-cli` brew install amazon-cli`

#### 9.4 Documentation
- [x] Create comprehensive `README.md`:
  - [ ] Project description and badges
  - [ ] Installation instructions (Homebrew, binary, source)
  - [ ] Quick start guide
  - [ ] Authentication setup
  - [ ] Command reference with examples
  - [ ] Configuration options
  - [ ] Contributing guidelines
- [x] Add inline help text to all commands (Cobra Long descriptions)
- [x] Create `CHANGELOG.md` with initial release notes
- [x] Add `LICENSE` file (MIT recommended)

---

### Phase 10: ClawdHub Integration

#### 10.1 Skills.md Structure Research
- [x] Review ClawdHub documentation for skills.md format
- [x] Study existing ClawdHub skills for examples
- [x] Identify required metadata fields

#### 10.2 Create skills.md
- [x] Create `skills.md` file with proper structure:
  - [ ] Skill metadata header:
    ```yaml
    ---
    name: amazon-cli
    description: CLI tool for managing Amazon orders, returns, purchases, and subscriptions
    version: 1.0.0
    author: michaelshimeles
    repository: https://github.com/michaelshimeles/amazon-cli
    ---
    ```
  - [ ] Overview section explaining what the skill does
  - [ ] Installation instructions

#### 10.3 Document All Commands as Actions
- [x] Document `auth` commands:
  - [ ] `auth login` - inputs: none, outputs: auth status
  - [ ] `auth status` - inputs: none, outputs: authentication state
  - [ ] `auth logout` - inputs: none, outputs: confirmation
- [x] Document `orders` commands:
  - [ ] `orders list` - inputs: limit, status; outputs: OrdersResponse schema
  - [ ] `orders get` - inputs: order_id; outputs: Order schema
  - [ ] `orders track` - inputs: order_id; outputs: Tracking schema
  - [ ] `orders history` - inputs: year; outputs: OrdersResponse schema
- [x] Document `returns` commands:
  - [ ] `returns list` - inputs: none; outputs: ReturnableItem[] schema
  - [ ] `returns options` - inputs: order_id, item_id; outputs: ReturnOption[] schema
  - [ ] `returns create` - inputs: order_id, item_id, reason, confirm; outputs: Return schema
  - [ ] `returns label` - inputs: return_id; outputs: ReturnLabel schema
  - [ ] `returns status` - inputs: return_id; outputs: Return schema
- [x] Document `search` command:
  - [ ] inputs: query, category, min_price, max_price, prime_only
  - [ ] outputs: SearchResponse schema
- [x] Document `product` commands:
  - [ ] `product get` - inputs: asin; outputs: Product schema
  - [ ] `product reviews` - inputs: asin, limit; outputs: ReviewsResponse schema
- [x] Document `cart` commands:
  - [ ] `cart add` - inputs: asin, quantity; outputs: Cart schema
  - [ ] `cart list` - inputs: none; outputs: Cart schema
  - [ ] `cart remove` - inputs: asin; outputs: Cart schema
  - [ ] `cart clear` - inputs: confirm; outputs: confirmation
  - [ ] `cart checkout` - inputs: confirm, address_id, payment_id; outputs: OrderConfirmation schema
- [x] Document `buy` command:
  - [ ] inputs: asin, quantity, confirm, address_id, payment_id
  - [ ] outputs: OrderConfirmation schema
- [x] Document `subscriptions` commands:
  - [ ] `subscriptions list` - inputs: none; outputs: SubscriptionsResponse schema
  - [ ] `subscriptions get` - inputs: subscription_id; outputs: Subscription schema
  - [ ] `subscriptions skip` - inputs: subscription_id, confirm; outputs: Subscription schema
  - [ ] `subscriptions frequency` - inputs: subscription_id, interval, confirm; outputs: Subscription schema
  - [ ] `subscriptions cancel` - inputs: subscription_id, confirm; outputs: Subscription schema
  - [ ] `subscriptions upcoming` - inputs: none; outputs: UpcomingDelivery[] schema

#### 10.4 Add AI Agent Usage Examples
- [x] Add example: "Check my recent orders"
  ```
  amazon-cli orders list --limit 5
  ```
- [x] Add example: "Find the tracking info for order X"
  ```
  amazon-cli orders track 123-4567890-1234567
  ```
- [x] Add example: "Search for wireless headphones under $100"
  ```
  amazon-cli search "wireless headphones" --max-price 100 --prime-only
  ```
- [x] Add example: "Return a defective item"
  ```
  amazon-cli returns create 123-4567890-1234567 ITEM123 --reason defective --confirm
  ```
- [x] Add example: "Skip next Subscribe - [ ] Add example: "Skip next Subscribe & Save delivery" Save delivery"
  ```
  amazon-cli subscriptions skip S01-1234567-8901234 --confirm
  ```
- [x] Add example: "Buy an item immediately"
  ```
  amazon-cli buy B08N5WRWNW --quantity 1 --confirm
  ```

#### 10.5 Safety & Error Documentation
- [x] Document all error codes and their meanings
- [x] Document the `--confirm` requirement for purchase operations
- [x] Add safety guidelines:
  - [ ] Always preview before purchasing (omit --confirm first)
  - [ ] Verify cart contents before checkout
  - [ ] Check subscription changes before confirming
- [x] Document rate limiting behavior
- [x] Document authentication expiry handling

#### 10.6 Final Review & Testing
- [x] Validate skills.md against ClawdHub schema
- [x] Test all example commands work as documented
- [x] Verify all JSON schemas match actual output
- [x] Submit to ClawdHub for review/publication

## Security Considerations

1. **Credentials**: Stored in plain text config file; users are responsible for file permissions
2. **--confirm flag**: Required for all purchase/modification actions to prevent accidental execution
3. **No credential logging**: Tokens never appear in verbose output
4. **HTTPS only**: All Amazon communication over TLS

## Success Metrics

1. All core commands functional against amazon.com
2. JSON output parseable by AI agents without errors
3. Rate limiting prevents account blocks
4. < 5 second response time for typical commands
5. Published on ClawdHub with working skills.md

## Open Questions

1. **Amazon API access**: Does Amazon provide official API access, or will this require web scraping? (Likely scraping for most features)
2. **2FA handling**: How to handle accounts with two-factor authentication enabled?
3. **CAPTCHA**: Strategy for handling CAPTCHA challenges?
4. **Terms of Service**: Review Amazon ToS for automation compliance

## Appendix: Command Reference Quick Sheet

```bash
# Authentication
amazon-cli auth login
amazon-cli auth status
amazon-cli auth logout

# Orders
amazon-cli orders list
amazon-cli orders get <order-id>
amazon-cli orders track <order-id>
amazon-cli orders history --year 2024

# Returns
amazon-cli returns list
amazon-cli returns options <order-id> <item-id>
amazon-cli returns create <order-id> <item-id> --reason defective --confirm
amazon-cli returns label <return-id>
amazon-cli returns status <return-id>

# Search & Products
amazon-cli search "query" --prime-only
amazon-cli product get <asin>
amazon-cli product reviews <asin>

# Cart & Checkout
amazon-cli cart add <asin>
amazon-cli cart list
amazon-cli cart remove <asin>
amazon-cli cart checkout --confirm
amazon-cli buy <asin> --confirm

# Subscriptions
amazon-cli subscriptions list
amazon-cli subscriptions get <subscription-id>
amazon-cli subscriptions skip <subscription-id> --confirm
amazon-cli subscriptions frequency <subscription-id> --interval 4 --confirm
amazon-cli subscriptions cancel <subscription-id> --confirm
amazon-cli subscriptions upcoming
```
