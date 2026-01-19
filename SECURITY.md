# Security Policy

## Supported Versions

We release patches for security vulnerabilities in the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |
| < 0.1   | :x:                |

We recommend always using the latest stable release to ensure you have all security patches and updates.

## Reporting Vulnerabilities

We take the security of amazon-cli seriously. If you discover a security vulnerability, please follow these steps:

### How to Report

**Email:** security@amazon-cli.io

Please include the following information in your report:

- Description of the vulnerability
- Steps to reproduce the issue
- Potential impact and severity
- Suggested fix (if you have one)
- Your contact information for follow-up

### What to Expect

1. **Acknowledgment:** We will acknowledge receipt of your vulnerability report within 48 hours.

2. **Initial Assessment:** We will provide an initial assessment of the report within 5 business days, including an expected timeline for a fix.

3. **Updates:** We will keep you informed of our progress throughout the investigation and remediation process.

4. **Resolution:** Once the vulnerability is fixed, we will:
   - Release a security patch
   - Publish a security advisory
   - Credit you for the discovery (unless you prefer to remain anonymous)

### Disclosure Policy

- Please do not publicly disclose the vulnerability until we have released a fix
- We aim to resolve critical vulnerabilities within 30 days
- We will coordinate the disclosure timeline with you

## Security Best Practices

Follow these security best practices when using amazon-cli:

### 1. Configuration File Permissions

**Critical:** Ensure your configuration file has restricted permissions to prevent unauthorized access to your credentials.

```bash
# Set correct permissions on config file
chmod 0600 ~/.amazon-cli/config.json
```

The configuration file at `~/.amazon-cli/config.json` contains sensitive authentication tokens. Setting permissions to `0600` ensures that only your user account can read and write this file.

**Verify permissions:**
```bash
ls -la ~/.amazon-cli/config.json
# Should show: -rw------- (0600)
```

### 2. Credential Rotation

Regularly rotate your authentication credentials to minimize the risk of compromised tokens:

- **Re-authenticate periodically:** Run `amazon-cli auth login` to generate fresh tokens
- **Logout when not in use:** Use `amazon-cli auth logout` to clear stored credentials when you're done using the CLI
- **Monitor token expiration:** Check `amazon-cli auth status` to see when your tokens expire
- **Revoke access:** If you suspect your credentials have been compromised, immediately logout and re-authenticate

```bash
# Logout and clear credentials
amazon-cli auth logout

# Re-authenticate to get fresh tokens
amazon-cli auth login
```

### 3. Use --confirm Flag for Sensitive Operations

amazon-cli requires the `--confirm` flag for all operations that modify data or make purchases. This is a safety mechanism to prevent accidental or unauthorized actions.

**Commands requiring --confirm:**

- `amazon-cli cart checkout --confirm` - Complete a purchase
- `amazon-cli cart clear --confirm` - Clear your shopping cart
- `amazon-cli buy <asin> --confirm` - Quick buy an item
- `amazon-cli returns create <order-id> <item-id> --confirm` - Initiate a return
- `amazon-cli subscriptions skip <subscription-id> --confirm` - Skip subscription delivery
- `amazon-cli subscriptions frequency <subscription-id> --interval <weeks> --confirm` - Change subscription frequency
- `amazon-cli subscriptions cancel <subscription-id> --confirm` - Cancel a subscription

**Without --confirm flag:** Commands show a preview of what would happen without executing the action.

```bash
# Safe: Preview checkout without executing
amazon-cli cart checkout

# Executes the purchase (requires explicit confirmation)
amazon-cli cart checkout --confirm
```

**Important:** Never script or automate commands with `--confirm` without proper safeguards and monitoring.

### 4. Additional Security Recommendations

- **Keep software updated:** Regularly update amazon-cli to the latest version
- **Use in trusted environments:** Avoid using amazon-cli on shared or untrusted systems
- **Secure your system:** Ensure your operating system and development environment are secure
- **Review automation scripts:** Carefully audit any scripts that use amazon-cli
- **Enable verbose logging cautiously:** Tokens never appear in verbose output, but avoid sharing logs publicly
- **Monitor account activity:** Regularly check your Amazon account for unauthorized actions

### 5. Security Features

amazon-cli implements several security features:

- **HTTPS only:** All communication with Amazon uses TLS encryption
- **No credential logging:** Authentication tokens are never written to logs
- **Rate limiting:** Built-in rate limiting helps prevent detection as automation
- **Token expiration:** Authentication tokens expire automatically
- **Confirmation required:** Destructive operations require explicit confirmation

## Reporting Security Issues in Dependencies

If you discover a vulnerability in one of our dependencies, please:

1. Report it to the dependency's maintainers directly
2. Also notify us at security@amazon-cli.io so we can track and update accordingly

## Security Update Notifications

To stay informed about security updates:

- Watch the GitHub repository for security advisories
- Subscribe to release notifications
- Follow the project changelog

Thank you for helping keep amazon-cli and its users secure.
