# amazon-cli

A command-line interface for Amazon shopping functionality, providing programmatic access to orders, returns, purchases, and subscriptions.

## Installation

### Homebrew (macOS/Linux)
```bash
brew tap zkwentz/amazon-cli
brew install amazon-cli
```

### Binary Releases
Download pre-compiled binaries for your platform from the [GitHub Releases](https://github.com/zkwentz/amazon-cli/releases) page.

Available platforms:
- macOS (arm64, amd64)
- Linux (arm64, amd64)
- Windows (amd64)

### Build from Source
```bash
go install github.com/zkwentz/amazon-cli@latest
```

## Authentication

### Authentication Method

**amazon-cli uses Browser-based OAuth authentication** via Amazon's Login with Amazon (LWA) API. This provides secure, user-authorized access to your Amazon account data.

#### How It Works

1. Run `amazon-cli auth login`
2. The CLI opens your default browser to Amazon's OAuth consent page
3. You log in with your Amazon credentials (supports 2FA)
4. After successful authentication, OAuth tokens are captured and stored locally
5. The CLI automatically refreshes tokens when they expire

#### Token Storage

Tokens are stored in plain text at `~/.amazon-cli/config.json` with the following structure:

```json
{
  "auth": {
    "access_token": "...",
    "refresh_token": "...",
    "expires_at": "2024-01-20T12:00:00Z"
  }
}
```

**Security Note:** The config file contains sensitive authentication tokens. Ensure it has appropriate file permissions (the CLI sets this to `0600` automatically).

#### Alternative: Browser Session Authentication

If the OAuth API does not provide sufficient access to certain Amazon features, the CLI can fall back to browser session-based authentication:

```bash
amazon-cli auth login --browser
```

This method:
- Opens Amazon's login page in your browser
- Uses browser automation to capture session cookies after you log in
- Stores cookies in the config file for subsequent requests
- Automatically detects when cookies expire and prompts for re-authentication

### Authentication Commands

```bash
# Log in to Amazon
amazon-cli auth login

# Check authentication status
amazon-cli auth status

# Log out and clear stored credentials
amazon-cli auth logout
```

## Quick Start

```bash
# Authenticate
amazon-cli auth login

# List recent orders
amazon-cli orders list --limit 5

# Search for products
amazon-cli search "wireless headphones" --prime-only

# View cart
amazon-cli cart list

# Get order details
amazon-cli orders get 123-4567890-1234567
```

## Core Features

### Orders Management
- List recent orders with filtering
- Get detailed order information
- Track shipments
- View order history by year

### Returns Management
- View returnable items
- Check return options
- Initiate returns with reason codes
- Download return labels
- Track return status

### Search & Products
- Search Amazon's catalog with filters
- Get detailed product information
- View product reviews

### Cart & Checkout
- Add/remove items from cart
- View cart contents
- Checkout with saved addresses and payment methods
- Quick buy for immediate purchases

### Subscribe & Save
- List all subscriptions
- Skip upcoming deliveries
- Change delivery frequency
- Cancel subscriptions
- View upcoming deliveries

## Configuration

Configuration file location: `~/.amazon-cli/config.json`

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

### Rate Limiting

The CLI implements rate limiting to avoid triggering Amazon's anti-automation measures:
- Minimum 1 second delay between requests (configurable)
- Random jitter added to each delay
- Exponential backoff on 429/503 responses
- User-agent rotation

## Safety Features

### Confirmation Requirements

All purchase and modification operations require the `--confirm` flag:

```bash
# Without --confirm: Shows what would happen (dry run)
amazon-cli cart checkout

# With --confirm: Actually executes the operation
amazon-cli cart checkout --confirm
```

This prevents accidental purchases or modifications.

## Output Formats

All commands output JSON by default for easy parsing by AI agents and scripts:

```bash
# JSON output (default)
amazon-cli orders list

# Table format for human readability
amazon-cli orders list --output table

# Raw format
amazon-cli orders list --output raw
```

## Error Handling

Errors are returned as JSON with consistent structure:

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
- `AUTH_REQUIRED` - Not logged in
- `AUTH_EXPIRED` - Token expired
- `NOT_FOUND` - Resource not found
- `RATE_LIMITED` - Too many requests
- `INVALID_INPUT` - Invalid command input
- `PURCHASE_FAILED` - Purchase could not be completed
- `NETWORK_ERROR` - Network connectivity issue
- `AMAZON_ERROR` - Amazon returned an error

## Global Flags

```bash
--output, -o      Output format: json (default), table, raw
--quiet, -q       Suppress non-essential output
--verbose, -v     Enable verbose logging
--config          Path to config file (default: ~/.amazon-cli/config.json)
--no-color        Disable colored output
```

## Examples

### Check recent orders
```bash
amazon-cli orders list --limit 10 --status delivered
```

### Search for products under $50
```bash
amazon-cli search "coffee maker" --max-price 50 --prime-only
```

### Return a defective item
```bash
# Preview the return
amazon-cli returns create 123-4567890-1234567 ITEM123 --reason defective

# Confirm the return
amazon-cli returns create 123-4567890-1234567 ITEM123 --reason defective --confirm
```

### Skip next Subscribe & Save delivery
```bash
amazon-cli subscriptions skip S01-1234567-8901234 --confirm
```

## Security Considerations

1. **Credentials Storage**: Tokens are stored in plain text in `~/.amazon-cli/config.json`. Users are responsible for protecting this file (the CLI sets permissions to `0600`).

2. **Confirmation Flags**: All purchase and modification operations require explicit `--confirm` flags to prevent accidents.

3. **No Credential Logging**: Authentication tokens never appear in logs, even in verbose mode.

4. **HTTPS Only**: All communication with Amazon uses TLS encryption.

## Privacy & Terms of Service

This tool interacts with Amazon's services on your behalf. Please review:
- Amazon's Terms of Service
- Amazon's API Usage Policies
- Your local laws regarding automation of web services

## Target Marketplace

Currently supports **US only** (amazon.com). Future expansion to other marketplaces may be considered.

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## License

MIT License - see LICENSE file for details.

## Support

For issues, questions, or feature requests, please visit the [GitHub Issues](https://github.com/zkwentz/amazon-cli/issues) page.
