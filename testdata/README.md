# Amazon CLI Test Data

This directory contains mock Amazon API responses for testing purposes.

## Directory Structure

```
testdata/
├── mocks/
│   ├── auth/              # Authentication responses
│   ├── orders/            # Order management responses
│   ├── returns/           # Return management responses
│   ├── products/          # Product search and details responses
│   ├── cart/              # Shopping cart and checkout responses
│   ├── subscriptions/     # Subscribe & Save responses
│   └── errors/            # Error responses
└── helpers.go             # Test helper utilities
```

## Mock Data Categories

### Authentication (`auth/`)
- `login_success.json` - Successful authentication response
- `auth_status_authenticated.json` - Status when authenticated
- `auth_status_unauthenticated.json` - Status when not authenticated
- `logout_success.json` - Successful logout response

### Orders (`orders/`)
- `list_response.json` - List of recent orders
- `get_order_response.json` - Detailed single order
- `track_order_response.json` - Order tracking information
- `history_response.json` - Order history for a specific year

### Returns (`returns/`)
- `list_returnable_items.json` - List of items eligible for return
- `return_options.json` - Available return methods
- `create_return_response.json` - Return initiation confirmation
- `return_label.json` - Return shipping label
- `return_status.json` - Return status with timeline

### Products (`products/`)
- `search_response.json` - Product search results
- `product_details.json` - Detailed product information
- `product_reviews.json` - Product reviews with ratings

### Cart (`cart/`)
- `cart_response.json` - Current cart contents
- `add_to_cart_response.json` - Response after adding item
- `checkout_preview.json` - Checkout preview (dry run)
- `order_confirmation.json` - Order placement confirmation
- `addresses.json` - Saved shipping addresses
- `payment_methods.json` - Saved payment methods

### Subscriptions (`subscriptions/`)
- `list_subscriptions.json` - List of all subscriptions
- `get_subscription.json` - Detailed subscription information
- `upcoming_deliveries.json` - Upcoming subscription deliveries
- `skip_delivery_response.json` - Response after skipping delivery
- `update_frequency_response.json` - Response after changing frequency
- `cancel_subscription_response.json` - Subscription cancellation confirmation

### Errors (`errors/`)
- `auth_required.json` - Not authenticated error
- `auth_expired.json` - Expired authentication error
- `not_found.json` - Resource not found error
- `rate_limited.json` - Rate limiting error
- `invalid_input.json` - Invalid input validation error
- `purchase_failed.json` - Purchase failure error
- `network_error.json` - Network connectivity error
- `amazon_error.json` - Amazon service error

## Usage in Tests

### Go Example

```go
import (
    "encoding/json"
    "os"
    "testing"
)

func loadMockResponse(t *testing.T, filepath string) []byte {
    data, err := os.ReadFile(filepath)
    if err != nil {
        t.Fatalf("Failed to read mock file %s: %v", filepath, err)
    }
    return data
}

func TestOrdersList(t *testing.T) {
    mockData := loadMockResponse(t, "testdata/mocks/orders/list_response.json")

    var response OrdersResponse
    err := json.Unmarshal(mockData, &response)
    if err != nil {
        t.Fatalf("Failed to unmarshal response: %v", err)
    }

    if len(response.Orders) != 3 {
        t.Errorf("Expected 3 orders, got %d", len(response.Orders))
    }
}
```

### Mock HTTP Server Example

```go
func TestOrdersAPI(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        mockData := loadMockResponse(t, "testdata/mocks/orders/list_response.json")
        w.Header().Set("Content-Type", "application/json")
        w.Write(mockData)
    }))
    defer server.Close()

    // Test your API client against the mock server
    client := NewClient(server.URL)
    orders, err := client.GetOrders(10, "")
    // ... assertions
}
```

## Updating Mock Data

When updating mock responses:

1. Ensure all JSON is properly formatted
2. Match the schema defined in the PRD
3. Use realistic data values
4. Include edge cases (empty lists, null fields, etc.)
5. Keep consistency across related responses (same order IDs, ASINs, etc.)

## Common Test Scenarios

### Happy Path
- Use the standard response files for successful operations

### Error Cases
- Use files from `errors/` directory for testing error handling

### Edge Cases
- Empty results: Modify responses to have empty arrays
- Missing fields: Remove optional fields from responses
- Expired sessions: Use `auth_expired.json`

## Adding New Mock Data

When adding new mock responses:

1. Create a descriptive filename (e.g., `feature_action_response.json`)
2. Follow the existing schema from the PRD
3. Include realistic data
4. Document it in this README
5. Add usage examples if needed
