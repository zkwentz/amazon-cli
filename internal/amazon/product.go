package amazon

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

// GetProduct fetches detailed product information for a given ASIN
func (c *Client) GetProduct(asin string) (*models.Product, error) {
	if err := validateASIN(asin); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://www.amazon.com/dp/%s", asin)

	resp, err := c.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("product not found: ASIN %s", asin)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	product := &models.Product{
		ASIN: asin,
	}

	// Extract title
	product.Title = strings.TrimSpace(doc.Find("#productTitle").Text())

	// Extract price
	product.Price = extractPrice(doc)

	// Extract original price (if on sale)
	if origPrice := extractOriginalPrice(doc); origPrice > 0 {
		product.OriginalPrice = &origPrice
	}

	// Extract rating
	product.Rating = extractRating(doc)

	// Extract review count
	product.ReviewCount = extractReviewCount(doc)

	// Check Prime eligibility
	product.Prime = doc.Find("#priceBadging_feature_div .a-icon-prime").Length() > 0 ||
		doc.Find("#deliveryMessageMirId .a-icon-prime").Length() > 0

	// Extract stock status
	product.InStock = extractStockStatus(doc)

	// Extract delivery estimate
	product.DeliveryEstimate = extractDeliveryEstimate(doc)

	// Extract description
	product.Description = extractDescription(doc)

	// Extract features
	product.Features = extractFeatures(doc)

	// Extract images
	product.Images = extractImages(doc)

	return product, nil
}

// validateASIN checks if the ASIN is in the correct format
func validateASIN(asin string) error {
	if len(asin) != 10 {
		return fmt.Errorf("invalid ASIN format: must be 10 characters, got %d", len(asin))
	}

	matched, _ := regexp.MatchString("^[A-Z0-9]{10}$", asin)
	if !matched {
		return fmt.Errorf("invalid ASIN format: must contain only uppercase letters and numbers")
	}

	return nil
}

// extractPrice extracts the product price from the document
func extractPrice(doc *goquery.Document) float64 {
	// Try multiple price selectors as Amazon's HTML varies
	priceSelectors := []string{
		".a-price .a-offscreen",
		"#priceblock_ourprice",
		"#priceblock_dealprice",
		".a-price-whole",
	}

	for _, selector := range priceSelectors {
		if priceText := doc.Find(selector).First().Text(); priceText != "" {
			price := parsePrice(priceText)
			if price > 0 {
				return price
			}
		}
	}

	return 0.0
}

// extractOriginalPrice extracts the original price if the product is on sale
func extractOriginalPrice(doc *goquery.Document) float64 {
	selectors := []string{
		".a-price.a-text-price .a-offscreen",
		"#priceblock_saleprice",
		".basisPrice .a-offscreen",
	}

	for _, selector := range selectors {
		if priceText := doc.Find(selector).First().Text(); priceText != "" {
			price := parsePrice(priceText)
			if price > 0 {
				return price
			}
		}
	}

	return 0.0
}

// parsePrice converts a price string to float64
func parsePrice(priceStr string) float64 {
	// Remove currency symbols and whitespace
	cleaned := strings.TrimSpace(priceStr)
	cleaned = strings.ReplaceAll(cleaned, "$", "")
	cleaned = strings.ReplaceAll(cleaned, ",", "")

	// Extract first number found
	re := regexp.MustCompile(`[\d.]+`)
	match := re.FindString(cleaned)
	if match == "" {
		return 0.0
	}

	price, err := strconv.ParseFloat(match, 64)
	if err != nil {
		return 0.0
	}

	return price
}

// extractRating extracts the product rating
func extractRating(doc *goquery.Document) float64 {
	ratingSelectors := []string{
		"#acrPopover",
		".a-icon-star .a-icon-alt",
		"#averageCustomerReviews .a-icon-alt",
	}

	for _, selector := range ratingSelectors {
		if ratingText := doc.Find(selector).First().Text(); ratingText != "" {
			// Extract number from text like "4.5 out of 5 stars"
			re := regexp.MustCompile(`([\d.]+)\s*out of`)
			matches := re.FindStringSubmatch(ratingText)
			if len(matches) > 1 {
				rating, err := strconv.ParseFloat(matches[1], 64)
				if err == nil {
					return rating
				}
			}
		}
	}

	return 0.0
}

// extractReviewCount extracts the number of reviews
func extractReviewCount(doc *goquery.Document) int {
	reviewSelectors := []string{
		"#acrCustomerReviewText",
		"#averageCustomerReviews span",
	}

	for _, selector := range reviewSelectors {
		if reviewText := doc.Find(selector).First().Text(); reviewText != "" {
			// Extract number from text like "1,234 ratings"
			cleaned := strings.ReplaceAll(reviewText, ",", "")
			re := regexp.MustCompile(`(\d+)`)
			matches := re.FindStringSubmatch(cleaned)
			if len(matches) > 1 {
				count, err := strconv.Atoi(matches[1])
				if err == nil {
					return count
				}
			}
		}
	}

	return 0
}

// extractStockStatus checks if the product is in stock
func extractStockStatus(doc *goquery.Document) bool {
	availabilityText := doc.Find("#availability").Text()
	availabilityText = strings.ToLower(strings.TrimSpace(availabilityText))

	// Check for out of stock indicators
	outOfStockKeywords := []string{
		"currently unavailable",
		"out of stock",
		"unavailable",
	}

	for _, keyword := range outOfStockKeywords {
		if strings.Contains(availabilityText, keyword) {
			return false
		}
	}

	// Check for in stock indicators
	inStockKeywords := []string{
		"in stock",
		"available",
		"add to cart",
	}

	for _, keyword := range inStockKeywords {
		if strings.Contains(availabilityText, keyword) {
			return true
		}
	}

	// Also check if "Add to Cart" button exists
	return doc.Find("#add-to-cart-button").Length() > 0
}

// extractDeliveryEstimate extracts the delivery estimate
func extractDeliveryEstimate(doc *goquery.Document) string {
	deliverySelectors := []string{
		"#deliveryMessageMirId",
		"#delivery-message",
		".a-color-success.a-text-bold",
	}

	for _, selector := range deliverySelectors {
		if deliveryText := doc.Find(selector).First().Text(); deliveryText != "" {
			return strings.TrimSpace(deliveryText)
		}
	}

	return ""
}

// extractDescription extracts the product description
func extractDescription(doc *goquery.Document) string {
	descSelectors := []string{
		"#productDescription p",
		"#feature-bullets",
		".a-unordered-list.a-vertical.a-spacing-mini",
	}

	for _, selector := range descSelectors {
		if desc := doc.Find(selector).First().Text(); desc != "" {
			return strings.TrimSpace(desc)
		}
	}

	return ""
}

// extractFeatures extracts product feature bullets
func extractFeatures(doc *goquery.Document) []string {
	features := []string{}

	// Try multiple selectors for feature bullets
	doc.Find("#feature-bullets li, .a-unordered-list.a-vertical li").Each(func(i int, s *goquery.Selection) {
		feature := strings.TrimSpace(s.Text())
		if feature != "" && !strings.Contains(feature, "See more product details") {
			features = append(features, feature)
		}
	})

	return features
}

// extractImages extracts product image URLs
func extractImages(doc *goquery.Document) []string {
	images := []string{}
	seen := make(map[string]bool)

	// Extract from image thumbnails
	doc.Find("#altImages img, #imageBlock img").Each(func(i int, s *goquery.Selection) {
		if src, exists := s.Attr("src"); exists {
			// Convert thumbnail to full-size image URL
			fullSizeURL := convertToFullSizeImage(src)
			if fullSizeURL != "" && !seen[fullSizeURL] {
				images = append(images, fullSizeURL)
				seen[fullSizeURL] = true
			}
		}
	})

	return images
}

// convertToFullSizeImage converts a thumbnail URL to full-size image URL
func convertToFullSizeImage(thumbnailURL string) string {
	// Amazon image URLs typically have size indicators like ._AC_UL160_SR160,160_
	// Replace with larger size or remove size restrictions
	re := regexp.MustCompile(`\._.*?_\.`)
	fullSizeURL := re.ReplaceAllString(thumbnailURL, ".")

	if strings.HasPrefix(fullSizeURL, "http") {
		return fullSizeURL
	}

	return ""
}
