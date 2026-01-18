# Amazon Login with Amazon (LWA) OAuth API Research

## Overview

Login with Amazon (LWA) is Amazon's OAuth 2.0-based authentication service that allows applications to leverage Amazon's authentication system. It enables users to sign in using their Amazon credentials and grant applications access to their profile information.

## Key Concepts

### OAuth 2.0 Standard
LWA is built on OAuth 2.0 principles, providing secure delegated access without exposing user credentials to third-party applications.

### Security Model
- Users authenticate directly with Amazon
- Applications receive temporary access tokens
- Users explicitly consent to data sharing via consent screen
- Tokens expire after 1 hour for enhanced security

## Authentication Flow (Authorization Code Grant)

### Step 1: Authorization Request
Direct users to Amazon's authorization endpoint to obtain an authorization code.

**Endpoint:** `https://www.amazon.com/ap/oa`

**Method:** GET (redirect)

**Required Parameters:**
- `client_id` - Your application identifier (max 100 bytes)
- `scope` - Requested permissions (e.g., "profile", "profile:user_id", "postal_code")
- `response_type` - Must be "code"
- `redirect_uri` - HTTPS callback URL for your application

**Recommended Parameters:**
- `state` - Random string for CSRF protection and request correlation
- `code_challenge` - For PKCE security (recommended for browser apps)
- `code_challenge_method` - Typically "S256" for SHA-256 hashing

**Example:**
```
https://www.amazon.com/ap/oa?client_id=YOUR_CLIENT_ID&scope=profile&response_type=code&redirect_uri=https://example.com/callback&state=RANDOM_STATE
```

### Step 2: User Authentication & Consent
1. User logs in with Amazon credentials (if not already logged in)
2. User sees consent screen with requested permissions
3. User approves or denies access

### Step 3: Authorization Response
Amazon redirects user back to your `redirect_uri` with:

**Success Response:**
- `code` - Authorization code (18-128 characters, valid for 5 minutes)
- `scope` - User-consented scopes
- `state` - Echo of the state parameter

**Error Response:**
- `error` - Error code (e.g., "access_denied", "invalid_request")
- `error_description` - Human-readable error description
- `state` - Echo of the state parameter

### Step 4: Exchange Authorization Code for Access Token
Exchange the authorization code for an access token using the token endpoint.

**Endpoint:** Regional token endpoints (applications are not region-specific)
- North America: `https://api.amazon.com/auth/o2/token`
- European Union: `https://api.amazon.co.uk/auth/o2/token`
- Far East: `https://api.amazon.co.jp/auth/o2/token`

**Method:** POST

**Content-Type:** `application/x-www-form-urlencoded`

**Required Parameters:**
- `grant_type` - Must be "authorization_code"
- `code` - The authorization code received
- `redirect_uri` - Must match the authorization request URI
- `client_id` - Your application identifier
- `client_secret` - Your application secret (optional for browser apps using PKCE)
- `code_verifier` - Required if code_challenge was used (PKCE)

**Example Request:**
```http
POST /auth/o2/token HTTP/1.1
Host: api.amazon.com
Content-Type: application/x-www-form-urlencoded

grant_type=authorization_code&code=AUTH_CODE&redirect_uri=https://example.com/callback&client_id=YOUR_CLIENT_ID&client_secret=YOUR_CLIENT_SECRET
```

**Success Response (HTTP 200):**
```json
{
  "access_token": "Atza|IQEBLjAsAhRmHjNgHpi0U...",
  "token_type": "bearer",
  "expires_in": 3600,
  "refresh_token": "Atzr|IQEBLjAsAhRmHjNgHpi0U..."
}
```

**Response Fields:**
- `access_token` - Bearer token for API access (max 2048 bytes, starts with "Atza|")
- `token_type` - Always "bearer"
- `expires_in` - Token lifetime in seconds (typically 3600 = 1 hour)
- `refresh_token` - Token for obtaining new access tokens (starts with "Atzr|")

**Error Response:**
- `error` - Error code (e.g., "invalid_grant", "invalid_client")
- `error_description` - Human-readable error description

## Access Tokens

### Format
- Alphanumeric strings
- Minimum length: 350 characters
- Maximum size: 2048 bytes
- Prefix: "Atza|"

### Characteristics
- Valid for 1 hour (3600 seconds)
- Bearer tokens (can be used by any client in possession)
- Must be transmitted over HTTPS
- Contain URL-unsafe characters (must be URL-encoded)

### Security Considerations
- Store securely and never expose to users
- Transmit only over secure channels
- Implement token expiration handling
- Use refresh tokens to obtain new access tokens

## Refresh Tokens

### Format
- Similar to access tokens
- Prefix: "Atzr|"

### Usage
Refresh tokens have extended validity and allow obtaining new access tokens without user re-authentication.

**Refresh Token Request:**

**Endpoint:** `https://api.amazon.com/auth/o2/token`

**Method:** POST

**Parameters:**
- `grant_type` - Must be "refresh_token"
- `refresh_token` - The refresh token
- `client_id` - Your application identifier
- `client_secret` - Your application secret

**Response:**
Returns a new access token and refresh token pair.

## Obtaining User Profile Information

### Profile Endpoint

**URL:** `https://api.amazon.com/user/profile`

**Method:** GET

**Authentication:** Three methods supported

1. **Query Parameter:**
```
https://api.amazon.com/user/profile?access_token=YOUR_ACCESS_TOKEN
```

2. **Bearer Token (Recommended):**
```http
GET /user/profile HTTP/1.1
Host: api.amazon.com
Authorization: Bearer Atza|IQEBLjAsAhRmHjNgHpi0U...
Accept: application/json
Accept-Language: en-US
```

3. **Custom Header:**
```http
GET /user/profile HTTP/1.1
Host: api.amazon.com
x-amz-access-token: Atza|IQEBLjAsAhRmHjNgHpi0U...
Accept: application/json
```

### Profile Response

**Success Response (HTTP 200):**
```json
{
  "user_id": "amznl.account.K2LI23KL2LK2",
  "email": "user@example.com",
  "name": "User Name",
  "postal_code": "98052"
}
```

**Available Fields (depends on requested scope):**
- `user_id` - Unique Amazon user identifier
- `email` - User's email address
- `name` - User's full name
- `postal_code` - User's postal/ZIP code

**Note:** Not all fields are guaranteed; returned data depends on the requested and consented scopes.

## Scopes

### Available Scopes
- `profile` - Access to user's name and email
- `profile:user_id` - Access to user_id only (minimal identification)
- `postal_code` - Access to user's postal/ZIP code

### Scope Combinations
Multiple scopes can be requested by separating them with spaces:
```
scope=profile postal_code
```

### Essential vs. Voluntary Scopes
LWA supports marking scopes as essential or voluntary, allowing users to selectively grant permissions.

## PKCE (Proof Key for Code Exchange)

### When to Use
- Required for browser-based apps without server-side component
- Recommended for enhanced security in all scenarios
- Necessary when client_secret cannot be securely stored

### Implementation

**Step 1: Generate Code Verifier**
- Random string (43-128 characters)
- Uses characters: A-Z, a-z, 0-9, and "-", ".", "_", "~"

**Step 2: Generate Code Challenge**
- SHA-256 hash of code_verifier
- Base64 URL-encoded

**Step 3: Authorization Request**
Include `code_challenge` and `code_challenge_method=S256` in authorization request.

**Step 4: Token Request**
Include `code_verifier` (original random string) in token request.

### Security Benefit
Prevents authorization code interception attacks by requiring proof of the original request.

## Error Handling

### Authorization Errors
- `invalid_request` - Missing or invalid parameters
- `unauthorized_client` - Client not authorized for this grant type
- `access_denied` - User denied authorization
- `unsupported_response_type` - Invalid response_type parameter
- `invalid_scope` - Invalid or unknown scope
- `server_error` - Authorization server error
- `temporarily_unavailable` - Server temporarily unavailable

### Token Errors
- `invalid_request` - Missing or invalid parameters
- `invalid_client` - Invalid client credentials
- `invalid_grant` - Invalid authorization code or refresh token
- `unauthorized_client` - Client not authorized for this grant type
- `unsupported_grant_type` - Invalid grant_type parameter

### Best Practices
- Always validate the `state` parameter to prevent CSRF attacks
- Handle all error responses gracefully
- Implement retry logic for temporary failures
- Log errors for debugging and monitoring

## Regional Considerations

### Token Endpoints
While regional endpoints exist, applications are NOT region-specific. You can create, verify, and refresh tokens using any regional LWA endpoint.

### Supported Regions
- North America: api.amazon.com
- European Union: api.amazon.co.uk
- Far East: api.amazon.co.jp

## Implementation Example (Conceptual Flow)

```javascript
// Step 1: Redirect to Amazon for authorization
const authUrl = 'https://www.amazon.com/ap/oa?' +
  'client_id=YOUR_CLIENT_ID&' +
  'scope=profile&' +
  'response_type=code&' +
  'redirect_uri=https://example.com/callback&' +
  'state=RANDOM_STATE';
window.location.href = authUrl;

// Step 2: Handle callback (server-side)
// Extract code from query parameters
const authCode = req.query.code;
const state = req.query.state;

// Validate state parameter (CSRF protection)
if (state !== expectedState) {
  throw new Error('Invalid state parameter');
}

// Step 3: Exchange code for token
const tokenResponse = await fetch('https://api.amazon.com/auth/o2/token', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/x-www-form-urlencoded'
  },
  body: new URLSearchParams({
    grant_type: 'authorization_code',
    code: authCode,
    redirect_uri: 'https://example.com/callback',
    client_id: 'YOUR_CLIENT_ID',
    client_secret: 'YOUR_CLIENT_SECRET'
  })
});

const tokenData = await tokenResponse.json();
const accessToken = tokenData.access_token;

// Step 4: Get user profile
const profileResponse = await fetch('https://api.amazon.com/user/profile', {
  headers: {
    'Authorization': `Bearer ${accessToken}`,
    'Accept': 'application/json'
  }
});

const profile = await profileResponse.json();
console.log('User:', profile.name, profile.email);
```

## Security Best Practices

### 1. State Parameter
- Always use the state parameter
- Generate a cryptographically random value
- Validate it matches on callback
- Prevents CSRF attacks

### 2. PKCE
- Use PKCE for browser-based applications
- Recommended even for server-side apps
- Prevents authorization code interception

### 3. Token Storage
- Never expose tokens to client-side code if avoidable
- Store refresh tokens securely (encrypted database)
- Use secure, httpOnly cookies for session management
- Implement token rotation

### 4. HTTPS
- All redirect URIs must use HTTPS
- All API calls must use HTTPS
- Prevents token interception

### 5. Token Validation
- Validate token expiration
- Implement refresh logic before expiration
- Handle token refresh failures gracefully

### 6. Scope Minimization
- Request only necessary scopes
- Use profile:user_id when full profile not needed
- Improves user trust and conversion rates

## Testing and Development

### Developer Console
Create and manage LWA applications at: https://developer.amazon.com

### Security Profile
- Create a security profile for your application
- Configure allowed return URLs
- Obtain client_id and client_secret
- Manage application settings

### Testing Considerations
- Test with different user consent scenarios (approve/deny)
- Test token expiration and refresh flows
- Test error handling for all error types
- Validate PKCE implementation
- Test across different browsers and devices

## Additional Resources

### Official Documentation
- Main Documentation: https://developer.amazon.com/docs/login-with-amazon/documentation-overview.html
- Authorization Code Grant: https://developer.amazon.com/docs/login-with-amazon/authorization-code-grant.html
- Conceptual Overview: https://developer.amazon.com/docs/login-with-amazon/conceptual-overview.html
- Access Tokens: https://developer.amazon.com/docs/login-with-amazon/access-token.html
- Customer Profile: https://developer.amazon.com/docs/login-with-amazon/obtain-customer-profile.html

### SDK Support
Amazon provides SDKs for multiple platforms:
- JavaScript (Web)
- iOS
- Android
- For other platforms, use direct REST API calls

## Summary

Amazon's Login with Amazon (LWA) OAuth API provides a robust, secure authentication solution based on OAuth 2.0 standards. Key points:

1. **Standard OAuth 2.0 Flow**: Authorization Code Grant is the primary flow
2. **Regional Flexibility**: Applications work across all regions
3. **Security Features**: PKCE support, state parameter, short-lived tokens
4. **Profile Access**: Easy access to user profile data with explicit consent
5. **Refresh Tokens**: Long-lived sessions without repeated user authentication
6. **Comprehensive Documentation**: Well-documented with examples in multiple languages

The implementation is straightforward, following standard OAuth 2.0 patterns, making it familiar to developers with OAuth experience.
