package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRunBuy_ASINValidation(t *testing.T) {
	tests := []struct {
		name        string
		asin        string
		wantErr     bool
		errContains string
	}{
		{
			name:        "valid ASIN",
			asin:        "B08N5WRWNW",
			wantErr:     false,
			errContains: "",
		},
		{
			name:        "ASIN too short",
			asin:        "B08N5",
			wantErr:     true,
			errContains: "invalid ASIN format",
		},
		{
			name:        "ASIN too long",
			asin:        "B08N5WRWNW123",
			wantErr:     true,
			errContains: "invalid ASIN format",
		},
		{
			name:        "empty ASIN",
			asin:        "",
			wantErr:     true,
			errContains: "invalid ASIN format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags
			buyConfirm = false
			buyQuantity = 1
			buyAddressID = ""
			buyPaymentID = ""

			// Create a new command instance for each test
			cmd := &cobra.Command{}

			err := runBuy(cmd, []string{tt.asin})

			if tt.wantErr {
				if err == nil {
					t.Errorf("runBuy() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("runBuy() error = %v, want error containing %q", err, tt.errContains)
				}
			} else if err != nil {
				t.Errorf("runBuy() unexpected error: %v", err)
			}
		})
	}
}

func TestRunBuy_QuantityValidation(t *testing.T) {
	tests := []struct {
		name        string
		quantity    int
		wantErr     bool
		errContains string
	}{
		{
			name:        "quantity 1",
			quantity:    1,
			wantErr:     false,
			errContains: "",
		},
		{
			name:        "quantity 5",
			quantity:    5,
			wantErr:     false,
			errContains: "",
		},
		{
			name:        "quantity 0",
			quantity:    0,
			wantErr:     true,
			errContains: "quantity must be positive",
		},
		{
			name:        "negative quantity",
			quantity:    -1,
			wantErr:     true,
			errContains: "quantity must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags
			buyConfirm = false
			buyQuantity = tt.quantity
			buyAddressID = ""
			buyPaymentID = ""

			cmd := &cobra.Command{}
			err := runBuy(cmd, []string{"B08N5WRWNW"})

			if tt.wantErr {
				if err == nil {
					t.Errorf("runBuy() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("runBuy() error = %v, want error containing %q", err, tt.errContains)
				}
			} else if err != nil {
				t.Errorf("runBuy() unexpected error: %v", err)
			}
		})
	}
}

func TestRunBuy_DryRun(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Reset flags for dry run
	buyConfirm = false
	buyQuantity = 2
	buyAddressID = "addr_123"
	buyPaymentID = "pay_456"

	cmd := &cobra.Command{}
	err := runBuy(cmd, []string{"B08N5WRWNW"})

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("runBuy() unexpected error: %v", err)
	}

	// Read captured output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Parse JSON output
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("failed to parse JSON output: %v", err)
	}

	// Verify dry run fields
	if dryRun, ok := result["dry_run"].(bool); !ok || !dryRun {
		t.Error("expected dry_run to be true")
	}

	if quantity, ok := result["quantity"].(float64); !ok || int(quantity) != 2 {
		t.Errorf("expected quantity to be 2, got %v", result["quantity"])
	}

	if addressID, ok := result["address_id"].(string); !ok || addressID != "addr_123" {
		t.Errorf("expected address_id to be 'addr_123', got %v", result["address_id"])
	}

	if paymentID, ok := result["payment_id"].(string); !ok || paymentID != "pay_456" {
		t.Errorf("expected payment_id to be 'pay_456', got %v", result["payment_id"])
	}

	if message, ok := result["message"].(string); !ok || !strings.Contains(message, "--confirm") {
		t.Error("expected message to mention --confirm flag")
	}
}

func TestRunBuy_Confirmed(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Reset flags for confirmed purchase
	buyConfirm = true
	buyQuantity = 1
	buyAddressID = ""
	buyPaymentID = ""
	quiet = true // Suppress stderr output

	cmd := &cobra.Command{}
	err := runBuy(cmd, []string{"B08N5WRWNW"})

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout
	quiet = false // Reset quiet flag

	if err != nil {
		t.Fatalf("runBuy() unexpected error: %v", err)
	}

	// Read captured output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Parse JSON output
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("failed to parse JSON output: %v", err)
	}

	// Verify order confirmation fields
	if orderID, ok := result["order_id"].(string); !ok || orderID == "" {
		t.Error("expected order_id to be non-empty")
	}

	if total, ok := result["total"].(float64); !ok || total <= 0 {
		t.Error("expected total to be greater than 0")
	}

	if estimatedDelivery, ok := result["estimated_delivery"].(string); !ok || estimatedDelivery == "" {
		t.Error("expected estimated_delivery to be non-empty")
	}

	// Make sure it's NOT a dry run
	if dryRun, ok := result["dry_run"].(bool); ok && dryRun {
		t.Error("expected dry_run to be false or absent for confirmed purchase")
	}
}

func TestRunBuy_ConfirmedWithCustomOptions(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Reset flags with custom options
	buyConfirm = true
	buyQuantity = 3
	buyAddressID = "custom_addr"
	buyPaymentID = "custom_pay"
	quiet = true

	cmd := &cobra.Command{}
	err := runBuy(cmd, []string{"B07XJ8C8F5"})

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout
	quiet = false

	if err != nil {
		t.Fatalf("runBuy() unexpected error: %v", err)
	}

	// Read captured output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Parse JSON output
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("failed to parse JSON output: %v", err)
	}

	// Verify order was created
	if orderID, ok := result["order_id"].(string); !ok || orderID == "" {
		t.Error("expected order_id to be non-empty")
	}
}

func TestBuyCmd_Flags(t *testing.T) {
	// Test that all expected flags are registered
	flags := buyCmd.Flags()

	if flags.Lookup("confirm") == nil {
		t.Error("expected --confirm flag to be registered")
	}

	if flags.Lookup("quantity") == nil {
		t.Error("expected --quantity flag to be registered")
	}

	if flags.Lookup("address-id") == nil {
		t.Error("expected --address-id flag to be registered")
	}

	if flags.Lookup("payment-id") == nil {
		t.Error("expected --payment-id flag to be registered")
	}
}

func TestBuyCmd_Args(t *testing.T) {
	// Test that command requires exactly 1 arg
	if buyCmd.Args == nil {
		t.Error("expected Args validation to be set")
		return
	}

	// Test with 0 args
	err := buyCmd.Args(buyCmd, []string{})
	if err == nil {
		t.Error("expected error with 0 arguments")
	}

	// Test with 1 arg (should pass)
	err = buyCmd.Args(buyCmd, []string{"B08N5WRWNW"})
	if err != nil {
		t.Errorf("unexpected error with 1 argument: %v", err)
	}

	// Test with 2 args
	err = buyCmd.Args(buyCmd, []string{"B08N5WRWNW", "extra"})
	if err == nil {
		t.Error("expected error with 2 arguments")
	}
}
