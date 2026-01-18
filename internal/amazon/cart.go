package amazon

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

const (
	amazonCartURL = "https://www.amazon.com/gp/cart/view.html"
)

// Client represents an Amazon API client
type Client struct {
	httpClient *http.Client
	// Additional fields for authentication, rate limiting, etc. will be added later
}

// NewClient creates a new Amazon client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
	}
}

// GetCart fetches the current shopping cart from Amazon
func (c *Client) GetCart() (*models.Cart, error) {
	// Create request to cart page
	req, err := http.NewRequest("GET", amazonCartURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cart request: %w", err)
	}

	// Set common headers to mimic browser request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cart: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cart request failed with status: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read cart response: %w", err)
	}

	// Parse cart from HTML response
	cart, err := parseCartHTML(string(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse cart: %w", err)
	}

	return cart, nil
}

// parseCartHTML parses the HTML response from Amazon cart page
func parseCartHTML(html string) (*models.Cart, error) {
	cart := &models.Cart{
		Items: []models.CartItem{},
	}

	// Check if cart is empty
	if strings.Contains(html, "Your Shopping Cart is empty") ||
	   strings.Contains(html, "Your Amazon Cart is empty") {
		return cart, nil
	}

	// Parse cart items
	// Note: This is a simplified parser. In production, you would use a proper HTML parser
	// like goquery (github.com/PuerkitoBio/goquery)
	items, err := extractCartItems(html)
	if err != nil {
		return nil, err
	}
	cart.Items = items

	// Calculate item count
	cart.ItemCount = len(items)

	// Extract totals
	subtotal, err := extractPrice(html, `subtotal`)
	if err == nil {
		cart.Subtotal = subtotal
	}

	estimatedTax, err := extractPrice(html, `estimated tax`)
	if err == nil {
		cart.EstimatedTax = estimatedTax
	}

	// Calculate total
	cart.Total = cart.Subtotal + cart.EstimatedTax

	return cart, nil
}

// extractCartItems extracts cart items from HTML
func extractCartItems(html string) ([]models.CartItem, error) {
	var items []models.CartItem

	// Pattern to find ASIN in cart items
	asinPattern := regexp.MustCompile(`data-asin="([A-Z0-9]{10})"`)
	asinMatches := asinPattern.FindAllStringSubmatch(html, -1)

	// For each ASIN found, extract item details
	for _, match := range asinMatches {
		if len(match) < 2 {
			continue
		}
		asin := match[1]

		// Extract item details (simplified - in production use proper HTML parsing)
		item := models.CartItem{
			ASIN:    asin,
			InStock: true,  // Default assumption
			Prime:   false, // Default assumption
		}

		// Try to extract title (looking for sc-product-title or similar classes)
		titlePattern := regexp.MustCompile(`data-asin="` + asin + `"[^>]*>[\s\S]*?<span class="[^"]*product-title[^"]*"[^>]*>(.*?)</span>`)
		if titleMatch := titlePattern.FindStringSubmatch(html); len(titleMatch) > 1 {
			item.Title = strings.TrimSpace(stripHTMLTags(titleMatch[1]))
		}

		// Try to extract price
		pricePattern := regexp.MustCompile(`data-asin="` + asin + `"[^>]*>[\s\S]*?\$([0-9]+\.[0-9]{2})`)
		if priceMatch := pricePattern.FindStringSubmatch(html); len(priceMatch) > 1 {
			if price, err := strconv.ParseFloat(priceMatch[1], 64); err == nil {
				item.Price = price
			}
		}

		// Try to extract quantity
		quantityPattern := regexp.MustCompile(`data-asin="` + asin + `"[^>]*>[\s\S]*?quantity["\s:]+([0-9]+)`)
		if qtyMatch := quantityPattern.FindStringSubmatch(html); len(qtyMatch) > 1 {
			if qty, err := strconv.Atoi(qtyMatch[1]); err == nil {
				item.Quantity = qty
			} else {
				item.Quantity = 1 // Default to 1
			}
		} else {
			item.Quantity = 1 // Default to 1
		}

		// Calculate subtotal
		item.Subtotal = item.Price * float64(item.Quantity)

		// Check for Prime
		if strings.Contains(html, asin) && strings.Contains(html, "prime-badge") {
			item.Prime = true
		}

		items = append(items, item)
	}

	return items, nil
}

// extractPrice extracts a price value from HTML based on a label
func extractPrice(html, label string) (float64, error) {
	pattern := regexp.MustCompile(`(?i)` + label + `[^$]*\$([0-9,]+\.[0-9]{2})`)
	matches := pattern.FindStringSubmatch(html)
	if len(matches) < 2 {
		return 0, fmt.Errorf("price not found for label: %s", label)
	}

	// Remove commas from price string
	priceStr := strings.ReplaceAll(matches[1], ",", "")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse price: %w", err)
	}

	return price, nil
}

// stripHTMLTags removes HTML tags from a string
func stripHTMLTags(s string) string {
	tagPattern := regexp.MustCompile(`<[^>]*>`)
	return tagPattern.ReplaceAllString(s, "")
}

// AddToCart adds an item to the cart
func (c *Client) AddToCart(asin string, quantity int) (*models.Cart, error) {
	// TODO: Implement add to cart functionality
	return nil, fmt.Errorf("not implemented")
}

// RemoveFromCart removes an item from the cart
func (c *Client) RemoveFromCart(asin string) (*models.Cart, error) {
	// TODO: Implement remove from cart functionality
	return nil, fmt.Errorf("not implemented")
}

// ClearCart removes all items from the cart
func (c *Client) ClearCart() error {
	// TODO: Implement clear cart functionality
	return fmt.Errorf("not implemented")
}

// GetAddresses fetches saved addresses
func (c *Client) GetAddresses() ([]models.Address, error) {
	// TODO: Implement get addresses functionality
	return nil, fmt.Errorf("not implemented")
}

// GetPaymentMethods fetches saved payment methods
func (c *Client) GetPaymentMethods() ([]models.PaymentMethod, error) {
	// TODO: Implement get payment methods functionality
	return nil, fmt.Errorf("not implemented")
}
