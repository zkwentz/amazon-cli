package amazon

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

// Client represents an Amazon API client
type Client struct {
	httpClient *http.Client
}

// NewClient creates a new Amazon API client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
	}
}

// GetOrderTracking fetches tracking information for a specific order
// orderID: Amazon order ID (e.g., "123-4567890-1234567")
// Returns tracking information or an error
func (c *Client) GetOrderTracking(orderID string) (*models.Tracking, error) {
	if orderID == "" {
		return nil, fmt.Errorf("order ID cannot be empty")
	}

	// Validate order ID format (Amazon order IDs are typically in format XXX-XXXXXXX-XXXXXXX)
	orderIDPattern := regexp.MustCompile(`^\d{3}-\d{7}-\d{7}$`)
	if !orderIDPattern.MatchString(orderID) {
		return nil, fmt.Errorf("invalid order ID format: %s (expected format: XXX-XXXXXXX-XXXXXXX)", orderID)
	}

	// Construct the tracking URL for the order
	trackingURL := fmt.Sprintf("https://www.amazon.com/progress-tracker/package/ref=ppx_yo_dt_b_track_package?_encoding=UTF8&itemId=%s&orderId=%s", orderID, orderID)

	// Create the HTTP request
	req, err := http.NewRequest("GET", trackingURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers to mimic a browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tracking information: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("authentication required: please log in to Amazon")
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("order not found: %s", orderID)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse the tracking information from the HTML response
	tracking, err := parseTrackingInfo(string(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse tracking information: %w", err)
	}

	return tracking, nil
}

// parseTrackingInfo extracts tracking information from the HTML response
func parseTrackingInfo(html string) (*models.Tracking, error) {
	tracking := &models.Tracking{}

	// Extract carrier name
	carrierPattern := regexp.MustCompile(`(?i)carrier["\s:]+([A-Za-z\s]+)`)
	if matches := carrierPattern.FindStringSubmatch(html); len(matches) > 1 {
		tracking.Carrier = strings.TrimSpace(matches[1])
	}

	// Extract tracking number (typically 10-40 alphanumeric characters)
	trackingNumPattern := regexp.MustCompile(`(?i)tracking[_\s]+number["\s:]+([A-Z0-9]{10,40})`)
	if matches := trackingNumPattern.FindStringSubmatch(html); len(matches) > 1 {
		tracking.TrackingNumber = strings.TrimSpace(matches[1])
	}

	// Extract status (delivered, in transit, out for delivery, etc.)
	statusPattern := regexp.MustCompile(`(?i)status["\s:]+([A-Za-z\s]+)`)
	if matches := statusPattern.FindStringSubmatch(html); len(matches) > 1 {
		tracking.Status = strings.ToLower(strings.TrimSpace(matches[1]))
	}

	// Extract delivery date (various date formats)
	datePattern := regexp.MustCompile(`(?i)delivered?["\s:]+([A-Za-z]+\s+\d{1,2},?\s+\d{4})`)
	if matches := datePattern.FindStringSubmatch(html); len(matches) > 1 {
		tracking.DeliveryDate = strings.TrimSpace(matches[1])
	}

	// If we didn't find any tracking information, return an error
	if tracking.Carrier == "" && tracking.TrackingNumber == "" && tracking.Status == "" {
		return nil, fmt.Errorf("no tracking information found in response")
	}

	return tracking, nil
}
