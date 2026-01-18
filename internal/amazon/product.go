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

// GetProductReviews fetches reviews for a product by ASIN
func (c *Client) GetProductReviews(asin string, limit int) (*models.ReviewsResponse, error) {
	if limit <= 0 {
		limit = 10
	}

	// Amazon reviews URL
	url := fmt.Sprintf("https://www.amazon.com/product-reviews/%s", asin)

	resp, err := c.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch reviews: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("amazon returned status %d: %s", resp.StatusCode, string(body))
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	response := &models.ReviewsResponse{
		ASIN:    asin,
		Reviews: []models.Review{},
	}

	// Extract average rating and total review count
	response.AverageRating = extractAverageRating(doc)
	response.TotalReviews = extractTotalReviews(doc)

	// Extract individual reviews
	reviewCount := 0
	doc.Find("div[data-hook='review']").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if reviewCount >= limit {
			return false
		}

		review := extractReview(s)
		if review.Title != "" || review.Body != "" {
			response.Reviews = append(response.Reviews, review)
			reviewCount++
		}

		return true
	})

	return response, nil
}

func extractAverageRating(doc *goquery.Document) float64 {
	// Try to find average rating
	ratingText := doc.Find("div[data-hook='rating-out-of-text']").First().Text()
	if ratingText == "" {
		ratingText = doc.Find("i[data-hook='average-star-rating'] span").First().Text()
	}

	// Extract number from text like "4.5 out of 5 stars"
	re := regexp.MustCompile(`([0-9.]+)`)
	matches := re.FindStringSubmatch(ratingText)
	if len(matches) > 1 {
		rating, _ := strconv.ParseFloat(matches[1], 64)
		return rating
	}

	return 0.0
}

func extractTotalReviews(doc *goquery.Document) int {
	// Try to find total review count
	countText := doc.Find("div[data-hook='cr-filter-info-review-rating-count']").First().Text()
	if countText == "" {
		countText = doc.Find("div[data-hook='total-review-count']").First().Text()
	}

	// Extract number from text like "1,234 global ratings"
	re := regexp.MustCompile(`([0-9,]+)`)
	matches := re.FindStringSubmatch(countText)
	if len(matches) > 1 {
		countStr := strings.ReplaceAll(matches[1], ",", "")
		count, _ := strconv.Atoi(countStr)
		return count
	}

	return 0
}

func extractReview(s *goquery.Selection) models.Review {
	review := models.Review{}

	// Extract rating
	ratingText := s.Find("i[data-hook='review-star-rating'] span").First().Text()
	if ratingText == "" {
		ratingText = s.Find("i[data-hook='cmps-review-star-rating'] span").First().Text()
	}
	re := regexp.MustCompile(`([0-9.]+)`)
	matches := re.FindStringSubmatch(ratingText)
	if len(matches) > 1 {
		review.Rating, _ = strconv.ParseFloat(matches[1], 64)
	}

	// Extract title
	review.Title = strings.TrimSpace(s.Find("a[data-hook='review-title'] span").Text())
	if review.Title == "" {
		review.Title = strings.TrimSpace(s.Find("a[data-hook='review-title']").Text())
	}

	// Extract body
	review.Body = strings.TrimSpace(s.Find("span[data-hook='review-body'] span").Text())
	if review.Body == "" {
		review.Body = strings.TrimSpace(s.Find("span[data-hook='review-body']").Text())
	}

	// Extract author
	review.Author = strings.TrimSpace(s.Find("span.a-profile-name").First().Text())

	// Extract date
	review.Date = strings.TrimSpace(s.Find("span[data-hook='review-date']").Text())

	// Check if verified purchase
	verifiedText := s.Find("span[data-hook='avp-badge']").Text()
	review.Verified = strings.Contains(verifiedText, "Verified Purchase")

	return review
}
