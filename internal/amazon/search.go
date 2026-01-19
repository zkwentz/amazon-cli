package amazon

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

// Search searches for products on Amazon
func (c *Client) Search(query string, opts models.SearchOptions) (*models.SearchResponse, error) {
	// Validate query
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	// Set default page if not provided
	if opts.Page <= 0 {
		opts.Page = 1
	}

	// Build search URL with query parameters
	searchURL := fmt.Sprintf("%s/s", c.baseURL)
	params := url.Values{}

	// Add search query
	params.Add("k", query)

	// Add category filter if provided
	if opts.Category != "" {
		params.Add("i", opts.Category)
	}

	// Add price range filters
	if opts.MinPrice > 0 {
		// Amazon uses price in cents for min price
		minPriceCents := int(opts.MinPrice * 100)
		params.Add("low-price", strconv.Itoa(minPriceCents))
	}
	if opts.MaxPrice > 0 {
		// Amazon uses price in cents for max price
		maxPriceCents := int(opts.MaxPrice * 100)
		params.Add("high-price", strconv.Itoa(maxPriceCents))
	}

	// Add Prime filter if requested
	if opts.PrimeOnly {
		params.Add("prime", "true")
	}

	// Add page number (Amazon uses 'page' parameter)
	if opts.Page > 1 {
		params.Add("page", strconv.Itoa(opts.Page))
	}

	// Construct full URL with query parameters
	fullURL := fmt.Sprintf("%s?%s", searchURL, params.Encode())

	// Create HTTP GET request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Execute the request with rate limiting and retries
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch search results: %w", err)
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for CAPTCHA
	if c.detectCAPTCHA(body.Bytes()) {
		return nil, fmt.Errorf("CAPTCHA detected - please try again later or use a different method")
	}

	// Parse the HTML response
	products, err := parseSearchResultsHTML(body.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to parse search results: %w", err)
	}

	return &models.SearchResponse{
		Query:        query,
		Results:      products,
		TotalResults: len(products),
		Page:         opts.Page,
	}, nil
}

// parseSearchResultsHTML parses Amazon search results HTML and extracts product information
func parseSearchResultsHTML(html []byte) ([]models.Product, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var products []models.Product

	// Amazon search results use various selectors depending on the page type
	// Common selectors: div[data-asin], .s-result-item, .sg-col-inner
	doc.Find("div[data-asin]").Each(func(i int, s *goquery.Selection) {
		// Skip items without ASIN or with empty ASIN
		asin, exists := s.Attr("data-asin")
		if !exists || asin == "" {
			return
		}

		product := models.Product{
			ASIN: asin,
		}

		// Extract title - multiple possible selectors
		titleEl := s.Find("h2 a span, .s-title-instructions-style, h2.s-line-clamp-2, h2 span")
		if titleEl.Length() > 0 {
			product.Title = strings.TrimSpace(titleEl.First().Text())
		}

		// Skip if no title found (likely not a valid product)
		if product.Title == "" {
			return
		}

		// Extract price - Amazon uses various price selectors
		priceEl := s.Find(".a-price .a-offscreen, .a-price-whole")
		if priceEl.Length() > 0 {
			priceText := priceEl.First().Text()
			product.Price = parsePriceFromText(priceText)
		}

		// Extract original price (if on sale)
		originalPriceEl := s.Find(".a-price.a-text-price .a-offscreen")
		if originalPriceEl.Length() > 0 {
			originalPriceText := originalPriceEl.First().Text()
			originalPrice := parsePriceFromText(originalPriceText)
			if originalPrice > 0 && originalPrice != product.Price {
				product.OriginalPrice = &originalPrice
			}
		}

		// Extract rating
		ratingEl := s.Find("span[aria-label*='out of'], .a-icon-star-small span.a-icon-alt")
		if ratingEl.Length() > 0 {
			ratingText := ratingEl.First().Text()
			if ariaLabel, exists := ratingEl.Attr("aria-label"); exists {
				ratingText = ariaLabel
			}
			product.Rating = parseRating(ratingText)
		}

		// Extract review count
		reviewEl := s.Find("span[aria-label*='ratings'], .a-size-base.s-underline-text")
		if reviewEl.Length() > 0 {
			reviewText := reviewEl.First().Text()
			if ariaLabel, exists := reviewEl.Attr("aria-label"); exists {
				reviewText = ariaLabel
			}
			product.ReviewCount = parseReviewCount(reviewText)
		}

		// Check if Prime eligible
		primeEl := s.Find("i.a-icon-prime, .s-prime, [aria-label*='Prime']")
		product.Prime = primeEl.Length() > 0

		// Check if in stock - assume in stock unless "unavailable" or "out of stock" found
		product.InStock = true
		availabilityText := strings.ToLower(s.Find(".a-size-base.a-color-secondary, .a-size-base.a-color-price").Text())
		if strings.Contains(availabilityText, "unavailable") ||
			strings.Contains(availabilityText, "out of stock") ||
			strings.Contains(availabilityText, "currently unavailable") {
			product.InStock = false
		}

		// Only add products with at least ASIN, title, and price
		if product.ASIN != "" && product.Title != "" && product.Price > 0 {
			products = append(products, product)
		}
	})

	return products, nil
}

// parsePriceFromText extracts a float64 price from a price string (e.g., "$29.99" -> 29.99)
func parsePriceFromText(priceStr string) float64 {
	// Remove whitespace
	priceStr = strings.TrimSpace(priceStr)

	// Use regex to extract numeric value (handles formats like $29.99, $1,299.99, etc.)
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

// parseRating extracts rating from text like "4.5 out of 5 stars"
func parseRating(ratingText string) float64 {
	ratingText = strings.TrimSpace(ratingText)

	// Try to extract the first number before "out of"
	re := regexp.MustCompile(`(\d+\.?\d*)\s*out of`)
	matches := re.FindStringSubmatch(ratingText)
	if len(matches) > 1 {
		rating, err := strconv.ParseFloat(matches[1], 64)
		if err == nil {
			return rating
		}
	}

	// Try to find any decimal number in the format X.X
	re = regexp.MustCompile(`\d+\.\d+`)
	match := re.FindString(ratingText)
	if match != "" {
		rating, err := strconv.ParseFloat(match, 64)
		if err == nil {
			return rating
		}
	}

	return 0.0
}

// parseReviewCount extracts review count from text like "1,234 ratings" or "1234"
func parseReviewCount(reviewText string) int {
	reviewText = strings.TrimSpace(reviewText)

	// Remove "ratings" or other text, extract just the number
	re := regexp.MustCompile(`([\d,]+)`)
	match := re.FindString(reviewText)
	if match == "" {
		return 0
	}

	// Remove commas
	match = strings.ReplaceAll(match, ",", "")

	// Parse to int
	count, err := strconv.Atoi(match)
	if err != nil {
		return 0
	}

	return count
}
