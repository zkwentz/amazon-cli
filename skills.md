---
name: amazon-cli
description: CLI tool for managing Amazon orders, returns, purchases, and subscriptions
version: 1.0.0
author: zkwentz
repository: https://github.com/zkwentz/amazon-cli
---

# amazon-cli

A command-line interface for Amazon shopping that outputs structured JSON, designed for AI agent integration.

## Installation

```bash
# Homebrew (macOS/Linux)
brew tap zkwentz/tap
brew install amazon-cli

# Or download binary from GitHub Releases
# Or build from source
go install github.com/zkwentz/amazon-cli@latest
```

## Authentication

Before using any commands, authenticate with Amazon:

```bash
amazon-cli auth login
```

This opens your browser for Amazon OAuth authentication. Tokens are stored locally.

---

## Actions

### auth login

Authenticate with Amazon via browser OAuth.

**Inputs:** None

**Output:**
```json
{
  "status": "authenticated",
  "expires_at": "2024-01-20T12:00:00Z"
}
```

---

### auth status

Check current authentication status.

**Inputs:** None

**Output:**
```json
{
  "authenticated": true,
  "expires_at": "2024-01-20T12:00:00Z",
  "expires_in_seconds": 3600
}
```

---

### auth logout

Clear stored credentials.

**Inputs:** None

**Output:**
```json
{
  "status": "logged_out"
}
```

---

### orders list

List recent orders.

**Inputs:**
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| --limit | int | No | 10 | Number of orders to return |
| --status | string | No | all | Filter: pending, delivered, returned |

**Command:**
```bash
amazon-cli orders list --limit 5 --status delivered
```

**Output:**
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

---

### orders get

Get details for a specific order.

**Inputs:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| order_id | string | Yes | Amazon order ID |

**Command:**
```bash
amazon-cli orders get 123-4567890-1234567
```

**Output:** Same as single order in `orders list`

---

### orders track

Get tracking information for an order.

**Inputs:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| order_id | string | Yes | Amazon order ID |

**Command:**
```bash
amazon-cli orders track 123-4567890-1234567
```

**Output:**
```json
{
  "carrier": "UPS",
  "tracking_number": "1Z999AA10123456784",
  "status": "in_transit",
  "delivery_date": "2024-01-17",
  "events": [
    {
      "timestamp": "2024-01-16T10:30:00Z",
      "location": "Local Facility",
      "status": "Out for delivery"
    }
  ]
}
```

---

### orders history

Get order history for a specific year.

**Inputs:**
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| --year | int | No | current | Year to fetch orders from |

**Command:**
```bash
amazon-cli orders history --year 2024
```

**Output:** Same schema as `orders list`

---

### returns list

List items eligible for return.

**Inputs:** None

**Command:**
```bash
amazon-cli returns list
```

**Output:**
```json
{
  "returnable_items": [
    {
      "order_id": "123-4567890-1234567",
      "item_id": "ITEM123",
      "asin": "B08N5WRWNW",
      "title": "Product Name",
      "price": 29.99,
      "purchase_date": "2024-01-15",
      "return_window": "2024-02-15"
    }
  ]
}
```

---

### returns options

Get available return options for an item.

**Inputs:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| order_id | string | Yes | Amazon order ID |
| item_id | string | Yes | Item ID within order |

**Command:**
```bash
amazon-cli returns options 123-4567890-1234567 ITEM123
```

**Output:**
```json
{
  "return_options": [
    {
      "method": "UPS Dropoff",
      "label": "Print label",
      "dropoff_location": "Any UPS Store",
      "fee": 0
    },
    {
      "method": "Amazon Locker",
      "label": "QR code",
      "dropoff_location": "Whole Foods - Downtown",
      "fee": 0
    }
  ]
}
```

---

### returns create

Initiate a return for an item. **Requires --confirm flag.**

**Inputs:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| order_id | string | Yes | Amazon order ID |
| item_id | string | Yes | Item ID within order |
| --reason | string | Yes | Reason code (see below) |
| --confirm | flag | Yes | Required to execute |

**Reason codes:** `defective`, `wrong_item`, `not_as_described`, `no_longer_needed`, `better_price`, `other`

**Command:**
```bash
# Preview (without --confirm)
amazon-cli returns create 123-4567890-1234567 ITEM123 --reason defective

# Execute (with --confirm)
amazon-cli returns create 123-4567890-1234567 ITEM123 --reason defective --confirm
```

**Output (preview):**
```json
{
  "dry_run": true,
  "would_return": {
    "order_id": "123-4567890-1234567",
    "item_id": "ITEM123",
    "title": "Product Name",
    "reason": "defective"
  },
  "message": "Add --confirm to execute"
}
```

**Output (confirmed):**
```json
{
  "return_id": "R01-1234567-8901234",
  "order_id": "123-4567890-1234567",
  "item_id": "ITEM123",
  "status": "initiated",
  "reason": "defective",
  "created_at": "2024-01-18T12:00:00Z"
}
```

---

### returns label

Get return shipping label.

**Inputs:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| return_id | string | Yes | Return ID |

**Command:**
```bash
amazon-cli returns label R01-1234567-8901234
```

**Output:**
```json
{
  "label_url": "https://...",
  "carrier": "UPS",
  "instructions": "Drop off at any UPS location"
}
```

---

### returns status

Check status of a return.

**Inputs:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| return_id | string | Yes | Return ID |

**Command:**
```bash
amazon-cli returns status R01-1234567-8901234
```

**Output:**
```json
{
  "return_id": "R01-1234567-8901234",
  "status": "refunded",
  "refund_amount": 29.99,
  "refund_date": "2024-01-25"
}
```

---

### search

Search for products on Amazon.

**Inputs:**
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| query | string | Yes | - | Search query |
| --category | string | No | all | Product category |
| --min-price | float | No | - | Minimum price filter |
| --max-price | float | No | - | Maximum price filter |
| --prime-only | flag | No | false | Only Prime items |
| --page | int | No | 1 | Results page |

**Command:**
```bash
amazon-cli search "wireless headphones" --max-price 100 --prime-only
```

**Output:**
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

---

### product get

Get detailed product information.

**Inputs:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| asin | string | Yes | Amazon product ASIN |

**Command:**
```bash
amazon-cli product get B08N5WRWNW
```

**Output:**
```json
{
  "asin": "B08N5WRWNW",
  "title": "Sony WH-1000XM4 Wireless Headphones",
  "price": 278.00,
  "original_price": 349.99,
  "rating": 4.7,
  "review_count": 52431,
  "prime": true,
  "in_stock": true,
  "delivery_estimate": "Tomorrow",
  "description": "Industry-leading noise canceling...",
  "features": [
    "30-hour battery life",
    "Touch sensor controls",
    "Speak-to-chat technology"
  ],
  "images": [
    "https://images-na.ssl-images-amazon.com/..."
  ]
}
```

---

### product reviews

Get product reviews.

**Inputs:**
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| asin | string | Yes | - | Amazon product ASIN |
| --limit | int | No | 10 | Number of reviews |

**Command:**
```bash
amazon-cli product reviews B08N5WRWNW --limit 5
```

**Output:**
```json
{
  "asin": "B08N5WRWNW",
  "average_rating": 4.7,
  "total_reviews": 52431,
  "reviews": [
    {
      "rating": 5,
      "title": "Best headphones I've owned",
      "body": "The noise canceling is incredible...",
      "author": "John D.",
      "date": "2024-01-10",
      "verified": true
    }
  ]
}
```

---

### cart add

Add item to shopping cart.

**Inputs:**
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| asin | string | Yes | - | Amazon product ASIN |
| --quantity | int | No | 1 | Quantity to add |

**Command:**
```bash
amazon-cli cart add B08N5WRWNW --quantity 2
```

**Output:**
```json
{
  "items": [
    {
      "asin": "B08N5WRWNW",
      "title": "Sony WH-1000XM4",
      "price": 278.00,
      "quantity": 2,
      "subtotal": 556.00,
      "prime": true,
      "in_stock": true
    }
  ],
  "subtotal": 556.00,
  "estimated_tax": 44.48,
  "total": 600.48,
  "item_count": 2
}
```

---

### cart list

View current cart contents.

**Inputs:** None

**Command:**
```bash
amazon-cli cart list
```

**Output:** Same schema as `cart add`

---

### cart remove

Remove item from cart.

**Inputs:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| asin | string | Yes | Amazon product ASIN |

**Command:**
```bash
amazon-cli cart remove B08N5WRWNW
```

**Output:** Updated cart (same schema as `cart add`)

---

### cart clear

Clear all items from cart. **Requires --confirm flag.**

**Inputs:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| --confirm | flag | Yes | Required to execute |

**Command:**
```bash
amazon-cli cart clear --confirm
```

**Output:**
```json
{
  "status": "cleared",
  "items_removed": 3
}
```

---

### cart checkout

Checkout current cart. **Requires --confirm flag to execute purchase.**

**Inputs:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| --confirm | flag | Yes | Required to complete purchase |
| --address-id | string | No | Shipping address ID |
| --payment-id | string | No | Payment method ID |

**Command:**
```bash
# Preview checkout
amazon-cli cart checkout

# Complete purchase
amazon-cli cart checkout --confirm
```

**Output (preview):**
```json
{
  "dry_run": true,
  "cart": {
    "items": [...],
    "subtotal": 278.00,
    "estimated_tax": 22.24,
    "total": 300.24
  },
  "address": {
    "id": "addr_default",
    "name": "John Doe",
    "street": "123 Main St",
    "city": "Seattle",
    "state": "WA",
    "zip": "98101"
  },
  "payment_method": {
    "type": "Visa",
    "last4": "1234"
  },
  "message": "Add --confirm to complete purchase"
}
```

**Output (confirmed):**
```json
{
  "order_id": "123-4567890-1234567",
  "total": 300.24,
  "estimated_delivery": "2024-01-20"
}
```

---

### buy

Quick purchase a single item. **Requires --confirm flag.**

**Inputs:**
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| asin | string | Yes | - | Amazon product ASIN |
| --quantity | int | No | 1 | Quantity to purchase |
| --confirm | flag | Yes | - | Required to execute |
| --address-id | string | No | default | Shipping address ID |
| --payment-id | string | No | default | Payment method ID |

**Command:**
```bash
# Preview
amazon-cli buy B08N5WRWNW --quantity 1

# Purchase
amazon-cli buy B08N5WRWNW --quantity 1 --confirm
```

**Output (preview):**
```json
{
  "dry_run": true,
  "product": {
    "asin": "B08N5WRWNW",
    "title": "Sony WH-1000XM4",
    "price": 278.00
  },
  "quantity": 1,
  "total": 300.24,
  "message": "Add --confirm to complete purchase"
}
```

**Output (confirmed):**
```json
{
  "order_id": "123-4567890-1234567",
  "total": 300.24,
  "estimated_delivery": "2024-01-20"
}
```

---

### subscriptions list

List all Subscribe & Save subscriptions.

**Inputs:** None

**Command:**
```bash
amazon-cli subscriptions list
```

**Output:**
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

---

### subscriptions get

Get details for a specific subscription.

**Inputs:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| subscription_id | string | Yes | Subscription ID |

**Command:**
```bash
amazon-cli subscriptions get S01-1234567-8901234
```

**Output:** Same as single subscription in `subscriptions list`

---

### subscriptions skip

Skip next delivery. **Requires --confirm flag.**

**Inputs:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| subscription_id | string | Yes | Subscription ID |
| --confirm | flag | Yes | Required to execute |

**Command:**
```bash
amazon-cli subscriptions skip S01-1234567-8901234 --confirm
```

**Output:**
```json
{
  "subscription_id": "S01-1234567-8901234",
  "status": "active",
  "skipped_delivery": "2024-02-01",
  "next_delivery": "2024-03-01"
}
```

---

### subscriptions frequency

Change delivery frequency. **Requires --confirm flag.**

**Inputs:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| subscription_id | string | Yes | Subscription ID |
| --interval | int | Yes | Frequency in weeks (1-26) |
| --confirm | flag | Yes | Required to execute |

**Command:**
```bash
amazon-cli subscriptions frequency S01-1234567-8901234 --interval 6 --confirm
```

**Output:**
```json
{
  "subscription_id": "S01-1234567-8901234",
  "previous_frequency_weeks": 4,
  "new_frequency_weeks": 6,
  "next_delivery": "2024-02-15"
}
```

---

### subscriptions cancel

Cancel a subscription. **Requires --confirm flag.**

**Inputs:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| subscription_id | string | Yes | Subscription ID |
| --confirm | flag | Yes | Required to execute |

**Command:**
```bash
amazon-cli subscriptions cancel S01-1234567-8901234 --confirm
```

**Output:**
```json
{
  "subscription_id": "S01-1234567-8901234",
  "status": "cancelled",
  "cancelled_at": "2024-01-18T12:00:00Z"
}
```

---

### subscriptions upcoming

View upcoming subscription deliveries.

**Inputs:** None

**Command:**
```bash
amazon-cli subscriptions upcoming
```

**Output:**
```json
{
  "upcoming_deliveries": [
    {
      "subscription_id": "S01-1234567-8901234",
      "asin": "B00EXAMPLE",
      "title": "Coffee Pods 100 Count",
      "delivery_date": "2024-02-01",
      "quantity": 1
    }
  ]
}
```

---

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

---

## Safety Guidelines

1. **Always preview before purchasing** - Run commands without `--confirm` first to see what will happen
2. **Verify cart contents** - Use `cart list` before `cart checkout --confirm`
3. **Check subscription changes** - Review output before confirming subscription modifications
4. **Keep credentials secure** - Config file at `~/.amazon-cli/config.json` contains tokens

---

## Examples

### Check recent orders
```bash
amazon-cli orders list --limit 5
```

### Find tracking info for an order
```bash
amazon-cli orders track 123-4567890-1234567
```

### Search for wireless headphones under $100
```bash
amazon-cli search "wireless headphones" --max-price 100 --prime-only
```

### Return a defective item
```bash
amazon-cli returns create 123-4567890-1234567 ITEM123 --reason defective --confirm
```

### Skip next Subscribe & Save delivery
```bash
amazon-cli subscriptions skip S01-1234567-8901234 --confirm
```

### Buy an item immediately
```bash
amazon-cli buy B08N5WRWNW --quantity 1 --confirm
```

---

## Rate Limiting

The CLI includes built-in rate limiting to avoid triggering Amazon's anti-automation:
- Minimum 1 second delay between requests
- Random jitter added to delays
- Exponential backoff on 429/503 responses
- Maximum 3 retries before failing

---

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--output` | `-o` | Output format: json (default), table, raw |
| `--quiet` | `-q` | Suppress non-essential output |
| `--verbose` | `-v` | Enable verbose logging |
| `--config` | | Path to config file |
| `--no-color` | | Disable colored output |
