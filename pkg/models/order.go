package models

// Order represents an Amazon order
type Order struct {
	OrderID  string      `json:"order_id"`
	Date     string      `json:"date"`
	Total    float64     `json:"total"`
	Status   string      `json:"status"`
	Items    []OrderItem `json:"items"`
	Tracking *Tracking   `json:"tracking,omitempty"`
}

// OrderItem represents an item within an order
type OrderItem struct {
	ASIN     string  `json:"asin"`
	Title    string  `json:"title"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

// Tracking represents shipment tracking information
type Tracking struct {
	Carrier        string          `json:"carrier"`
	TrackingNumber string          `json:"tracking_number"`
	Status         string          `json:"status"`
	DeliveryDate   string          `json:"delivery_date,omitempty"`
	Events         []TrackingEvent `json:"events,omitempty"`
}

// TrackingEvent represents a single tracking event
type TrackingEvent struct {
	Timestamp string `json:"timestamp"`
	Location  string `json:"location"`
	Status    string `json:"status"`
}

// OrdersResponse represents the response for listing orders
type OrdersResponse struct {
	Orders     []Order `json:"orders"`
	TotalCount int     `json:"total_count"`
}
