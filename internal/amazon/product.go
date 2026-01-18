package amazon

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/michaelshimeles/amazon-cli/pkg/models"
)

// GetProduct fetches detailed product information for a given ASIN
func (c *Client) GetProduct(asin string) (*models.Product, error) {
	// Validate ASIN format (10 alphanumeric characters)
	if !isValidASIN(asin) {
		return nil, fmt.Errorf("invalid ASIN format: %s", asin)
	}

	// Build product detail page URL
	productURL := fmt.Sprintf("https://www.amazon.com/dp/%s", asin)

	// Fetch the product page
	resp, err := c.Get(productURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product page: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("product not found: %s", asin)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse HTML response
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Extract product information
	product := &models.Product{
		ASIN: asin,
	}

	// Extract title
	product.Title = strings.TrimSpace(doc.Find("#productTitle").First().Text())

	// Extract price
	priceText := doc.Find(".a-price .a-offscreen").First().Text()
	if priceText != "" {
		product.Price = parsePrice(priceText)
	}

	// Extract original price (if on sale)
	originalPriceText := doc.Find(".a-text-price .a-offscreen").First().Text()
	if originalPriceText != "" && originalPriceText != priceText {
		originalPrice := parsePrice(originalPriceText)
		product.OriginalPrice = &originalPrice
	}

	// Extract rating
	ratingText := doc.Find("#acrPopover").AttrOr("title", "")
	if ratingText != "" {
		product.Rating = parseRating(ratingText)
	}

	// Extract review count
	reviewCountText := doc.Find("#acrCustomerReviewText").First().Text()
	product.ReviewCount = parseReviewCount(reviewCountText)

	// Check for Prime eligibility
	product.Prime = doc.Find("#priceBadging_feature_div").Text() != "" ||
		doc.Find("i.a-icon-prime").Length() > 0

	// Extract stock status
	availabilityText := doc.Find("#availability").Text()
	product.InStock = !strings.Contains(strings.ToLower(availabilityText), "out of stock") &&
		!strings.Contains(strings.ToLower(availabilityText), "unavailable")

	// Extract delivery estimate
	deliveryText := doc.Find("#mir-layout-DELIVERY_BLOCK-slot-PRIMARY_DELIVERY_MESSAGE_LARGE").Text()
	if deliveryText == "" {
		deliveryText = doc.Find("#deliveryMessageMirId").Text()
	}
	product.DeliveryEstimate = strings.TrimSpace(deliveryText)

	// Extract description
	product.Description = strings.TrimSpace(doc.Find("#productDescription p").First().Text())

	// Extract feature bullets
	doc.Find("#feature-bullets ul li").Each(func(i int, s *goquery.Selection) {
		feature := strings.TrimSpace(s.Text())
		if feature != "" && !strings.Contains(feature, "Make sure") {
			product.Features = append(product.Features, feature)
		}
	})

	// Extract images
	doc.Find("#altImages ul li img").Each(func(i int, s *goquery.Selection) {
		if imgURL, exists := s.Attr("src"); exists {
			// Convert thumbnail to full size
			imgURL = strings.ReplaceAll(imgURL, "._SS40_", "")
			product.Images = append(product.Images, imgURL)
		}
	})

	return product, nil
}

// GetProductReviews fetches product reviews for a given ASIN
func (c *Client) GetProductReviews(asin string, limit int) (*models.ReviewsResponse, error) {
	// Validate ASIN format
	if !isValidASIN(asin) {
		return nil, fmt.Errorf("invalid ASIN format: %s", asin)
	}

	// Build reviews page URL
	reviewsURL := fmt.Sprintf("https://www.amazon.com/product-reviews/%s", asin)

	// Fetch the reviews page
	resp, err := c.Get(reviewsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch reviews page: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("reviews not found for ASIN: %s", asin)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse HTML response
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Initialize response
	reviewsResp := &models.ReviewsResponse{
		ASIN:    asin,
		Reviews: make([]models.Review, 0),
	}

	// Extract average rating
	avgRatingText := doc.Find("[data-hook='rating-out-of-text']").First().Text()
	reviewsResp.AverageRating = parseRating(avgRatingText)

	// Extract total review count
	totalReviewsText := doc.Find("[data-hook='total-review-count']").First().Text()
	reviewsResp.TotalReviews = parseReviewCount(totalReviewsText)

	// Extract individual reviews
	count := 0
	doc.Find("[data-hook='review']").Each(func(i int, s *goquery.Selection) {
		if limit > 0 && count >= limit {
			return
		}

		review := models.Review{}

		// Extract rating
		ratingText := s.Find("[data-hook='review-star-rating']").AttrOr("class", "")
		review.Rating = parseStarRating(ratingText)

		// Extract title
		review.Title = strings.TrimSpace(s.Find("[data-hook='review-title']").Text())

		// Extract body
		review.Body = strings.TrimSpace(s.Find("[data-hook='review-body']").Text())

		// Extract author
		review.Author = strings.TrimSpace(s.Find(".a-profile-name").Text())

		// Extract date
		review.Date = strings.TrimSpace(s.Find("[data-hook='review-date']").Text())

		// Check if verified purchase
		review.Verified = s.Find("[data-hook='avp-badge']").Length() > 0

		reviewsResp.Reviews = append(reviewsResp.Reviews, review)
		count++
	})

	return reviewsResp, nil
}

// isValidASIN validates ASIN format (10 alphanumeric characters)
func isValidASIN(asin string) bool {
	match, _ := regexp.MatchString(`^[A-Z0-9]{10}$`, asin)
	return match
}

// parsePrice extracts numeric price from price string (e.g., "$29.99" -> 29.99)
func parsePrice(priceStr string) float64 {
	// Remove currency symbols and whitespace
	priceStr = strings.TrimSpace(priceStr)
	priceStr = strings.ReplaceAll(priceStr, "$", "")
	priceStr = strings.ReplaceAll(priceStr, ",", "")

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return 0.0
	}
	return price
}

// parseRating extracts numeric rating from rating text (e.g., "4.5 out of 5 stars" -> 4.5)
func parseRating(ratingStr string) float64 {
	// Look for pattern like "4.5 out of 5"
	re := regexp.MustCompile(`([\d.]+)\s*out\s*of`)
	matches := re.FindStringSubmatch(ratingStr)
	if len(matches) >= 2 {
		rating, err := strconv.ParseFloat(matches[1], 64)
		if err == nil {
			return rating
		}
	}
	return 0.0
}

// parseStarRating extracts rating from CSS class (e.g., "a-star-4-5" -> 4.5)
func parseStarRating(classStr string) float64 {
	re := regexp.MustCompile(`a-star-(\d+)-?(\d*)`)
	matches := re.FindStringSubmatch(classStr)
	if len(matches) >= 2 {
		whole := matches[1]
		decimal := "0"
		if len(matches) >= 3 && matches[2] != "" {
			decimal = matches[2]
		}
		ratingStr := fmt.Sprintf("%s.%s", whole, decimal)
		rating, err := strconv.ParseFloat(ratingStr, 64)
		if err == nil {
			return rating
		}
	}
	return 0.0
}

// parseReviewCount extracts numeric count from review count text (e.g., "1,234 ratings" -> 1234)
func parseReviewCount(countStr string) int {
	// Extract numbers and remove commas
	re := regexp.MustCompile(`[\d,]+`)
	match := re.FindString(countStr)
	if match == "" {
		return 0
	}

	match = strings.ReplaceAll(match, ",", "")
	count, err := strconv.Atoi(match)
	if err != nil {
		return 0
	}
	return count
}

// Get is a convenience method for making GET requests with rate limiting
func (c *Client) Get(url string) (*http.Response, error) {
	// Wait according to rate limiter
	if c.rateLimiter != nil {
		if err := c.rateLimiter.Wait(); err != nil {
			return nil, fmt.Errorf("rate limiter error: %w", err)
		}
	}

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Execute request with retry logic
	return c.Do(req)
}

// Do executes an HTTP request with retry logic for rate limiting
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	// Set random User-Agent
	if len(c.userAgents) > 0 {
		req.Header.Set("User-Agent", c.userAgents[c.nextUserAgentIndex()])
	}

	// Set common headers
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Check if we need to retry due to rate limiting
	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
		resp.Body.Close()

		// Check if we should retry
		if c.rateLimiter != nil && c.rateLimiter.ShouldRetry(resp.StatusCode, c.retryCount) {
			c.retryCount++
			if err := c.rateLimiter.WaitWithBackoff(c.retryCount); err != nil {
				return nil, fmt.Errorf("backoff wait failed: %w", err)
			}
			// Retry the request
			return c.Do(req)
		}

		return nil, fmt.Errorf("rate limited, max retries exceeded")
	}

	// Reset retry count on successful request
	c.retryCount = 0

	return resp, nil
}

// Client stub - will be defined in client.go
type Client struct {
	httpClient  *http.Client
	rateLimiter *RateLimiter
	userAgents  []string
	retryCount  int
	uaIndex     int
}

// RateLimiter stub - will be defined in ratelimit package
type RateLimiter struct{}

func (r *RateLimiter) Wait() error                              { return nil }
func (r *RateLimiter) WaitWithBackoff(attempt int) error        { return nil }
func (r *RateLimiter) ShouldRetry(statusCode, attempt int) bool { return false }

func (c *Client) nextUserAgentIndex() int {
	c.uaIndex = (c.uaIndex + 1) % len(c.userAgents)
	return c.uaIndex
}
