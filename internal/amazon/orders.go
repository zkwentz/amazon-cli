package amazon

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

// GetOrder fetches individual order details by order ID
func (c *Client) GetOrder(orderID string) (*models.Order, error) {
	if orderID == "" {
		return nil, fmt.Errorf("order ID cannot be empty")
	}

	// Construct the Amazon order details URL
	url := fmt.Sprintf("https://www.amazon.com/gp/your-account/order-details?orderID=%s", orderID)

	// Create the HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("User-Agent", c.userAgents[0])
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	// Execute the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order details: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("order not found: %s", orderID)
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("authentication required")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the HTML response and extract order details
	order, err := c.parseOrderDetails(orderID, string(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse order details: %w", err)
	}

	return order, nil
}

// parseOrderDetails parses the HTML response and extracts order information
func (c *Client) parseOrderDetails(orderID string, htmlContent string) (*models.Order, error) {
	// This is a placeholder implementation
	// In a real implementation, this would use a proper HTML parser like goquery
	// to extract the actual order details from the Amazon HTML

	// For now, we'll create a basic structure to demonstrate the function
	// In production, this would parse:
	// - Order date
	// - Total amount
	// - Order status (delivered, pending, shipped, etc.)
	// - Items list with ASIN, title, quantity, price
	// - Tracking information if available

	order := &models.Order{
		OrderID: orderID,
		Date:    "", // Would be extracted from HTML
		Total:   0.0, // Would be extracted from HTML
		Status:  "unknown", // Would be extracted from HTML
		Items:   []models.OrderItem{},
	}

	// Check if we received an error page or login redirect
	if strings.Contains(htmlContent, "Sign in") || strings.Contains(htmlContent, "sign-in") {
		return nil, fmt.Errorf("authentication required to access order details")
	}

	// Check if order not found
	if strings.Contains(htmlContent, "order could not be found") ||
	   strings.Contains(htmlContent, "We cannot find this order") {
		return nil, fmt.Errorf("order not found: %s", orderID)
	}

	// NOTE: In a complete implementation, we would:
	// 1. Use goquery or a similar library to parse the HTML
	// 2. Extract order date from elements like: <div class="order-date-invoice-item">
	// 3. Extract total from elements containing price information
	// 4. Extract status from order status banners
	// 5. Parse item details from the items table
	// 6. Extract tracking info if present

	return order, nil
}
