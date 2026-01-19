# Amazon CLI Tasks

Each task is independent and self-contained. Execute in order for best results.

---

- [ ] Create file `docs/amazon-api-research.md` that documents Amazon's authentication options including Login with Amazon (LWA) OAuth, session cookies, and which approach works best for accessing order history, cart, and search. Include example curl commands.

- [ ] Create file `docs/rate-limiting-strategy.md` documenting Amazon's rate limiting behavior: requests per minute before blocking, CAPTCHA trigger conditions, and recommended delays between requests (minimum 1-2 seconds with jitter).

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/internal/config/config.go`, create a complete configuration management system with `AuthConfig` struct (AccessToken, RefreshToken, ExpiresAt), `LoadConfig(path string) (*Config, error)` that reads from `~/.amazon-cli/config.json`, and `SaveConfig(config *Config, path string) error` that writes with 0600 permissions.

- [ ] Create file `internal/config/config_test.go` with unit tests for LoadConfig and SaveConfig including: test loading valid config, test loading missing file returns empty config, test saving creates directory if needed, test file permissions are 0600.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/internal/ratelimit/limiter.go`, create a complete rate limiter with: `RateLimiter` struct, `NewRateLimiter(minDelay, maxDelay time.Duration, maxRetries int)`, `Wait()` that enforces minimum delay with random jitter 0-500ms, `WaitWithBackoff(attempt int)` with exponential backoff capped at 60 seconds, `ShouldRetry(statusCode, attempt int) bool` that returns true for 429/503 within retry limit.

- [ ] Create file `internal/ratelimit/limiter_test.go` with tests: TestWait_EnforcesMinimumDelay, TestWait_AddsJitter, TestWaitWithBackoff_ExponentialIncrease, TestShouldRetry_Returns_True_For_429, TestShouldRetry_Returns_False_After_MaxRetries. Achieve 90%+ coverage.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/internal/amazon/client.go`, add a slice of 10 common browser User-Agent strings (Chrome, Firefox, Safari on Windows/macOS/mobile) and implement `getRandomUserAgent() string` that returns a random one.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/internal/amazon/client.go`, implement `Do(req *http.Request) (*http.Response, error)` method that: calls rate limiter Wait(), sets random User-Agent, sets Accept and Accept-Language headers, executes request, retries with WaitWithBackoff on 429/503 up to maxRetries times.

- [x] In file `/Users/zacharywentz/Development/amazon-cli/internal/amazon/client.go`, add circuit breaker with: failureCount, threshold (5), resetTimeout (60s). Add `checkCircuitBreaker() error` that errors if too many failures, `recordFailure()` and `recordSuccess()` methods. Integrate into Do() method.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/internal/amazon/client.go`, implement `detectCAPTCHA(body []byte) bool` that checks response body for CAPTCHA indicators like "captcha" or Amazon's specific CAPTCHA HTML patterns.

- [x] Create file `internal/testutil/mock_server.go` with `MockAmazonServer` struct and `NewMockAmazonServer() *MockAmazonServer` that creates an httptest.Server. Add `ServeFixture(path, fixtureFile string)` to configure which fixture file to serve for each URL path.

- [x] Create file `internal/amazon/client_test.go` with integration tests using MockAmazonServer: TestDo_Success, TestDo_Retry_On_429, TestDo_Retry_On_503, TestDo_Stops_After_MaxRetries, TestDo_CircuitBreaker_Opens_After_Failures, TestDo_Detects_CAPTCHA. Achieve 90%+ coverage.

- [x] Create file `internal/amazon/auth.go` with: `AuthTokens` struct (AccessToken, RefreshToken, ExpiresAt), `IsExpired() bool` method, `ExpiresWithin(duration time.Duration) bool` method, and `RefreshTokens(refreshToken string) (*AuthTokens, error)` placeholder that returns mock tokens for now.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/cmd/auth.go`, implement `authLoginCmd` Run function that: outputs `{"status": "login_required", "message": "Browser-based login not yet implemented"}` as JSON. This is a placeholder for future OAuth implementation.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/cmd/auth.go`, implement `authStatusCmd` Run function that: loads config, checks if tokens exist and not expired, outputs JSON `{"authenticated": true/false, "expires_at": "...", "expires_in_seconds": N}`.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/cmd/auth.go`, implement `authLogoutCmd` Run function that: loads config, clears auth tokens, saves config, outputs JSON `{"status": "logged_out"}`.

- [x] Create file `internal/amazon/auth_test.go` with tests: TestAuthTokens_IsExpired_True_When_Past, TestAuthTokens_IsExpired_False_When_Future, TestAuthTokens_ExpiresWithin_True, TestAuthTokens_ExpiresWithin_False. Achieve 85%+ coverage.

- [x] Create directory `testdata/orders/` and create file `testdata/orders/order_list_sample.html` containing sample HTML structure that mimics Amazon's order history page with 3 orders, each having order_id, date, total, status, and item details. Anonymize all data.

- [x] Create file `testdata/orders/order_detail_sample.html` containing sample HTML for a single Amazon order detail page with: order header info, 2 items with ASIN/title/price/quantity, shipping address, payment method, and tracking information.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/internal/amazon/orders.go`, add `import "github.com/PuerkitoBio/goquery"` and implement `parseOrdersHTML(html []byte) ([]models.Order, error)` that parses order list HTML and extracts order_id, date, total, status for each order. Use CSS selectors appropriate for the fixture structure.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/internal/amazon/orders.go`, implement `parseOrderDetailHTML(html []byte) (*models.Order, error)` that parses single order HTML and extracts complete order with items array (ASIN, title, price, quantity) and tracking info if present.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/internal/amazon/orders.go`, replace the mock `GetOrders(limit int, status string)` implementation with: make HTTP GET to order history URL, parse response with parseOrdersHTML, filter by status if provided, limit results, return OrdersResponse.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/internal/amazon/orders.go`, replace the mock `GetOrder(orderID string)` implementation with: validate orderID format, make HTTP GET to order detail URL, parse response with parseOrderDetailHTML, return Order.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/internal/amazon/orders.go`, replace the mock `GetOrderTracking(orderID string)` implementation with: make HTTP GET to tracking URL, parse tracking info (carrier, tracking_number, status, delivery_date), return Tracking struct.

- [x] Create file `internal/amazon/orders_test.go` with parser tests: TestParseOrdersHTML_ReturnsCorrectCount (load fixture, verify 3 orders), TestParseOrdersHTML_ExtractsAllFields (verify order_id, date, total, status not empty), TestParseOrderDetailHTML_ExtractsItems (verify items array populated).

- [x] Create file `internal/amazon/orders_test.go` integration tests using MockAmazonServer: TestGetOrders_Integration (mock serves fixture, verify response), TestGetOrder_Integration, TestGetOrders_EmptyHistory (fixture with no orders returns empty array).

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/cmd/orders.go`, update `ordersListCmd` Run function to: create Amazon client, call GetOrders with limit and status flags, handle errors with proper error codes, output JSON result.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/cmd/orders.go`, update `ordersGetCmd` Run function to: validate orderID argument provided, create client, call GetOrder, handle NOT_FOUND error, output JSON result.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/cmd/orders.go`, update `ordersTrackCmd` Run function to: validate orderID, call GetOrderTracking, output tracking JSON.

- [x] Create file `pkg/models/return.go` with structs: `ReturnableItem` (OrderID, ItemID, ASIN, Title, Price, PurchaseDate, ReturnWindow), `ReturnOption` (Method, Label, DropoffLocation, Fee), `Return` (ReturnID, OrderID, ItemID, Status, Reason, CreatedAt), `ReturnLabel` (URL, Carrier, Instructions).

- [x] Create file `testdata/returns/returnable_items_sample.html` with sample HTML mimicking Amazon's returns center showing 2 returnable items with order_id, item_id, title, price, purchase_date, return_window.

- [x] Create file `internal/amazon/returns.go` with: `GetReturnableItems() ([]models.ReturnableItem, error)` that makes GET request and parses HTML, `GetReturnOptions(orderID, itemID string) ([]models.ReturnOption, error)` placeholder returning mock options.

- [x] Create file `internal/amazon/returns.go` with: `CreateReturn(orderID, itemID, reason string) (*models.Return, error)` that validates reason against allowed list (defective, wrong_item, not_as_described, no_longer_needed, better_price, other), returns mock Return with generated returnID.

- [x] Create file `internal/amazon/returns.go` with: `GetReturnLabel(returnID string) (*models.ReturnLabel, error)` returning mock label data, `GetReturnStatus(returnID string) (*models.Return, error)` returning mock status.

- [x] Create file `cmd/returns.go` with `returns` parent command and subcommands: `returns list` (calls GetReturnableItems, outputs JSON), `returns options <order-id> <item-id>` (calls GetReturnOptions, outputs JSON).

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/cmd/returns.go`, implement `returns create <order-id> <item-id>` command with --reason flag (required) and --confirm flag. Without --confirm output dry_run preview. With --confirm call CreateReturn and output result.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/cmd/returns.go`, implement `returns label <return-id>` and `returns status <return-id>` commands that call respective client methods and output JSON.

- [x] Create file `internal/amazon/returns_test.go` with tests: TestCreateReturn_InvalidReason_ReturnsError, TestCreateReturn_ValidReason_Succeeds, TestGetReturnableItems_ParsesHTML (using fixture).

- [x] Create file `testdata/search/search_results_sample.html` with sample HTML mimicking Amazon search results page showing 5 products with ASIN, title, price, rating, review_count, prime badge, stock status.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/internal/amazon/search.go`, implement `parseSearchResultsHTML(html []byte) ([]models.Product, error)` that parses search HTML and extracts ASIN, title, price, rating, review_count, prime, in_stock for each product.

- [x] In file `/Users/zacharywentz/Development/amazon-cli/internal/amazon/search.go`, update `Search(query string, opts SearchOptions)` to: build search URL with query params for category/minPrice/maxPrice/primeOnly, make HTTP GET, parse with parseSearchResultsHTML, return SearchResponse.

- [x] Create file `internal/amazon/search_test.go` with tests: TestParseSearchResultsHTML_ExtractsProducts, TestSearch_WithPrimeFilter, TestSearch_WithPriceRange. Use fixtures and MockAmazonServer.

- [x] In file `/Users/zacharywentz/Development/amazon-cli/cmd/search.go`, update Run function to: get all flags (category, min-price, max-price, prime-only), create SearchOptions, call client.Search, output JSON result.

- [x] Create file `testdata/products/product_detail_sample.html` with sample HTML mimicking Amazon product page with: ASIN, title, price, original_price, rating, review_count, prime badge, stock status, description, feature bullets, image URLs.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/internal/amazon/product.go`, implement `parseProductDetailHTML(html []byte) (*models.Product, error)` that extracts all product fields from detail page HTML, handling missing optional fields gracefully.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/internal/amazon/product.go`, update `GetProduct(asin string)` to: validate ASIN format, make HTTP GET to product URL, parse with parseProductDetailHTML, return Product.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/internal/amazon/product.go`, update `GetProductReviews(asin string, limit int)` to: make HTTP GET to reviews URL, parse reviews (rating, title, body, author, date, verified), limit results, return ReviewsResponse.

- [ ] Create file `internal/amazon/product_test.go` with tests: TestParseProductDetailHTML_AllFields, TestParseProductDetailHTML_MissingOptionalFields, TestGetProduct_InvalidASIN_ReturnsError.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/cmd/product.go`, update `productGetCmd` Run function to: validate ASIN argument, call client.GetProduct, output JSON. Update `productReviewsCmd` to call GetProductReviews with --limit flag.

- [x] Create file `internal/validation/validators.go` with: `ValidateASIN(asin string) error` (must be 10 alphanumeric chars), `ValidateOrderID(id string) error` (format XXX-XXXXXXX-XXXXXXX), `ValidateQuantity(qty int) error` (1-999), `ValidatePriceRange(min, max float64) error` (min >= 0, max > min).

- [x] Create file `internal/validation/validators_test.go` with exhaustive tests: valid inputs pass, empty strings fail, wrong lengths fail, special characters fail, boundary values tested. Achieve 100% coverage on validators.go.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/internal/amazon/cart.go`, update `AddToCart(asin string, quantity int)` to: call ValidateASIN, call ValidateQuantity, then proceed with existing mock implementation that adds to in-memory cart.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/internal/amazon/cart.go`, update `RemoveFromCart(asin string)` to: call ValidateASIN, then actually remove the item from the in-memory cart.Items slice, recalculate totals.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/internal/amazon/cart.go`, update `ClearCart()` to: reset cart.Items to empty slice, set cart.Subtotal, cart.EstimatedTax, cart.Total, cart.ItemCount all to 0.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/internal/amazon/cart.go`, update `PreviewCheckout` to return more realistic preview data including: current cart contents, mock address with all fields populated, mock payment method with all fields, delivery options array.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/internal/amazon/cart.go`, add comments to `CompleteCheckout` function clearly stating: "MOCK IMPLEMENTATION - Never test against production Amazon. Real implementation requires sandbox environment." Keep mock behavior.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/internal/amazon/cart_test.go`, add test `TestCompleteCheckout_NeverMakesRealHTTPPost` that verifies the mock implementation doesn't make external HTTP calls by checking no network requests are made.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/internal/amazon/cart_test.go`, add test `TestRemoveFromCart_ActuallyRemovesItem` that: adds item, verifies count is 1, removes item, verifies count is 0.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/internal/amazon/cart_test.go`, add test `TestClearCart_ResetsAllTotals` that: adds multiple items, clears cart, verifies ItemCount=0, Subtotal=0, Total=0.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/cmd/cart.go`, update `cartCheckoutCmd` to verify --confirm flag is checked BEFORE any checkout logic runs, ensuring preview mode is the default.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/cmd/buy.go`, update `buyCmd` Run function to: validate ASIN, if no --confirm show product preview with price and quantity, if --confirm call AddToCart then CompleteCheckout, output result.

- [ ] Create file `pkg/models/subscription.go` with structs: `Subscription` (SubscriptionID, ASIN, Title, Price, DiscountPercent, FrequencyWeeks, NextDelivery, Status, Quantity), `SubscriptionsResponse` (Subscriptions []Subscription), `UpcomingDelivery` (SubscriptionID, ASIN, Title, DeliveryDate, Quantity).

- [ ] Create file `internal/amazon/subscriptions.go` with mock implementations: `GetSubscriptions() (*models.SubscriptionsResponse, error)` returning 2 mock subscriptions, `GetSubscription(id string) (*models.Subscription, error)` returning mock subscription.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/internal/amazon/subscriptions.go`, add: `SkipDelivery(id string) (*models.Subscription, error)` that returns subscription with NextDelivery advanced by FrequencyWeeks, `CancelSubscription(id string) (*models.Subscription, error)` that returns subscription with Status="cancelled".

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/internal/amazon/subscriptions.go`, add: `UpdateFrequency(id string, weeks int) (*models.Subscription, error)` that validates weeks is 1-26 and returns subscription with updated FrequencyWeeks, `GetUpcomingDeliveries() ([]models.UpcomingDelivery, error)` returning mock deliveries sorted by date.

- [ ] Create file `cmd/subscriptions.go` with `subscriptions` parent command and subcommands: `list` (calls GetSubscriptions), `get <id>` (calls GetSubscription), `upcoming` (calls GetUpcomingDeliveries). All output JSON.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/cmd/subscriptions.go`, add `skip <id>` command with --confirm flag. Without --confirm show what would be skipped. With --confirm call SkipDelivery.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/cmd/subscriptions.go`, add `frequency <id>` command with --interval flag (required, weeks) and --confirm flag. Validate interval 1-26. Without --confirm show preview. With --confirm call UpdateFrequency.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/cmd/subscriptions.go`, add `cancel <id>` command with --confirm flag. Without --confirm show cancellation preview. With --confirm call CancelSubscription.

- [ ] Create file `internal/amazon/subscriptions_test.go` with tests: TestGetSubscriptions_ReturnsList, TestSkipDelivery_AdvancesDate, TestUpdateFrequency_InvalidWeeks_ReturnsError, TestCancelSubscription_SetsStatusCancelled.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/pkg/models/errors.go`, ensure these error codes exist as constants: AUTH_REQUIRED, AUTH_EXPIRED, NOT_FOUND, RATE_LIMITED, INVALID_INPUT, PURCHASE_FAILED, NETWORK_ERROR, AMAZON_ERROR, CAPTCHA_REQUIRED. Add any missing ones.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/pkg/models/errors.go`, ensure `CLIError` struct has Code, Message, Details fields and implements `Error() string` method that returns JSON formatted error string.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/pkg/models/errors.go`, add `NewCLIError(code, message string, details map[string]interface{}) *CLIError` constructor function and exit code constants: ExitSuccess=0, ExitGeneralError=1, ExitInvalidArgs=2, ExitAuthError=3, ExitNetworkError=4, ExitRateLimited=5, ExitNotFound=6.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/internal/output/output.go`, ensure `Error(code, message string, details map[string]interface{})` function exists that outputs JSON formatted error to stderr: `{"error": {"code": "...", "message": "...", "details": {...}}}`.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/internal/amazon/client.go`, wrap all errors with context using fmt.Errorf with %w: network errors get "network request failed: %w", parse errors get "failed to parse response: %w", auth errors get "authentication failed: %w".

- [ ] Create file `internal/output/output_test.go` with tests: TestJSON_OutputsValidJSON, TestError_OutputsErrorSchema, TestError_IncludesAllFields. Verify JSON is parseable.

- [ ] Create file `.github/workflows/ci.yml` with: name "CI", triggers on push to main and pull_request to main, job "test" with matrix (ubuntu-latest, macos-latest, windows-latest) and go 1.21, steps to checkout, setup-go, run `go test -v -race -coverprofile=coverage.txt ./...`.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/.github/workflows/ci.yml`, add job "lint" that runs on ubuntu-latest, uses golangci/golangci-lint-action@v4 with version latest.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/.github/workflows/ci.yml`, add step to upload coverage to codecov using codecov/codecov-action@v4 with file ./coverage.txt.

- [ ] Create file `.github/workflows/release.yml` with: name "Release", trigger on push tags 'v*', job "release" on ubuntu-latest, steps to checkout with fetch-depth 0, setup-go, run goreleaser/goreleaser-action@v5 with args "release --clean" and env GITHUB_TOKEN.

- [ ] Create file `.github/PULL_REQUEST_TEMPLATE.md` with sections: Description, Changes (bullet list), Testing checklist (unit tests, integration tests, manual testing), Checklist (tests pass, coverage maintained, docs updated).

- [ ] Create file `.github/ISSUE_TEMPLATE/bug_report.md` with fields: Description, Steps to Reproduce, Expected Behavior, Actual Behavior, Environment (OS, CLI version), Logs/Output.

- [ ] Create file `.github/ISSUE_TEMPLATE/feature_request.md` with fields: Feature Description, Use Case, Proposed Solution, Alternatives Considered.

- [ ] Create file `CONTRIBUTING.md` with sections: Development Setup (go install, clone repo), Running Tests (go test ./...), Building (go build), Code Style (gofmt, golangci-lint), Submitting PRs (branch naming, commit messages, PR template).

- [ ] Create file `SECURITY.md` with sections: Supported Versions (table), Reporting Vulnerabilities (email/process), Security Best Practices (config file permissions 0600, credential rotation, --confirm flag usage).

- [ ] Create file `CHANGELOG.md` with header following Keep a Changelog format, [Unreleased] section, and [1.0.0] section listing all features: Authentication, Orders, Returns, Search, Products, Cart, Checkout, Subscriptions, Rate Limiting.

- [ ] Create file `docs/TROUBLESHOOTING.md` with common issues and solutions: "Authentication Failed" (re-run auth login), "Rate Limited" (wait and retry), "CAPTCHA Required" (complete in browser), "Command Not Found" (check PATH), "Permission Denied on Config" (check file permissions).

- [ ] Update file `/Users/zacharywentz/Development/amazon-cli/README.md` to add badges at top: CI status badge from GitHub Actions, Go version badge, License badge. Use shields.io format.

- [ ] Update file `/Users/zacharywentz/Development/amazon-cli/README.md` to ensure all command examples in Quick Start section actually work with current implementation. Test each command.

- [ ] Create file `skills.md` with YAML frontmatter: name: amazon-cli, description: CLI for Amazon shopping automation, version: 1.0.0, author: zkwentz, repository URL, tags: [shopping, e-commerce, amazon, automation].

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/skills.md`, add Installation section with Homebrew instructions (brew tap, brew install) and binary download instructions.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/skills.md`, add Authentication section documenting: auth login (inputs: none, output: status JSON), auth status (inputs: none, output: auth state JSON), auth logout (inputs: none, output: confirmation JSON).

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/skills.md`, add Orders section documenting each command with inputs, outputs, and example: orders list (--limit, --status), orders get (order-id), orders track (order-id), orders history (--year).

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/skills.md`, add Search section documenting: search command with all flags (query, --category, --min-price, --max-price, --prime-only), output schema showing results array structure.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/skills.md`, add Products section documenting: product get (asin input, full product output), product reviews (asin, --limit inputs, reviews array output).

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/skills.md`, add Cart section documenting: cart add, cart list, cart remove, cart clear (--confirm required), cart checkout (--confirm required, preview without flag).

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/skills.md`, add Safety section with bold warning: --confirm flag required for all purchase operations, always preview before confirming, cart checkout without --confirm shows preview only.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/skills.md`, add Returns section documenting: returns list, returns options, returns create (--reason codes list, --confirm required), returns label, returns status.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/skills.md`, add Subscriptions section documenting: subscriptions list, get, skip (--confirm), frequency (--interval, --confirm), cancel (--confirm), upcoming.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/skills.md`, add Error Codes section with table: AUTH_REQUIRED, AUTH_EXPIRED, NOT_FOUND, RATE_LIMITED, INVALID_INPUT, PURCHASE_FAILED, NETWORK_ERROR, AMAZON_ERROR and their descriptions.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/skills.md`, add Examples section with 5 common AI agent tasks: "Check recent orders" (orders list --limit 5), "Track a package" (orders track <id>), "Search products" (search "query" --prime-only), "Add to cart" (cart add <asin>), "Preview checkout" (cart checkout without --confirm).

- [ ] Create file `test/e2e/cli_test.go` with test `TestCLI_Help_Works` that builds the binary with `go build`, runs `./amazon-cli --help`, and verifies exit code is 0 and output contains "amazon-cli".

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/test/e2e/cli_test.go`, add test `TestCLI_OrdersList_OutputsJSON` that runs `./amazon-cli orders list --limit 1` and verifies output is valid JSON with "orders" key.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/test/e2e/cli_test.go`, add test `TestCLI_CartCheckout_RequiresConfirm` that runs `./amazon-cli cart checkout` without --confirm and verifies output contains "dry_run": true.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/test/e2e/cli_test.go`, add test `TestCLI_Search_ReturnsResults` that runs `./amazon-cli search "test"` and verifies output is valid JSON with "query" and "results" keys.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/test/e2e/cli_test.go`, add test `TestCLI_InvalidASIN_ReturnsError` that runs `./amazon-cli product get "invalid"` and verifies exit code is non-zero and output contains error.

- [ ] Run command `go test -v -coverprofile=coverage.out ./...` in project root and verify all tests pass. Run `go tool cover -func=coverage.out | grep total` and document the coverage percentage.

- [ ] Run command `golangci-lint run` in project root and fix any linter errors or warnings. Ensure zero issues reported.

- [ ] Run command `go build -o amazon-cli .` in project root and verify binary is created. Test `./amazon-cli --help` works.

- [ ] Run command `go build -ldflags "-X main.version=1.0.0" -o amazon-cli .` to verify version embedding works. If main.version variable doesn't exist, add it to main.go.

- [ ] In file `/Users/zacharywentz/Development/amazon-cli/main.go`, add `var version = "dev"` at package level and update to print version when `--version` flag is passed or add version command if not exists.

- [ ] Create file `Makefile` with targets: `build` (go build -o amazon-cli .), `test` (go test -v ./...), `cover` (go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out), `lint` (golangci-lint run), `clean` (rm -f amazon-cli coverage.out).

- [ ] Run command `make test` to verify Makefile works and all tests pass.

- [ ] Run command `make lint` to verify linting passes with zero issues.

- [ ] Run command `make build` to verify build succeeds and produces amazon-cli binary.

- [ ] Delete file `/Users/zacharywentz/Development/amazon-cli/IMPLEMENTATION_AND_TEST_PLAN.md` as it has been replaced by this tasks.md file.

- [ ] Run command `go mod tidy` to clean up go.mod and go.sum files, removing any unused dependencies.

- [ ] Run command `gofmt -w .` to format all Go files in the project.

- [ ] Commit all changes with message "Complete amazon-cli implementation with comprehensive tests and documentation".
