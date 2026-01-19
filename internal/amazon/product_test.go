package amazon

import (
	"testing"
)

func TestParseProductDetailHTML(t *testing.T) {
	// Sample HTML that mimics Amazon product detail page structure
	html := []byte(`
		<html>
		<head>
			<link rel="canonical" href="https://www.amazon.com/dp/B08N5WRWNW">
		</head>
		<body>
			<div id="dp-container">
				<input type="hidden" name="ASIN" value="B08N5WRWNW">
				<span id="productTitle">Sony WH-1000XM4 Wireless Premium Noise Canceling Headphones</span>

				<div id="corePrice_feature_div">
					<span class="a-price" data-a-color="price">
						<span class="a-offscreen">$278.00</span>
					</span>
					<span class="a-price a-text-price" data-a-strike="true">
						<span class="a-offscreen">$349.99</span>
					</span>
				</div>

				<div id="averageCustomerReviews">
					<span id="acrPopover" title="4.7 out of 5 stars">
						<span class="a-icon-alt">4.7 out of 5 stars</span>
					</span>
					<span id="acrCustomerReviewText">52,431 ratings</span>
				</div>

				<div id="priceBadging_feature_div">
					<i class="a-icon-prime"></i>
				</div>

				<div id="availability">
					<span class="a-size-medium a-color-success">In Stock</span>
				</div>

				<div id="deliveryMessageMirId">
					<span class="a-text-bold">Tomorrow</span>
				</div>

				<div id="feature-bullets">
					<ul class="a-unordered-list">
						<li><span class="a-list-item">Industry-leading noise cancellation</span></li>
						<li><span class="a-list-item">30-hour battery life</span></li>
						<li><span class="a-list-item">Touch sensor controls</span></li>
						<li><span class="a-list-item">Speak-to-chat technology</span></li>
					</ul>
				</div>

				<div id="productDescription">
					<p>Industry-leading noise canceling with Dual Noise Sensor technology. Next-level music with Edge-AI.</p>
				</div>

				<div id="altImages">
					<ul>
						<li class="imageThumbnail">
							<img src="https://images-na.ssl-images-amazon.com/images/I/71o8Q5XJS5L._AC_SL1500_.jpg"
								 data-old-hires="https://images-na.ssl-images-amazon.com/images/I/71o8Q5XJS5L._AC_SL1500_.jpg">
						</li>
						<li class="imageThumbnail">
							<img src="https://images-na.ssl-images-amazon.com/images/I/81WpXBD4uWL._AC_SL1500_.jpg"
								 data-old-hires="https://images-na.ssl-images-amazon.com/images/I/81WpXBD4uWL._AC_SL1500_.jpg">
						</li>
					</ul>
				</div>
			</div>
		</body>
		</html>
	`)

	product, err := parseProductDetailHTML(html)
	if err != nil {
		t.Fatalf("parseProductDetailHTML failed: %v", err)
	}

	// Verify ASIN
	if product.ASIN != "B08N5WRWNW" {
		t.Errorf("Expected ASIN B08N5WRWNW, got %s", product.ASIN)
	}

	// Verify title
	if product.Title != "Sony WH-1000XM4 Wireless Premium Noise Canceling Headphones" {
		t.Errorf("Expected Sony headphones title, got %s", product.Title)
	}

	// Verify price
	if product.Price != 278.00 {
		t.Errorf("Expected price 278.00, got %f", product.Price)
	}

	// Verify original price
	if product.OriginalPrice == nil || *product.OriginalPrice != 349.99 {
		t.Errorf("Expected original price 349.99, got %v", product.OriginalPrice)
	}

	// Verify rating
	if product.Rating != 4.7 {
		t.Errorf("Expected rating 4.7, got %f", product.Rating)
	}

	// Verify review count
	if product.ReviewCount != 52431 {
		t.Errorf("Expected review count 52431, got %d", product.ReviewCount)
	}

	// Verify Prime
	if !product.Prime {
		t.Error("Expected Prime to be true")
	}

	// Verify in stock
	if !product.InStock {
		t.Error("Expected InStock to be true")
	}

	// Verify delivery estimate
	if product.DeliveryEstimate != "Tomorrow" {
		t.Errorf("Expected delivery estimate 'Tomorrow', got %s", product.DeliveryEstimate)
	}

	// Verify description
	if product.Description == "" {
		t.Error("Expected description to be populated")
	}

	// Verify features
	if len(product.Features) != 4 {
		t.Errorf("Expected 4 features, got %d", len(product.Features))
	}

	// Verify images
	if len(product.Images) != 2 {
		t.Errorf("Expected 2 images, got %d", len(product.Images))
	}
}

func TestParseProductDetailHTML_MinimalFields(t *testing.T) {
	// Product with only required fields
	html := []byte(`
		<html>
		<body>
			<input type="hidden" name="ASIN" value="B12345TEST">
			<span id="productTitle">Test Product Title</span>
		</body>
		</html>
	`)

	product, err := parseProductDetailHTML(html)
	if err != nil {
		t.Fatalf("parseProductDetailHTML failed: %v", err)
	}

	if product.ASIN != "B12345TEST" {
		t.Errorf("Expected ASIN B12345TEST, got %s", product.ASIN)
	}

	if product.Title != "Test Product Title" {
		t.Errorf("Expected title 'Test Product Title', got %s", product.Title)
	}

	// Optional fields should have zero/default values
	if product.Price != 0 {
		t.Errorf("Expected price 0, got %f", product.Price)
	}

	if product.OriginalPrice != nil {
		t.Errorf("Expected original price to be nil, got %v", product.OriginalPrice)
	}

	if product.Rating != 0 {
		t.Errorf("Expected rating 0, got %f", product.Rating)
	}

	if product.ReviewCount != 0 {
		t.Errorf("Expected review count 0, got %d", product.ReviewCount)
	}

	if product.Prime {
		t.Error("Expected Prime to be false")
	}

	if !product.InStock {
		t.Error("Expected InStock to be true (default)")
	}

	if product.DeliveryEstimate != "" {
		t.Errorf("Expected delivery estimate to be empty, got %s", product.DeliveryEstimate)
	}

	if product.Description != "" {
		t.Errorf("Expected description to be empty, got %s", product.Description)
	}

	if len(product.Features) != 0 {
		t.Errorf("Expected 0 features, got %d", len(product.Features))
	}

	if len(product.Images) != 0 {
		t.Errorf("Expected 0 images, got %d", len(product.Images))
	}
}

func TestParseProductDetailHTML_MissingASIN(t *testing.T) {
	html := []byte(`
		<html>
		<body>
			<span id="productTitle">Product Without ASIN</span>
			<div class="a-price"><span class="a-offscreen">$99.99</span></div>
		</body>
		</html>
	`)

	_, err := parseProductDetailHTML(html)
	if err == nil {
		t.Error("Expected error for missing ASIN, got nil")
	}
}

func TestParseProductDetailHTML_MissingTitle(t *testing.T) {
	html := []byte(`
		<html>
		<body>
			<input type="hidden" name="ASIN" value="B12345TEST">
			<div class="a-price"><span class="a-offscreen">$99.99</span></div>
		</body>
		</html>
	`)

	_, err := parseProductDetailHTML(html)
	if err == nil {
		t.Error("Expected error for missing title, got nil")
	}
}

func TestParseProductDetailHTML_OutOfStock(t *testing.T) {
	tests := []struct {
		name          string
		availText     string
		expectedStock bool
	}{
		{"In Stock", "In Stock", true},
		{"Out of Stock", "Out of Stock", false},
		{"Currently Unavailable", "Currently unavailable", false},
		{"Not Available", "This item is not available", false},
		{"Unavailable", "Unavailable", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html := []byte(`
				<html>
				<body>
					<input type="hidden" name="ASIN" value="B12345TEST">
					<span id="productTitle">Test Product</span>
					<div id="availability">
						<span>` + tt.availText + `</span>
					</div>
				</body>
				</html>
			`)

			product, err := parseProductDetailHTML(html)
			if err != nil {
				t.Fatalf("parseProductDetailHTML failed: %v", err)
			}

			if product.InStock != tt.expectedStock {
				t.Errorf("Expected InStock=%v for %q, got %v", tt.expectedStock, tt.availText, product.InStock)
			}
		})
	}
}

func TestParseProductDetailHTML_ASINExtraction(t *testing.T) {
	tests := []struct {
		name         string
		html         []byte
		expectedASIN string
	}{
		{
			name: "From input field",
			html: []byte(`
				<html>
				<body>
					<input type="hidden" name="ASIN" value="B08N5WRWNW">
					<span id="productTitle">Test Product</span>
				</body>
				</html>
			`),
			expectedASIN: "B08N5WRWNW",
		},
		{
			name: "From canonical link",
			html: []byte(`
				<html>
				<head>
					<link rel="canonical" href="https://www.amazon.com/dp/B08N5WRWNW">
				</head>
				<body>
					<span id="productTitle">Test Product</span>
				</body>
				</html>
			`),
			expectedASIN: "B08N5WRWNW",
		},
		{
			name: "From data-asin attribute",
			html: []byte(`
				<html>
				<body>
					<div data-asin="B08N5WRWNW">
						<span id="productTitle">Test Product</span>
					</div>
				</body>
				</html>
			`),
			expectedASIN: "B08N5WRWNW",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product, err := parseProductDetailHTML(tt.html)
			if err != nil {
				t.Fatalf("parseProductDetailHTML failed: %v", err)
			}

			if product.ASIN != tt.expectedASIN {
				t.Errorf("Expected ASIN %s, got %s", tt.expectedASIN, product.ASIN)
			}
		})
	}
}

func TestParseProductDetailHTML_PriceVariations(t *testing.T) {
	tests := []struct {
		name          string
		priceHTML     string
		expectedPrice float64
	}{
		{
			name: "Standard price format",
			priceHTML: `
				<div class="a-price" data-a-color="price">
					<span class="a-offscreen">$99.99</span>
				</div>
			`,
			expectedPrice: 99.99,
		},
		{
			name: "Priceblock ourprice",
			priceHTML: `
				<span id="priceblock_ourprice">$149.50</span>
			`,
			expectedPrice: 149.50,
		},
		{
			name: "Deal price",
			priceHTML: `
				<span id="priceblock_dealprice">$79.99</span>
			`,
			expectedPrice: 79.99,
		},
		{
			name: "Price whole",
			priceHTML: `
				<span class="a-price-whole">29</span>
			`,
			expectedPrice: 29.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html := []byte(`
				<html>
				<body>
					<input type="hidden" name="ASIN" value="B12345TEST">
					<span id="productTitle">Test Product</span>
					` + tt.priceHTML + `
				</body>
				</html>
			`)

			product, err := parseProductDetailHTML(html)
			if err != nil {
				t.Fatalf("parseProductDetailHTML failed: %v", err)
			}

			if product.Price != tt.expectedPrice {
				t.Errorf("Expected price %f, got %f", tt.expectedPrice, product.Price)
			}
		})
	}
}

func TestParseProductDetailHTML_RatingVariations(t *testing.T) {
	tests := []struct {
		name           string
		ratingHTML     string
		expectedRating float64
	}{
		{
			name: "acrPopover with title",
			ratingHTML: `
				<span id="acrPopover" title="4.7 out of 5 stars">
					<span class="a-icon-alt">4.7 out of 5 stars</span>
				</span>
			`,
			expectedRating: 4.7,
		},
		{
			name: "Icon alt text",
			ratingHTML: `
				<i class="a-icon-star">
					<span class="a-icon-alt">4.5 out of 5 stars</span>
				</i>
			`,
			expectedRating: 4.5,
		},
		{
			name: "Simple decimal",
			ratingHTML: `
				<span class="a-icon-alt">3.8</span>
			`,
			expectedRating: 3.8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html := []byte(`
				<html>
				<body>
					<input type="hidden" name="ASIN" value="B12345TEST">
					<span id="productTitle">Test Product</span>
					` + tt.ratingHTML + `
				</body>
				</html>
			`)

			product, err := parseProductDetailHTML(html)
			if err != nil {
				t.Fatalf("parseProductDetailHTML failed: %v", err)
			}

			if product.Rating != tt.expectedRating {
				t.Errorf("Expected rating %f, got %f", tt.expectedRating, product.Rating)
			}
		})
	}
}

func TestParseProductDetailHTML_PrimeIndicators(t *testing.T) {
	tests := []struct {
		name      string
		primeHTML string
		isPrime   bool
	}{
		{
			name: "Icon prime in price badging",
			primeHTML: `
				<div id="priceBadging_feature_div">
					<i class="a-icon-prime"></i>
				</div>
			`,
			isPrime: true,
		},
		{
			name: "Generic icon prime",
			primeHTML: `
				<i class="a-icon-prime"></i>
			`,
			isPrime: true,
		},
		{
			name: "Prime badge span",
			primeHTML: `
				<span class="prime-badge">Prime</span>
			`,
			isPrime: true,
		},
		{
			name: "Aria label Prime",
			primeHTML: `
				<span aria-label="Amazon Prime">Prime eligible</span>
			`,
			isPrime: true,
		},
		{
			name:      "No prime indicator",
			primeHTML: ``,
			isPrime:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html := []byte(`
				<html>
				<body>
					<input type="hidden" name="ASIN" value="B12345TEST">
					<span id="productTitle">Test Product</span>
					` + tt.primeHTML + `
				</body>
				</html>
			`)

			product, err := parseProductDetailHTML(html)
			if err != nil {
				t.Fatalf("parseProductDetailHTML failed: %v", err)
			}

			if product.Prime != tt.isPrime {
				t.Errorf("Expected Prime=%v, got %v", tt.isPrime, product.Prime)
			}
		})
	}
}

func TestParseProductDetailHTML_FeaturesExtraction(t *testing.T) {
	html := []byte(`
		<html>
		<body>
			<input type="hidden" name="ASIN" value="B12345TEST">
			<span id="productTitle">Test Product</span>
			<div id="feature-bullets">
				<ul class="a-unordered-list">
					<li><span class="a-list-item">Feature 1</span></li>
					<li><span class="a-list-item">Feature 2</span></li>
					<li><span class="a-list-item">Feature 3</span></li>
					<li><span class="a-list-item">See more product details</span></li>
				</ul>
			</div>
		</body>
		</html>
	`)

	product, err := parseProductDetailHTML(html)
	if err != nil {
		t.Fatalf("parseProductDetailHTML failed: %v", err)
	}

	// Should have 3 features (excluding "See more product details" which starts with "See more")
	expectedFeatures := 3
	if len(product.Features) != expectedFeatures {
		t.Errorf("Expected %d features, got %d", expectedFeatures, len(product.Features))
	}

	// Verify feature content
	if len(product.Features) > 0 && product.Features[0] != "Feature 1" {
		t.Errorf("Expected first feature to be 'Feature 1', got %s", product.Features[0])
	}
}

func TestParseProductDetailHTML_ImagesDeduplication(t *testing.T) {
	html := []byte(`
		<html>
		<body>
			<input type="hidden" name="ASIN" value="B12345TEST">
			<span id="productTitle">Test Product</span>
			<div id="altImages">
				<ul>
					<li class="imageThumbnail">
						<img src="https://example.com/image1.jpg"
							 data-old-hires="https://example.com/image1-hires.jpg">
					</li>
					<li class="imageThumbnail">
						<img src="https://example.com/image1.jpg"
							 data-old-hires="https://example.com/image1-hires.jpg">
					</li>
					<li class="imageThumbnail">
						<img src="https://example.com/image2.jpg"
							 data-old-hires="https://example.com/image2-hires.jpg">
					</li>
				</ul>
			</div>
		</body>
		</html>
	`)

	product, err := parseProductDetailHTML(html)
	if err != nil {
		t.Fatalf("parseProductDetailHTML failed: %v", err)
	}

	// Should deduplicate images, expecting 2 unique images
	if len(product.Images) != 2 {
		t.Errorf("Expected 2 unique images, got %d", len(product.Images))
	}
}

func TestParseProductDetailHTML_ImageFiltering(t *testing.T) {
	html := []byte(`
		<html>
		<body>
			<input type="hidden" name="ASIN" value="B12345TEST">
			<span id="productTitle">Test Product</span>
			<div id="altImages">
				<ul>
					<li class="imageThumbnail">
						<img src="https://example.com/1x1.gif">
					</li>
					<li class="imageThumbnail">
						<img src="https://example.com/pixel.png">
					</li>
					<li class="imageThumbnail">
						<img src="https://example.com/transparent.gif">
					</li>
					<li class="imageThumbnail">
						<img src="https://example.com/valid-image.jpg">
					</li>
				</ul>
			</div>
		</body>
		</html>
	`)

	product, err := parseProductDetailHTML(html)
	if err != nil {
		t.Fatalf("parseProductDetailHTML failed: %v", err)
	}

	// Should filter out placeholder images, expecting only 1 valid image
	if len(product.Images) != 1 {
		t.Errorf("Expected 1 valid image, got %d", len(product.Images))
	}

	if len(product.Images) > 0 && product.Images[0] != "https://example.com/valid-image.jpg" {
		t.Errorf("Expected valid-image.jpg, got %s", product.Images[0])
	}
}

func TestParseProductDetailHTML_EmptyHTML(t *testing.T) {
	html := []byte(`<html><body></body></html>`)

	_, err := parseProductDetailHTML(html)
	if err == nil {
		t.Error("Expected error for HTML without required fields, got nil")
	}
}

func TestParseProductDetailHTML_InvalidHTML(t *testing.T) {
	html := []byte(`not valid html at all`)

	_, err := parseProductDetailHTML(html)
	if err == nil {
		t.Error("Expected error for invalid HTML, got nil")
	}
}

func TestParseProductDetailHTML_DescriptionExtraction(t *testing.T) {
	html := []byte(`
		<html>
		<body>
			<input type="hidden" name="ASIN" value="B12345TEST">
			<span id="productTitle">Test Product</span>
			<div id="productDescription">
				<p>This is the first paragraph of the description.</p>
				<p>This is the second paragraph.</p>
			</div>
		</body>
		</html>
	`)

	product, err := parseProductDetailHTML(html)
	if err != nil {
		t.Fatalf("parseProductDetailHTML failed: %v", err)
	}

	if product.Description == "" {
		t.Error("Expected description to be populated")
	}

	// Description should contain both paragraphs
	expectedSubstring := "first paragraph"
	if !stringContains(product.Description, expectedSubstring) {
		t.Errorf("Expected description to contain %q", expectedSubstring)
	}
}

func TestParseProductDetailHTML_DeliveryEstimate(t *testing.T) {
	tests := []struct {
		name             string
		deliveryHTML     string
		expectedDelivery string
	}{
		{
			name: "Tomorrow delivery",
			deliveryHTML: `
				<div id="deliveryMessageMirId">
					<span class="a-text-bold">Tomorrow</span>
				</div>
			`,
			expectedDelivery: "Tomorrow",
		},
		{
			name: "Specific date",
			deliveryHTML: `
				<div id="mir-layout-DELIVERY_BLOCK">
					<span class="a-text-bold">Friday, Jan 24</span>
				</div>
			`,
			expectedDelivery: "Friday, Jan 24",
		},
		{
			name: "Feature name delivery",
			deliveryHTML: `
				<div data-feature-name="deliveryMessage">
					<span>Arrives Mon, Jan 27</span>
				</div>
			`,
			expectedDelivery: "Arrives Mon, Jan 27",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html := []byte(`
				<html>
				<body>
					<input type="hidden" name="ASIN" value="B12345TEST">
					<span id="productTitle">Test Product</span>
					` + tt.deliveryHTML + `
				</body>
				</html>
			`)

			product, err := parseProductDetailHTML(html)
			if err != nil {
				t.Fatalf("parseProductDetailHTML failed: %v", err)
			}

			if product.DeliveryEstimate != tt.expectedDelivery {
				t.Errorf("Expected delivery %q, got %q", tt.expectedDelivery, product.DeliveryEstimate)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestGetProduct_EmptyASIN(t *testing.T) {
	client := NewClient()
	_, err := client.GetProduct("")
	if err == nil {
		t.Error("Expected error for empty ASIN, got nil")
	}
	expectedMsg := "ASIN cannot be empty"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message %q, got %q", expectedMsg, err.Error())
	}
}

func TestGetProduct_InvalidASINFormat(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name string
		asin string
	}{
		{"Too short", "B08N5"},
		{"Too long", "B08N5WRWNW123"},
		{"Contains lowercase", "b08n5wrwnw"},
		{"Contains special chars", "B08N5WRWN!"},
		{"Contains spaces", "B08N5 WRWN"},
		{"Only 9 chars", "B08N5WRWN"},
		{"Only 11 chars", "B08N5WRWNW1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.GetProduct(tt.asin)
			if err == nil {
				t.Errorf("Expected error for invalid ASIN %q, got nil", tt.asin)
			}
			expectedSubstring := "invalid ASIN format"
			if !stringContains(err.Error(), expectedSubstring) {
				t.Errorf("Expected error to contain %q, got %q", expectedSubstring, err.Error())
			}
		})
	}
}

func TestGetProduct_ValidASINFormat(t *testing.T) {
	tests := []struct {
		name string
		asin string
	}{
		{"Standard ASIN", "B08N5WRWNW"},
		{"All uppercase letters", "ABCDEFGHIJ"},
		{"All numbers", "1234567890"},
		{"Mixed alphanumeric", "A1B2C3D4E5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We can't test the full flow without mocking HTTP, but we can verify
			// the ASIN validation passes and we get a network/parsing error
			client := NewClient()
			_, err := client.GetProduct(tt.asin)

			// Should not be an ASIN format error
			if err != nil && stringContains(err.Error(), "invalid ASIN format") {
				t.Errorf("ASIN %q should be valid format, got error: %v", tt.asin, err)
			}
		})
	}
}
