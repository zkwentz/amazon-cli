package models

// CartItem represents an item in the shopping cart
type CartItem struct {
	ASIN     string  `json:"asin"`
	Title    string  `json:"title"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Subtotal float64 `json:"subtotal"`
	Prime    bool    `json:"prime"`
	InStock  bool    `json:"in_stock"`
}

// Cart represents the shopping cart with all items and totals
type Cart struct {
	Items        []CartItem `json:"items"`
	Subtotal     float64    `json:"subtotal"`
	EstimatedTax float64    `json:"estimated_tax"`
	Total        float64    `json:"total"`
	ItemCount    int        `json:"item_count"`
}

// Address represents a shipping address
type Address struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
	Country string `json:"country"`
	Default bool   `json:"default"`
}

// PaymentMethod represents a payment method
type PaymentMethod struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Last4   string `json:"last4"`
	Default bool   `json:"default"`
}

// DeliveryOption represents a delivery option during checkout
type DeliveryOption struct {
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Date        string  `json:"date"`
}

// CheckoutPreview represents a preview of the checkout
type CheckoutPreview struct {
	Cart            *Cart            `json:"cart"`
	Address         *Address         `json:"address"`
	PaymentMethod   *PaymentMethod   `json:"payment_method"`
	DeliveryOptions []DeliveryOption `json:"delivery_options"`
}

// OrderConfirmation represents a confirmed order
type OrderConfirmation struct {
	OrderID           string  `json:"order_id"`
	Total             float64 `json:"total"`
	EstimatedDelivery string  `json:"estimated_delivery"`
}
