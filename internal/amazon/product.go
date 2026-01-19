package amazon

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

// GetProduct retrieves detailed product information
func (c *Client) GetProduct(asin string) (*models.Product, error) {
	if asin == "" {
		return nil, fmt.Errorf("ASIN cannot be empty")
	}

	// TODO: Implement actual Amazon API call
	// For now, return mock data

	originalPrice := 349.99

	return &models.Product{
		ASIN:             asin,
		Title:            "Sony WH-1000XM4 Wireless Premium Noise Canceling Overhead Headphones",
		Price:            278.00,
		OriginalPrice:    &originalPrice,
		Rating:           4.7,
		ReviewCount:      52431,
		Prime:            true,
		InStock:          true,
		DeliveryEstimate: "Tomorrow",
		Description:      "Industry-leading noise canceling with Dual Noise Sensor technology. Next-level music with Edge-AI, co-developed with Sony Music Studios Tokyo. Up to 30-hour battery life with quick charging (10 min charge for 5 hours of playback).",
		Features: []string{
			"Industry-leading noise cancellation",
			"30-hour battery life",
			"Touch sensor controls",
			"Speak-to-chat technology",
			"Wearing detection",
			"Multipoint connection",
		},
		Images: []string{
			"https://images-na.ssl-images-amazon.com/images/I/71o8Q5XJS5L._AC_SL1500_.jpg",
			"https://images-na.ssl-images-amazon.com/images/I/81WpXBD4uWL._AC_SL1500_.jpg",
		},
	}, nil
}

// GetProductReviews retrieves reviews for a product
func (c *Client) GetProductReviews(asin string, limit int) (*models.ReviewsResponse, error) {
	if asin == "" {
		return nil, fmt.Errorf("ASIN cannot be empty")
	}

	if limit <= 0 {
		limit = 10
	}

	// TODO: Implement actual Amazon API call
	// For now, return mock data

	reviews := []models.Review{
		{
			Rating:   5,
			Title:    "Best headphones I've ever owned",
			Body:     "The noise canceling is incredible. I use these daily for work calls and music. Battery life is exactly as advertised.",
			Author:   "John D.",
			Date:     time.Now().AddDate(0, 0, -10).Format("2006-01-02"),
			Verified: true,
		},
		{
			Rating:   4,
			Title:    "Great but pricey",
			Body:     "Sound quality is excellent and the ANC is top-notch. Only complaint is the price, but you get what you pay for.",
			Author:   "Sarah M.",
			Date:     time.Now().AddDate(0, 0, -25).Format("2006-01-02"),
			Verified: true,
		},
		{
			Rating:   5,
			Title:    "Perfect for travel",
			Body:     "Used these on a 12-hour flight and they were amazing. The noise canceling blocked out all the engine noise.",
			Author:   "Mike R.",
			Date:     time.Now().AddDate(0, -1, -5).Format("2006-01-02"),
			Verified: true,
		},
	}

	if len(reviews) > limit {
		reviews = reviews[:limit]
	}

	return &models.ReviewsResponse{
		ASIN:          asin,
		AverageRating: 4.7,
		TotalReviews:  52431,
		Reviews:       reviews,
	}, nil
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
