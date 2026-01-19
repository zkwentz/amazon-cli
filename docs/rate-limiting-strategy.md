# Rate Limiting Strategy

## Overview

This document outlines Amazon's rate limiting behavior and the recommended strategies to avoid triggering anti-automation measures when using the amazon-cli tool. Understanding and respecting these limits is crucial for maintaining reliable access to Amazon's services without encountering CAPTCHAs or IP blocks.

## Amazon's Rate Limiting Behavior

### Requests Per Minute (RPM) Limits

Amazon implements dynamic rate limiting that varies based on several factors:

- **Conservative threshold:** 20-30 requests per minute per IP address
- **Aggressive threshold:** 60+ requests per minute will likely trigger blocking
- **Recommended safe limit:** 15-20 requests per minute to maintain a safety margin

These limits can vary based on:
- Account history and reputation
- Time of day (higher limits during peak hours)
- Type of requests (search vs. checkout vs. order history)
- Geographic location
- Previous automation detection patterns

### CAPTCHA Trigger Conditions

Amazon's anti-bot systems will present CAPTCHA challenges or block access when detecting:

1. **High Request Velocity**
   - More than 30 requests per minute consistently
   - Sudden bursts of requests (e.g., 10 requests in 10 seconds)
   - Sustained high-volume traffic over several minutes

2. **Suspicious Request Patterns**
   - Requests with identical timing intervals (e.g., exactly 1.000s between each request)
   - Sequential product browsing without natural pauses
   - Immediate cart actions without typical "thinking time"
   - Searches followed immediately by purchases without viewing product details

3. **Missing Browser Characteristics**
   - Lack of typical browser headers (User-Agent, Accept-Language, etc.)
   - Missing cookies or session tokens
   - No JavaScript execution fingerprint
   - Absence of mouse movement or interaction patterns

4. **Account-Level Indicators**
   - Multiple rapid login attempts
   - Accessing orders from different geographic locations simultaneously
   - Unusual access patterns (e.g., viewing 100s of orders in quick succession)

5. **Network-Level Signals**
   - Requests from known datacenter IP ranges
   - Traffic from VPN or proxy services
   - Multiple accounts accessed from the same IP
   - Lack of TLS fingerprint diversity

## Recommended Rate Limiting Strategy

### Minimum Delays Between Requests

To safely interact with Amazon's APIs without triggering rate limits:

1. **Base Delay: 1-2 seconds**
   - Minimum delay: **1000ms** (1 second)
   - Recommended delay: **1500ms** (1.5 seconds)
   - Conservative delay: **2000ms** (2 seconds)

2. **Request-Specific Delays**
   - **Search queries:** 1.5-3 seconds between searches
   - **Product detail views:** 1-2 seconds between products
   - **Cart operations:** 2-3 seconds between cart modifications
   - **Checkout actions:** 3-5 seconds between checkout steps
   - **Order history requests:** 2-4 seconds between page loads

3. **Jitter Implementation**
   - Add random jitter to prevent detectable patterns
   - Jitter range: **0-500ms** added to base delay
   - Example: 1500ms base + 0-500ms jitter = 1500-2000ms actual delay

### Implementation Guidelines

#### Basic Rate Limiter Configuration

```json
{
  "rate_limiting": {
    "min_delay_ms": 1000,
    "max_delay_ms": 2000,
    "jitter_ms": 500,
    "max_retries": 3,
    "backoff_multiplier": 2,
    "max_backoff_ms": 60000
  }
}
```

#### Request Flow with Rate Limiting

1. **Before Request**
   ```
   base_delay = min_delay_ms + random(0, jitter_ms)
   wait(base_delay)
   ```

2. **Execute Request**
   - Send HTTP request with proper headers
   - Include session cookies and tokens
   - Maintain connection pooling

3. **After Request**
   - If `200 OK`: Continue normal flow
   - If `429 Too Many Requests`: Apply exponential backoff
   - If `503 Service Unavailable`: Apply exponential backoff
   - If CAPTCHA detected: Pause and notify user

#### Exponential Backoff Strategy

When rate limiting is detected (429, 503, or CAPTCHA):

```
retry_delay = min(
    max_backoff_ms,
    base_delay * (backoff_multiplier ^ retry_attempt)
)

Example:
- Attempt 1: 2 seconds
- Attempt 2: 4 seconds
- Attempt 3: 8 seconds
- Attempt 4+: 60 seconds (capped at max_backoff_ms)
```

Maximum retries: **3 attempts** before returning error to user

### Advanced Techniques

#### Request Prioritization

When rate limiting is a concern, prioritize requests:

1. **High Priority** (minimal delay acceptable)
   - Purchase confirmation
   - Authentication token refresh
   - Critical order updates

2. **Medium Priority** (standard delays)
   - Product searches
   - Cart operations
   - Order history retrieval

3. **Low Priority** (extended delays acceptable)
   - Bulk order history scraping
   - Product detail crawling
   - Review collection

#### Burst Protection

Prevent rapid bursts that trigger detection:

```
- Track request timestamps in sliding window
- If requests_in_last_60s > 20: enforce cooling period
- Cooling period: 10-30 seconds with no requests
- Resume with increased base delay (e.g., 2-3 seconds)
```

#### Session Management

Maintain realistic session behavior:

- Keep sessions alive with periodic heartbeats
- Don't create new sessions for every request
- Reuse cookies and authentication tokens
- Respect session expiration times

### Error Handling

When rate limiting occurs:

1. **Detection**
   - Monitor for HTTP 429 (Too Many Requests)
   - Check for HTTP 503 (Service Unavailable)
   - Detect CAPTCHA pages in HTML responses
   - Watch for unusual redirect patterns

2. **Response**
   - Log the rate limit event
   - Increment retry counter
   - Apply exponential backoff
   - If max retries exceeded: return error to user

3. **User Notification**
   ```json
   {
     "error": {
       "code": "RATE_LIMITED",
       "message": "Request rate limit exceeded. Please wait before retrying.",
       "details": {
         "retry_after_seconds": 60,
         "requests_made": 35,
         "window_seconds": 60
       }
     }
   }
   ```

### Testing Rate Limits

To verify rate limiting implementation:

1. **Gradual Load Test**
   - Start with 10 requests/minute
   - Gradually increase to 20, 30, 40 requests/minute
   - Monitor for CAPTCHA or blocking responses
   - Identify safe threshold for your use case

2. **Pattern Detection Test**
   - Test with fixed intervals (no jitter)
   - Test with randomized jitter
   - Compare success rates
   - Measure time-to-CAPTCHA

3. **Burst Test**
   - Send 5-10 rapid requests
   - Measure recovery time needed
   - Validate backoff implementation

## Best Practices Summary

1. **Always use jitter** (0-500ms) to randomize request timing
2. **Maintain minimum 1-2 second delays** between requests
3. **Implement exponential backoff** for retries (max 3 attempts)
4. **Monitor for rate limit signals** (429, 503, CAPTCHA)
5. **Use realistic session patterns** (keep sessions alive, reuse cookies)
6. **Prioritize requests** when approaching limits
7. **Implement burst protection** to prevent rapid request clusters
8. **Test your limits** gradually to find safe thresholds
9. **Respect cooling-off periods** after hitting rate limits
10. **Log rate limit events** for monitoring and adjustment

## Configuration Examples

### Conservative (Safest)
```json
{
  "min_delay_ms": 2000,
  "max_delay_ms": 3000,
  "jitter_ms": 500,
  "max_rpm": 15
}
```

### Balanced (Recommended)
```json
{
  "min_delay_ms": 1500,
  "max_delay_ms": 2000,
  "jitter_ms": 500,
  "max_rpm": 20
}
```

### Aggressive (Risk of blocking)
```json
{
  "min_delay_ms": 1000,
  "max_delay_ms": 1500,
  "jitter_ms": 300,
  "max_rpm": 30
}
```

## Monitoring and Adjustment

Continuously monitor your rate limiting effectiveness:

- **Track success rates** for different delay configurations
- **Log CAPTCHA encounters** and analyze patterns
- **Measure request latency** to detect slowdowns
- **Adjust parameters** based on observed behavior
- **Document changes** in rate limiting patterns over time

Amazon's rate limiting behavior may change over time. Regularly review and adjust your strategy based on observed patterns and any blocking incidents.

## References

- Amazon Rate Limiting: Dynamic and account-specific
- Recommended base delay: 1-2 seconds with 0-500ms jitter
- CAPTCHA triggers: >30 RPM, pattern detection, missing browser signals
- Exponential backoff: 2^n with 60-second maximum
- Maximum retries: 3 attempts before failure
