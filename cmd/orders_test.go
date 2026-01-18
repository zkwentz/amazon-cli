package cmd

import (
	"testing"
	"time"

	"github.com/zkwentz/amazon-cli/internal/amazon"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

// TestOrdersHistoryCommandLogic tests the core logic of the orders history command
func TestOrdersHistoryCommandLogic(t *testing.T) {
	// Test that the Amazon client can be created
	client, err := amazon.NewClient()
	if err != nil {
		t.Fatalf("Failed to create Amazon client: %v", err)
	}

	// Test with current year
	currentYear := time.Now().Year()
	response, err := client.GetOrderHistory(currentYear)
	if err != nil {
		t.Errorf("Failed to get order history for current year: %v", err)
	}

	if response == nil {
		t.Error("Expected response but got nil")
	}

	if response.Orders == nil {
		t.Error("Expected orders slice but got nil")
	}

	if response.TotalCount != len(response.Orders) {
		t.Errorf("TotalCount (%d) doesn't match orders length (%d)",
			response.TotalCount, len(response.Orders))
	}
}

// TestOrdersHistoryInvalidYear tests error handling for invalid years
func TestOrdersHistoryInvalidYear(t *testing.T) {
	client, err := amazon.NewClient()
	if err != nil {
		t.Fatalf("Failed to create Amazon client: %v", err)
	}

	// Test with invalid year
	_, err = client.GetOrderHistory(1990)
	if err == nil {
		t.Error("Expected error for invalid year but got none")
		return
	}

	cliErr, ok := err.(*models.CLIError)
	if !ok {
		t.Errorf("Expected CLIError but got: %T", err)
		return
	}

	if cliErr.Code != models.ErrorCodeInvalidInput {
		t.Errorf("Expected error code %s but got %s",
			models.ErrorCodeInvalidInput, cliErr.Code)
	}
}

// TestOrdersHistoryValidYears tests various valid years
func TestOrdersHistoryValidYears(t *testing.T) {
	client, err := amazon.NewClient()
	if err != nil {
		t.Fatalf("Failed to create Amazon client: %v", err)
	}

	validYears := []int{1995, 2000, 2020, 2023, time.Now().Year()}

	for _, year := range validYears {
		t.Run(string(rune(year/1000)+'0')+string(rune((year/100)%10)+'0')+string(rune((year/10)%10)+'0')+string(rune(year%10)+'0'), func(t *testing.T) {
			response, err := client.GetOrderHistory(year)
			if err != nil {
				t.Errorf("Failed to get order history for year %d: %v", year, err)
				return
			}

			if response == nil {
				t.Error("Expected response but got nil")
			}
		})
	}
}
