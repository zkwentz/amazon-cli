# Amazon CLI - Detailed Task List for Ralphy

> **NOTE**: This task list is designed for AI-powered execution with ralphy (https://github.com/michaelshimeles/ralphy)
> Each task is granular, specific, and independently executable.

---

## Phase 1: Core Infrastructure Completion

### 1.1 Amazon API Integration Research

#### 1.1.1 Research Amazon Shopping API
- [ ] Search for "Amazon Shopping API documentation" and read official docs
- [ ] Check if Amazon MWS (Marketplace Web Service) provides needed functionality
- [ ] Research Amazon Advertising API for product search capabilities
- [ ] Document API endpoints, rate limits, and authentication requirements in `docs/amazon-api-research.md`
- [ ] Test API access with example requests using curl
- [ ] Document pricing/cost implications of API usage

#### 1.1.2 Research Amazon OAuth/Authentication
- [ ] Visit https://developer.amazon.com/ and explore Login with Amazon (LWA)
- [ ] Document OAuth 2.0 flow for Amazon login
- [ ] Identify required OAuth scopes for order access
- [ ] Test OAuth flow with a test application
- [ ] Document refresh token mechanics
- [ ] Research token expiration policies
- [ ] Create `docs/authentication-strategy.md` with detailed auth flow diagrams

#### 1.1.3 Research Web Scraping Approach (Fallback)
- [ ] Analyze Amazon.com order history page HTML structure
- [ ] Document CSS selectors for order information
- [ ] Analyze cart page structure at amazon.com/cart
- [ ] Document product search results page structure
- [ ] Identify CSRF token locations and extraction methods
- [ ] Document session cookie requirements
- [ ] Test basic HTTP requests to Amazon with curl
- [ ] Identify anti-scraping measures (rate limits, CAPTCHAs)
- [ ] Document User-Agent requirements

#### 1.1.4 Rate Limiting Analysis
- [ ] Make test requests to Amazon at various frequencies
- [ ] Document when rate limiting kicks in (requests per minute)
- [ ] Test different User-Agent strings
- [ ] Identify IP-based vs session-based rate limiting
- [ ] Document CAPTCHA trigger conditions
- [ ] Create rate limiting strategy in `docs/rate-limiting-strategy.md`

#### 1.1.5 Decision Document
- [ ] Create `docs/implementation-decision.md`
- [ ] Document chosen approach (API vs scraping vs hybrid)
- [ ] List pros/cons of chosen approach
- [ ] Document fallback strategies
- [ ] Get decision reviewed and approved

---

### 1.2 Authentication System Implementation

#### 1.2.1 Setup Authentication Package Structure
- [ ] Create `internal/amazon/auth.go` file
- [ ] Define `AuthTokens` struct with `AccessToken`, `RefreshToken`, `ExpiresAt` fields
- [ ] Define OAuth constants: `authURL`, `tokenURL`, `redirectURI`
- [ ] Add `clientID` and `clientSecret` configuration support
- [ ] Create `internal/amazon/auth_test.go` file

#### 1.2.2 Implement OAuth Flow (if using OAuth)
- [ ] Implement `GenerateAuthURL(state string) string` function
- [ ] Implement `StartLocalServer(port int) (*http.Server, chan string, error)` for callback
- [ ] Implement `HandleOAuthCallback(code, state string) (*AuthTokens, error)` function
- [ ] Implement `ExchangeCodeForTokens(code string) (*AuthTokens, error)` function
- [ ] Add CSRF state validation logic
- [ ] Add error handling for auth failures
- [ ] Test with real Amazon developer account

#### 1.2.3 Implement Browser-Based Auth (Alternative)
- [ ] Research using https://github.com/go-rod/rod for browser automation
- [ ] Implement `OpenBrowserLogin() error` function
- [ ] Implement cookie extraction from browser session
- [ ] Store session cookies securely in config
- [ ] Implement cookie validation check
- [ ] Add timeout handling (2 minutes)

#### 1.2.4 Implement Token Refresh Logic
- [ ] Create `RefreshTokenIfNeeded(config *Config) error` function in `internal/amazon/auth.go`
- [ ] Check if token expires within 5 minutes
- [ ] Implement refresh token exchange HTTP request
- [ ] Update config with new tokens
- [ ] Save updated config to disk
- [ ] Add retry logic for refresh failures
- [ ] Write unit tests for refresh logic

#### 1.2.5 Implement Auth Commands
- [ ] Open `cmd/auth.go`
- [ ] Implement `authLoginCmd` Run function:
  - [ ] Call authentication flow
  - [ ] Handle browser opening
  - [ ] Wait for callback/cookies
  - [ ] Save credentials to config
  - [ ] Output success JSON
- [ ] Implement `authStatusCmd` Run function:
  - [ ] Load config
  - [ ] Check for tokens
  - [ ] Calculate expiry time
  - [ ] Output status JSON
- [ ] Implement `authLogoutCmd` Run function:
  - [ ] Load config
  - [ ] Clear auth section
  - [ ] Save config
  - [ ] Output confirmation JSON

#### 1.2.6 Implement Secure Credential Storage
- [ ] Open `internal/config/config.go`
- [ ] Create `AuthConfig` struct with `AccessToken`, `RefreshToken`, `ExpiresAt` fields
- [ ] Implement `SaveConfig(config *Config, path string) error` function
- [ ] Set file permissions to 0600 after writing
- [ ] Verify parent directory exists, create if needed with 0700
- [ ] Add validation for config structure
- [ ] Implement `LoadConfig(path string) (*Config, error)` function
- [ ] Handle missing config file (first run)
- [ ] Write unit tests for config load/save

#### 1.2.7 Authentication Testing
- [ ] Create `internal/amazon/auth_test.go`
- [ ] Write test `TestGenerateAuthURL` - verify URL format
- [ ] Write test `TestTokenRefresh_NotExpired` - should not refresh
- [ ] Write test `TestTokenRefresh_Expired` - should refresh
- [ ] Write test `TestTokenRefresh_InvalidRefreshToken` - should error
- [ ] Write test `TestOAuthCallback_ValidState` - should succeed
- [ ] Write test `TestOAuthCallback_InvalidState` - should fail
- [ ] Create integration test with mock OAuth server
- [ ] Test full login flow end-to-end
- [ ] Verify coverage >= 85% with `go test -cover ./internal/amazon`

---

### 1.3 HTTP Client & Rate Limiting Enhancement

#### 1.3.1 Implement Rate Limiter
- [ ] Create `internal/ratelimit/limiter.go`
- [ ] Define `RateLimiter` struct with fields:
  - [ ] `minDelay time.Duration`
  - [ ] `maxDelay time.Duration`
  - [ ] `lastRequest time.Time`
  - [ ] `mutex sync.Mutex`
  - [ ] `maxRetries int`
- [ ] Implement `NewRateLimiter(minDelay, maxDelay time.Duration, maxRetries int) *RateLimiter`
- [ ] Implement `Wait() error` function:
  - [ ] Lock mutex
  - [ ] Calculate time since last request
  - [ ] If < minDelay, sleep for difference
  - [ ] Add jitter: `rand.Intn(500) * time.Millisecond`
  - [ ] Update lastRequest time
  - [ ] Unlock mutex
- [ ] Implement `WaitWithBackoff(attempt int) error` function:
  - [ ] Calculate backoff: `min(2^attempt * 1000ms, 60000ms)`
  - [ ] Sleep for calculated duration
  - [ ] Log backoff if verbose
- [ ] Implement `ShouldRetry(statusCode int, attempt int) bool` function:
  - [ ] Return true if statusCode == 429 or 503
  - [ ] Return false if attempt >= maxRetries
  - [ ] Return false otherwise

#### 1.3.2 Test Rate Limiter
- [ ] Create `internal/ratelimit/limiter_test.go`
- [ ] Write test `TestRateLimiter_Wait_EnforcesMinDelay`:
  - [ ] Create limiter with 1 second min delay
  - [ ] Call Wait() twice rapidly
  - [ ] Verify second call takes >= 1 second
- [ ] Write test `TestRateLimiter_Jitter_IsRandom`:
  - [ ] Call Wait() multiple times
  - [ ] Verify delays are not identical (jitter working)
- [ ] Write test `TestRateLimiter_WaitWithBackoff_Exponential`:
  - [ ] Verify attempt 0 = 1s, attempt 1 = 2s, attempt 2 = 4s, etc.
- [ ] Write test `TestRateLimiter_ShouldRetry_429`:
  - [ ] Verify returns true for 429 within retry limit
- [ ] Write test `TestRateLimiter_ShouldRetry_MaxRetriesExceeded`:
  - [ ] Verify returns false when attempts >= maxRetries
- [ ] Run tests: `go test -v ./internal/ratelimit`
- [ ] Verify coverage >= 90%

#### 1.3.3 Implement Production HTTP Client
- [ ] Open `internal/amazon/client.go`
- [ ] Add User-Agent list (10+ common browsers):
  - [ ] Chrome on Windows
  - [ ] Chrome on macOS
  - [ ] Firefox on Windows
  - [ ] Firefox on macOS
  - [ ] Safari on macOS
  - [ ] Edge on Windows
  - [ ] Mobile Chrome
  - [ ] Mobile Safari
- [ ] Implement `NewClient() *Client`:
  - [ ] Create http.Client with 30 second timeout
  - [ ] Create cookie jar
  - [ ] Initialize rate limiter
  - [ ] Set baseURL to "https://www.amazon.com"
- [ ] Implement `Do(req *http.Request) (*http.Response, error)`:
  - [ ] Call rate limiter Wait()
  - [ ] Set random User-Agent from list
  - [ ] Set Accept header
  - [ ] Set Accept-Language header
  - [ ] Execute request
  - [ ] Check status code
  - [ ] If 429/503, check ShouldRetry
  - [ ] If should retry, call WaitWithBackoff and retry
  - [ ] Return response or error
- [ ] Implement `Get(url string) (*http.Response, error)` wrapper
- [ ] Implement `PostForm(url string, data url.Values) (*http.Response, error)` wrapper

#### 1.3.4 Add Request/Response Logging
- [ ] Add `logRequest(req *http.Request)` function:
  - [ ] Log method, URL
  - [ ] Log headers (redact auth)
  - [ ] Only log in verbose mode
- [ ] Add `logResponse(resp *http.Response)` function:
  - [ ] Log status code
  - [ ] Log response headers
  - [ ] Only log in verbose mode
- [ ] Integrate logging into `Do()` method

#### 1.3.5 Implement Circuit Breaker
- [ ] Add `circuitBreaker` struct to Client:
  - [ ] `failureCount int`
  - [ ] `lastFailure time.Time`
  - [ ] `threshold int` (default 5)
  - [ ] `resetTimeout time.Duration` (default 60s)
- [ ] Implement `checkCircuitBreaker() error`:
  - [ ] If failureCount >= threshold, return error
  - [ ] If time since lastFailure > resetTimeout, reset counter
- [ ] Implement `recordFailure()`:
  - [ ] Increment failureCount
  - [ ] Update lastFailure time
- [ ] Implement `recordSuccess()`:
  - [ ] Reset failureCount to 0
- [ ] Integrate into `Do()` method

#### 1.3.6 Add CAPTCHA Detection
- [ ] Implement `detectCAPTCHA(resp *http.Response) bool`:
  - [ ] Check if response contains "captcha" in body
  - [ ] Check for specific Amazon CAPTCHA HTML markers
- [ ] Return custom error when CAPTCHA detected
- [ ] Add CAPTCHA error to `pkg/models/errors.go`

#### 1.3.7 Test HTTP Client
- [ ] Create `internal/amazon/client_test.go`
- [ ] Create `internal/testutil/mock_server.go` with mock HTTP server
- [ ] Write test `TestClient_Do_Success`:
  - [ ] Mock server returns 200
  - [ ] Verify request succeeds
- [ ] Write test `TestClient_Do_Retry_429`:
  - [ ] Mock server returns 429, then 200
  - [ ] Verify retry happens
- [ ] Write test `TestClient_Do_Retry_503`:
  - [ ] Mock server returns 503, then 200
  - [ ] Verify retry happens
- [ ] Write test `TestClient_Do_MaxRetries`:
  - [ ] Mock server always returns 429
  - [ ] Verify stops after max retries
- [ ] Write test `TestClient_Do_CircuitBreaker`:
  - [ ] Simulate 5 consecutive failures
  - [ ] Verify circuit breaker opens
- [ ] Write test `TestClient_Do_UserAgentRotation`:
  - [ ] Make multiple requests
  - [ ] Verify User-Agent changes
- [ ] Write test `TestClient_DetectCAPTCHA`:
  - [ ] Mock server returns CAPTCHA HTML
  - [ ] Verify error returned
- [ ] Run tests: `go test -v ./internal/amazon`
- [ ] Verify coverage >= 90%

---

## Phase 2: Orders & Returns Features

### 2.1 Orders Data Models

#### 2.1.1 Create Order Models
- [ ] Open `pkg/models/order.go`
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
- [ ] Define `OrderItem` struct:
  ```go
  type OrderItem struct {
      ASIN      string  `json:"asin"`
      Title     string  `json:"title"`
      Quantity  int     `json:"quantity"`
      Price     float64 `json:"price"`
  }
  ```
- [ ] Define `Tracking` struct:
  ```go
  type Tracking struct {
      Carrier        string `json:"carrier"`
      TrackingNumber string `json:"tracking_number"`
      Status         string `json:"status"`
      DeliveryDate   string `json:"delivery_date"`
  }
  ```
- [ ] Define `OrdersResponse` struct:
  ```go
  type OrdersResponse struct {
      Orders     []Order `json:"orders"`
      TotalCount int     `json:"total_count"`
  }
  ```

---

### 2.2 Orders Implementation

#### 2.2.1 Research Amazon Order Pages
- [ ] Login to Amazon.com manually
- [ ] Navigate to order history page
- [ ] Save HTML source to `testdata/orders/order_list.html`
- [ ] Anonymize personal data (names, addresses, order IDs)
- [ ] Document URL pattern for order history
- [ ] Navigate to single order detail page
- [ ] Save HTML to `testdata/orders/order_detail.html`
- [ ] Document CSS selectors for:
  - [ ] Order ID
  - [ ] Order date
  - [ ] Order total
  - [ ] Order status
  - [ ] Item ASIN
  - [ ] Item title
  - [ ] Item price
  - [ ] Item quantity
- [ ] Document pagination mechanism

#### 2.2.2 Implement Order HTML Parser
- [ ] Open `internal/amazon/orders.go`
- [ ] Install goquery: `go get github.com/PuerkitoBio/goquery`
- [ ] Implement `parseOrderList(doc *goquery.Document) ([]models.Order, error)`:
  - [ ] Find order containers using CSS selector
  - [ ] For each order, extract:
    - [ ] Order ID from data attribute or link
    - [ ] Date from date element
    - [ ] Total from price element
    - [ ] Status from status badge
  - [ ] Return slice of orders
- [ ] Implement `parseOrderDetail(doc *goquery.Document) (*models.Order, error)`:
  - [ ] Extract order header info (ID, date, total, status)
  - [ ] Find items list
  - [ ] For each item, extract ASIN, title, price, quantity
  - [ ] Extract tracking info if present
  - [ ] Return complete Order struct
- [ ] Implement `parseTracking(doc *goquery.Document) (*models.Tracking, error)`:
  - [ ] Extract carrier name
  - [ ] Extract tracking number
  - [ ] Extract delivery status
  - [ ] Extract estimated/actual delivery date
  - [ ] Return Tracking struct

#### 2.2.3 Implement GetOrders
- [ ] In `internal/amazon/orders.go`, implement `GetOrders(limit int, status string) (*models.OrdersResponse, error)`:
  - [ ] Build order history URL
  - [ ] Make GET request using client
  - [ ] Parse response body with goquery
  - [ ] Call parseOrderList
  - [ ] Filter by status if provided
  - [ ] Limit results to requested count
  - [ ] Return OrdersResponse

#### 2.2.4 Implement GetOrder
- [ ] Implement `GetOrder(orderID string) (*models.Order, error)`:
  - [ ] Build order detail URL with orderID
  - [ ] Make GET request
  - [ ] Parse response with goquery
  - [ ] Call parseOrderDetail
  - [ ] Return Order

#### 2.2.5 Implement GetOrderTracking
- [ ] Implement `GetOrderTracking(orderID string) (*models.Tracking, error)`:
  - [ ] Build tracking URL with orderID
  - [ ] Make GET request
  - [ ] Parse response
  - [ ] Call parseTracking
  - [ ] Return Tracking

#### 2.2.6 Implement GetOrderHistory
- [ ] Implement `GetOrderHistory(year int) (*models.OrdersResponse, error)`:
  - [ ] Build URL with year parameter
  - [ ] Make GET request
  - [ ] Handle pagination:
    - [ ] Parse first page
    - [ ] Check for next page link
    - [ ] Request next pages
    - [ ] Combine results
  - [ ] Return all orders for year

#### 2.2.7 Wire Up Orders Commands
- [ ] Open `cmd/orders.go`
- [ ] In `ordersListCmd` Run function:
  - [ ] Get limit and status flags
  - [ ] Create Amazon client
  - [ ] Call client.GetOrders(limit, status)
  - [ ] Handle errors
  - [ ] Output JSON using output.JSON()
- [ ] In `ordersGetCmd` Run function:
  - [ ] Validate orderID argument provided
  - [ ] Create client
  - [ ] Call client.GetOrder(orderID)
  - [ ] Output JSON
- [ ] In `ordersTrackCmd` Run function:
  - [ ] Validate orderID
  - [ ] Create client
  - [ ] Call client.GetOrderTracking(orderID)
  - [ ] Output JSON
- [ ] In `ordersHistoryCmd` Run function:
  - [ ] Get year flag (default current year)
  - [ ] Create client
  - [ ] Call client.GetOrderHistory(year)
  - [ ] Output JSON

---

### 2.3 Orders Testing

#### 2.3.1 Create Test Fixtures
- [ ] Create directory `testdata/orders/`
- [ ] Save real Amazon order list HTML to `testdata/orders/order_list.html`
- [ ] Anonymize all personal data
- [ ] Create `testdata/orders/order_detail.html` with single order
- [ ] Create `testdata/orders/empty_orders.html` for zero orders
- [ ] Create `testdata/orders/order_tracking.html` with tracking info
- [ ] Create `testdata/orders/order_canceled.html` for canceled order

#### 2.3.2 Write Parser Tests
- [ ] Create `internal/amazon/orders_test.go`
- [ ] Write test `TestParseOrderList`:
  - [ ] Load fixture `testdata/orders/order_list.html`
  - [ ] Parse with goquery
  - [ ] Call parseOrderList
  - [ ] Verify order count
  - [ ] Verify first order has all fields
  - [ ] Verify order ID not empty
  - [ ] Verify total > 0
- [ ] Write test `TestParseOrderList_Empty`:
  - [ ] Load empty_orders.html
  - [ ] Verify returns empty slice, no error
- [ ] Write test `TestParseOrderDetail`:
  - [ ] Load order_detail.html
  - [ ] Call parseOrderDetail
  - [ ] Verify all order fields populated
  - [ ] Verify items array has elements
  - [ ] Verify each item has ASIN, title, price
- [ ] Write test `TestParseTracking`:
  - [ ] Load order_tracking.html
  - [ ] Call parseTracking
  - [ ] Verify carrier, tracking number, status, date

#### 2.3.3 Write Integration Tests
- [ ] Write test `TestGetOrders_Integration`:
  - [ ] Create mock HTTP server
  - [ ] Serve order_list.html fixture
  - [ ] Create client pointed at mock server
  - [ ] Call GetOrders(10, "")
  - [ ] Verify response structure
- [ ] Write test `TestGetOrder_Integration`:
  - [ ] Mock server serves order_detail.html
  - [ ] Call GetOrder("123-456-789")
  - [ ] Verify order details
- [ ] Write test `TestGetOrderTracking_Integration`:
  - [ ] Mock server serves order_tracking.html
  - [ ] Call GetOrderTracking("123-456-789")
  - [ ] Verify tracking info

#### 2.3.4 Write Error Tests
- [ ] Write test `TestGetOrders_NetworkError`:
  - [ ] Mock server returns error/timeout
  - [ ] Verify error returned
- [ ] Write test `TestGetOrders_AuthRequired`:
  - [ ] Mock server returns 401
  - [ ] Verify AUTH_REQUIRED error
- [ ] Write test `TestGetOrder_NotFound`:
  - [ ] Mock server returns 404
  - [ ] Verify NOT_FOUND error
- [ ] Write test `TestGetOrders_MalformedHTML`:
  - [ ] Mock server returns invalid HTML
  - [ ] Verify graceful error handling

#### 2.3.5 Run Orders Tests
- [ ] Run: `go test -v ./internal/amazon -run Order`
- [ ] Verify all tests pass
- [ ] Run: `go test -cover ./internal/amazon -run Order`
- [ ] Verify coverage >= 80%

---

### 2.4 Returns Implementation

#### 2.4.1 Create Returns Models
- [ ] Create `pkg/models/return.go`
- [ ] Define `ReturnableItem` struct:
  ```go
  type ReturnableItem struct {
      OrderID      string  `json:"order_id"`
      ItemID       string  `json:"item_id"`
      ASIN         string  `json:"asin"`
      Title        string  `json:"title"`
      Price        float64 `json:"price"`
      PurchaseDate string  `json:"purchase_date"`
      ReturnWindow string  `json:"return_window"`
  }
  ```
- [ ] Define `ReturnOption` struct:
  ```go
  type ReturnOption struct {
      Method          string  `json:"method"`
      Label           string  `json:"label"`
      DropoffLocation string  `json:"dropoff_location"`
      Fee             float64 `json:"fee"`
  }
  ```
- [ ] Define `Return` struct:
  ```go
  type Return struct {
      ReturnID  string `json:"return_id"`
      OrderID   string `json:"order_id"`
      ItemID    string `json:"item_id"`
      Status    string `json:"status"`
      Reason    string `json:"reason"`
      CreatedAt string `json:"created_at"`
  }
  ```
- [ ] Define `ReturnLabel` struct:
  ```go
  type ReturnLabel struct {
      URL          string `json:"url"`
      Carrier      string `json:"carrier"`
      Instructions string `json:"instructions"`
  }
  ```

#### 2.4.2 Research Amazon Returns Pages
- [ ] Navigate to Amazon returns center
- [ ] Save HTML to `testdata/returns/returnable_items.html`
- [ ] Click on a returnable item
- [ ] Save return options page to `testdata/returns/return_options.html`
- [ ] Document CSS selectors for returnable items
- [ ] Document return reason codes Amazon uses
- [ ] Document return options structure
- [ ] Document return confirmation page structure

#### 2.4.3 Implement Returns Client
- [ ] Create `internal/amazon/returns.go`
- [ ] Implement `GetReturnableItems() ([]models.ReturnableItem, error)`:
  - [ ] Build returns center URL
  - [ ] Make GET request
  - [ ] Parse HTML with goquery
  - [ ] Extract returnable items
  - [ ] Return slice of ReturnableItem
- [ ] Implement `GetReturnOptions(orderID, itemID string) ([]models.ReturnOption, error)`:
  - [ ] Build return options URL
  - [ ] Make GET request with orderID, itemID
  - [ ] Parse return options
  - [ ] Return slice of ReturnOption
- [ ] Implement `CreateReturn(orderID, itemID, reason string) (*models.Return, error)`:
  - [ ] Validate reason code
  - [ ] Build return creation form data
  - [ ] Make POST request
  - [ ] Parse confirmation response
  - [ ] Extract return ID
  - [ ] Return Return struct
- [ ] Implement `GetReturnLabel(returnID string) (*models.ReturnLabel, error)`:
  - [ ] Build label URL
  - [ ] Make GET request
  - [ ] Parse label info
  - [ ] Return ReturnLabel
- [ ] Implement `GetReturnStatus(returnID string) (*models.Return, error)`:
  - [ ] Build status URL
  - [ ] Make GET request
  - [ ] Parse status
  - [ ] Return Return

#### 2.4.4 Implement Returns Commands
- [ ] Create `cmd/returns.go`
- [ ] Add `returns` parent command
- [ ] Implement `returnsListCmd`:
  - [ ] Call client.GetReturnableItems()
  - [ ] Output JSON
- [ ] Implement `returnsOptionsCmd`:
  - [ ] Validate orderID and itemID args
  - [ ] Call client.GetReturnOptions
  - [ ] Output JSON
- [ ] Implement `returnsCreateCmd`:
  - [ ] Add --reason flag (required)
  - [ ] Add --confirm flag (required)
  - [ ] Validate reason in allowed list
  - [ ] If no --confirm, output dry run preview
  - [ ] If --confirm, call client.CreateReturn
  - [ ] Output confirmation JSON
- [ ] Implement `returnsLabelCmd`:
  - [ ] Validate returnID
  - [ ] Call client.GetReturnLabel
  - [ ] Output JSON
- [ ] Implement `returnsStatusCmd`:
  - [ ] Validate returnID
  - [ ] Call client.GetReturnStatus
  - [ ] Output JSON
- [ ] Wire up commands in `init()`:
  - [ ] Add to rootCmd
  - [ ] Add flags

#### 2.4.5 Test Returns
- [ ] Create `internal/amazon/returns_test.go`
- [ ] Write parser tests with fixtures
- [ ] Write integration tests with mock server
- [ ] Write test for invalid reason codes
- [ ] Write test for --confirm flag validation
- [ ] Run: `go test -v ./internal/amazon -run Return`
- [ ] Verify coverage >= 80%

---

## Phase 3: Search & Product Features

### 3.1 Product Models
- [ ] Open `pkg/models/product.go`
- [ ] Define `Product` struct with all fields
- [ ] Define `SearchResponse` struct
- [ ] Define `Review` struct
- [ ] Define `ReviewsResponse` struct

### 3.2 Search Implementation

#### 3.2.1 Research Amazon Search
- [ ] Navigate to Amazon search page
- [ ] Search for "wireless headphones"
- [ ] Save HTML to `testdata/search/search_results.html`
- [ ] Document URL parameters (k, i, rh, etc.)
- [ ] Document CSS selectors for products
- [ ] Test different filters (price, prime)
- [ ] Document pagination structure

#### 3.2.2 Implement Search Parser
- [ ] Create `internal/amazon/search.go`
- [ ] Define `SearchOptions` struct
- [ ] Implement `parseSearchResults(doc *goquery.Document) ([]models.Product, error)`
- [ ] Extract: ASIN, title, price, rating, prime badge, stock

#### 3.2.3 Implement Search Function
- [ ] Implement `Search(query string, opts SearchOptions) (*models.SearchResponse, error)`:
  - [ ] Build search URL with parameters
  - [ ] Add filters (category, price range, prime)
  - [ ] Make GET request
  - [ ] Parse results
  - [ ] Return SearchResponse

#### 3.2.4 Test Search
- [ ] Create test fixtures in `testdata/search/`
- [ ] Write parser tests
- [ ] Write integration tests
- [ ] Test filter combinations
- [ ] Run: `go test -v ./internal/amazon -run Search`

---

### 3.3 Product Implementation

#### 3.3.1 Research Product Pages
- [ ] Navigate to product detail page
- [ ] Save HTML to `testdata/products/product_detail.html`
- [ ] Document selectors for all fields
- [ ] Navigate to reviews page
- [ ] Save to `testdata/products/reviews.html`

#### 3.3.2 Implement Product Parser
- [ ] Create `internal/amazon/product.go`
- [ ] Implement `parseProductDetail(doc *goquery.Document) (*models.Product, error)`
- [ ] Extract all product fields
- [ ] Handle missing optional fields gracefully

#### 3.3.3 Implement Product Functions
- [ ] Implement `GetProduct(asin string) (*models.Product, error)`
- [ ] Implement `GetProductReviews(asin string, limit int) (*models.ReviewsResponse, error)`

#### 3.3.4 Test Products
- [ ] Create fixtures
- [ ] Write parser tests
- [ ] Test with missing data
- [ ] Test out-of-stock products
- [ ] Run: `go test -v ./internal/amazon -run Product`

---

### 3.4 Search & Product Commands

#### 3.4.1 Implement Search Command
- [ ] Create `cmd/search.go`
- [ ] Add flags: --category, --min-price, --max-price, --prime-only, --page
- [ ] Implement Run function
- [ ] Call client.Search
- [ ] Output JSON

#### 3.4.2 Implement Product Commands
- [ ] Open `cmd/product.go`
- [ ] Implement `product get` command
- [ ] Add ASIN validation
- [ ] Implement `product reviews` command
- [ ] Add --limit flag

#### 3.4.3 Test Commands
- [ ] Test search with various queries
- [ ] Test all filter combinations
- [ ] Test ASIN validation
- [ ] Verify JSON output format

---

## Phase 4: Cart & Checkout

### 4.1 Cart Models
- [ ] Verify `pkg/models/cart.go` has all needed structs
- [ ] Add `CheckoutPreview` struct if missing
- [ ] Add `OrderConfirmation` struct if missing

### 4.2 Cart Implementation

#### 4.2.1 Research Amazon Cart
- [ ] Navigate to amazon.com/cart
- [ ] Save HTML to `testdata/cart/cart_with_items.html`
- [ ] Document add to cart endpoint
- [ ] Document CSRF token location
- [ ] Test add to cart flow manually
- [ ] Save empty cart to `testdata/cart/cart_empty.html`

#### 4.2.2 Implement Real Cart Operations
- [ ] Open `internal/amazon/cart.go`
- [ ] Replace mock `AddToCart` with real implementation:
  - [ ] Extract CSRF token
  - [ ] Build form data
  - [ ] POST to add-to-cart endpoint
  - [ ] Parse response
  - [ ] Return updated cart
- [ ] Replace mock `GetCart`:
  - [ ] GET cart page
  - [ ] Parse cart items
  - [ ] Calculate totals
  - [ ] Return Cart struct
- [ ] Replace mock `RemoveFromCart`:
  - [ ] Build remove request
  - [ ] POST to remove endpoint
  - [ ] Return updated cart
- [ ] Replace mock `ClearCart`:
  - [ ] Iterate over cart items
  - [ ] Remove each item
  - [ ] Return success

#### 4.2.3 Implement Address/Payment Methods
- [ ] Replace mock `GetAddresses`:
  - [ ] Navigate to address management
  - [ ] Parse saved addresses
  - [ ] Return slice of Address
- [ ] Replace mock `GetPaymentMethods`:
  - [ ] Navigate to payment methods
  - [ ] Parse saved payment methods
  - [ ] Return slice of PaymentMethod

#### 4.2.4 Implement Checkout Preview
- [ ] Replace mock `PreviewCheckout`:
  - [ ] Build checkout initiation request
  - [ ] Submit with addressID and paymentID
  - [ ] Parse checkout preview page
  - [ ] Extract totals, delivery estimates
  - [ ] Return CheckoutPreview

#### 4.2.5 Implement Complete Checkout
- [ ] **CRITICAL**: Review mock implementation in `CompleteCheckout`
- [ ] **DO NOT** implement real checkout against production
- [ ] Keep mock implementation for safety
- [ ] Add extensive validation:
  - [ ] Verify --confirm flag handled in command layer
  - [ ] Verify cart not empty
  - [ ] Verify address exists
  - [ ] Verify payment method exists
- [ ] Add TODO comment for production implementation
- [ ] Document that real implementation requires:
  - [ ] Test Amazon account
  - [ ] Sandbox environment
  - [ ] Or complete mock server

---

### 4.3 Cart Testing (CRITICAL)

#### 4.3.1 Expand Cart Tests
- [ ] Open `internal/amazon/cart_test.go`
- [ ] Keep existing tests (they're good)
- [ ] Add test `TestAddToCart_RealParser`:
  - [ ] Load fixture with cart HTML
  - [ ] Test parsing actual cart response
- [ ] Add test `TestGetCart_WithMultipleItems`:
  - [ ] Parse cart with 5+ items
  - [ ] Verify totals calculated correctly
- [ ] Add test `TestCart_CSRFTokenExtraction`:
  - [ ] Mock cart page with CSRF token
  - [ ] Verify token extracted correctly
- [ ] Add test `TestCart_QuantityLimits`:
  - [ ] Test adding quantity > 10
  - [ ] Verify error or capping

#### 4.3.2 Checkout Safety Tests
- [ ] Write test `TestCheckout_RequiresConfirmFlag`:
  - [ ] Verify command layer checks --confirm
  - [ ] Verify preview mode when --confirm missing
- [ ] Write test `TestCheckout_EmptyCart`:
  - [ ] Verify error when cart empty
- [ ] Write test `TestCheckout_PreviewNeverSubmits`:
  - [ ] Mock server tracks if order submitted
  - [ ] Call PreviewCheckout
  - [ ] Verify no order submission
- [ ] Write test `TestCheckout_MockOnly`:
  - [ ] Document that CompleteCheckout is mock-only
  - [ ] Verify it returns mock order ID
  - [ ] Verify it never makes real HTTP POST

#### 4.3.3 Integration Tests with Mock Server
- [ ] Create mock server in `internal/testutil/mock_amazon.go`
- [ ] Add routes for:
  - [ ] /cart (GET - return cart)
  - [ ] /cart/add (POST - add item)
  - [ ] /cart/remove (POST - remove item)
  - [ ] /checkout/preview (GET - preview)
- [ ] Write test `TestCartFlow_Integration`:
  - [ ] Add item
  - [ ] Get cart
  - [ ] Remove item
  - [ ] Verify cart empty
- [ ] Write test `TestCheckoutFlow_Integration`:
  - [ ] Add items
  - [ ] Preview checkout
  - [ ] Verify preview data

#### 4.3.4 Run Cart Tests
- [ ] Run: `go test -v ./internal/amazon -run Cart`
- [ ] Run: `go test -v ./internal/amazon -run Checkout`
- [ ] Verify all tests pass
- [ ] Run: `go test -cover ./internal/amazon -run "Cart|Checkout"`
- [ ] Verify coverage >= 90%

---

### 4.4 Buy Command

#### 4.4.1 Implement Buy Command
- [ ] Open `cmd/buy.go`
- [ ] Implement `buyCmd` Run function:
  - [ ] Validate ASIN
  - [ ] Get quantity, addressID, paymentID flags
  - [ ] If no --confirm:
    - [ ] Get product details
    - [ ] Calculate total
    - [ ] Output dry run preview JSON
  - [ ] If --confirm:
    - [ ] Add to cart
    - [ ] Complete checkout
    - [ ] Output OrderConfirmation JSON

#### 4.4.2 Test Buy Command
- [ ] Test without --confirm shows preview
- [ ] Test with --confirm calls checkout
- [ ] Test invalid ASIN rejected

---

## Phase 5: Subscriptions

### 5.1 Subscriptions Models
- [ ] Create `pkg/models/subscription.go`
- [ ] Define `Subscription` struct
- [ ] Define `SubscriptionsResponse` struct
- [ ] Define `UpcomingDelivery` struct

### 5.2 Research Subscribe & Save
- [ ] Navigate to Amazon Subscribe & Save dashboard
- [ ] Save HTML to `testdata/subscriptions/subscription_list.html`
- [ ] Document page structure
- [ ] Test skipping delivery
- [ ] Test changing frequency
- [ ] Document form submissions

### 5.3 Implement Subscriptions Client
- [ ] Create `internal/amazon/subscriptions.go`
- [ ] Implement `GetSubscriptions()`
- [ ] Implement `GetSubscription(id)`
- [ ] Implement `SkipDelivery(id)`
- [ ] Implement `UpdateFrequency(id, weeks)`
- [ ] Implement `CancelSubscription(id)`
- [ ] Implement `GetUpcomingDeliveries()`

### 5.4 Implement Subscriptions Commands
- [ ] Create `cmd/subscriptions.go`
- [ ] Implement all subscription commands
- [ ] Add --confirm flags where needed
- [ ] Wire up to rootCmd

### 5.5 Test Subscriptions
- [ ] Create test fixtures
- [ ] Write parser tests
- [ ] Write integration tests
- [ ] Test frequency validation
- [ ] Run: `go test -v ./internal/amazon -run Subscription`
- [ ] Verify coverage >= 75%

---

## Phase 6: Error Handling & Polish

### 6.1 Input Validation

#### 6.1.1 Create Validation Package
- [ ] Create `internal/validation/validators.go`
- [ ] Implement `ValidateASIN(asin string) error`:
  - [ ] Check length == 10
  - [ ] Check alphanumeric only
  - [ ] Return error if invalid
- [ ] Implement `ValidateOrderID(id string) error`:
  - [ ] Check format: XXX-XXXXXXX-XXXXXXX
  - [ ] Return error if invalid
- [ ] Implement `ValidateQuantity(qty int) error`:
  - [ ] Check > 0
  - [ ] Check <= 999
  - [ ] Return error if invalid
- [ ] Implement `ValidatePriceRange(min, max float64) error`:
  - [ ] Check min >= 0
  - [ ] Check max > min
  - [ ] Return error if invalid
- [ ] Implement `ValidateReturnReason(reason string) error`:
  - [ ] Check against allowed reasons
  - [ ] Return error if invalid

#### 6.1.2 Test Validators
- [ ] Create `internal/validation/validators_test.go`
- [ ] Write test `TestValidateASIN`:
  - [ ] Test valid ASIN
  - [ ] Test too short
  - [ ] Test too long
  - [ ] Test special characters
  - [ ] Test empty string
- [ ] Write test `TestValidateOrderID`
- [ ] Write test `TestValidateQuantity`
- [ ] Write test `TestValidatePriceRange`
- [ ] Write test `TestValidateReturnReason`
- [ ] Run: `go test -v ./internal/validation`
- [ ] Verify coverage == 100%

#### 6.1.3 Integrate Validators
- [ ] In `cmd/cart.go`, validate ASIN before AddToCart
- [ ] In `cmd/orders.go`, validate orderID before GetOrder
- [ ] In `cmd/product.go`, validate ASIN before GetProduct
- [ ] In `cmd/search.go`, validate price range
- [ ] In `cmd/returns.go`, validate reason code

---

### 6.2 Comprehensive Error Handling

#### 6.2.1 Enhance Error Types
- [ ] Open `pkg/models/errors.go`
- [ ] Verify all error codes defined:
  - [ ] AUTH_REQUIRED
  - [ ] AUTH_EXPIRED
  - [ ] NOT_FOUND
  - [ ] RATE_LIMITED
  - [ ] INVALID_INPUT
  - [ ] PURCHASE_FAILED
  - [ ] NETWORK_ERROR
  - [ ] AMAZON_ERROR
  - [ ] CAPTCHA_REQUIRED
- [ ] Add `NewCLIError(code, message string, details map[string]interface{}) *CLIError` constructor
- [ ] Implement `Error() string` method

#### 6.2.2 Add Error Wrapping
- [ ] In `internal/amazon/client.go`, wrap all errors with context:
  - [ ] Network errors → `fmt.Errorf("network request failed: %w", err)`
  - [ ] Parse errors → `fmt.Errorf("failed to parse response: %w", err)`
- [ ] Use `errors.Is()` and `errors.As()` for error checking

#### 6.2.3 Improve Error Messages
- [ ] For AUTH_REQUIRED, suggest: "Run 'amazon-cli auth login' to authenticate"
- [ ] For RATE_LIMITED, suggest: "Wait a few minutes and try again"
- [ ] For CAPTCHA_REQUIRED, suggest: "Visit amazon.com in browser to complete CAPTCHA"
- [ ] For NOT_FOUND, include the resource ID in error

#### 6.2.4 Test Error Handling
- [ ] Write test for each error type
- [ ] Verify error messages are helpful
- [ ] Verify error codes are correct
- [ ] Test error wrapping preserves original error

---

### 6.3 Logging

#### 6.3.1 Implement Structured Logging
- [ ] Add `log/slog` to imports
- [ ] Create logger in `internal/amazon/client.go`:
  - [ ] `logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{}))`
- [ ] Add log levels based on --verbose flag
- [ ] Log requests: `logger.Debug("making request", "url", url)`
- [ ] Log responses: `logger.Debug("received response", "status", resp.StatusCode)`
- [ ] Log rate limiting: `logger.Info("rate limit wait", "duration", delay)`

#### 6.3.2 Redact Sensitive Data
- [ ] Implement `redactHeaders(headers http.Header) http.Header`:
  - [ ] Redact Authorization header
  - [ ] Redact Cookie header
  - [ ] Return safe headers
- [ ] Use in logging

#### 6.3.3 Add Request Timing
- [ ] Before request: `start := time.Now()`
- [ ] After request: `duration := time.Since(start)`
- [ ] Log: `logger.Info("request completed", "duration", duration, "url", url)`

---

## Phase 7: Integration Testing

### 7.1 Mock Amazon Server

#### 7.1.1 Create Mock Server Package
- [ ] Create `internal/testutil/mock_amazon.go`
- [ ] Define `MockAmazonServer` struct
- [ ] Implement `NewMockAmazonServer() *MockAmazonServer`
- [ ] Implement `handler(w http.ResponseWriter, r *http.Request)` with routing:
  - [ ] /orders → serve order list fixture
  - [ ] /orders/:id → serve order detail
  - [ ] /cart → serve cart
  - [ ] /cart/add → handle add item
  - [ ] /search → serve search results
  - [ ] /products/:asin → serve product detail
- [ ] Implement `WithFixture(path, fixture string)` for custom fixtures
- [ ] Implement `Close()` to shutdown server

#### 7.1.2 Create Comprehensive Fixtures
- [ ] Organize `testdata/` directory:
  ```
  testdata/
  ├── auth/
  ├── orders/
  ├── cart/
  ├── search/
  ├── products/
  ├── returns/
  └── subscriptions/
  ```
- [ ] Create fixture for each API response
- [ ] Anonymize all personal data
- [ ] Document what each fixture tests

---

### 7.2 End-to-End Tests

#### 7.2.1 Setup E2E Test Framework
- [ ] Create `test/e2e/` directory
- [ ] Create `test/e2e/e2e_test.go`
- [ ] Add build helper: `buildCLI() (string, error)` to compile binary
- [ ] Add cleanup helper: `defer os.Remove(binaryPath)`

#### 7.2.2 Write E2E Test: Auth Flow
- [ ] Write `TestE2E_AuthFlow`:
  - [ ] Start mock server
  - [ ] Run: `amazon-cli auth status` (expect not authenticated)
  - [ ] Mock login success
  - [ ] Run: `amazon-cli auth status` (expect authenticated)
  - [ ] Run: `amazon-cli auth logout`
  - [ ] Run: `amazon-cli auth status` (expect not authenticated)

#### 7.2.3 Write E2E Test: Order Flow
- [ ] Write `TestE2E_OrderFlow`:
  - [ ] Setup mock with order fixtures
  - [ ] Run: `amazon-cli orders list --limit 5`
  - [ ] Parse JSON output
  - [ ] Verify structure matches schema
  - [ ] Extract first order ID
  - [ ] Run: `amazon-cli orders get <order-id>`
  - [ ] Verify order details
  - [ ] Run: `amazon-cli orders track <order-id>`
  - [ ] Verify tracking info

#### 7.2.4 Write E2E Test: Search to Cart Flow
- [ ] Write `TestE2E_SearchToCartFlow`:
  - [ ] Run: `amazon-cli search "headphones"`
  - [ ] Parse results
  - [ ] Extract first ASIN
  - [ ] Run: `amazon-cli product get <asin>`
  - [ ] Verify product details
  - [ ] Run: `amazon-cli cart add <asin>`
  - [ ] Run: `amazon-cli cart list`
  - [ ] Verify item in cart
  - [ ] Run: `amazon-cli cart checkout` (no --confirm)
  - [ ] Verify preview mode
  - [ ] Run: `amazon-cli cart clear --confirm`

#### 7.2.5 Write E2E Test: Return Flow
- [ ] Write `TestE2E_ReturnFlow`:
  - [ ] Run: `amazon-cli returns list`
  - [ ] Parse returnable items
  - [ ] Extract orderID and itemID
  - [ ] Run: `amazon-cli returns options <order-id> <item-id>`
  - [ ] Run: `amazon-cli returns create <order-id> <item-id> --reason defective` (no --confirm)
  - [ ] Verify dry run

#### 7.2.6 Write E2E Test: Error Recovery
- [ ] Write `TestE2E_ErrorRecovery`:
  - [ ] Mock server returns 401 (auth required)
  - [ ] Verify CLI returns correct error code
  - [ ] Verify error message suggests login
  - [ ] Mock auth refresh
  - [ ] Retry request
  - [ ] Verify success

#### 7.2.7 Write E2E Test: Rate Limiting
- [ ] Write `TestE2E_RateLimiting`:
  - [ ] Mock server tracks request timestamps
  - [ ] Make 5 rapid requests
  - [ ] Verify rate limiting enforced (min 1 sec between)
  - [ ] Mock 429 response
  - [ ] Verify retry with backoff

#### 7.2.8 Run E2E Tests
- [ ] Run: `go test -v ./test/e2e`
- [ ] Verify all workflows pass
- [ ] Run on macOS, Linux, Windows

---

### 7.3 Performance Testing

#### 7.3.1 Create Benchmarks
- [ ] In `internal/amazon/orders_test.go`, add:
  ```go
  func BenchmarkParseOrders_10(b *testing.B) {
      html := loadFixture("order_list_10.html")
      for i := 0; i < b.N; i++ {
          parseOrders(html)
      }
  }
  ```
- [ ] Create benchmarks for 10, 100, 1000 orders
- [ ] Benchmark search result parsing
- [ ] Benchmark cart operations
- [ ] Benchmark JSON marshaling

#### 7.3.2 Run Benchmarks
- [ ] Run: `go test -bench=. -benchmem ./internal/amazon`
- [ ] Document results in `docs/performance.md`
- [ ] Verify:
  - [ ] Order parsing (10 items): < 1ms
  - [ ] Order parsing (100 items): < 10ms
  - [ ] Order parsing (1000 items): < 100ms
  - [ ] Memory allocations reasonable

#### 7.3.3 Profile Performance
- [ ] Run: `go test -cpuprofile=cpu.prof -bench=. ./internal/amazon`
- [ ] Analyze: `go tool pprof cpu.prof`
- [ ] Identify bottlenecks
- [ ] Run: `go test -memprofile=mem.prof -bench=. ./internal/amazon`
- [ ] Analyze memory usage
- [ ] Optimize hot paths if needed

---

## Phase 8: CI/CD & Automation

### 8.1 GitHub Actions - CI Workflow

#### 8.1.1 Create CI Workflow File
- [ ] Create `.github/workflows/ci.yml`
- [ ] Add workflow name: `CI`
- [ ] Add triggers:
  ```yaml
  on:
    push:
      branches: [ main ]
    pull_request:
      branches: [ main ]
  ```

#### 8.1.2 Add Test Job
- [ ] Add job `test` with matrix strategy:
  ```yaml
  strategy:
    matrix:
      os: [ubuntu-latest, macos-latest, windows-latest]
      go-version: [1.25.x]
  ```
- [ ] Add steps:
  - [ ] Checkout code
  - [ ] Setup Go
  - [ ] Cache Go modules
  - [ ] Run tests with coverage
  - [ ] Upload coverage to codecov

#### 8.1.3 Add Lint Job
- [ ] Add job `lint`
- [ ] Use `golangci-lint-action@v4`
- [ ] Configure to fail on warnings

#### 8.1.4 Add Build Job
- [ ] Add job `build`
- [ ] Build for all platforms:
  - [ ] darwin/amd64
  - [ ] darwin/arm64
  - [ ] linux/amd64
  - [ ] linux/arm64
  - [ ] windows/amd64
- [ ] Upload artifacts

#### 8.1.5 Test CI Workflow
- [ ] Commit workflow file
- [ ] Push to branch
- [ ] Open PR
- [ ] Verify CI runs
- [ ] Verify all jobs pass

---

### 8.2 GitHub Actions - Release Workflow

#### 8.2.1 Create Release Workflow
- [ ] Create `.github/workflows/release.yml`
- [ ] Add trigger:
  ```yaml
  on:
    push:
      tags:
        - 'v*'
  ```

#### 8.2.2 Add GoReleaser Job
- [ ] Checkout code
- [ ] Setup Go
- [ ] Run GoReleaser:
  ```yaml
  - uses: goreleaser/goreleaser-action@v5
    with:
      version: latest
      args: release --clean
    env:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      HOMEBREW_TAP_GITHUB_TOKEN: ${{ secrets.HOMEBREW_TAP_GITHUB_TOKEN }}
  ```

#### 8.2.3 Configure GoReleaser
- [ ] Verify `.goreleaser.yml` exists
- [ ] Update version template
- [ ] Configure changelog generation
- [ ] Configure Homebrew formula update

#### 8.2.4 Setup Secrets
- [ ] Create `HOMEBREW_TAP_GITHUB_TOKEN` in repo settings
- [ ] Generate GitHub personal access token
- [ ] Add token to repository secrets

#### 8.2.5 Test Release
- [ ] Create test tag: `git tag v0.1.0-test`
- [ ] Push tag: `git push origin v0.1.0-test`
- [ ] Verify release workflow runs
- [ ] Verify binaries created
- [ ] Verify Homebrew formula updated
- [ ] Delete test tag and release

---

### 8.3 Quality Gates

#### 8.3.1 Create Pre-commit Hook
- [ ] Create `.git/hooks/pre-commit`:
  ```bash
  #!/bin/bash
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
- [ ] Make executable: `chmod +x .git/hooks/pre-commit`

#### 8.3.2 Create Branch Protection Rules
- [ ] Go to GitHub repo Settings → Branches
- [ ] Add rule for `main` branch:
  - [ ] Require PR before merging
  - [ ] Require status checks: `test`, `lint`, `build`
  - [ ] Require branches be up to date
  - [ ] Require conversation resolution
  - [ ] Do not allow bypassing

#### 8.3.3 Create PR Template
- [ ] Create `.github/PULL_REQUEST_TEMPLATE.md`:
  ```markdown
  ## Description

  ## Changes
  -

  ## Testing
  - [ ] Unit tests added/updated
  - [ ] Integration tests added/updated
  - [ ] E2E tests added/updated
  - [ ] Manual testing completed

  ## Checklist
  - [ ] All tests pass
  - [ ] Coverage >= 80%
  - [ ] Documentation updated
  - [ ] CHANGELOG.md updated
  - [ ] No new linter warnings
  ```

#### 8.3.4 Create CODEOWNERS
- [ ] Create `.github/CODEOWNERS`:
  ```
  * @zkwentz

  # Require additional review for critical files
  /internal/amazon/cart.go @zkwentz
  /cmd/cart.go @zkwentz
  /.github/ @zkwentz
  ```

---

### 8.4 Code Coverage

#### 8.4.1 Setup Codecov
- [ ] Sign up at https://codecov.io
- [ ] Connect GitHub repository
- [ ] Get upload token
- [ ] Add token to GitHub secrets: `CODECOV_TOKEN`

#### 8.4.2 Add Coverage Upload to CI
- [ ] In `.github/workflows/ci.yml`, add:
  ```yaml
  - name: Upload coverage
    uses: codecov/codecov-action@v4
    with:
      file: ./coverage.txt
      token: ${{ secrets.CODECOV_TOKEN }}
  ```

#### 8.4.3 Add Coverage Badge
- [ ] Get badge markdown from Codecov
- [ ] Add to `README.md`:
  ```markdown
  [![codecov](https://codecov.io/gh/zkwentz/amazon-cli/branch/main/graph/badge.svg)](https://codecov.io/gh/zkwentz/amazon-cli)
  ```

#### 8.4.4 Configure Coverage Requirements
- [ ] In Codecov settings, set:
  - [ ] Target coverage: 80%
  - [ ] Fail PR if coverage drops
  - [ ] Show coverage diff in PR comments

---

## Phase 9: Documentation & Release

### 9.1 Documentation Files

#### 9.1.1 Create CONTRIBUTING.md
- [ ] Create `CONTRIBUTING.md`:
  - [ ] Development setup instructions
  - [ ] How to run tests
  - [ ] How to build locally
  - [ ] Code style guidelines
  - [ ] How to submit PRs
  - [ ] Where to get help

#### 9.1.2 Create CHANGELOG.md
- [ ] Create `CHANGELOG.md`
- [ ] Add header:
  ```markdown
  # Changelog

  All notable changes to this project will be documented in this file.

  The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
  and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
  ```
- [ ] Add `[Unreleased]` section
- [ ] Document all features in v1.0.0 section

#### 9.1.3 Create SECURITY.md
- [ ] Create `SECURITY.md`:
  - [ ] Supported versions
  - [ ] How to report vulnerabilities
  - [ ] Security best practices:
    - [ ] Config file permissions
    - [ ] Credential rotation
    - [ ] --confirm flag usage
  - [ ] Scope of security (not responsible for Amazon account compromise)

#### 9.1.4 Create Development Guide
- [ ] Create `docs/DEVELOPMENT.md`:
  - [ ] Project structure explanation
  - [ ] How to add new features
  - [ ] Testing guidelines
  - [ ] Debugging tips
  - [ ] Performance profiling
  - [ ] Common issues

#### 9.1.5 Create Troubleshooting Guide
- [ ] Create `docs/TROUBLESHOOTING.md`:
  - [ ] Auth issues
  - [ ] CAPTCHA problems
  - [ ] Rate limiting
  - [ ] Network errors
  - [ ] Common error messages and fixes

#### 9.1.6 Update README.md
- [ ] Add badges:
  - [ ] CI status
  - [ ] Code coverage
  - [ ] Go version
  - [ ] License
- [ ] Verify all examples work
- [ ] Add screenshots/GIFs if possible
- [ ] Add "Star History" graph
- [ ] Add "Contributors" section

---

### 9.2 skills.md for ClawdHub

#### 9.2.1 Research ClawdHub Format
- [ ] Visit ClawdHub documentation
- [ ] Review example skills.md files
- [ ] Identify required metadata fields
- [ ] Document JSON schema format

#### 9.2.2 Create skills.md Header
- [ ] Create `skills.md`
- [ ] Add metadata:
  ```yaml
  ---
  name: amazon-cli
  description: CLI tool for Amazon shopping automation - orders, returns, cart, subscriptions
  version: 1.0.0
  author: zkwentz
  repository: https://github.com/zkwentz/amazon-cli
  tags: [shopping, e-commerce, amazon, automation, ai-agent]
  license: MIT
  ---
  ```

#### 9.2.3 Document Installation
- [ ] Add installation section:
  - [ ] Homebrew instructions
  - [ ] Binary download instructions
  - [ ] Build from source
- [ ] Add quick start example

#### 9.2.4 Document Authentication Actions
- [ ] Document `auth login`:
  - [ ] Purpose
  - [ ] Inputs: none
  - [ ] Output schema
  - [ ] Example command
- [ ] Document `auth status`
- [ ] Document `auth logout`

#### 9.2.5 Document Orders Actions
- [ ] For each orders command, document:
  - [ ] Purpose
  - [ ] Required inputs
  - [ ] Optional inputs with defaults
  - [ ] Output JSON schema
  - [ ] Example command
  - [ ] Example output
- [ ] Commands: list, get, track, history

#### 9.2.6 Document Search & Product Actions
- [ ] Document `search`:
  - [ ] All filter options
  - [ ] Output schema with example
- [ ] Document `product get`
- [ ] Document `product reviews`

#### 9.2.7 Document Cart & Checkout Actions
- [ ] Document `cart add`, `cart list`, `cart remove`, `cart clear`
- [ ] Document `cart checkout`:
  - [ ] **Emphasize --confirm requirement**
  - [ ] Show dry run example
  - [ ] Show confirmed example
  - [ ] Warn about real purchases
- [ ] Document `buy` command with safety warnings

#### 9.2.8 Document Returns Actions
- [ ] Document all returns commands
- [ ] Include reason code reference
- [ ] Show --confirm examples

#### 9.2.9 Document Subscriptions Actions
- [ ] Document all subscription commands
- [ ] Include frequency options

#### 9.2.10 Add AI Agent Examples
- [ ] Add "Common Tasks" section:
  - [ ] "Check my recent orders"
  - [ ] "Track package for order X"
  - [ ] "Search for wireless headphones under $100"
  - [ ] "Return a defective item"
  - [ ] "Skip next subscription delivery"
- [ ] Show complete workflow examples

#### 9.2.11 Add Safety & Error Documentation
- [ ] Document all error codes
- [ ] Explain --confirm flag purpose
- [ ] Add safety guidelines:
  - [ ] Always preview before purchasing
  - [ ] Verify cart before checkout
  - [ ] Check subscription changes
- [ ] Document rate limiting behavior

#### 9.2.12 Validate skills.md
- [ ] Check YAML frontmatter is valid
- [ ] Verify all JSON schemas match actual output
- [ ] Test all example commands
- [ ] Validate against ClawdHub schema (if available)

---

### 9.3 Release Preparation

#### 9.3.1 Version Bump
- [ ] Update version in `main.go` to `1.0.0`
- [ ] Update version in `.goreleaser.yml`
- [ ] Update version in `skills.md`
- [ ] Update version in `README.md` examples

#### 9.3.2 Update CHANGELOG.md
- [ ] Move all items from `[Unreleased]` to `[1.0.0]` section
- [ ] Add release date
- [ ] Organize by categories:
  - [ ] Added
  - [ ] Changed
  - [ ] Fixed
  - [ ] Security
- [ ] Review for completeness

#### 9.3.3 Create Release Checklist
- [ ] Create `docs/RELEASE_CHECKLIST.md`:
  ```markdown
  ## v1.0.0 Release Checklist

  ### Pre-release
  - [ ] All tests passing on main
  - [ ] Coverage >= 80%
  - [ ] All documentation updated
  - [ ] CHANGELOG.md complete
  - [ ] Version bumped
  - [ ] No known critical bugs

  ### Testing
  - [ ] Manual test on macOS
  - [ ] Manual test on Linux
  - [ ] Manual test on Windows
  - [ ] Test Homebrew installation (local)
  - [ ] Test binary downloads work

  ### Security
  - [ ] Run: go list -m -json all | nancy sleuth
  - [ ] Review dependencies for vulnerabilities
  - [ ] Verify credentials not logged
  - [ ] Verify config file permissions

  ### Performance
  - [ ] Run benchmarks
  - [ ] Verify no regressions
  - [ ] Memory usage < 50MB

  ### Release
  - [ ] Create tag: git tag v1.0.0
  - [ ] Push tag: git push origin v1.0.0
  - [ ] Verify release workflow succeeds
  - [ ] Verify binaries uploaded
  - [ ] Verify Homebrew formula updated
  - [ ] Test Homebrew install

  ### Post-release
  - [ ] Announce on social media
  - [ ] Submit to ClawdHub
  - [ ] Create discussion thread
  - [ ] Monitor for issues
  ```

#### 9.3.4 Security Audit
- [ ] Install nancy: `go install github.com/sonatype-nexus-community/nancy@latest`
- [ ] Run: `go list -m -json all | nancy sleuth`
- [ ] Review vulnerabilities
- [ ] Update dependencies if needed
- [ ] Run: `go mod tidy`

#### 9.3.5 Performance Validation
- [ ] Run: `go test -bench=. -benchmem ./...`
- [ ] Compare with baseline
- [ ] Verify no regressions
- [ ] Document results

#### 9.3.6 Manual Testing
- [ ] Build binary: `go build -o amazon-cli .`
- [ ] Test on macOS:
  - [ ] Run auth flow
  - [ ] List orders
  - [ ] Search products
  - [ ] Add to cart
  - [ ] Preview checkout (no --confirm)
- [ ] Test on Linux (Docker or VM)
- [ ] Test on Windows (VM)

#### 9.3.7 Create Release
- [ ] Commit all changes
- [ ] Push to main
- [ ] Create tag: `git tag -a v1.0.0 -m "Release v1.0.0"`
- [ ] Push tag: `git push origin v1.0.0`
- [ ] Monitor GitHub Actions
- [ ] Verify release created
- [ ] Download and test binaries

#### 9.3.8 Test Homebrew Installation
- [ ] Wait for Homebrew formula update
- [ ] Run: `brew tap zkwentz/tap`
- [ ] Run: `brew install amazon-cli`
- [ ] Verify installation: `amazon-cli --version`
- [ ] Test basic commands

#### 9.3.9 Submit to ClawdHub
- [ ] Follow ClawdHub submission process
- [ ] Upload skills.md
- [ ] Verify listing looks correct
- [ ] Test installation from ClawdHub

---

## Phase 10: Monitoring & Maintenance

### 10.1 Post-Release Monitoring

#### 10.1.1 Setup Issue Templates
- [ ] Create `.github/ISSUE_TEMPLATE/bug_report.md`:
  - [ ] Description field
  - [ ] Steps to reproduce
  - [ ] Expected behavior
  - [ ] Actual behavior
  - [ ] Environment (OS, version)
  - [ ] Logs
- [ ] Create `.github/ISSUE_TEMPLATE/feature_request.md`
- [ ] Create `.github/ISSUE_TEMPLATE/question.md`

#### 10.1.2 Monitor Issues
- [ ] Enable GitHub notifications
- [ ] Respond to issues within 48 hours
- [ ] Triage and label issues
- [ ] Create milestones for fixes

#### 10.1.3 Track Metrics
- [ ] Monitor GitHub stars
- [ ] Track Homebrew installs (if metrics available)
- [ ] Monitor ClawdHub usage
- [ ] Track common error reports

---

### 10.2 Future Enhancements (v1.1.0+)

#### 10.2.1 Additional Features Backlog
- [ ] Table output format (instead of just JSON)
- [ ] Bash/Zsh completion scripts
- [ ] Fish shell completion
- [ ] Configuration profiles (multiple Amazon accounts)
- [ ] Export orders to CSV
- [ ] Price tracking alerts
- [ ] Deal notifications
- [ ] International marketplace support (amazon.co.uk, etc.)

#### 10.2.2 Performance Optimizations
- [ ] Concurrent order fetching
- [ ] HTTP/2 support
- [ ] Connection pooling
- [ ] Response caching
- [ ] Incremental updates

#### 10.2.3 Security Enhancements
- [ ] Encrypted credential storage
- [ ] Keychain integration (macOS)
- [ ] Credential manager integration (Windows)
- [ ] Secret service integration (Linux)
- [ ] OAuth token rotation

---

## Appendix: Quick Reference

### Essential Commands

#### Run Tests
```bash
# All tests
go test ./...

# With coverage
go test -cover ./...

# Specific package
go test ./internal/amazon

# Verbose
go test -v ./...

# Short (skip slow tests)
go test -short ./...

# E2E tests
go test ./test/e2e/...

# Benchmarks
go test -bench=. -benchmem ./...

# Coverage HTML
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

#### Build Commands
```bash
# Build for current platform
go build -o amazon-cli .

# Build with version info
go build -ldflags "-X main.version=1.0.0" -o amazon-cli .

# Build for all platforms (using GoReleaser)
goreleaser build --snapshot --clean

# Cross-compile manually
GOOS=linux GOARCH=amd64 go build -o amazon-cli-linux .
GOOS=darwin GOARCH=arm64 go build -o amazon-cli-darwin-arm64 .
GOOS=windows GOARCH=amd64 go build -o amazon-cli.exe .
```

#### Linting
```bash
# Run linter
golangci-lint run

# Fix auto-fixable issues
golangci-lint run --fix

# Specific linters
golangci-lint run --enable-all
```

#### Coverage Analysis
```bash
# Generate coverage
go test -coverprofile=coverage.out ./...

# View in terminal
go tool cover -func=coverage.out

# View in browser
go tool cover -html=coverage.out

# Coverage by package
go tool cover -func=coverage.out | sort -k3 -n

# Find untested code
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## Task Completion Tracking

### Phase 1: Core Infrastructure
- Total tasks: 87
- Completed: 0
- In progress: 0
- Blocked: 0

### Phase 2: Orders & Returns
- Total tasks: 56
- Completed: 0
- In progress: 0
- Blocked: 0

### Phase 3: Search & Products
- Total tasks: 32
- Completed: 0
- In progress: 0
- Blocked: 0

### Phase 4: Cart & Checkout
- Total tasks: 41
- Completed: 0
- In progress: 0
- Blocked: 0

### Phase 5: Subscriptions
- Total tasks: 18
- Completed: 0
- In progress: 0
- Blocked: 0

### Phase 6: Error Handling
- Total tasks: 35
- Completed: 0
- In progress: 0
- Blocked: 0

### Phase 7: Integration Testing
- Total tasks: 28
- Completed: 0
- In progress: 0
- Blocked: 0

### Phase 8: CI/CD
- Total tasks: 47
- Completed: 0
- In progress: 0
- Blocked: 0

### Phase 9: Documentation
- Total tasks: 52
- Completed: 0
- In progress: 0
- Blocked: 0

### Phase 10: Monitoring
- Total tasks: 12
- Completed: 0
- In progress: 0
- Blocked: 0

**GRAND TOTAL: 408 tasks**

---

## Notes for Ralphy

This task list is designed for AI-powered execution. Each task is:

1. **Atomic**: Can be completed independently
2. **Specific**: Has clear acceptance criteria
3. **Actionable**: Includes exact file paths and code snippets where applicable
4. **Testable**: Has verification steps
5. **Ordered**: Respects dependencies (e.g., create file before editing it)

### Task Format
- `[ ]` = Not started
- File paths are absolute or relative from project root
- Code blocks show expected implementation
- Tests should be written alongside implementation

### Safety Notes
- **NEVER** test real checkout against production Amazon
- Keep mock implementations for checkout
- Respect rate limits during manual testing
- Anonymize any real data in fixtures

### Ralphy Execution Recommendations
1. Start with Phase 1 (core infrastructure)
2. Complete all tasks in a section before moving to next
3. Run tests after each implementation task
4. Verify coverage after each phase
5. Commit frequently with descriptive messages
6. Open PRs for major features for review

Good luck! 🚀
