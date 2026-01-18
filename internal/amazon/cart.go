package amazon

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

// CartService handles cart-related operations
type CartService struct {
	client *Client
}

// NewCartService creates a new cart service
func NewCartService(client *Client) *CartService {
	return &CartService{
		client: client,
	}
}

// AddToCart adds an item to the shopping cart
// Returns the updated cart after adding the item
func (s *CartService) AddToCart(asin string, quantity int) (*models.Cart, error) {
	if asin == "" {
		return nil, fmt.Errorf("ASIN cannot be empty")
	}
	if quantity <= 0 {
		return nil, fmt.Errorf("quantity must be greater than 0")
	}

	// Prepare form data for add-to-cart request
	formData := url.Values{}
	formData.Set("ASIN", asin)
	formData.Set("quantity", fmt.Sprintf("%d", quantity))

	// Build the add-to-cart URL
	addToCartURL := "https://www.amazon.com/gp/aws/cart/add.html"

	req, err := http.NewRequest("POST", addToCartURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create add-to-cart request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to add item to cart: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusFound {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("add to cart failed with status %d: %s", resp.StatusCode, string(body))
	}

	// After adding, fetch the updated cart
	return s.GetCart()
}

// GetCart retrieves the current shopping cart contents
func (s *CartService) GetCart() (*models.Cart, error) {
	cartURL := "https://www.amazon.com/gp/cart/view.html"

	req, err := http.NewRequest("GET", cartURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cart request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get cart failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the cart page HTML
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read cart response: %w", err)
	}

	cart, err := s.parseCartHTML(string(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse cart: %w", err)
	}

	return cart, nil
}

// RemoveFromCart removes an item from the shopping cart
func (s *CartService) RemoveFromCart(asin string) (*models.Cart, error) {
	if asin == "" {
		return nil, fmt.Errorf("ASIN cannot be empty")
	}

	// First, get the cart to find the item
	cart, err := s.GetCart()
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Check if item exists in cart
	found := false
	for _, item := range cart.Items {
		if item.ASIN == asin {
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("item %s not found in cart", asin)
	}

	// Prepare form data for delete request
	formData := url.Values{}
	formData.Set("ASIN", asin)
	formData.Set("quantity", "0") // Setting quantity to 0 removes the item

	deleteURL := "https://www.amazon.com/gp/cart/ajax-update.html"

	req, err := http.NewRequest("POST", deleteURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create remove request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to remove item from cart: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("remove from cart failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Get updated cart
	return s.GetCart()
}

// ClearCart removes all items from the shopping cart
func (s *CartService) ClearCart() error {
	cart, err := s.GetCart()
	if err != nil {
		return fmt.Errorf("failed to get cart: %w", err)
	}

	// Remove each item
	for _, item := range cart.Items {
		_, err := s.RemoveFromCart(item.ASIN)
		if err != nil {
			return fmt.Errorf("failed to remove item %s: %w", item.ASIN, err)
		}
	}

	return nil
}

// GetAddresses retrieves all saved shipping addresses
func (s *CartService) GetAddresses() ([]models.Address, error) {
	addressURL := "https://www.amazon.com/a/addresses"

	req, err := http.NewRequest("GET", addressURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create addresses request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get addresses: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get addresses failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read addresses response: %w", err)
	}

	addresses, err := s.parseAddressesHTML(string(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse addresses: %w", err)
	}

	return addresses, nil
}

// GetPaymentMethods retrieves all saved payment methods
func (s *CartService) GetPaymentMethods() ([]models.PaymentMethod, error) {
	paymentURL := "https://www.amazon.com/cpe/yourpayments/wallet"

	req, err := http.NewRequest("GET", paymentURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment methods request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment methods: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get payment methods failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read payment methods response: %w", err)
	}

	methods, err := s.parsePaymentMethodsHTML(string(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse payment methods: %w", err)
	}

	return methods, nil
}

// PreviewCheckout initiates checkout and returns a preview without completing the purchase
func (s *CartService) PreviewCheckout(addressID, paymentID string) (*models.CheckoutPreview, error) {
	checkoutURL := "https://www.amazon.com/gp/buy/spc/handlers/display.html"

	req, err := http.NewRequest("GET", checkoutURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create checkout preview request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get checkout preview: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("checkout preview failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read checkout preview response: %w", err)
	}

	preview, err := s.parseCheckoutPreviewHTML(string(body), addressID, paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse checkout preview: %w", err)
	}

	return preview, nil
}

// CompleteCheckout completes the checkout process and places the order
func (s *CartService) CompleteCheckout(addressID, paymentID string) (*models.OrderConfirmation, error) {
	// First, preview the checkout to ensure everything is valid
	_, err := s.PreviewCheckout(addressID, paymentID)
	if err != nil {
		return nil, fmt.Errorf("checkout preview failed: %w", err)
	}

	// Submit the purchase
	formData := url.Values{}
	if addressID != "" {
		formData.Set("addressID", addressID)
	}
	if paymentID != "" {
		formData.Set("paymentMethodID", paymentID)
	}

	purchaseURL := "https://www.amazon.com/gp/buy/spc/handlers/static-submit-decoupled.html"

	req, err := http.NewRequest("POST", purchaseURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create checkout request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to complete checkout: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusFound {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("checkout failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read checkout response: %w", err)
	}

	confirmation, err := s.parseOrderConfirmation(string(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse order confirmation: %w", err)
	}

	return confirmation, nil
}

// parseCartHTML parses the cart HTML and extracts cart data
// This is a placeholder implementation that would need to be completed with actual HTML parsing
func (s *CartService) parseCartHTML(html string) (*models.Cart, error) {
	// TODO: Implement actual HTML parsing using goquery or similar
	// This is a stub implementation for now
	cart := &models.Cart{
		Items:        []models.CartItem{},
		Subtotal:     0.0,
		EstimatedTax: 0.0,
		Total:        0.0,
		ItemCount:    0,
	}

	// In a real implementation, this would parse the HTML to extract:
	// - Individual cart items with ASIN, title, price, quantity
	// - Subtotal, tax, and total amounts
	// - Stock status and Prime eligibility for each item

	return cart, nil
}

// parseAddressesHTML parses the addresses HTML page
func (s *CartService) parseAddressesHTML(html string) ([]models.Address, error) {
	// TODO: Implement actual HTML parsing
	// This is a stub implementation
	addresses := []models.Address{}

	// In a real implementation, this would parse the HTML to extract:
	// - Address ID, name, street, city, state, zip, country
	// - Default address flag

	return addresses, nil
}

// parsePaymentMethodsHTML parses the payment methods HTML page
func (s *CartService) parsePaymentMethodsHTML(html string) ([]models.PaymentMethod, error) {
	// TODO: Implement actual HTML parsing
	// This is a stub implementation
	methods := []models.PaymentMethod{}

	// In a real implementation, this would parse the HTML to extract:
	// - Payment method ID, type (Visa, Mastercard, etc.)
	// - Last 4 digits of card number
	// - Default payment method flag

	return methods, nil
}

// parseCheckoutPreviewHTML parses the checkout preview page
func (s *CartService) parseCheckoutPreviewHTML(html string, addressID, paymentID string) (*models.CheckoutPreview, error) {
	// TODO: Implement actual HTML parsing
	// This is a stub implementation
	preview := &models.CheckoutPreview{
		Cart:            &models.Cart{},
		Address:         &models.Address{},
		PaymentMethod:   &models.PaymentMethod{},
		DeliveryOptions: []models.DeliveryOption{},
	}

	// In a real implementation, this would parse the HTML to extract:
	// - Cart summary with totals
	// - Selected or default address
	// - Selected or default payment method
	// - Available delivery options with dates and prices

	return preview, nil
}

// parseOrderConfirmation parses the order confirmation page
func (s *CartService) parseOrderConfirmation(html string) (*models.OrderConfirmation, error) {
	// TODO: Implement actual HTML parsing
	// This is a stub implementation
	confirmation := &models.OrderConfirmation{
		OrderID:           "",
		Total:             0.0,
		EstimatedDelivery: "",
	}

	// In a real implementation, this would parse the HTML to extract:
	// - Order ID
	// - Total amount charged
	// - Estimated delivery date

	return confirmation, nil
}

// Client is a placeholder for the HTTP client
// In the real implementation, this would be defined in client.go
type Client struct {
	httpClient *http.Client
}

// Do performs an HTTP request with rate limiting and retry logic
// This is a placeholder - the real implementation would be in client.go
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if c.httpClient == nil {
		c.httpClient = &http.Client{}
	}
	return c.httpClient.Do(req)
}

// MarshalJSON custom marshaling for better JSON output
func (c *Cart) MarshalJSON() ([]byte, error) {
	type Alias models.Cart
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	})
}
