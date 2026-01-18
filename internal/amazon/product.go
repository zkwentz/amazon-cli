package amazon

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

// GetProduct fetches product details by ASIN
func (c *Client) GetProduct(asin string) (*models.Product, error) {
	// Validate ASIN format (10 alphanumeric characters)
	if !isValidASIN(asin) {
		return nil, models.NewCLIError(models.ErrCodeInvalidInput, "Invalid ASIN format. ASIN must be 10 alphanumeric characters")
	}

	url := fmt.Sprintf("https://www.amazon.com/dp/%s", asin)
	resp, err := c.Get(url)
	if err != nil {
		return nil, models.NewCLIError(models.ErrCodeNetworkError, fmt.Sprintf("Failed to fetch product: %v", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, models.NewCLIError(models.ErrCodeNotFound, fmt.Sprintf("Product with ASIN %s not found", asin))
	}

	if resp.StatusCode != 200 {
		return nil, models.NewCLIError(models.ErrCodeAmazonError, fmt.Sprintf("Amazon returned status code %d", resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, models.NewCLIError(models.ErrCodeNetworkError, fmt.Sprintf("Failed to read response: %v", err))
	}

	product, err := parseProductPage(asin, body)
	if err != nil {
		return nil, models.NewCLIError(models.ErrCodeAmazonError, fmt.Sprintf("Failed to parse product page: %v", err))
	}

	return product, nil
}

// isValidASIN validates ASIN format
func isValidASIN(asin string) bool {
	match, _ := regexp.MatchString(`^[A-Z0-9]{10}$`, asin)
	return match
}

// parseProductPage parses the Amazon product page HTML
func parseProductPage(asin string, htmlData []byte) (*models.Product, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(htmlData)))
	if err != nil {
		return nil, err
	}

	product := &models.Product{
		ASIN: asin,
	}

	// Extract title
	product.Title = strings.TrimSpace(doc.Find("#productTitle").Text())
	if product.Title == "" {
		product.Title = strings.TrimSpace(doc.Find("h1.a-size-large").Text())
	}

	// Extract price
	priceText := doc.Find(".a-price .a-offscreen").First().Text()
	if priceText == "" {
		priceText = doc.Find("#priceblock_ourprice").Text()
	}
	if priceText == "" {
		priceText = doc.Find("#priceblock_dealprice").Text()
	}
	product.Price = parsePrice(priceText)

	// Extract original price (if on sale)
	originalPriceText := doc.Find(".a-price.a-text-price .a-offscreen").Text()
	if originalPriceText != "" {
		originalPrice := parsePrice(originalPriceText)
		if originalPrice > product.Price {
			product.OriginalPrice = &originalPrice
		}
	}

	// Extract rating
	ratingText := doc.Find("span.a-icon-alt").First().Text()
	product.Rating = parseRating(ratingText)

	// Extract review count
	reviewText := doc.Find("#acrCustomerReviewText").Text()
	product.ReviewCount = parseReviewCount(reviewText)

	// Check for Prime
	product.Prime = doc.Find("#priceBadging_feature_div i.a-icon-prime").Length() > 0

	// Check stock status
	availText := strings.ToLower(doc.Find("#availability span").Text())
	product.InStock = !strings.Contains(availText, "out of stock") &&
	                  !strings.Contains(availText, "currently unavailable")

	// Extract delivery estimate
	deliveryText := doc.Find("#mir-layout-DELIVERY_BLOCK-slot-PRIMARY_DELIVERY_MESSAGE_LARGE").Text()
	if deliveryText == "" {
		deliveryText = doc.Find("#deliveryMessageMirId").Text()
	}
	product.DeliveryEstimate = strings.TrimSpace(deliveryText)

	// Extract description
	product.Description = strings.TrimSpace(doc.Find("#productDescription p").Text())
	if product.Description == "" {
		product.Description = strings.TrimSpace(doc.Find("#feature-bullets").Text())
	}

	// Extract features
	doc.Find("#feature-bullets li span.a-list-item").Each(func(i int, s *goquery.Selection) {
		feature := strings.TrimSpace(s.Text())
		if feature != "" && !strings.HasPrefix(feature, "›") {
			product.Features = append(product.Features, feature)
		}
	})

	// Extract images
	doc.Find("#altImages li.imageThumbnail img").Each(func(i int, s *goquery.Selection) {
		if src, exists := s.Attr("src"); exists {
			// Convert thumbnail to full-size image URL
			fullSizeURL := strings.Replace(src, "._AC_US40_", "", 1)
			product.Images = append(product.Images, fullSizeURL)
		}
	})

	// If no images from thumbnails, try main image
	if len(product.Images) == 0 {
		if mainImg, exists := doc.Find("#landingImage").Attr("src"); exists {
			product.Images = append(product.Images, mainImg)
		}
	}

	return product, nil
}

// parsePrice extracts numeric price from text
func parsePrice(text string) float64 {
	// Remove all non-numeric characters except decimal point
	re := regexp.MustCompile(`[^0-9.]`)
	cleaned := re.ReplaceAllString(text, "")

	price, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0
	}
	return price
}

// parseRating extracts rating from text like "4.5 out of 5 stars"
func parseRating(text string) float64 {
	re := regexp.MustCompile(`([0-9.]+)\s+out of`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		rating, err := strconv.ParseFloat(matches[1], 64)
		if err == nil {
			return rating
		}
	}
	return 0
}

// parseReviewCount extracts review count from text
func parseReviewCount(text string) int {
	// Remove commas and extract number
	re := regexp.MustCompile(`[0-9,]+`)
	matches := re.FindString(text)
	cleaned := strings.ReplaceAll(matches, ",", "")

	count, err := strconv.Atoi(cleaned)
	if err != nil {
		return 0
	}
	return count
}
