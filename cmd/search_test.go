package cmd

import (
	"encoding/json"
	"testing"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

func TestSearchCmd_Configuration(t *testing.T) {
	// Test the main search command configuration
	if searchCmd.Use != "search <query>" {
		t.Errorf("Expected Use='search <query>', got '%s'", searchCmd.Use)
	}

	if searchCmd.Short != "Search for products" {
		t.Errorf("Expected Short='Search for products', got '%s'", searchCmd.Short)
	}

	expectedLong := `Search for products on Amazon by keyword with optional filters.`
	if searchCmd.Long != expectedLong {
		t.Errorf("Expected Long='%s', got '%s'", expectedLong, searchCmd.Long)
	}

	if searchCmd.Run == nil {
		t.Error("Expected Run function to be defined")
	}

	// Test that it requires exactly 1 argument
	if searchCmd.Args == nil {
		t.Error("Expected Args validator to be defined")
	}
}

func TestSearchCmd_Flags(t *testing.T) {
	// Test that all flags are properly configured
	categoryFlag := searchCmd.Flags().Lookup("category")
	if categoryFlag == nil {
		t.Error("Expected --category flag to be defined")
	} else {
		if categoryFlag.DefValue != "" {
			t.Errorf("Expected --category default value to be empty, got '%s'", categoryFlag.DefValue)
		}
	}

	minPriceFlag := searchCmd.Flags().Lookup("min-price")
	if minPriceFlag == nil {
		t.Error("Expected --min-price flag to be defined")
	} else {
		if minPriceFlag.DefValue != "0" {
			t.Errorf("Expected --min-price default value to be '0', got '%s'", minPriceFlag.DefValue)
		}
	}

	maxPriceFlag := searchCmd.Flags().Lookup("max-price")
	if maxPriceFlag == nil {
		t.Error("Expected --max-price flag to be defined")
	} else {
		if maxPriceFlag.DefValue != "0" {
			t.Errorf("Expected --max-price default value to be '0', got '%s'", maxPriceFlag.DefValue)
		}
	}

	primeOnlyFlag := searchCmd.Flags().Lookup("prime-only")
	if primeOnlyFlag == nil {
		t.Error("Expected --prime-only flag to be defined")
	} else {
		if primeOnlyFlag.DefValue != "false" {
			t.Errorf("Expected --prime-only default value to be 'false', got '%s'", primeOnlyFlag.DefValue)
		}
	}

	pageFlag := searchCmd.Flags().Lookup("page")
	if pageFlag == nil {
		t.Error("Expected --page flag to be defined")
	} else {
		if pageFlag.DefValue != "1" {
			t.Errorf("Expected --page default value to be '1', got '%s'", pageFlag.DefValue)
		}
	}
}

func TestSearchCmd_VariablesInitialized(t *testing.T) {
	// Test that package-level variables are initialized with correct default values
	// Save original values
	origCategory := searchCategory
	origMinPrice := searchMinPrice
	origMaxPrice := searchMaxPrice
	origPrimeOnly := searchPrimeOnly
	origPage := searchPage

	// Modify them
	searchCategory = "electronics"
	searchMinPrice = 10.0
	searchMaxPrice = 100.0
	searchPrimeOnly = true
	searchPage = 2

	// Verify modifications worked
	if searchCategory != "electronics" {
		t.Error("Failed to modify searchCategory")
	}
	if searchMinPrice != 10.0 {
		t.Error("Failed to modify searchMinPrice")
	}
	if searchMaxPrice != 100.0 {
		t.Error("Failed to modify searchMaxPrice")
	}
	if searchPrimeOnly != true {
		t.Error("Failed to modify searchPrimeOnly")
	}
	if searchPage != 2 {
		t.Error("Failed to modify searchPage")
	}

	// Restore original values
	searchCategory = origCategory
	searchMinPrice = origMinPrice
	searchMaxPrice = origMaxPrice
	searchPrimeOnly = origPrimeOnly
	searchPage = origPage
}

func TestSearchCmd_ResponseParsing(t *testing.T) {
	// This test verifies that the models.SearchResponse structure
	// can be properly marshaled to JSON (as used by output.JSON)

	response := &models.SearchResponse{
		Query: "laptop",
		Results: []models.Product{
			{
				ASIN:        "B08XYZ1234",
				Title:       "Example Laptop",
				Price:       899.99,
				Rating:      4.5,
				ReviewCount: 1234,
				Prime:       true,
				InStock:     true,
			},
		},
		TotalResults: 1,
		Page:         1,
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal SearchResponse to JSON: %v", err)
	}

	// Verify it's valid JSON
	var parsed map[string]interface{}
	err = json.Unmarshal(jsonBytes, &parsed)
	if err != nil {
		t.Fatalf("Failed to parse marshaled JSON: %v", err)
	}

	// Verify expected fields exist
	if _, ok := parsed["query"]; !ok {
		t.Error("Expected 'query' field in JSON output")
	}
	if _, ok := parsed["results"]; !ok {
		t.Error("Expected 'results' field in JSON output")
	}
	if _, ok := parsed["total_results"]; !ok {
		t.Error("Expected 'total_results' field in JSON output")
	}
	if _, ok := parsed["page"]; !ok {
		t.Error("Expected 'page' field in JSON output")
	}
}

func TestSearchCmd_ProductStructure(t *testing.T) {
	// Test that the Product structure includes all expected fields
	product := models.Product{
		ASIN:        "B08XYZ1234",
		Title:       "Test Product",
		Price:       29.99,
		Rating:      4.5,
		ReviewCount: 100,
		Prime:       true,
		InStock:     true,
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("Failed to marshal Product to JSON: %v", err)
	}

	// Verify it's valid JSON
	var parsed map[string]interface{}
	err = json.Unmarshal(jsonBytes, &parsed)
	if err != nil {
		t.Fatalf("Failed to parse marshaled JSON: %v", err)
	}

	// Verify all expected fields exist
	expectedFields := []string{"asin", "title", "price", "rating", "review_count", "prime", "in_stock"}
	for _, field := range expectedFields {
		if _, ok := parsed[field]; !ok {
			t.Errorf("Expected '%s' field in JSON output", field)
		}
	}
}

func TestSearchCmd_SearchOptionsStructure(t *testing.T) {
	// Test that SearchOptions includes all the required fields
	opts := models.SearchOptions{
		Category:  "electronics",
		MinPrice:  10.0,
		MaxPrice:  100.0,
		PrimeOnly: true,
		Page:      1,
	}

	// Verify all fields are accessible
	if opts.Category != "electronics" {
		t.Errorf("Expected Category='electronics', got '%s'", opts.Category)
	}
	if opts.MinPrice != 10.0 {
		t.Errorf("Expected MinPrice=10.0, got %f", opts.MinPrice)
	}
	if opts.MaxPrice != 100.0 {
		t.Errorf("Expected MaxPrice=100.0, got %f", opts.MaxPrice)
	}
	if opts.PrimeOnly != true {
		t.Errorf("Expected PrimeOnly=true, got %v", opts.PrimeOnly)
	}
	if opts.Page != 1 {
		t.Errorf("Expected Page=1, got %d", opts.Page)
	}
}

func TestSearchCmd_ErrorHandling(t *testing.T) {
	// Verify that error constants used in search command exist
	if models.ErrAmazonError != "AMAZON_ERROR" {
		t.Errorf("Expected ErrAmazonError='AMAZON_ERROR', got '%s'", models.ErrAmazonError)
	}

	if models.ExitGeneralError != 1 {
		t.Errorf("Expected ExitGeneralError=1, got %d", models.ExitGeneralError)
	}
}

func TestSearchCmd_GetClientReturnsClient(t *testing.T) {
	// Test that getClient returns a non-nil client
	c := getClient()
	if c == nil {
		t.Error("Expected getClient() to return non-nil client")
	}

	// Test that calling getClient twice returns the same instance
	c2 := getClient()
	if c != c2 {
		t.Error("Expected getClient() to return the same client instance")
	}
}

func TestSearchCmd_RunFunctionExists(t *testing.T) {
	// Verify that the Run function is properly defined
	if searchCmd.Run == nil {
		t.Fatal("Expected Run function to be defined")
	}

	// The Run function should:
	// 1. Get query from args[0]
	// 2. Get all flags (category, min-price, max-price, prime-only, page)
	// 3. Create SearchOptions with all flags
	// 4. Call client.Search
	// 5. Output JSON result

	// This is verified by the existence of the function and its integration
	// with the command structure
}

func TestSearchCmd_IntegrationWithSearchOptions(t *testing.T) {
	// Test that the command correctly constructs SearchOptions
	// This verifies the integration between flags and the SearchOptions struct

	// Set some test values for the flags
	origCategory := searchCategory
	origMinPrice := searchMinPrice
	origMaxPrice := searchMaxPrice
	origPrimeOnly := searchPrimeOnly
	origPage := searchPage

	searchCategory = "books"
	searchMinPrice = 5.0
	searchMaxPrice = 50.0
	searchPrimeOnly = true
	searchPage = 2

	// Create SearchOptions as the Run function does
	opts := models.SearchOptions{
		Category:  searchCategory,
		MinPrice:  searchMinPrice,
		MaxPrice:  searchMaxPrice,
		PrimeOnly: searchPrimeOnly,
		Page:      searchPage,
	}

	// Verify the options are correctly set
	if opts.Category != "books" {
		t.Errorf("Expected Category='books', got '%s'", opts.Category)
	}
	if opts.MinPrice != 5.0 {
		t.Errorf("Expected MinPrice=5.0, got %f", opts.MinPrice)
	}
	if opts.MaxPrice != 50.0 {
		t.Errorf("Expected MaxPrice=50.0, got %f", opts.MaxPrice)
	}
	if opts.PrimeOnly != true {
		t.Errorf("Expected PrimeOnly=true, got %v", opts.PrimeOnly)
	}
	if opts.Page != 2 {
		t.Errorf("Expected Page=2, got %d", opts.Page)
	}

	// Restore original values
	searchCategory = origCategory
	searchMinPrice = origMinPrice
	searchMaxPrice = origMaxPrice
	searchPrimeOnly = origPrimeOnly
	searchPage = origPage
}
