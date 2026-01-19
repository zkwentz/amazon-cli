# amazon-cli

[![CI](https://img.shields.io/github/actions/workflow/status/zkwentz/amazon-cli/ci.yml?branch=main&label=CI&logo=github)](https://github.com/zkwentz/amazon-cli/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.25.5-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/zkwentz/amazon-cli/blob/main/LICENSE)

A command-line interface for Amazon shopping, designed for AI agent integration and programmatic access to Amazon.com.

## Overview

**amazon-cli** replaces the Amazon web interface with a terminal-based workflow, outputting structured JSON for seamless automation. Built with Go and the Cobra framework, it provides access to orders, returns, purchases, and subscriptions.

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap amazon-cli/tap
brew install amazon-cli
```

### Binary Releases

Download pre-compiled binaries from the [GitHub Releases](https://github.com/zkwentz/amazon-cli/releases) page:

- macOS (arm64, amd64)
- Linux (arm64, amd64)
- Windows (amd64)

### Build from Source

```bash
go install github.com/zkwentz/amazon-cli@latest
```

Or clone and build:

```bash
git clone https://github.com/zkwentz/amazon-cli.git
cd amazon-cli
go build -o amazon-cli .
```

## Quick Start

```bash
# Authenticate with Amazon
amazon-cli auth login

# Check authentication status
amazon-cli auth status

# List recent orders
amazon-cli orders list --limit 5

# Search for products
amazon-cli search "wireless headphones" --prime-only

# Add item to cart
amazon-cli cart add B08N5WRWNW --quantity 1

# View cart
amazon-cli cart list

# Checkout (preview first, then confirm)
amazon-cli cart checkout              # Preview
amazon-cli cart checkout --confirm    # Execute
```

## Authentication

amazon-cli uses browser-based OAuth for authentication.

### Login

```bash
amazon-cli auth login
```

This opens your default browser to Amazon's login page. After authenticating, tokens are stored locally in `~/.amazon-cli/config.json`.

### Check Status

```bash
amazon-cli auth status
```

Output:
```json
{
  "authenticated": true,
  "expires_at": "2024-01-20T12:00:00Z",
  "expires_in_seconds": 3600
}
```

### Logout

```bash
amazon-cli auth logout
```

## Commands

### Orders

```bash
# List recent orders
amazon-cli orders list [--limit N] [--status pending|delivered|returned]

# Get order details
amazon-cli orders get <order-id>

# Track shipment
amazon-cli orders track <order-id>

# Get order history for a specific year
amazon-cli orders history [--year YYYY]
```

**Example output (orders list):**
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

### Returns

```bash
# List returnable items
amazon-cli returns list

# Get return options for an item
amazon-cli returns options <order-id> <item-id>

# Initiate a return (requires --confirm)
amazon-cli returns create <order-id> <item-id> --reason <reason-code> --confirm

# Get return label
amazon-cli returns label <return-id>

# Check return status
amazon-cli returns status <return-id>
```

**Return reason codes:**
| Code | Description |
|------|-------------|
| `defective` | Item is defective or doesn't work |
| `wrong_item` | Received wrong item |
| `not_as_described` | Item not as described |
| `no_longer_needed` | No longer needed |
| `better_price` | Found better price elsewhere |
| `other` | Other reason |

### Search & Products

```bash
# Search products
amazon-cli search "<query>" [--category <cat>] [--min-price N] [--max-price N] [--prime-only]

# Get product details
amazon-cli product get <asin>

# Get product reviews
amazon-cli product reviews <asin> [--limit N]
```

**Example output (search):**
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

### Cart & Checkout

```bash
# Add to cart
amazon-cli cart add <asin> [--quantity N]

# View cart
amazon-cli cart list

# Remove from cart
amazon-cli cart remove <asin>

# Clear cart (requires --confirm)
amazon-cli cart clear --confirm

# Checkout (requires --confirm to execute)
amazon-cli cart checkout --confirm [--address-id <id>] [--payment-id <id>]

# Quick buy (requires --confirm)
amazon-cli buy <asin> --confirm [--quantity N] [--address-id <id>]
```

**Safety:** Purchase commands require the `--confirm` flag. Without it, the command shows a preview of what would happen.

**Example (cart list):**
```json
{
  "items": [
    {
      "asin": "B08N5WRWNW",
      "title": "Sony WH-1000XM4",
      "price": 278.00,
      "quantity": 1,
      "subtotal": 278.00,
      "prime": true,
      "in_stock": true
    }
  ],
  "subtotal": 278.00,
  "estimated_tax": 22.24,
  "total": 300.24,
  "item_count": 1
}
```

### Subscriptions (Subscribe & Save)

```bash
# List all subscriptions
amazon-cli subscriptions list

# Get subscription details
amazon-cli subscriptions get <subscription-id>

# Skip next delivery (requires --confirm)
amazon-cli subscriptions skip <subscription-id> --confirm

# Change frequency (requires --confirm)
amazon-cli subscriptions frequency <subscription-id> --interval <weeks> --confirm

# Cancel subscription (requires --confirm)
amazon-cli subscriptions cancel <subscription-id> --confirm

# View upcoming deliveries
amazon-cli subscriptions upcoming
```

**Example output (subscriptions list):**
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

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--output` | `-o` | Output format: json, table, raw | json |
| `--quiet` | `-q` | Suppress non-essential output | false |
| `--verbose` | `-v` | Enable verbose logging | false |
| `--config` | | Path to config file | ~/.amazon-cli/config.json |
| `--no-color` | | Disable colored output | false |

## Configuration

Configuration is stored in `~/.amazon-cli/config.json`:

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

## Error Handling

All errors return JSON with a consistent schema:

```json
{
  "error": {
    "code": "AUTH_EXPIRED",
    "message": "Authentication token has expired. Run 'amazon-cli auth login' to re-authenticate.",
    "details": {}
  }
}
```

### Error Codes

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

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Invalid arguments |
| 3 | Authentication error |
| 4 | Network error |
| 5 | Rate limited |
| 6 | Not found |

## Rate Limiting

To avoid triggering Amazon's anti-automation measures, amazon-cli implements:

- **Minimum delay:** 1 second between requests (configurable)
- **Jitter:** Random 0-500ms added to each delay
- **Exponential backoff:** On 429/503 responses, wait 2^n seconds (max 60s)
- **Max retries:** 3 attempts before failing (configurable)

## Project Structure

```
amazon-cli/
├── cmd/
│   └── root.go              # Root command and global flags
├── internal/
│   ├── amazon/              # Amazon API client
│   │   ├── cart.go
│   │   └── cart_test.go
│   ├── config/              # Configuration management
│   ├── output/              # Output formatting
│   └── ratelimit/           # Rate limiting logic
├── pkg/
│   └── models/              # Shared data models
│       └── cart.go
├── main.go
├── go.mod
├── go.sum
└── README.md
```

## Development

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o amazon-cli .
```

### Linting

```bash
golangci-lint run
```

## Security Considerations

1. **Credentials:** Stored in `~/.amazon-cli/config.json` with restricted permissions
2. **--confirm flag:** Required for all purchase/modification actions to prevent accidental execution
3. **No credential logging:** Tokens never appear in verbose output
4. **HTTPS only:** All Amazon communication over TLS

## License

MIT

## Contributing

Contributions are welcome. Please open an issue or submit a pull request.
