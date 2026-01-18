package amazon

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Order represents an Amazon order with all its details
type Order struct {
	OrderID  string      `json:"order_id"`
	Date     string      `json:"date"`
	Total    float64     `json:"total"`
	Status   string      `json:"status"`
	Items    []OrderItem `json:"items"`
	Tracking *Tracking   `json:"tracking,omitempty"`
}

// OrderItem represents a single item in an order
type OrderItem struct {
	ASIN     string  `json:"asin"`
	Title    string  `json:"title"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

// Tracking represents shipment tracking information
type Tracking struct {
	Carrier        string `json:"carrier"`
	TrackingNumber string `json:"tracking_number"`
	Status         string `json:"status"`
	DeliveryDate   string `json:"delivery_date"`
}

// OrdersResponse represents the response for orders list API
type OrdersResponse struct {
	Orders     []Order `json:"orders"`
	TotalCount int     `json:"total_count"`
}

// GetOrders retrieves a list of orders with optional filtering
// limit: maximum number of orders to return
// status: filter by status (pending, delivered, returned, or empty for all)
func (c *Client) GetOrders(limit int, status string) (*OrdersResponse, error) {
	if c == nil {
		return nil, fmt.Errorf("client is nil")
	}

	// Build the order history URL
	orderURL := "https://www.amazon.com/gp/your-account/order-history"
	params := url.Values{}

	if status != "" {
		// Map status to Amazon's filter parameter
		statusMap := map[string]string{
			"pending":   "open",
			"delivered": "completed",
			"returned":  "returned",
		}
		if amazonStatus, ok := statusMap[status]; ok {
			params.Set("orderFilter", amazonStatus)
		}
	}

	if len(params) > 0 {
		orderURL += "?" + params.Encode()
	}

	// Make request to Amazon
	resp, err := c.Get(orderURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse HTML response
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	orders := []Order{}

	// Parse order cards from the page
	// Note: This is a simplified implementation. Real implementation would need
	// to handle Amazon's actual HTML structure which can vary
	doc.Find(".order").Each(func(i int, s *goquery.Selection) {
		if limit > 0 && len(orders) >= limit {
			return
		}

		order := Order{
			Items: []OrderItem{},
		}

		// Extract order ID
		if orderID, exists := s.Find(".order-id").First().Attr("data-order-id"); exists {
			order.OrderID = orderID
		} else {
			orderIDText := s.Find(".order-id").First().Text()
			order.OrderID = strings.TrimSpace(orderIDText)
		}

		// Extract order date
		dateText := s.Find(".order-date").First().Text()
		order.Date = strings.TrimSpace(dateText)

		// Extract total price
		totalText := s.Find(".order-total").First().Text()
		order.Total = parsePrice(totalText)

		// Extract status
		statusText := s.Find(".order-status").First().Text()
		order.Status = strings.ToLower(strings.TrimSpace(statusText))

		// Extract items
		s.Find(".order-item").Each(func(j int, item *goquery.Selection) {
			orderItem := OrderItem{}

			// Extract ASIN
			if asin, exists := item.Attr("data-asin"); exists {
				orderItem.ASIN = asin
			}

			// Extract title
			orderItem.Title = strings.TrimSpace(item.Find(".product-title").Text())

			// Extract quantity
			qtyText := item.Find(".item-quantity").Text()
			orderItem.Quantity = parseQuantity(qtyText)

			// Extract price
			priceText := item.Find(".item-price").Text()
			orderItem.Price = parsePrice(priceText)

			if orderItem.ASIN != "" || orderItem.Title != "" {
				order.Items = append(order.Items, orderItem)
			}
		})

		if order.OrderID != "" {
			orders = append(orders, order)
		}
	})

	return &OrdersResponse{
		Orders:     orders,
		TotalCount: len(orders),
	}, nil
}

// GetOrder retrieves detailed information about a specific order
func (c *Client) GetOrder(orderID string) (*Order, error) {
	if c == nil {
		return nil, fmt.Errorf("client is nil")
	}

	if orderID == "" {
		return nil, fmt.Errorf("order ID is required")
	}

	// Build order details URL
	orderURL := fmt.Sprintf("https://www.amazon.com/gp/your-account/order-details?orderID=%s", url.QueryEscape(orderID))

	// Make request
	resp, err := c.Get(orderURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order details: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("order not found: %s", orderID)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse HTML response
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	order := &Order{
		OrderID: orderID,
		Items:   []OrderItem{},
	}

	// Extract order details
	order.Date = strings.TrimSpace(doc.Find(".order-date-invoice-item").First().Text())
	order.Total = parsePrice(doc.Find(".grand-total-price").First().Text())
	order.Status = strings.ToLower(strings.TrimSpace(doc.Find(".delivery-box__primary-text").First().Text()))

	// Extract items
	doc.Find(".shipment-item").Each(func(i int, s *goquery.Selection) {
		item := OrderItem{}

		if asin, exists := s.Attr("data-asin"); exists {
			item.ASIN = asin
		}

		item.Title = strings.TrimSpace(s.Find(".product-title").Text())
		item.Quantity = parseQuantity(s.Find(".item-quantity").Text())
		item.Price = parsePrice(s.Find(".item-price").Text())

		if item.ASIN != "" || item.Title != "" {
			order.Items = append(order.Items, item)
		}
	})

	// Try to get tracking information if available
	tracking, err := c.GetOrderTracking(orderID)
	if err == nil && tracking != nil {
		order.Tracking = tracking
	}

	return order, nil
}

// GetOrderTracking retrieves tracking information for a specific order
func (c *Client) GetOrderTracking(orderID string) (*Tracking, error) {
	if c == nil {
		return nil, fmt.Errorf("client is nil")
	}

	if orderID == "" {
		return nil, fmt.Errorf("order ID is required")
	}

	// Build tracking URL
	trackingURL := fmt.Sprintf("https://www.amazon.com/gp/css/shiptrack/view.html?orderID=%s", url.QueryEscape(orderID))

	// Make request
	resp, err := c.Get(trackingURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tracking info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("tracking info not found for order: %s", orderID)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse HTML response
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	tracking := &Tracking{}

	// Extract tracking details
	tracking.Carrier = strings.TrimSpace(doc.Find(".carrier-info").First().Text())
	tracking.TrackingNumber = strings.TrimSpace(doc.Find(".tracking-number").First().Text())
	tracking.Status = strings.ToLower(strings.TrimSpace(doc.Find(".tracking-status").First().Text()))
	tracking.DeliveryDate = strings.TrimSpace(doc.Find(".delivery-date").First().Text())

	// If we didn't find any tracking info, return nil
	if tracking.TrackingNumber == "" && tracking.Carrier == "" {
		return nil, nil
	}

	return tracking, nil
}

// GetOrderHistory retrieves order history for a specific year
func (c *Client) GetOrderHistory(year int) (*OrdersResponse, error) {
	if c == nil {
		return nil, fmt.Errorf("client is nil")
	}

	// Default to current year if not specified
	if year == 0 {
		year = time.Now().Year()
	}

	// Validate year
	if year < 1995 || year > time.Now().Year() {
		return nil, fmt.Errorf("invalid year: %d", year)
	}

	// Build URL with year filter
	orderURL := fmt.Sprintf("https://www.amazon.com/gp/your-account/order-history?orderFilter=year-%d", year)

	// Make request
	resp, err := c.Get(orderURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order history: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse HTML response
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	orders := []Order{}

	// Parse all orders from the year
	doc.Find(".order").Each(func(i int, s *goquery.Selection) {
		order := Order{
			Items: []OrderItem{},
		}

		// Extract order ID
		if orderID, exists := s.Find(".order-id").First().Attr("data-order-id"); exists {
			order.OrderID = orderID
		} else {
			orderIDText := s.Find(".order-id").First().Text()
			order.OrderID = strings.TrimSpace(orderIDText)
		}

		// Extract order date
		dateText := s.Find(".order-date").First().Text()
		order.Date = strings.TrimSpace(dateText)

		// Extract total
		totalText := s.Find(".order-total").First().Text()
		order.Total = parsePrice(totalText)

		// Extract status
		statusText := s.Find(".order-status").First().Text()
		order.Status = strings.ToLower(strings.TrimSpace(statusText))

		// Extract items (basic info only for list view)
		s.Find(".order-item").Each(func(j int, item *goquery.Selection) {
			orderItem := OrderItem{}

			if asin, exists := item.Attr("data-asin"); exists {
				orderItem.ASIN = asin
			}

			orderItem.Title = strings.TrimSpace(item.Find(".product-title").Text())
			orderItem.Quantity = parseQuantity(item.Find(".item-quantity").Text())
			orderItem.Price = parsePrice(item.Find(".item-price").Text())

			if orderItem.ASIN != "" || orderItem.Title != "" {
				order.Items = append(order.Items, orderItem)
			}
		})

		if order.OrderID != "" {
			orders = append(orders, order)
		}
	})

	// Handle pagination if necessary
	// This is a simplified version - real implementation would need to follow
	// pagination links to get all orders from the year

	return &OrdersResponse{
		Orders:     orders,
		TotalCount: len(orders),
	}, nil
}

// Helper function to parse price strings (e.g., "$29.99" -> 29.99)
func parsePrice(priceStr string) float64 {
	// Remove common currency symbols and whitespace
	cleaned := strings.TrimSpace(priceStr)
	cleaned = strings.ReplaceAll(cleaned, "$", "")
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	cleaned = strings.TrimSpace(cleaned)

	// Parse to float
	price, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0.0
	}

	return price
}

// Helper function to parse quantity strings (e.g., "Qty: 2" -> 2)
func parseQuantity(qtyStr string) int {
	// Remove common quantity prefixes
	cleaned := strings.TrimSpace(qtyStr)
	cleaned = strings.ReplaceAll(cleaned, "Qty:", "")
	cleaned = strings.ReplaceAll(cleaned, "Quantity:", "")
	cleaned = strings.TrimSpace(cleaned)

	// Parse to int
	qty, err := strconv.Atoi(cleaned)
	if err != nil {
		return 1 // Default to 1 if parsing fails
	}

	return qty
}

// Client placeholder - this should be defined in client.go
// Including a minimal version here for compilation
type Client struct {
	httpClient *http.Client
}

// Get performs a GET request (placeholder - should be in client.go)
func (c *Client) Get(url string) (*http.Response, error) {
	if c.httpClient == nil {
		c.httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}
	return c.httpClient.Get(url)
}

// MarshalJSON provides custom JSON serialization for better API output
func (o *Order) MarshalJSON() ([]byte, error) {
	type Alias Order
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	})
}

// UnmarshalJSON provides custom JSON deserialization
func (o *Order) UnmarshalJSON(data []byte) error {
	type Alias Order
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	}
	return json.Unmarshal(data, &aux)
}
