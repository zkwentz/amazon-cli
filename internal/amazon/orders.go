package amazon

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

// GetOrders retrieves a list of orders with optional filtering
func (c *Client) GetOrders(limit int, status string) (*models.OrdersResponse, error) {
	// TODO: Implement actual Amazon API call
	// For now, return mock data for testing

	if limit <= 0 {
		limit = 10
	}

	orders := []models.Order{
		{
			OrderID: "123-4567890-1234567",
			Date:    time.Now().AddDate(0, 0, -5).Format("2006-01-02"),
			Total:   29.99,
			Status:  "delivered",
			Items: []models.OrderItem{
				{
					ASIN:     "B08N5WRWNW",
					Title:    "Wireless Bluetooth Headphones",
					Quantity: 1,
					Price:    29.99,
				},
			},
			Tracking: &models.Tracking{
				Carrier:        "UPS",
				TrackingNumber: "1Z999AA10123456784",
				Status:         "delivered",
				DeliveryDate:   time.Now().AddDate(0, 0, -3).Format("2006-01-02"),
			},
		},
		{
			OrderID: "123-7654321-9876543",
			Date:    time.Now().AddDate(0, 0, -2).Format("2006-01-02"),
			Total:   54.99,
			Status:  "pending",
			Items: []models.OrderItem{
				{
					ASIN:     "B09XYZ1234",
					Title:    "USB-C Charging Cable 3-Pack",
					Quantity: 2,
					Price:    14.99,
				},
				{
					ASIN:     "B07ABC5678",
					Title:    "Phone Case - Clear",
					Quantity: 1,
					Price:    24.99,
				},
			},
			Tracking: &models.Tracking{
				Carrier:        "AMZL",
				TrackingNumber: "TBA123456789000",
				Status:         "in_transit",
			},
		},
	}

	// Filter by status if provided
	if status != "" {
		filtered := []models.Order{}
		for _, order := range orders {
			if order.Status == status {
				filtered = append(filtered, order)
			}
		}
		orders = filtered
	}

	// Apply limit
	if len(orders) > limit {
		orders = orders[:limit]
	}

	return &models.OrdersResponse{
		Orders:     orders,
		TotalCount: len(orders),
	}, nil
}

// GetOrder retrieves details for a specific order
func (c *Client) GetOrder(orderID string) (*models.Order, error) {
	if orderID == "" {
		return nil, fmt.Errorf("order ID cannot be empty")
	}

	// TODO: Implement actual Amazon API call
	// For now, return mock data

	return &models.Order{
		OrderID: orderID,
		Date:    time.Now().AddDate(0, 0, -5).Format("2006-01-02"),
		Total:   29.99,
		Status:  "delivered",
		Items: []models.OrderItem{
			{
				ASIN:     "B08N5WRWNW",
				Title:    "Wireless Bluetooth Headphones",
				Quantity: 1,
				Price:    29.99,
			},
		},
		Tracking: &models.Tracking{
			Carrier:        "UPS",
			TrackingNumber: "1Z999AA10123456784",
			Status:         "delivered",
			DeliveryDate:   time.Now().AddDate(0, 0, -3).Format("2006-01-02"),
		},
	}, nil
}

// GetOrderTracking retrieves tracking information for an order
func (c *Client) GetOrderTracking(orderID string) (*models.Tracking, error) {
	if orderID == "" {
		return nil, fmt.Errorf("order ID cannot be empty")
	}

	// TODO: Implement actual Amazon API call
	// For now, return mock data

	return &models.Tracking{
		Carrier:        "UPS",
		TrackingNumber: "1Z999AA10123456784",
		Status:         "in_transit",
		DeliveryDate:   time.Now().AddDate(0, 0, 2).Format("2006-01-02"),
		Events: []models.TrackingEvent{
			{
				Timestamp: time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
				Location:  "Local Distribution Center",
				Status:    "Out for delivery",
			},
			{
				Timestamp: time.Now().Add(-8 * time.Hour).Format(time.RFC3339),
				Location:  "Regional Facility",
				Status:    "In transit",
			},
			{
				Timestamp: time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
				Location:  "Origin Facility",
				Status:    "Shipped",
			},
		},
	}, nil
}

// GetOrderHistory retrieves order history for a specific year
func (c *Client) GetOrderHistory(year int) (*models.OrdersResponse, error) {
	if year <= 0 {
		year = time.Now().Year()
	}

	// TODO: Implement actual Amazon API call
	// For now, return mock data

	orders := []models.Order{
		{
			OrderID: "123-1111111-1111111",
			Date:    fmt.Sprintf("%d-06-15", year),
			Total:   149.99,
			Status:  "delivered",
			Items: []models.OrderItem{
				{
					ASIN:     "B08XYZ9876",
					Title:    "Kindle Paperwhite",
					Quantity: 1,
					Price:    149.99,
				},
			},
		},
		{
			OrderID: "123-2222222-2222222",
			Date:    fmt.Sprintf("%d-03-20", year),
			Total:   35.50,
			Status:  "delivered",
			Items: []models.OrderItem{
				{
					ASIN:     "B07DEF4567",
					Title:    "Book: The Go Programming Language",
					Quantity: 1,
					Price:    35.50,
				},
			},
		},
	}

	return &models.OrdersResponse{
		Orders:     orders,
		TotalCount: len(orders),
	}, nil
}

// parseOrdersHTML parses order list HTML and extracts order information
func parseOrdersHTML(html []byte) ([]models.Order, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var orders []models.Order

	// Find all order elements
	doc.Find(".order").Each(func(i int, s *goquery.Selection) {
		order := models.Order{}

		// Extract order ID from data attribute
		if orderID, exists := s.Attr("data-order-id"); exists {
			order.OrderID = orderID
		} else {
			// Try to extract from order number text
			orderNumText := s.Find(".order-number").Text()
			re := regexp.MustCompile(`\d{3}-\d{7}-\d{7}`)
			if match := re.FindString(orderNumText); match != "" {
				order.OrderID = match
			}
		}

		// Extract order date
		dateText := s.Find(".order-date").Text()
		order.Date = strings.TrimSpace(dateText)

		// Extract order total
		totalText := s.Find(".order-total").Text()
		order.Total = parsePrice(totalText)

		// Extract order status from delivery status text
		statusText := strings.ToLower(strings.TrimSpace(s.Find(".delivery-status").Text()))
		if strings.Contains(statusText, "delivered") {
			order.Status = "delivered"
		} else if strings.Contains(statusText, "arriving") || strings.Contains(statusText, "shipping") {
			order.Status = "pending"
		} else if strings.Contains(statusText, "cancelled") {
			order.Status = "cancelled"
		} else if strings.Contains(statusText, "returned") {
			order.Status = "returned"
		} else {
			order.Status = "unknown"
		}

		// Only add orders that have at least an order ID
		if order.OrderID != "" {
			orders = append(orders, order)
		}
	})

	return orders, nil
}

// parseOrderDetailHTML parses Amazon order detail HTML and extracts complete order information
func parseOrderDetailHTML(html []byte) (*models.Order, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	order := &models.Order{
		Items: []models.OrderItem{},
	}

	// Extract order ID
	orderIDText := doc.Find(".order-id-value").Text()
	if orderIDText == "" {
		// Try alternative selector
		doc.Find(".order-info").Each(func(i int, s *goquery.Selection) {
			if strings.Contains(s.Text(), "Order #") {
				orderIDText = strings.TrimSpace(strings.TrimPrefix(s.Text(), "Order #"))
			}
		})
	}
	order.OrderID = strings.TrimSpace(orderIDText)

	// Extract order date
	dateText := doc.Find(".order-date .value").Text()
	if dateText != "" {
		// Try to parse the date and convert to YYYY-MM-DD format
		parsedDate, err := time.Parse("January 2, 2006", strings.TrimSpace(dateText))
		if err == nil {
			order.Date = parsedDate.Format("2006-01-02")
		} else {
			order.Date = strings.TrimSpace(dateText)
		}
	}

	// Extract total
	totalText := doc.Find(".order-total .value").Text()
	if totalText != "" {
		order.Total = parsePrice(totalText)
	}

	// Extract status
	statusText := doc.Find(".order-status .status-badge").Text()
	if statusText != "" {
		order.Status = strings.ToLower(strings.TrimSpace(statusText))
	}

	// Extract order items
	doc.Find(".order-item").Each(func(i int, s *goquery.Selection) {
		item := models.OrderItem{}

		// Extract ASIN from data attribute or text
		if asin, exists := s.Attr("data-asin"); exists {
			item.ASIN = asin
		} else {
			// Try to extract from ASIN label/value
			s.Find(".item-asin .value").Each(func(i int, val *goquery.Selection) {
				item.ASIN = strings.TrimSpace(val.Text())
			})
		}

		// Extract title
		titleText := s.Find(".item-title").Text()
		if titleText == "" {
			titleText = s.Find(".item-title a").Text()
		}
		item.Title = strings.TrimSpace(titleText)

		// Extract price
		priceText := s.Find(".item-price .value").Text()
		if priceText != "" {
			item.Price = parsePrice(priceText)
		}

		// Extract quantity
		quantityText := s.Find(".item-quantity .value").Text()
		if quantityText != "" {
			quantity, err := strconv.Atoi(strings.TrimSpace(quantityText))
			if err == nil {
				item.Quantity = quantity
			}
		}

		// Only add item if we have at least ASIN and title
		if item.ASIN != "" && item.Title != "" {
			order.Items = append(order.Items, item)
		}
	})

	// Extract tracking information if present
	trackingSection := doc.Find(".tracking-section")
	if trackingSection.Length() > 0 {
		tracking := &models.Tracking{}

		carrier := trackingSection.Find(".tracking-carrier .value").Text()
		if carrier != "" {
			tracking.Carrier = strings.TrimSpace(carrier)
		}

		trackingNumber := trackingSection.Find(".tracking-number .value").Text()
		if trackingNumber != "" {
			tracking.TrackingNumber = strings.TrimSpace(trackingNumber)
		}

		status := trackingSection.Find(".tracking-status .value").Text()
		if status != "" {
			tracking.Status = strings.ToLower(strings.TrimSpace(status))
		}

		deliveryDate := trackingSection.Find(".delivery-date .value").Text()
		if deliveryDate != "" {
			// Try to parse the date and convert to YYYY-MM-DD format
			parsedDate, err := time.Parse("January 2, 2006", strings.TrimSpace(deliveryDate))
			if err == nil {
				tracking.DeliveryDate = parsedDate.Format("2006-01-02")
			} else {
				tracking.DeliveryDate = strings.TrimSpace(deliveryDate)
			}
		}

		// Only set tracking if we have at least a tracking number or carrier
		if tracking.TrackingNumber != "" || tracking.Carrier != "" {
			order.Tracking = tracking
		}
	}

	// Validate that we extracted essential information
	if order.OrderID == "" {
		return nil, fmt.Errorf("failed to extract order ID from HTML")
	}

	return order, nil
}

// parsePrice extracts a float64 price from a price string (e.g., "$29.99" -> 29.99)
func parsePrice(priceStr string) float64 {
	// Remove currency symbols and whitespace
	priceStr = strings.TrimSpace(priceStr)

	// Use regex to extract numeric value
	re := regexp.MustCompile(`[\d,]+\.?\d*`)
	match := re.FindString(priceStr)
	if match == "" {
		return 0.0
	}

	// Remove commas
	match = strings.ReplaceAll(match, ",", "")

	// Parse to float
	price, err := strconv.ParseFloat(match, 64)
	if err != nil {
		return 0.0
	}

	return price
}
