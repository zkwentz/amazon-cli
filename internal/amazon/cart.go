package amazon

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

// AddToCart adds an item to the shopping cart
func (c *Client) AddToCart(asin string, quantity int) (*models.Cart, error) {
	// TODO: Implement add to cart functionality
	// This would make a POST request to Amazon's add-to-cart endpoint
	return nil, fmt.Errorf("AddToCart not yet implemented")
}

// GetCart retrieves the current shopping cart contents
func (c *Client) GetCart() (*models.Cart, error) {
	// TODO: Implement get cart functionality
	// This would fetch the cart page and parse the items
	return nil, fmt.Errorf("GetCart not yet implemented")
}

// RemoveFromCart removes an item from the shopping cart
func (c *Client) RemoveFromCart(asin string) (*models.Cart, error) {
	// TODO: Implement remove from cart functionality
	return nil, fmt.Errorf("RemoveFromCart not yet implemented")
}

// ClearCart removes all items from the shopping cart
func (c *Client) ClearCart() error {
	// TODO: Implement clear cart functionality
	return fmt.Errorf("ClearCart not yet implemented")
}

// GetAddresses retrieves all saved shipping addresses
func (c *Client) GetAddresses() ([]models.Address, error) {
	// TODO: Implement get addresses functionality
	// This would fetch the addresses from the user's account
	return nil, fmt.Errorf("GetAddresses not yet implemented")
}

// GetPaymentMethods retrieves all saved payment methods
func (c *Client) GetPaymentMethods() ([]models.PaymentMethod, error) {
	// TODO: Implement get payment methods functionality
	// This would fetch payment methods from the user's account
	return nil, fmt.Errorf("GetPaymentMethods not yet implemented")
}

// PreviewCheckout initiates the checkout flow and returns a preview without completing the purchase
// It takes addressID and paymentID as parameters to specify which address and payment method to use
func (c *Client) PreviewCheckout(addressID, paymentID string) (*models.CheckoutPreview, error) {
	// Step 1: Get current cart contents
	cart, err := c.getCartForCheckout()
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Validate cart is not empty
	if len(cart.Items) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	// Step 2: Get the specified address or default address
	address, err := c.getAddressForCheckout(addressID)
	if err != nil {
		return nil, fmt.Errorf("failed to get address: %w", err)
	}

	// Step 3: Get the specified payment method or default payment method
	paymentMethod, err := c.getPaymentMethodForCheckout(paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment method: %w", err)
	}

	// Step 4: Initiate checkout flow to get preview information
	preview, err := c.fetchCheckoutPreview(cart, address, paymentMethod)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch checkout preview: %w", err)
	}

	return preview, nil
}

// getCartForCheckout retrieves the cart for checkout purposes
func (c *Client) getCartForCheckout() (*models.Cart, error) {
	// This would make a request to Amazon's cart page
	// For now, returning a mock implementation structure
	// TODO: Implement actual cart fetching logic with HTML parsing

	req, err := http.NewRequest("GET", "https://www.amazon.com/gp/cart/view.html", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cart request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cart: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code fetching cart: %d", resp.StatusCode)
	}

	// TODO: Parse cart HTML/JSON response
	// This is a placeholder that would be replaced with actual parsing logic
	return nil, fmt.Errorf("cart parsing not yet implemented")
}

// getAddressForCheckout retrieves the address to use for checkout
func (c *Client) getAddressForCheckout(addressID string) (*models.Address, error) {
	// Fetch all addresses
	addresses, err := c.GetAddresses()
	if err != nil {
		return nil, err
	}

	if len(addresses) == 0 {
		return nil, fmt.Errorf("no addresses found on account")
	}

	// If addressID is specified, find and return that specific address
	if addressID != "" {
		for _, addr := range addresses {
			if addr.ID == addressID {
				return &addr, nil
			}
		}
		return nil, fmt.Errorf("address with ID %s not found", addressID)
	}

	// Otherwise, return the default address
	for _, addr := range addresses {
		if addr.Default {
			return &addr, nil
		}
	}

	// If no default is set, return the first address
	return &addresses[0], nil
}

// getPaymentMethodForCheckout retrieves the payment method to use for checkout
func (c *Client) getPaymentMethodForCheckout(paymentID string) (*models.PaymentMethod, error) {
	// Fetch all payment methods
	paymentMethods, err := c.GetPaymentMethods()
	if err != nil {
		return nil, err
	}

	if len(paymentMethods) == 0 {
		return nil, fmt.Errorf("no payment methods found on account")
	}

	// If paymentID is specified, find and return that specific payment method
	if paymentID != "" {
		for _, pm := range paymentMethods {
			if pm.ID == paymentID {
				return &pm, nil
			}
		}
		return nil, fmt.Errorf("payment method with ID %s not found", paymentID)
	}

	// Otherwise, return the default payment method
	for _, pm := range paymentMethods {
		if pm.Default {
			return &pm, nil
		}
	}

	// If no default is set, return the first payment method
	return &paymentMethods[0], nil
}

// fetchCheckoutPreview makes the actual request to Amazon to get checkout preview information
func (c *Client) fetchCheckoutPreview(cart *models.Cart, address *models.Address, paymentMethod *models.PaymentMethod) (*models.CheckoutPreview, error) {
	// Build form data for checkout preview request
	formData := url.Values{}
	formData.Set("addressID", address.ID)
	formData.Set("paymentMethodID", paymentMethod.ID)

	// Make request to Amazon's checkout preview endpoint
	// This would typically be something like /gp/buy/spc/handlers/display.html
	req, err := http.NewRequest("POST", "https://www.amazon.com/gp/buy/spc/handlers/display.html", strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create checkout preview request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch checkout preview: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code fetching checkout preview: %d", resp.StatusCode)
	}

	// TODO: Parse the response to extract delivery options, tax, shipping, and totals
	// This would involve parsing HTML or JSON response from Amazon

	// For now, construct a basic preview structure
	preview := &models.CheckoutPreview{
		Cart:          cart,
		Address:       address,
		PaymentMethod: paymentMethod,
		DeliveryOptions: []models.DeliveryOption{
			{
				Method:        "Standard Shipping",
				EstimatedDate: "Jan 25-27",
				Cost:          0.0,
				BusinessDays:  5,
			},
			{
				Method:         "Two-Day Shipping",
				EstimatedDate:  "Jan 22",
				Cost:           9.99,
				BusinessDays:   2,
				GuaranteedDate: "Jan 22",
			},
		},
		Subtotal: cart.Subtotal,
		Tax:      calculateEstimatedTax(cart.Subtotal, address.State),
		Shipping: 0.0, // Assuming Prime/free shipping for standard
		Total:    0.0, // Will be calculated below
	}

	// Calculate total
	preview.Total = preview.Subtotal + preview.Tax + preview.Shipping

	return preview, nil
}

// calculateEstimatedTax calculates estimated tax based on subtotal and state
func calculateEstimatedTax(subtotal float64, state string) float64 {
	// Tax rates by state (simplified - actual implementation would use precise rates)
	taxRates := map[string]float64{
		"CA": 0.0725, // California
		"NY": 0.08,   // New York
		"TX": 0.0625, // Texas
		"FL": 0.06,   // Florida
		"WA": 0.065,  // Washington
		// Add more states as needed
	}

	rate, ok := taxRates[state]
	if !ok {
		rate = 0.07 // Default average tax rate
	}

	tax := subtotal * rate
	// Round to 2 decimal places
	return float64(int(tax*100+0.5)) / 100
}

// CompleteCheckout completes the checkout process and places the order
func (c *Client) CompleteCheckout(addressID, paymentID string) (*models.OrderConfirmation, error) {
	// TODO: Implement complete checkout functionality
	// This would finalize the purchase and return order confirmation
	return nil, fmt.Errorf("CompleteCheckout not yet implemented")
}

// parseFloat parses a price string (e.g., "$29.99") to a float64
func parseFloat(s string) (float64, error) {
	// Remove currency symbols and commas
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "$")
	s = strings.ReplaceAll(s, ",", "")

	return strconv.ParseFloat(s, 64)
}
