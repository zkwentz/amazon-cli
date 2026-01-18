package amazon

import (
	"testing"

	"github.com/michaelshimeles/amazon-cli/pkg/models"
)

func TestIsValidFrequency(t *testing.T) {
	tests := []struct {
		name     string
		weeks    int
		expected bool
	}{
		{"Valid 1 week", 1, true},
		{"Valid 2 weeks", 2, true},
		{"Valid 3 weeks", 3, true},
		{"Valid 4 weeks", 4, true},
		{"Valid 5 weeks", 5, true},
		{"Valid 6 weeks", 6, true},
		{"Valid 8 weeks", 8, true},
		{"Valid 10 weeks", 10, true},
		{"Valid 12 weeks", 12, true},
		{"Valid 16 weeks", 16, true},
		{"Valid 20 weeks", 20, true},
		{"Valid 24 weeks", 24, true},
		{"Valid 26 weeks", 26, true},
		{"Invalid 0 weeks", 0, false},
		{"Invalid 7 weeks", 7, false},
		{"Invalid 9 weeks", 9, false},
		{"Invalid 11 weeks", 11, false},
		{"Invalid 13 weeks", 13, false},
		{"Invalid 15 weeks", 15, false},
		{"Invalid 25 weeks", 25, false},
		{"Invalid 27 weeks", 27, false},
		{"Invalid 30 weeks", 30, false},
		{"Invalid negative", -1, false},
		{"Invalid 100 weeks", 100, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidFrequency(tt.weeks)
			if result != tt.expected {
				t.Errorf("IsValidFrequency(%d) = %v, want %v", tt.weeks, result, tt.expected)
			}
		})
	}
}

func TestUpdateFrequency_Success(t *testing.T) {
	client := NewSubscriptionClient()

	tests := []struct {
		name           string
		subscriptionID string
		weeks          int
	}{
		{"Update to 1 week", "S01-1234567-8901234", 1},
		{"Update to 2 weeks", "S01-1234567-8901234", 2},
		{"Update to 4 weeks", "S01-1234567-8901234", 4},
		{"Update to 8 weeks", "S01-1234567-8901234", 8},
		{"Update to 12 weeks", "S01-1234567-8901234", 12},
		{"Update to 26 weeks", "S01-1234567-8901234", 26},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := client.UpdateFrequency(tt.subscriptionID, tt.weeks)
			if err != nil {
				t.Fatalf("UpdateFrequency returned unexpected error: %v", err)
			}
			if result == nil {
				t.Fatal("UpdateFrequency returned nil subscription")
			}
			if result.FrequencyWeeks != tt.weeks {
				t.Errorf("FrequencyWeeks = %d, want %d", result.FrequencyWeeks, tt.weeks)
			}
			if result.SubscriptionID != tt.subscriptionID {
				t.Errorf("SubscriptionID = %s, want %s", result.SubscriptionID, tt.subscriptionID)
			}
			if result.Status != "active" {
				t.Errorf("Status = %s, want active", result.Status)
			}
		})
	}
}

func TestUpdateFrequency_EmptySubscriptionID(t *testing.T) {
	client := NewSubscriptionClient()

	_, err := client.UpdateFrequency("", 4)
	if err == nil {
		t.Fatal("Expected error for empty subscription ID, got nil")
	}

	cliErr, ok := err.(*models.CLIError)
	if !ok {
		t.Fatalf("Expected CLIError, got %T", err)
	}

	if cliErr.Code != models.ErrorCodeInvalidInput {
		t.Errorf("Error code = %s, want %s", cliErr.Code, models.ErrorCodeInvalidInput)
	}

	if cliErr.Message != "subscription ID cannot be empty" {
		t.Errorf("Error message = %s, want 'subscription ID cannot be empty'", cliErr.Message)
	}
}

func TestUpdateFrequency_FrequencyTooLow(t *testing.T) {
	client := NewSubscriptionClient()

	tests := []struct {
		name  string
		weeks int
	}{
		{"Zero weeks", 0},
		{"Negative weeks", -1},
		{"Negative 10 weeks", -10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.UpdateFrequency("S01-1234567-8901234", tt.weeks)
			if err == nil {
				t.Fatal("Expected error for frequency too low, got nil")
			}

			cliErr, ok := err.(*models.CLIError)
			if !ok {
				t.Fatalf("Expected CLIError, got %T", err)
			}

			if cliErr.Code != models.ErrorCodeInvalidInput {
				t.Errorf("Error code = %s, want %s", cliErr.Code, models.ErrorCodeInvalidInput)
			}

			if cliErr.Message != "frequency must be at least 1 week" {
				t.Errorf("Error message = %s, want 'frequency must be at least 1 week'", cliErr.Message)
			}

			// Check that details contain the provided weeks
			if details, ok := cliErr.Details["provided_weeks"].(int); !ok || details != tt.weeks {
				t.Errorf("Details provided_weeks = %v, want %d", cliErr.Details["provided_weeks"], tt.weeks)
			}

			// Check that details contain minimum weeks
			if min, ok := cliErr.Details["minimum_weeks"].(int); !ok || min != 1 {
				t.Errorf("Details minimum_weeks = %v, want 1", cliErr.Details["minimum_weeks"])
			}
		})
	}
}

func TestUpdateFrequency_FrequencyTooHigh(t *testing.T) {
	client := NewSubscriptionClient()

	tests := []struct {
		name  string
		weeks int
	}{
		{"27 weeks", 27},
		{"30 weeks", 30},
		{"52 weeks", 52},
		{"100 weeks", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.UpdateFrequency("S01-1234567-8901234", tt.weeks)
			if err == nil {
				t.Fatal("Expected error for frequency too high, got nil")
			}

			cliErr, ok := err.(*models.CLIError)
			if !ok {
				t.Fatalf("Expected CLIError, got %T", err)
			}

			if cliErr.Code != models.ErrorCodeInvalidInput {
				t.Errorf("Error code = %s, want %s", cliErr.Code, models.ErrorCodeInvalidInput)
			}

			if cliErr.Message != "frequency cannot exceed 26 weeks" {
				t.Errorf("Error message = %s, want 'frequency cannot exceed 26 weeks'", cliErr.Message)
			}

			// Check that details contain the provided weeks
			if details, ok := cliErr.Details["provided_weeks"].(int); !ok || details != tt.weeks {
				t.Errorf("Details provided_weeks = %v, want %d", cliErr.Details["provided_weeks"], tt.weeks)
			}

			// Check that details contain maximum weeks
			if max, ok := cliErr.Details["maximum_weeks"].(int); !ok || max != 26 {
				t.Errorf("Details maximum_weeks = %v, want 26", cliErr.Details["maximum_weeks"])
			}
		})
	}
}

func TestUpdateFrequency_InvalidFrequencyInterval(t *testing.T) {
	client := NewSubscriptionClient()

	// Test frequencies that are within range (1-26) but not in the valid list
	tests := []struct {
		name  string
		weeks int
	}{
		{"7 weeks - not supported", 7},
		{"9 weeks - not supported", 9},
		{"11 weeks - not supported", 11},
		{"13 weeks - not supported", 13},
		{"14 weeks - not supported", 14},
		{"15 weeks - not supported", 15},
		{"17 weeks - not supported", 17},
		{"18 weeks - not supported", 18},
		{"19 weeks - not supported", 19},
		{"21 weeks - not supported", 21},
		{"22 weeks - not supported", 22},
		{"23 weeks - not supported", 23},
		{"25 weeks - not supported", 25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.UpdateFrequency("S01-1234567-8901234", tt.weeks)
			if err == nil {
				t.Fatalf("Expected error for invalid frequency interval %d, got nil", tt.weeks)
			}

			cliErr, ok := err.(*models.CLIError)
			if !ok {
				t.Fatalf("Expected CLIError, got %T", err)
			}

			if cliErr.Code != models.ErrorCodeInvalidInput {
				t.Errorf("Error code = %s, want %s", cliErr.Code, models.ErrorCodeInvalidInput)
			}

			expectedMsg := "frequency must be one of the supported intervals: [1 2 3 4 5 6 8 10 12 16 20 24 26] weeks"
			if cliErr.Message != expectedMsg {
				t.Errorf("Error message = %s, want %s", cliErr.Message, expectedMsg)
			}

			// Check that details contain the provided weeks
			if details, ok := cliErr.Details["provided_weeks"].(int); !ok || details != tt.weeks {
				t.Errorf("Details provided_weeks = %v, want %d", cliErr.Details["provided_weeks"], tt.weeks)
			}

			// Check that details contain valid frequencies
			if _, ok := cliErr.Details["valid_frequencies"]; !ok {
				t.Error("Details should contain valid_frequencies")
			}
		})
	}
}

func TestUpdateFrequency_BoundaryConditions(t *testing.T) {
	client := NewSubscriptionClient()

	// Test exact boundaries
	tests := []struct {
		name        string
		weeks       int
		shouldError bool
		errorMsg    string
	}{
		{"Minimum valid (1 week)", 1, false, ""},
		{"Just below minimum (0 weeks)", 0, true, "frequency must be at least 1 week"},
		{"Maximum valid (26 weeks)", 26, false, ""},
		{"Just above maximum (27 weeks)", 27, true, "frequency cannot exceed 26 weeks"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := client.UpdateFrequency("S01-1234567-8901234", tt.weeks)

			if tt.shouldError {
				if err == nil {
					t.Fatalf("Expected error, got nil")
				}
				cliErr, ok := err.(*models.CLIError)
				if !ok {
					t.Fatalf("Expected CLIError, got %T", err)
				}
				if cliErr.Message != tt.errorMsg {
					t.Errorf("Error message = %s, want %s", cliErr.Message, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}
				if result.FrequencyWeeks != tt.weeks {
					t.Errorf("FrequencyWeeks = %d, want %d", result.FrequencyWeeks, tt.weeks)
				}
			}
		})
	}
}

func TestUpdateFrequency_AllValidFrequencies(t *testing.T) {
	client := NewSubscriptionClient()

	// Test that all frequencies in ValidFrequencies work
	for _, weeks := range ValidFrequencies {
		t.Run(string(rune(weeks))+" weeks", func(t *testing.T) {
			result, err := client.UpdateFrequency("S01-1234567-8901234", weeks)
			if err != nil {
				t.Fatalf("UpdateFrequency(%d) returned unexpected error: %v", weeks, err)
			}
			if result == nil {
				t.Fatal("UpdateFrequency returned nil subscription")
			}
			if result.FrequencyWeeks != weeks {
				t.Errorf("FrequencyWeeks = %d, want %d", result.FrequencyWeeks, weeks)
			}
		})
	}
}

func TestGetSubscriptions(t *testing.T) {
	client := NewSubscriptionClient()

	result, err := client.GetSubscriptions()
	if err != nil {
		t.Fatalf("GetSubscriptions returned unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("GetSubscriptions returned nil")
	}

	if len(result.Subscriptions) == 0 {
		t.Error("Expected at least one subscription in mock response")
	}
}

func TestGetSubscription(t *testing.T) {
	client := NewSubscriptionClient()

	result, err := client.GetSubscription("S01-1234567-8901234")
	if err != nil {
		t.Fatalf("GetSubscription returned unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("GetSubscription returned nil")
	}

	if result.SubscriptionID != "S01-1234567-8901234" {
		t.Errorf("SubscriptionID = %s, want S01-1234567-8901234", result.SubscriptionID)
	}
}

func TestGetSubscription_EmptyID(t *testing.T) {
	client := NewSubscriptionClient()

	_, err := client.GetSubscription("")
	if err == nil {
		t.Fatal("Expected error for empty subscription ID, got nil")
	}

	cliErr, ok := err.(*models.CLIError)
	if !ok {
		t.Fatalf("Expected CLIError, got %T", err)
	}

	if cliErr.Code != models.ErrorCodeInvalidInput {
		t.Errorf("Error code = %s, want %s", cliErr.Code, models.ErrorCodeInvalidInput)
	}
}

func TestSkipDelivery(t *testing.T) {
	client := NewSubscriptionClient()

	result, err := client.SkipDelivery("S01-1234567-8901234")
	if err != nil {
		t.Fatalf("SkipDelivery returned unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("SkipDelivery returned nil")
	}

	if result.Status != "active" {
		t.Errorf("Status = %s, want active", result.Status)
	}
}

func TestSkipDelivery_EmptyID(t *testing.T) {
	client := NewSubscriptionClient()

	_, err := client.SkipDelivery("")
	if err == nil {
		t.Fatal("Expected error for empty subscription ID, got nil")
	}

	cliErr, ok := err.(*models.CLIError)
	if !ok {
		t.Fatalf("Expected CLIError, got %T", err)
	}

	if cliErr.Code != models.ErrorCodeInvalidInput {
		t.Errorf("Error code = %s, want %s", cliErr.Code, models.ErrorCodeInvalidInput)
	}
}

func TestCancelSubscription(t *testing.T) {
	client := NewSubscriptionClient()

	result, err := client.CancelSubscription("S01-1234567-8901234")
	if err != nil {
		t.Fatalf("CancelSubscription returned unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("CancelSubscription returned nil")
	}

	if result.Status != "cancelled" {
		t.Errorf("Status = %s, want cancelled", result.Status)
	}
}

func TestCancelSubscription_EmptyID(t *testing.T) {
	client := NewSubscriptionClient()

	_, err := client.CancelSubscription("")
	if err == nil {
		t.Fatal("Expected error for empty subscription ID, got nil")
	}

	cliErr, ok := err.(*models.CLIError)
	if !ok {
		t.Fatalf("Expected CLIError, got %T", err)
	}

	if cliErr.Code != models.ErrorCodeInvalidInput {
		t.Errorf("Error code = %s, want %s", cliErr.Code, models.ErrorCodeInvalidInput)
	}
}

func TestGetUpcomingDeliveries(t *testing.T) {
	client := NewSubscriptionClient()

	result, err := client.GetUpcomingDeliveries()
	if err != nil {
		t.Fatalf("GetUpcomingDeliveries returned unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("GetUpcomingDeliveries returned nil")
	}

	if len(result) == 0 {
		t.Error("Expected at least one upcoming delivery in mock response")
	}
}
