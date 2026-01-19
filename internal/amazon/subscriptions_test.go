package amazon

import (
	"testing"
	"time"
)

func TestSkipDelivery(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		wantErr     bool
		errContains string
	}{
		{
			name:        "valid subscription ID",
			id:          "sub123",
			wantErr:     false,
			errContains: "",
		},
		{
			name:        "empty subscription ID should fail",
			id:          "",
			wantErr:     true,
			errContains: "subscription ID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()

			// Get the current time for comparison
			beforeSkip := time.Now()

			subscription, err := client.SkipDelivery(tt.id)

			if tt.wantErr {
				if err == nil {
					t.Error("SkipDelivery() expected error but got none")
					return
				}
				if err.Error() != tt.errContains {
					t.Errorf("SkipDelivery() error = %v, want %v", err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("SkipDelivery() unexpected error: %v", err)
				return
			}

			if subscription == nil {
				t.Error("SkipDelivery() returned nil subscription")
				return
			}

			// Verify subscription ID matches
			if subscription.ID != tt.id {
				t.Errorf("SkipDelivery() subscription.ID = %v, want %v", subscription.ID, tt.id)
			}

			// Verify NextDelivery was advanced
			// NextDelivery should be at least FrequencyWeeks in the future from now
			expectedMinDelivery := beforeSkip.AddDate(0, 0, subscription.FrequencyWeeks*7)
			if subscription.NextDelivery.Before(expectedMinDelivery) {
				t.Errorf("SkipDelivery() NextDelivery = %v, expected at least %v", subscription.NextDelivery, expectedMinDelivery)
			}

			// Verify status is still active
			if subscription.Status != "active" {
				t.Errorf("SkipDelivery() status = %v, want active", subscription.Status)
			}
		})
	}
}

func TestSkipDelivery_AdvancesDeliveryByFrequencyWeeks(t *testing.T) {
	client := NewClient()

	subscription, err := client.SkipDelivery("sub123")
	if err != nil {
		t.Fatalf("SkipDelivery() unexpected error: %v", err)
	}

	// Verify that NextDelivery is at least FrequencyWeeks * 7 days from now
	// The mock implementation sets initial NextDelivery to 14 days from now,
	// then advances it by FrequencyWeeks (28 days), so total is ~42 days from now
	minExpectedDelivery := time.Now().AddDate(0, 0, subscription.FrequencyWeeks*7)
	if subscription.NextDelivery.Before(minExpectedDelivery) {
		t.Errorf("SkipDelivery() NextDelivery = %v, expected at least %v", subscription.NextDelivery, minExpectedDelivery)
	}

	// Verify NextDelivery is not too far in the future (within 60 days as a sanity check)
	maxExpectedDelivery := time.Now().AddDate(0, 0, 60)
	if subscription.NextDelivery.After(maxExpectedDelivery) {
		t.Errorf("SkipDelivery() NextDelivery = %v, expected at most %v", subscription.NextDelivery, maxExpectedDelivery)
	}
}

func TestCancelSubscription(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		wantErr     bool
		errContains string
	}{
		{
			name:        "valid subscription ID",
			id:          "sub123",
			wantErr:     false,
			errContains: "",
		},
		{
			name:        "empty subscription ID should fail",
			id:          "",
			wantErr:     true,
			errContains: "subscription ID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()

			subscription, err := client.CancelSubscription(tt.id)

			if tt.wantErr {
				if err == nil {
					t.Error("CancelSubscription() expected error but got none")
					return
				}
				if err.Error() != tt.errContains {
					t.Errorf("CancelSubscription() error = %v, want %v", err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("CancelSubscription() unexpected error: %v", err)
				return
			}

			if subscription == nil {
				t.Error("CancelSubscription() returned nil subscription")
				return
			}

			// Verify subscription ID matches
			if subscription.ID != tt.id {
				t.Errorf("CancelSubscription() subscription.ID = %v, want %v", subscription.ID, tt.id)
			}

			// Verify status is set to cancelled
			if subscription.Status != "cancelled" {
				t.Errorf("CancelSubscription() status = %v, want cancelled", subscription.Status)
			}
		})
	}
}

func TestCancelSubscription_SetsStatusToCancelled(t *testing.T) {
	client := NewClient()

	subscription, err := client.CancelSubscription("sub456")
	if err != nil {
		t.Fatalf("CancelSubscription() unexpected error: %v", err)
	}

	if subscription.Status != "cancelled" {
		t.Errorf("CancelSubscription() status = %v, want cancelled", subscription.Status)
	}

	// Verify other fields are still populated correctly
	if subscription.ID != "sub456" {
		t.Errorf("CancelSubscription() ID = %v, want sub456", subscription.ID)
	}

	if subscription.ASIN == "" {
		t.Error("CancelSubscription() ASIN should not be empty")
	}

	if subscription.Title == "" {
		t.Error("CancelSubscription() Title should not be empty")
	}
}
