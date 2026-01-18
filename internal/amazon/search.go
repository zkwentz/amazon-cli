package amazon

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

// SearchOptions contains parameters for product search
type SearchOptions struct {
	Category  string
	MinPrice  float64
	MaxPrice  float64
	PrimeOnly bool
	Page      int
}

// Client represents the Amazon API client
// This is a minimal definition for the search module
// The full client implementation will be in client.go
type Client struct {
	httpClient *http.Client
	userAgent  string
}

// NewClient creates a new Amazon client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
		userAgent:  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
	}
}

// Search performs a product search on Amazon
func (c *Client) Search(query string, opts SearchOptions) (*models.SearchResponse, error) {
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	// Default page to 1 if not set
	if opts.Page < 1 {
		opts.Page = 1
	}

	// Build search URL
	searchURL, err := c.buildSearchURL(query, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to build search URL: %w", err)
	}

	// Create request
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search request failed with status: %d", resp.StatusCode)
	}

	// Parse response
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Extract products from the page
	products, totalResults := c.parseSearchResults(doc)

	return &models.SearchResponse{
		Query:        query,
		Results:      products,
		TotalResults: totalResults,
		Page:         opts.Page,
	}, nil
}

// buildSearchURL constructs the Amazon search URL with all parameters
func (c *Client) buildSearchURL(query string, opts SearchOptions) (string, error) {
	baseURL := "https://www.amazon.com/s"
	params := url.Values{}

	// Add search query
	params.Add("k", query)

	// Add category if specified
	if opts.Category != "" {
		params.Add("i", opts.Category)
	}

	// Add price range filters
	if opts.MinPrice > 0 {
		params.Add("low-price", fmt.Sprintf("%.2f", opts.MinPrice))
	}
	if opts.MaxPrice > 0 {
		params.Add("high-price", fmt.Sprintf("%.2f", opts.MaxPrice))
	}

	// Add Prime filter
	if opts.PrimeOnly {
		params.Add("prime", "prime")
	}

	// Add page number (Amazon uses page parameter)
	if opts.Page > 1 {
		params.Add("page", strconv.Itoa(opts.Page))
	}

	return fmt.Sprintf("%s?%s", baseURL, params.Encode()), nil
}

// parseSearchResults extracts product information from the search results page
func (c *Client) parseSearchResults(doc *goquery.Document) ([]models.Product, int) {
	var products []models.Product

	// Amazon search results are typically in divs with data-component-type="s-search-result"
	doc.Find("[data-component-type='s-search-result']").Each(func(i int, s *goquery.Selection) {
		product := c.parseProductCard(s)
		if product.ASIN != "" {
			products = append(products, product)
		}
	})

	// Try to extract total results count
	totalResults := c.parseTotalResults(doc)

	return products, totalResults
}

// parseProductCard extracts product details from a single search result card
func (c *Client) parseProductCard(s *goquery.Selection) models.Product {
	product := models.Product{}

	// Extract ASIN
	if asin, exists := s.Attr("data-asin"); exists {
		product.ASIN = asin
	}

	// Extract title
	titleSel := s.Find("h2 a span")
	if titleSel.Length() > 0 {
		product.Title = strings.TrimSpace(titleSel.First().Text())
	}

	// Extract price
	priceSel := s.Find(".a-price .a-offscreen").First()
	if priceSel.Length() > 0 {
		priceText := strings.TrimSpace(priceSel.Text())
		product.Price = c.parsePrice(priceText)
	}

	// Extract original price (if on sale)
	originalPriceSel := s.Find(".a-price[data-a-strike='true'] .a-offscreen")
	if originalPriceSel.Length() > 0 {
		originalPriceText := strings.TrimSpace(originalPriceSel.Text())
		originalPrice := c.parsePrice(originalPriceText)
		if originalPrice > 0 {
			product.OriginalPrice = &originalPrice
		}
	}

	// Extract rating
	ratingSel := s.Find(".a-icon-alt")
	if ratingSel.Length() > 0 {
		ratingText := ratingSel.First().Text()
		product.Rating = c.parseRating(ratingText)
	}

	// Extract review count
	reviewCountSel := s.Find("[aria-label*='stars']")
	if reviewCountSel.Length() > 0 {
		if ariaLabel, exists := reviewCountSel.Attr("aria-label"); exists {
			product.ReviewCount = c.parseReviewCount(ariaLabel)
		}
	}

	// Check for Prime badge
	primeSel := s.Find("i[aria-label='Amazon Prime']")
	product.Prime = primeSel.Length() > 0

	// Assume in stock if price is present (simplified logic)
	product.InStock = product.Price > 0

	// Extract delivery estimate
	deliverySel := s.Find("[aria-label*='delivery'], .a-color-success")
	if deliverySel.Length() > 0 {
		product.DeliveryEstimate = strings.TrimSpace(deliverySel.First().Text())
	}

	return product
}

// parsePrice extracts a numeric price from a price string like "$29.99"
func (c *Client) parsePrice(priceStr string) float64 {
	// Remove currency symbols and commas
	priceStr = strings.TrimSpace(priceStr)
	priceStr = strings.ReplaceAll(priceStr, "$", "")
	priceStr = strings.ReplaceAll(priceStr, ",", "")

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return 0
	}
	return price
}

// parseRating extracts a numeric rating from a string like "4.5 out of 5 stars"
func (c *Client) parseRating(ratingStr string) float64 {
	parts := strings.Split(ratingStr, " ")
	if len(parts) > 0 {
		rating, err := strconv.ParseFloat(parts[0], 64)
		if err == nil {
			return rating
		}
	}
	return 0
}

// parseReviewCount extracts the number of reviews from an aria-label
func (c *Client) parseReviewCount(ariaLabel string) int {
	// Look for numbers in the aria label
	parts := strings.Fields(ariaLabel)
	for _, part := range parts {
		// Remove commas from numbers like "1,234"
		cleanPart := strings.ReplaceAll(part, ",", "")
		if count, err := strconv.Atoi(cleanPart); err == nil {
			return count
		}
	}
	return 0
}

// parseTotalResults attempts to extract the total number of results
func (c *Client) parseTotalResults(doc *goquery.Document) int {
	// Look for result count in the page header
	resultTextSel := doc.Find(".s-result-count, [data-component-type='s-result-info-bar']")
	if resultTextSel.Length() > 0 {
		text := resultTextSel.Text()
		// Try to extract a number from text like "1-48 of over 50,000 results"
		words := strings.Fields(text)
		for i, word := range words {
			if word == "of" && i+1 < len(words) {
				// Next word might be the total
				nextWord := strings.ReplaceAll(words[i+1], ",", "")
				nextWord = strings.ReplaceAll(nextWord, "over", "")
				nextWord = strings.TrimSpace(nextWord)
				if count, err := strconv.Atoi(nextWord); err == nil {
					return count
				}
			}
		}
	}
	return 0
}

// GetProduct fetches detailed information about a specific product by ASIN
func (c *Client) GetProduct(asin string) (*models.Product, error) {
	if asin == "" {
		return nil, fmt.Errorf("ASIN cannot be empty")
	}

	// Validate ASIN format (10 alphanumeric characters)
	if len(asin) != 10 {
		return nil, fmt.Errorf("invalid ASIN format: must be 10 characters")
	}

	// Build product URL
	productURL := fmt.Sprintf("https://www.amazon.com/dp/%s", asin)

	// Create request
	req, err := http.NewRequest("GET", productURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("product not found: %s", asin)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	// Parse response
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Extract product details
	product := c.parseProductDetail(doc, asin)

	return &product, nil
}

// parseProductDetail extracts detailed product information from the product page
func (c *Client) parseProductDetail(doc *goquery.Document, asin string) models.Product {
	product := models.Product{
		ASIN: asin,
	}

	// Extract title
	titleSel := doc.Find("#productTitle")
	if titleSel.Length() > 0 {
		product.Title = strings.TrimSpace(titleSel.Text())
	}

	// Extract price
	priceSel := doc.Find(".a-price .a-offscreen").First()
	if priceSel.Length() > 0 {
		priceText := strings.TrimSpace(priceSel.Text())
		product.Price = c.parsePrice(priceText)
	}

	// Extract rating
	ratingSel := doc.Find(".a-icon-star .a-icon-alt")
	if ratingSel.Length() > 0 {
		ratingText := ratingSel.Text()
		product.Rating = c.parseRating(ratingText)
	}

	// Extract review count
	reviewCountSel := doc.Find("#acrCustomerReviewText")
	if reviewCountSel.Length() > 0 {
		reviewText := reviewCountSel.Text()
		product.ReviewCount = c.parseReviewCountFromText(reviewText)
	}

	// Check Prime eligibility
	primeSel := doc.Find("#priceBadging_feature_div, [data-csa-c-type='element'][data-csa-c-content-id*='prime']")
	product.Prime = primeSel.Length() > 0 || strings.Contains(doc.Text(), "Prime FREE")

	// Check stock status
	availabilitySel := doc.Find("#availability")
	if availabilitySel.Length() > 0 {
		availText := strings.ToLower(availabilitySel.Text())
		product.InStock = strings.Contains(availText, "in stock") || strings.Contains(availText, "available")
	} else {
		product.InStock = product.Price > 0
	}

	// Extract description
	descSel := doc.Find("#feature-bullets, #productDescription")
	if descSel.Length() > 0 {
		product.Description = strings.TrimSpace(descSel.First().Text())
	}

	// Extract feature bullets
	doc.Find("#feature-bullets li").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" && !strings.Contains(text, "See more") {
			product.Features = append(product.Features, text)
		}
	})

	// Extract images
	doc.Find("#altImages img, #main-image-container img").Each(func(i int, s *goquery.Selection) {
		if src, exists := s.Attr("src"); exists {
			// Convert thumbnail to larger version if possible
			largeURL := strings.ReplaceAll(src, "._SS40_", "")
			largeURL = strings.ReplaceAll(largeURL, "._AC_US40_", "")
			product.Images = append(product.Images, largeURL)
		}
	})

	// Extract delivery estimate
	deliverySel := doc.Find("#mir-layout-DELIVERY_BLOCK")
	if deliverySel.Length() > 0 {
		product.DeliveryEstimate = strings.TrimSpace(deliverySel.Text())
	}

	return product
}

// parseReviewCountFromText extracts review count from text like "1,234 ratings"
func (c *Client) parseReviewCountFromText(text string) int {
	// Remove commas and extract the first number
	text = strings.TrimSpace(text)
	parts := strings.Fields(text)
	if len(parts) > 0 {
		numStr := strings.ReplaceAll(parts[0], ",", "")
		if count, err := strconv.Atoi(numStr); err == nil {
			return count
		}
	}
	return 0
}
