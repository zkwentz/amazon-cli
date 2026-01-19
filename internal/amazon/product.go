package amazon

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

// GetProduct retrieves detailed product information
func (c *Client) GetProduct(asin string) (*models.Product, error) {
	// Validate ASIN is not empty
	if asin == "" {
		return nil, fmt.Errorf("ASIN cannot be empty")
	}

	// Validate ASIN format - ASINs are typically 10 characters (alphanumeric)
	asinRegex := regexp.MustCompile(`^[A-Z0-9]{10}$`)
	if !asinRegex.MatchString(asin) {
		return nil, fmt.Errorf("invalid ASIN format: must be 10 alphanumeric characters")
	}

	// Construct product detail URL
	productURL := fmt.Sprintf("%s/dp/%s", c.baseURL, asin)

	// Create HTTP GET request
	req, err := http.NewRequest("GET", productURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Execute the request with rate limiting and retries
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product details: %w", err)
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
	product, err := parseProductDetailHTML(body.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to parse product details: %w", err)
	}

	return product, nil
}

// GetProductReviews retrieves reviews for a product
func (c *Client) GetProductReviews(asin string, limit int) (*models.ReviewsResponse, error) {
	if asin == "" {
		return nil, fmt.Errorf("ASIN cannot be empty")
	}

	if limit <= 0 {
		limit = 10
	}

	// Build reviews URL
	reviewsURL := fmt.Sprintf("%s/product-reviews/%s", c.baseURL, asin)

	// Create HTTP GET request
	req, err := http.NewRequest("GET", reviewsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Execute the request with rate limiting and retries
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch reviews: %w", err)
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
	reviewsResponse, err := parseReviewsHTML(body.Bytes(), asin, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to parse reviews: %w", err)
	}

	return reviewsResponse, nil
}

// parseProductDetailHTML parses Amazon product detail page HTML and extracts product information
func parseProductDetailHTML(html []byte) (*models.Product, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	product := &models.Product{}

	// Extract ASIN - multiple possible locations
	// 1. From data-asin attribute
	asinEl := doc.Find("[data-asin]").First()
	if asin, exists := asinEl.Attr("data-asin"); exists && asin != "" {
		product.ASIN = asin
	}

	// 2. From canonical link or URL pattern
	if product.ASIN == "" {
		canonicalEl := doc.Find("link[rel='canonical']")
		if href, exists := canonicalEl.Attr("href"); exists {
			asinRegex := regexp.MustCompile(`/dp/([A-Z0-9]{10})`)
			if matches := asinRegex.FindStringSubmatch(href); len(matches) > 1 {
				product.ASIN = matches[1]
			}
		}
	}

	// 3. From input field with name="ASIN"
	if product.ASIN == "" {
		asinInput := doc.Find("input[name='ASIN']")
		if asin, exists := asinInput.Attr("value"); exists && asin != "" {
			product.ASIN = asin
		}
	}

	// Extract title
	titleSelectors := []string{
		"#productTitle",
		"#title",
		"h1.product-title",
		"span#productTitle",
	}
	for _, selector := range titleSelectors {
		titleEl := doc.Find(selector)
		if titleEl.Length() > 0 {
			product.Title = strings.TrimSpace(titleEl.First().Text())
			if product.Title != "" {
				break
			}
		}
	}

	// Extract price - Amazon has multiple price containers
	priceSelectors := []string{
		".a-price[data-a-color='price'] .a-offscreen",
		"#priceblock_ourprice",
		"#priceblock_dealprice",
		".a-price .a-offscreen",
		"span.a-price-whole",
	}
	for _, selector := range priceSelectors {
		priceEl := doc.Find(selector)
		if priceEl.Length() > 0 {
			priceText := priceEl.First().Text()
			price := parsePriceFromText(priceText)
			if price > 0 {
				product.Price = price
				break
			}
		}
	}

	// Extract original price (list price / was price)
	originalPriceSelectors := []string{
		".a-price[data-a-strike='true'] .a-offscreen",
		"#priceblock_listprice",
		".a-text-price .a-offscreen",
		"span.a-price.a-text-price span.a-offscreen",
	}
	for _, selector := range originalPriceSelectors {
		originalPriceEl := doc.Find(selector)
		if originalPriceEl.Length() > 0 {
			originalPriceText := originalPriceEl.First().Text()
			originalPrice := parsePriceFromText(originalPriceText)
			if originalPrice > 0 && originalPrice != product.Price {
				product.OriginalPrice = &originalPrice
				break
			}
		}
	}

	// Extract rating
	ratingSelectors := []string{
		"#acrPopover",
		"span.a-icon-alt",
		"i.a-icon-star span.a-icon-alt",
	}
	for _, selector := range ratingSelectors {
		ratingEl := doc.Find(selector)
		if ratingEl.Length() > 0 {
			ratingText := ratingEl.First().Text()
			if title, exists := ratingEl.Attr("title"); exists {
				ratingText = title
			}
			rating := parseRating(ratingText)
			if rating > 0 {
				product.Rating = rating
				break
			}
		}
	}

	// Extract review count
	reviewSelectors := []string{
		"#acrCustomerReviewText",
		"span#acrCustomerReviewText",
		"a#acrCustomerReviewLink span",
	}
	for _, selector := range reviewSelectors {
		reviewEl := doc.Find(selector)
		if reviewEl.Length() > 0 {
			reviewText := reviewEl.First().Text()
			count := parseReviewCount(reviewText)
			if count > 0 {
				product.ReviewCount = count
				break
			}
		}
	}

	// Check if Prime eligible
	primeSelectors := []string{
		"#priceBadging_feature_div i.a-icon-prime",
		"i.a-icon-prime",
		"span.prime-badge",
		"[aria-label*='Prime']",
	}
	for _, selector := range primeSelectors {
		primeEl := doc.Find(selector)
		if primeEl.Length() > 0 {
			product.Prime = true
			break
		}
	}

	// Check availability/stock status
	product.InStock = true // Default to in stock
	availabilitySelectors := []string{
		"#availability span",
		"#availability",
		"div#availability",
	}
	for _, selector := range availabilitySelectors {
		availEl := doc.Find(selector)
		if availEl.Length() > 0 {
			availText := strings.ToLower(strings.TrimSpace(availEl.First().Text()))
			if strings.Contains(availText, "out of stock") ||
				strings.Contains(availText, "unavailable") ||
				strings.Contains(availText, "currently unavailable") ||
				strings.Contains(availText, "not available") {
				product.InStock = false
				break
			}
		}
	}

	// Extract delivery estimate
	deliverySelectors := []string{
		"#deliveryMessageMirId span",
		"#mir-layout-DELIVERY_BLOCK span.a-text-bold",
		"div[data-feature-name='deliveryMessage'] span",
	}
	for _, selector := range deliverySelectors {
		deliveryEl := doc.Find(selector)
		if deliveryEl.Length() > 0 {
			deliveryText := strings.TrimSpace(deliveryEl.First().Text())
			if deliveryText != "" {
				product.DeliveryEstimate = deliveryText
				break
			}
		}
	}

	// Extract description
	descriptionSelectors := []string{
		"#productDescription p",
		"#feature-bullets ul li span.a-list-item",
		"div#productDescription",
	}
	for _, selector := range descriptionSelectors {
		descEl := doc.Find(selector)
		if descEl.Length() > 0 {
			var descParts []string
			descEl.Each(func(i int, s *goquery.Selection) {
				text := strings.TrimSpace(s.Text())
				if text != "" {
					descParts = append(descParts, text)
				}
			})
			if len(descParts) > 0 {
				product.Description = strings.Join(descParts, " ")
				break
			}
		}
	}

	// Extract features/bullet points
	featureSelectors := []string{
		"#feature-bullets ul li span.a-list-item",
		"div#feature-bullets ul.a-unordered-list li span",
	}
	for _, selector := range featureSelectors {
		featureEl := doc.Find(selector)
		if featureEl.Length() > 0 {
			var features []string
			featureEl.Each(func(i int, s *goquery.Selection) {
				text := strings.TrimSpace(s.Text())
				// Filter out empty strings and common non-feature text
				if text != "" && !strings.HasPrefix(text, "See more") {
					features = append(features, text)
				}
			})
			if len(features) > 0 {
				product.Features = features
				break
			}
		}
	}

	// Extract images
	imageSelectors := []string{
		"#altImages ul li.imageThumbnail img",
		"#imageBlock img[data-old-hires]",
		"#landingImage",
		"#imgTagWrapperId img",
	}

	imageSet := make(map[string]bool) // Use map to avoid duplicates

	for _, selector := range imageSelectors {
		imgEl := doc.Find(selector)
		if imgEl.Length() > 0 {
			imgEl.Each(func(i int, s *goquery.Selection) {
				// Try to get high-res image first
				if hiresURL, exists := s.Attr("data-old-hires"); exists && hiresURL != "" {
					imageSet[hiresURL] = true
				} else if hiresURL, exists := s.Attr("data-a-hires"); exists && hiresURL != "" {
					imageSet[hiresURL] = true
				} else if srcURL, exists := s.Attr("src"); exists && srcURL != "" {
					// Filter out tiny thumbnails and placeholder images
					if !strings.Contains(srcURL, "1x1") &&
					   !strings.Contains(srcURL, "pixel") &&
					   !strings.Contains(srcURL, "transparent") {
						imageSet[srcURL] = true
					}
				}
			})
		}
	}

	// Convert map to slice
	for imgURL := range imageSet {
		product.Images = append(product.Images, imgURL)
	}

	// Validate that we at least have ASIN and title
	if product.ASIN == "" || product.Title == "" {
		return nil, fmt.Errorf("failed to extract required fields (ASIN or title)")
	}

	return product, nil
}

// parseReviewsHTML parses Amazon product reviews page HTML and extracts review information
func parseReviewsHTML(html []byte, asin string, limit int) (*models.ReviewsResponse, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	response := &models.ReviewsResponse{
		ASIN:    asin,
		Reviews: []models.Review{},
	}

	// Extract average rating from the page header
	avgRatingSelectors := []string{
		"div[data-hook='rating-out-of-text']",
		"span[data-hook='rating-out-of-text']",
		"i[data-hook='average-star-rating'] span.a-icon-alt",
	}
	for _, selector := range avgRatingSelectors {
		ratingEl := doc.Find(selector)
		if ratingEl.Length() > 0 {
			ratingText := strings.TrimSpace(ratingEl.First().Text())
			rating := parseRating(ratingText)
			if rating > 0 {
				response.AverageRating = rating
				break
			}
		}
	}

	// Extract total review count
	totalReviewSelectors := []string{
		"div[data-hook='total-review-count']",
		"span[data-hook='total-review-count']",
		"div[data-hook='cr-filter-info-review-rating-count']",
	}
	for _, selector := range totalReviewSelectors {
		countEl := doc.Find(selector)
		if countEl.Length() > 0 {
			countText := strings.TrimSpace(countEl.First().Text())
			count := parseReviewCount(countText)
			if count > 0 {
				response.TotalReviews = count
				break
			}
		}
	}

	// Parse individual reviews
	reviewCount := 0
	doc.Find("div[data-hook='review']").Each(func(i int, s *goquery.Selection) {
		if reviewCount >= limit {
			return
		}

		review := models.Review{}

		// Extract rating - look for star rating
		ratingEl := s.Find("i[data-hook='review-star-rating'] span.a-icon-alt, i[data-hook='cmps-review-star-rating'] span.a-icon-alt")
		if ratingEl.Length() > 0 {
			ratingText := strings.TrimSpace(ratingEl.First().Text())
			ratingFloat := parseRating(ratingText)
			review.Rating = int(ratingFloat)
		}

		// Extract title
		titleEl := s.Find("a[data-hook='review-title'] span, span[data-hook='review-title']")
		if titleEl.Length() > 0 {
			review.Title = strings.TrimSpace(titleEl.First().Text())
		}

		// Extract body
		bodyEl := s.Find("span[data-hook='review-body'] span")
		if bodyEl.Length() > 0 {
			review.Body = strings.TrimSpace(bodyEl.First().Text())
		}

		// Extract author
		authorEl := s.Find("span.a-profile-name")
		if authorEl.Length() > 0 {
			review.Author = strings.TrimSpace(authorEl.First().Text())
		}

		// Extract date
		dateEl := s.Find("span[data-hook='review-date']")
		if dateEl.Length() > 0 {
			dateText := strings.TrimSpace(dateEl.First().Text())
			review.Date = parseDateFromReview(dateText)
		}

		// Check if verified purchase
		verifiedEl := s.Find("span[data-hook='avp-badge']")
		review.Verified = verifiedEl.Length() > 0

		// Only add review if it has at least title or body
		if review.Title != "" || review.Body != "" {
			response.Reviews = append(response.Reviews, review)
			reviewCount++
		}
	})

	return response, nil
}

// parseDateFromReview extracts and formats date from review date text
// Amazon review dates are typically in format "Reviewed in [Country] on [Date]"
func parseDateFromReview(dateText string) string {
	// Remove "Reviewed in [Country] on " prefix
	dateText = strings.TrimSpace(dateText)

	// Try to extract date after "on "
	if idx := strings.LastIndex(dateText, " on "); idx != -1 {
		dateText = dateText[idx+4:]
	}

	// Try to parse common date formats
	dateFormats := []string{
		"January 2, 2006",
		"Jan 2, 2006",
		"2 January 2006",
		"2006-01-02",
	}

	for _, format := range dateFormats {
		if t, err := time.Parse(format, dateText); err == nil {
			return t.Format("2006-01-02")
		}
	}

	// If parsing fails, return the original text
	return dateText
}
