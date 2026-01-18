package amazon

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

// ParseOrderFromJSON parses an order from JSON data
func ParseOrderFromJSON(data []byte) (*models.Order, error) {
	var order models.Order
	if err := json.Unmarshal(data, &order); err != nil {
		return nil, fmt.Errorf("failed to parse order JSON: %w", err)
	}
	return &order, nil
}

// ParseOrdersResponse parses a list of orders from JSON
func ParseOrdersResponse(data []byte) (*models.OrdersResponse, error) {
	var response models.OrdersResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse orders response JSON: %w", err)
	}
	return &response, nil
}

// ParseOrderFromHTML parses order data from HTML scraped content
// This simulates parsing order data from Amazon's HTML pages
func ParseOrderFromHTML(html string) (*models.Order, error) {
	// This is a simplified parser for demonstration
	// In a real implementation, this would use a proper HTML parser like goquery

	if html == "" {
		return nil, fmt.Errorf("empty HTML content")
	}

	order := &models.Order{
		Items: []models.OrderItem{},
	}

	// Extract order ID (e.g., "Order ID: 123-4567890-1234567")
	if idx := strings.Index(html, "Order ID:"); idx != -1 {
		start := idx + len("Order ID:")
		end := strings.Index(html[start:], "\n")
		if end == -1 {
			end = len(html[start:])
		}
		order.OrderID = strings.TrimSpace(html[start : start+end])
	}

	// Extract date
	if idx := strings.Index(html, "Date:"); idx != -1 {
		start := idx + len("Date:")
		end := strings.Index(html[start:], "\n")
		if end == -1 {
			end = len(html[start:])
		}
		order.Date = strings.TrimSpace(html[start : start+end])
	}

	// Extract total
	if idx := strings.Index(html, "Total:"); idx != -1 {
		start := idx + len("Total:")
		end := strings.Index(html[start:], "\n")
		if end == -1 {
			end = len(html[start:])
		}
		totalStr := strings.TrimSpace(html[start : start+end])
		totalStr = strings.TrimPrefix(totalStr, "$")
		totalStr = strings.ReplaceAll(totalStr, ",", "")
		if total, err := strconv.ParseFloat(totalStr, 64); err == nil {
			order.Total = total
		}
	}

	// Extract status
	if idx := strings.Index(html, "Status:"); idx != -1 {
		start := idx + len("Status:")
		end := strings.Index(html[start:], "\n")
		if end == -1 {
			end = len(html[start:])
		}
		order.Status = strings.TrimSpace(html[start : start+end])
	}

	// Validate required fields
	if order.OrderID == "" {
		return nil, fmt.Errorf("order ID not found in HTML")
	}

	return order, nil
}

// ValidateOrder validates that an order has all required fields
func ValidateOrder(order *models.Order) error {
	if order == nil {
		return fmt.Errorf("order is nil")
	}
	if order.OrderID == "" {
		return fmt.Errorf("order ID is required")
	}
	if !IsValidOrderID(order.OrderID) {
		return fmt.Errorf("invalid order ID format: %s", order.OrderID)
	}
	if order.Date == "" {
		return fmt.Errorf("order date is required")
	}
	if order.Total < 0 {
		return fmt.Errorf("order total cannot be negative")
	}
	if order.Status == "" {
		return fmt.Errorf("order status is required")
	}
	return nil
}

// IsValidOrderID checks if an order ID matches Amazon's format
// Amazon order IDs are typically in the format: XXX-XXXXXXX-XXXXXXX (19 characters)
func IsValidOrderID(orderID string) bool {
	parts := strings.Split(orderID, "-")
	if len(parts) != 3 {
		return false
	}
	if len(parts[0]) != 3 || len(parts[1]) != 7 || len(parts[2]) != 7 {
		return false
	}
	return true
}

// FilterOrdersByStatus filters orders by their status
func FilterOrdersByStatus(orders []models.Order, status string) []models.Order {
	if status == "" {
		return orders
	}

	filtered := make([]models.Order, 0)
	for _, order := range orders {
		if strings.EqualFold(order.Status, status) {
			filtered = append(filtered, order)
		}
	}
	return filtered
}

// CalculateOrderTotal calculates the total from order items
func CalculateOrderTotal(items []models.OrderItem) float64 {
	total := 0.0
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}
