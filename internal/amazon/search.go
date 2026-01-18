package amazon

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

// Search performs a product search on Amazon with the given query and options
func (c *Client) Search(query string, opts models.SearchOptions) (*models.SearchResponse, error) {
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	// Build the search URL with parameters
	searchURL, err := c.buildSearchURL(query, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to build search URL: %w", err)
	}

	// Create and configure the HTTP request
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers to mimic a real browser
	c.setRequestHeaders(req)

	// Execute the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search request failed with status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the HTML response
	searchResponse, err := c.parseSearchResults(string(body), query, opts.Page)
	if err != nil {
		return nil, fmt.Errorf("failed to parse search results: %w", err)
	}

	return searchResponse, nil
}

// buildSearchURL constructs the Amazon search URL with query parameters
func (c *Client) buildSearchURL(query string, opts models.SearchOptions) (string, error) {
	baseURL := c.baseURL + "/s"

	params := url.Values{}
	params.Set("k", query)

	// Add category filter if specified
	if opts.Category != "" {
		params.Set("i", opts.Category)
	}

	// Add price range filters if specified
	if opts.MinPrice > 0 {
		params.Set("low-price", fmt.Sprintf("%.2f", opts.MinPrice))
	}
	if opts.MaxPrice > 0 {
		params.Set("high-price", fmt.Sprintf("%.2f", opts.MaxPrice))
	}

	// Add Prime filter if specified
	if opts.PrimeOnly {
		params.Set("prime", "true")
	}

	// Add page parameter if specified (Amazon uses 'page' parameter)
	if opts.Page > 1 {
		params.Set("page", strconv.Itoa(opts.Page))
	}

	return baseURL + "?" + params.Encode(), nil
}

// setRequestHeaders sets common browser headers to avoid detection
func (c *Client) setRequestHeaders(req *http.Request) {
	// Select a random user agent
	userAgent := c.userAgents[rand.Intn(len(c.userAgents))]

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
}

// parseSearchResults extracts product information from the HTML response
func (c *Client) parseSearchResults(html string, query string, page int) (*models.SearchResponse, error) {
	response := &models.SearchResponse{
		Query:   query,
		Results: []models.Product{},
		Page:    page,
	}

	if page == 0 {
		response.Page = 1
	}

	// Extract total results count
	response.TotalResults = c.extractTotalResults(html)

	// Parse individual product items
	products := c.extractProducts(html)
	response.Results = products

	return response, nil
}

// extractTotalResults attempts to extract the total number of results from the page
func (c *Client) extractTotalResults(html string) int {
	// Look for result count patterns in Amazon's search results
	// Example: "1-48 of over 10,000 results"
	patterns := []string{
		`(\d+(?:,\d+)*)\s+results?`,
		`of\s+over\s+(\d+(?:,\d+)*)`,
		`of\s+(\d+(?:,\d+)*)\s+results?`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(html)
		if len(matches) > 1 {
			// Remove commas and parse the number
			numStr := strings.ReplaceAll(matches[1], ",", "")
			if num, err := strconv.Atoi(numStr); err == nil {
				return num
			}
		}
	}

	return 0
}

// extractProducts parses product data from the HTML
func (c *Client) extractProducts(html string) []models.Product {
	var products []models.Product

	// Amazon uses data-asin attributes for product identification
	// This regex finds product containers with ASIN
	asinPattern := regexp.MustCompile(`data-asin="([A-Z0-9]{10})"`)
	asinMatches := asinPattern.FindAllStringSubmatch(html, -1)

	// Keep track of unique ASINs
	seenASINs := make(map[string]bool)

	for _, match := range asinMatches {
		if len(match) < 2 {
			continue
		}

		asin := match[1]

		// Skip empty ASINs or duplicates
		if asin == "" || seenASINs[asin] {
			continue
		}
		seenASINs[asin] = true

		// Extract the product block for this ASIN
		product := c.extractProductByASIN(html, asin)
		if product != nil {
			products = append(products, *product)
		}

		// Limit to reasonable number of results
		if len(products) >= 50 {
			break
		}
	}

	return products
}

// extractProductByASIN extracts product details for a specific ASIN
func (c *Client) extractProductByASIN(html string, asin string) *models.Product {
	// Find the product block containing this ASIN
	// Look for data-asin="ASIN" and extract surrounding content
	startPattern := fmt.Sprintf(`data-asin="%s"`, asin)
	startIdx := strings.Index(html, startPattern)
	if startIdx == -1 {
		return nil
	}

	// Extract a reasonable chunk of HTML around this product
	// Search backwards for div opening and forwards for div closing
	chunkStart := startIdx - 500
	if chunkStart < 0 {
		chunkStart = 0
	}
	chunkEnd := startIdx + 3000
	if chunkEnd > len(html) {
		chunkEnd = len(html)
	}
	productHTML := html[chunkStart:chunkEnd]

	product := &models.Product{
		ASIN:    asin,
		InStock: true, // Assume in stock unless we find evidence otherwise
	}

	// Extract title
	product.Title = c.extractTitle(productHTML)

	// Extract price
	product.Price = c.extractPrice(productHTML)

	// Extract original price (for discounts)
	if originalPrice := c.extractOriginalPrice(productHTML); originalPrice > 0 {
		product.OriginalPrice = &originalPrice
	}

	// Extract rating
	product.Rating = c.extractRating(productHTML)

	// Extract review count
	product.ReviewCount = c.extractReviewCount(productHTML)

	// Check for Prime eligibility
	product.Prime = c.checkPrime(productHTML)

	// Extract delivery estimate
	product.DeliveryEstimate = c.extractDeliveryEstimate(productHTML)

	// Check stock status
	if strings.Contains(productHTML, "out of stock") || strings.Contains(productHTML, "Currently unavailable") {
		product.InStock = false
	}

	// Only return product if we extracted at least a title
	if product.Title == "" {
		return nil
	}

	return product
}

// extractTitle extracts the product title
func (c *Client) extractTitle(html string) string {
	patterns := []string{
		`<span class="[^"]*">([^<]+)</span>`,
		`<h2[^>]*>.*?<span[^>]*>([^<]+)</span>`,
		`aria-label="([^"]+)"`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(html)
		if len(matches) > 1 && len(matches[1]) > 10 {
			title := strings.TrimSpace(matches[1])
			// Clean up HTML entities
			title = strings.ReplaceAll(title, "&amp;", "&")
			title = strings.ReplaceAll(title, "&#39;", "'")
			title = strings.ReplaceAll(title, "&quot;", "\"")
			if title != "" {
				return title
			}
		}
	}

	return ""
}

// extractPrice extracts the product price
func (c *Client) extractPrice(html string) float64 {
	patterns := []string{
		`\$(\d+(?:\.\d{2})?)`,
		`<span[^>]*>(\d+(?:\.\d{2})?)</span>`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(html, -1)
		for _, match := range matches {
			if len(match) > 1 {
				if price, err := strconv.ParseFloat(match[1], 64); err == nil && price > 0 {
					return price
				}
			}
		}
	}

	return 0.0
}

// extractOriginalPrice extracts the original price (before discount)
func (c *Client) extractOriginalPrice(html string) float64 {
	// Look for strikethrough prices or "was" prices
	patterns := []string{
		`<span[^>]*text-decoration-line-through[^>]*>\$(\d+(?:\.\d{2})?)</span>`,
		`was[^$]*\$(\d+(?:\.\d{2})?)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(html)
		if len(matches) > 1 {
			if price, err := strconv.ParseFloat(matches[1], 64); err == nil && price > 0 {
				return price
			}
		}
	}

	return 0.0
}

// extractRating extracts the product rating
func (c *Client) extractRating(html string) float64 {
	patterns := []string{
		`(\d+(?:\.\d+)?)\s+out of 5`,
		`rating[^>]*>(\d+(?:\.\d+)?)</`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(html)
		if len(matches) > 1 {
			if rating, err := strconv.ParseFloat(matches[1], 64); err == nil {
				return rating
			}
		}
	}

	return 0.0
}

// extractReviewCount extracts the number of reviews
func (c *Client) extractReviewCount(html string) int {
	patterns := []string{
		`(\d+(?:,\d+)*)\s+ratings?`,
		`(\d+(?:,\d+)*)\s+reviews?`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(html)
		if len(matches) > 1 {
			numStr := strings.ReplaceAll(matches[1], ",", "")
			if count, err := strconv.Atoi(numStr); err == nil {
				return count
			}
		}
	}

	return 0
}

// checkPrime checks if the product is Prime eligible
func (c *Client) checkPrime(html string) bool {
	primeIndicators := []string{
		"prime",
		"Prime",
		"FREE delivery",
		"FREE Delivery",
	}

	htmlLower := strings.ToLower(html)
	for _, indicator := range primeIndicators {
		if strings.Contains(htmlLower, strings.ToLower(indicator)) {
			return true
		}
	}

	return false
}

// extractDeliveryEstimate extracts the delivery estimate
func (c *Client) extractDeliveryEstimate(html string) string {
	patterns := []string{
		`Get it by ([^<]+)`,
		`Delivery ([^<]+)`,
		`Arrives ([^<]+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(html)
		if len(matches) > 1 {
			estimate := strings.TrimSpace(matches[1])
			if estimate != "" {
				return estimate
			}
		}
	}

	return ""
}
