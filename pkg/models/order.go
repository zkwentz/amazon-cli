package models

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

// OrdersResponse represents the response from listing orders
type OrdersResponse struct {
	Orders     []Order `json:"orders"`
	TotalCount int     `json:"total_count"`
}
