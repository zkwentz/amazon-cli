# Amazon API Research: Authentication & Access Methods

## Overview

This document provides a comprehensive overview of Amazon's authentication options and API access methods for programmatic interaction with Amazon services, specifically focusing on order history, shopping cart, and product search functionality.

## Authentication Methods

### 1. Login with Amazon (LWA) OAuth

**Type:** Official OAuth 2.0 authentication service

**Description:**
Login with Amazon (LWA) is Amazon's official OAuth 2.0 implementation that allows third-party applications to request permission to access customer data on their behalf.

**Use Cases:**
- Third-party applications needing user authorization
- Accessing customer profile information
- Limited access to order data (through specific APIs)

**Endpoints:**
- Authorization: `https://www.amazon.com/ap/oa`
- Token: `https://api.amazon.com/auth/o2/token`

**Flow:**
1. Redirect user to Amazon's authorization page
2. User grants permissions
3. Receive authorization code
4. Exchange code for access token
5. Use access token for API requests

**Example - Step 1: Authorization Request**
```bash
# Redirect user to this URL (browser-based, not curl)
https://www.amazon.com/ap/oa?client_id=YOUR_CLIENT_ID&scope=profile&response_type=code&redirect_uri=YOUR_REDIRECT_URI
```

**Example - Step 2: Exchange Authorization Code for Token**
```bash
curl -X POST https://api.amazon.com/auth/o2/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=authorization_code" \
  -d "code=YOUR_AUTH_CODE" \
  -d "client_id=YOUR_CLIENT_ID" \
  -d "client_secret=YOUR_CLIENT_SECRET" \
  -d "redirect_uri=YOUR_REDIRECT_URI"
```

**Example - Step 3: Use Access Token**
```bash
curl -X GET https://api.amazon.com/user/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Limitations:**
- Requires app registration with Amazon
- Limited scope of accessible data
- Does not provide direct access to order history, cart, or detailed shopping data
- Primarily designed for customer profile information

**Best For:**
- Mobile apps and web applications
- User identity verification
- Basic profile information access

---

### 2. Amazon Advertising API (via LWA)

**Type:** Official API for advertising partners

**Description:**
Uses LWA for authentication but provides access to Amazon Advertising data, not consumer shopping data.

**Example Authorization:**
```bash
curl -X POST https://api.amazon.com/auth/o2/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=refresh_token" \
  -d "refresh_token=YOUR_REFRESH_TOKEN" \
  -d "client_id=YOUR_CLIENT_ID" \
  -d "client_secret=YOUR_CLIENT_SECRET"
```

**Limitations:**
- Only for advertising partners
- Does not provide access to order history or cart

---

### 3. Amazon MWS / SP-API (Selling Partner API)

**Type:** Official API for sellers and vendors

**Description:**
Amazon's marketplace web service for sellers to programmatically exchange data on listings, orders, payments, reports, and more.

**Authentication:** LWA-based with restricted access tokens

**Example - Get Orders (SP-API):**
```bash
curl -X GET "https://sellingpartnerapi-na.amazon.com/orders/v0/orders?MarketplaceIds=ATVPDKIKX0DER&CreatedAfter=2024-01-01T00:00:00Z" \
  -H "x-amz-access-token: YOUR_LWA_ACCESS_TOKEN" \
  -H "x-amz-date: 20240118T120000Z" \
  -H "Authorization: AWS4-HMAC-SHA256 Credential=YOUR_ACCESS_KEY/20240118/us-east-1/execute-api/aws4_request, SignedHeaders=host;x-amz-date, Signature=YOUR_SIGNATURE"
```

**Limitations:**
- Only accessible to registered sellers/vendors
- Requires AWS signature
- Not for end-consumer order access
- Focused on seller/vendor operations

---

### 4. Session Cookies Authentication

**Type:** Unofficial browser session-based authentication

**Description:**
Using authenticated browser session cookies to make requests to Amazon's internal APIs and web endpoints. This mimics browser behavior.

**Required Cookies:**
- `session-id`
- `session-id-time`
- `ubid-main`
- `at-main`
- `sess-at-main`
- `x-main`

**How to Obtain Cookies:**
1. Log in to Amazon via browser
2. Extract cookies from browser's developer tools (Application > Cookies)
3. Use cookies in API requests

**Example - Access Order History:**
```bash
curl -X GET "https://www.amazon.com/gp/your-account/order-history" \
  -H "Cookie: session-id=YOUR_SESSION_ID; session-id-time=YOUR_SESSION_TIME; ubid-main=YOUR_UBID; at-main=YOUR_AT_TOKEN; sess-at-main=YOUR_SESS_AT; x-main=YOUR_X_MAIN" \
  -H "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36" \
  -H "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8" \
  -H "Accept-Language: en-US,en;q=0.9"
```

**Example - Search Products:**
```bash
curl -X GET "https://www.amazon.com/s?k=laptop" \
  -H "Cookie: session-id=YOUR_SESSION_ID; ubid-main=YOUR_UBID" \
  -H "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36" \
  -H "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"
```

**Example - Access Cart:**
```bash
curl -X GET "https://www.amazon.com/gp/cart/view.html" \
  -H "Cookie: session-id=YOUR_SESSION_ID; session-id-time=YOUR_SESSION_TIME; ubid-main=YOUR_UBID; at-main=YOUR_AT_TOKEN" \
  -H "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36" \
  -H "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"
```

**Example - Add Item to Cart:**
```bash
curl -X POST "https://www.amazon.com/gp/aws/cart/add.html" \
  -H "Cookie: session-id=YOUR_SESSION_ID; session-id-time=YOUR_SESSION_TIME; ubid-main=YOUR_UBID; at-main=YOUR_AT_TOKEN; csrf=YOUR_CSRF_TOKEN" \
  -H "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -H "anti-csrftoken-a2z: YOUR_CSRF_TOKEN" \
  -d "ASIN.1=B08N5WRWNW&Quantity.1=1"
```

**Advantages:**
- Full access to all user-visible features
- Order history, cart, search, product details
- No API registration required

**Limitations:**
- Against Amazon's Terms of Service
- Cookies expire and require re-authentication
- Requires HTML parsing (no structured JSON responses in many cases)
- Risk of account suspension
- Frequent changes to HTML structure
- May trigger bot detection/CAPTCHA
- No official support or documentation

**Best For:**
- Personal automation scripts
- Research and development
- Prototyping

---

### 5. Amazon Product Advertising API (PA-API)

**Type:** Official API for product information and affiliate marketing

**Description:**
Provides access to Amazon's product catalog, including search, item details, and browse information. Designed for affiliates and content creators.

**Authentication:** AWS Signature Version 4

**Example - Search Products:**
```bash
curl -X POST "https://webservices.amazon.com/paapi5/searchitems" \
  -H "Content-Type: application/json; charset=utf-8" \
  -H "X-Amz-Target: com.amazon.paapi5.v1.ProductAdvertisingAPIv1.SearchItems" \
  -H "X-Amz-Date: 20240118T120000Z" \
  -H "Authorization: AWS4-HMAC-SHA256 Credential=YOUR_ACCESS_KEY/20240118/us-east-1/ProductAdvertisingAPI/aws4_request, SignedHeaders=content-type;host;x-amz-date;x-amz-target, Signature=YOUR_SIGNATURE" \
  -d '{
    "Keywords": "laptop",
    "Resources": [
      "Images.Primary.Large",
      "ItemInfo.Title",
      "Offers.Listings.Price"
    ],
    "PartnerTag": "YOUR_PARTNER_TAG",
    "PartnerType": "Associates",
    "Marketplace": "www.amazon.com"
  }'
```

**Example - Get Item Details:**
```bash
curl -X POST "https://webservices.amazon.com/paapi5/getitems" \
  -H "Content-Type: application/json; charset=utf-8" \
  -H "X-Amz-Target: com.amazon.paapi5.v1.ProductAdvertisingAPIv1.GetItems" \
  -H "X-Amz-Date: 20240118T120000Z" \
  -H "Authorization: AWS4-HMAC-SHA256 Credential=YOUR_ACCESS_KEY/20240118/us-east-1/ProductAdvertisingAPI/aws4_request, SignedHeaders=content-type;host;x-amz-date;x-amz-target, Signature=YOUR_SIGNATURE" \
  -d '{
    "ItemIds": ["B08N5WRWNW"],
    "Resources": [
      "Images.Primary.Large",
      "ItemInfo.Title",
      "ItemInfo.Features",
      "Offers.Listings.Price"
    ],
    "PartnerTag": "YOUR_PARTNER_TAG",
    "PartnerType": "Associates",
    "Marketplace": "www.amazon.com"
  }'
```

**Limitations:**
- Requires Amazon Associates account
- Does NOT provide access to order history or cart
- Rate limited
- Only for product catalog information

**Best For:**
- Product search and recommendations
- Price comparison
- Affiliate marketing
- Content creation

---

## Comparison Matrix

| Feature | LWA OAuth | SP-API | Session Cookies | PA-API |
|---------|-----------|--------|-----------------|---------|
| **Order History** | ❌ No | ✅ Yes (seller orders only) | ✅ Yes (full access) | ❌ No |
| **Shopping Cart** | ❌ No | ❌ No | ✅ Yes (full access) | ❌ No |
| **Product Search** | ❌ No | ❌ No | ✅ Yes (full access) | ✅ Yes (best for this) |
| **Official Support** | ✅ Yes | ✅ Yes | ❌ No | ✅ Yes |
| **TOS Compliant** | ✅ Yes | ✅ Yes | ❌ No | ✅ Yes |
| **Registration Required** | ✅ Yes | ✅ Yes (seller account) | ❌ No | ✅ Yes (Associates) |
| **Setup Complexity** | Medium | High | Low | Medium |
| **Response Format** | JSON | JSON | HTML | JSON |
| **Rate Limits** | Yes | Yes | Varies | Yes |

---

## Recommended Approaches by Use Case

### For Order History Access

**Best Option: Session Cookies** (only viable option for consumer orders)

**Why:**
- Amazon does not provide an official consumer-facing API for order history
- LWA OAuth does not include order data in its scope
- SP-API only works for seller/vendor orders, not consumer purchases

**Implementation Approach:**
1. Authenticate via browser
2. Extract session cookies
3. Make requests with cookies to order history endpoints
4. Parse HTML responses

**Risks:**
- Violates Amazon TOS
- Requires maintenance as HTML structure changes
- Risk of account suspension

---

### For Shopping Cart Access

**Best Option: Session Cookies** (only option)

**Why:**
- No official API provides cart access for consumers
- Cart functionality is only available through authenticated web sessions

**Implementation Approach:**
1. Use authenticated session cookies
2. Make GET requests to cart endpoints
3. Parse HTML for cart contents
4. Use POST requests with CSRF tokens to modify cart

**Risks:**
- Same as order history (TOS violation, maintenance burden)

---

### For Product Search

**Best Option: Product Advertising API (PA-API)**

**Why:**
- Official, supported API
- Returns structured JSON data
- Comprehensive product information
- TOS compliant

**Fallback: Session Cookies**
- Use if PA-API approval is not possible
- Provides access to all search features
- Requires HTML parsing

---

## Security Considerations

### Session Cookies Method

1. **Cookie Storage:** Never commit cookies to version control
2. **Encryption:** Encrypt cookies at rest
3. **Expiration:** Implement cookie refresh logic
4. **User-Agent:** Rotate user agents to avoid detection
5. **Rate Limiting:** Implement delays between requests
6. **CAPTCHA Handling:** Have fallback for CAPTCHA challenges

### Official APIs

1. **Credential Management:** Store API keys securely (environment variables, secrets manager)
2. **Token Rotation:** Implement refresh token flow for LWA
3. **AWS Signatures:** Use AWS SDK for proper signature generation
4. **HTTPS Only:** Always use encrypted connections

---

## Legal and Compliance

**Official APIs (LWA, SP-API, PA-API):**
- ✅ Terms of Service compliant
- ✅ Legally defensible
- ✅ Receive updates and support

**Session Cookies:**
- ❌ Violates Amazon Terms of Service (Section 3: "You may not... use any robot, spider, scraper, or other automated means to access the Services")
- ❌ Risk of account termination
- ❌ No legal recourse if access is blocked
- ⚠️ Consider for personal use only, never for commercial products

---

## Implementation Recommendations

### For Personal/Research Use
1. Start with session cookies for quick prototyping
2. Implement proper error handling for cookie expiration
3. Add rate limiting and delays
4. Monitor for CAPTCHA challenges

### For Production Applications
1. Use official APIs where possible:
   - PA-API for product search
   - SP-API if you're a seller
   - LWA for user authentication
2. Avoid session cookie approach due to TOS violations
3. Consider alternative data sources (Amazon's official widgets, RSS feeds)

### For Sellers/Vendors
1. Use SP-API for all seller operations
2. Implement OAuth 2.0 flow with refresh tokens
3. Use AWS signature libraries for authentication
4. Follow API rate limits and best practices

---

## Additional Resources

- [Login with Amazon Documentation](https://developer.amazon.com/docs/login-with-amazon/documentation-overview.html)
- [Selling Partner API Documentation](https://developer-docs.amazon.com/sp-api/)
- [Product Advertising API Documentation](https://webservices.amazon.com/paapi5/documentation/)
- [Amazon Associates Program](https://affiliate-program.amazon.com/)

---

## Conclusion

**The Reality Check:**

Amazon intentionally does not provide consumer-facing APIs for order history and shopping cart access. For these use cases, the only technical solution is session cookie-based scraping, which violates their Terms of Service.

**Recommended Path Forward:**

- **For Product Search:** Use PA-API (official, supported)
- **For Order History/Cart:** Either:
  - Accept limitations and don't implement (safest)
  - Use session cookies for personal use only (risky)
  - Contact Amazon Business Development for partnership opportunities
- **For Seller Operations:** Use SP-API (official, supported)

Always prioritize official APIs when available, and carefully weigh the legal and practical risks of unofficial methods.
