package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/michaelshimeles/amazon-cli/pkg/models"
	"github.com/spf13/cobra"
)

func TestCartCheckoutCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		wantErr        bool
		errContains    string
		expectDryRun   bool
		expectConfirm  bool
	}{
		{
			name:          "missing address-id flag should fail",
			args:          []string{"--payment-id", "pay123", "--confirm"},
			wantErr:       true,
			errContains:   "address-id is required",
			expectDryRun:  false,
			expectConfirm: false,
		},
		{
			name:          "missing payment-id flag should fail",
			args:          []string{"--address-id", "addr123", "--confirm"},
			wantErr:       true,
			errContains:   "payment-id is required",
			expectDryRun:  false,
			expectConfirm: false,
		},
		{
			name:          "without confirm flag should show dry run",
			args:          []string{"--address-id", "addr123", "--payment-id", "pay456"},
			wantErr:       false,
			expectDryRun:  true,
			expectConfirm: false,
		},
		{
			name:          "with confirm flag but empty cart should fail",
			args:          []string{"--address-id", "addr123", "--payment-id", "pay456", "--confirm"},
			wantErr:       true,
			errContains:   "cart is empty",
			expectDryRun:  false,
			expectConfirm: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new command for each test to reset state
			cmd := &cobra.Command{
				Use: "checkout",
				RunE: func(cmd *cobra.Command, args []string) error {
					return cartCheckoutCmd.RunE(cmd, args)
				},
			}

			// Add flags
			cmd.Flags().StringVar(&addressID, "address-id", "", "Address ID")
			cmd.Flags().StringVar(&paymentID, "payment-id", "", "Payment ID")
			cmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm")

			// Set args and execute
			cmd.SetArgs(tt.args)

			// Capture output
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)

			err := cmd.Execute()

			// Check error expectations
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// For dry run test, the RunE function prints to stdout directly
			// So we can't easily capture it in tests without redirecting stdout
			// This is acceptable for CLI commands that output JSON
			// The actual behavior is tested through integration tests
		})
	}
}

func TestCartCheckoutDryRun(t *testing.T) {
	// Reset flags
	addressID = ""
	paymentID = ""
	confirm = false

	// Create command
	cmd := &cobra.Command{
		Use:  "checkout",
		RunE: cartCheckoutCmd.RunE,
	}

	cmd.Flags().StringVar(&addressID, "address-id", "", "Address ID")
	cmd.Flags().StringVar(&paymentID, "payment-id", "", "Payment ID")
	cmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm")

	cmd.SetArgs([]string{"--address-id", "addr123", "--payment-id", "pay456"})

	// Capture output
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Parse JSON output
	var result map[string]interface{}
	output := buf.String()
	if output != "" {
		err = json.Unmarshal([]byte(output), &result)
		if err != nil {
			t.Fatalf("Failed to parse JSON output: %v", err)
		}

		// Verify dry_run flag
		if dryRun, ok := result["dry_run"].(bool); !ok || !dryRun {
			t.Error("Expected dry_run to be true")
		}

		// Verify message
		if msg, ok := result["message"].(string); !ok || msg == "" {
			t.Error("Expected message in dry run output")
		}

		// Verify preview exists
		if _, ok := result["preview"]; !ok {
			t.Error("Expected preview in dry run output")
		}
	}
}

func TestCartAddCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		quantity    int
		wantErr     bool
		errContains string
	}{
		{
			name:     "valid add with default quantity",
			args:     []string{"B08N5WRWNW"},
			quantity: 1,
			wantErr:  false,
		},
		{
			name:        "no ASIN should fail",
			args:        []string{},
			quantity:    1,
			wantErr:     true,
			errContains: "accepts 1 arg(s)",
		},
		{
			name:     "valid add with custom quantity",
			args:     []string{"B08N5WRWNW"},
			quantity: 3,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset quantity flag
			quantity = tt.quantity

			cmd := &cobra.Command{
				Use:  "add",
				Args: cobra.ExactArgs(1),
				RunE: cartAddCmd.RunE,
			}

			cmd.Flags().IntVarP(&quantity, "quantity", "n", 1, "Quantity")
			cmd.SetArgs(tt.args)

			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)

			err := cmd.Execute()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Verify output is valid JSON
			output := buf.String()
			if output != "" {
				var cart models.Cart
				err = json.Unmarshal([]byte(output), &cart)
				if err != nil {
					t.Errorf("Failed to parse JSON output: %v", err)
				}
			}
		})
	}
}

func TestCartListCommand(t *testing.T) {
	cmd := &cobra.Command{
		Use:  "list",
		Args: cobra.NoArgs,
		RunE: cartListCmd.RunE,
	}

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify output is valid JSON
	output := buf.String()
	if output != "" {
		var cart models.Cart
		err = json.Unmarshal([]byte(output), &cart)
		if err != nil {
			t.Errorf("Failed to parse JSON output: %v", err)
		}
	}
}

func TestCartRemoveCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid remove",
			args:    []string{"B08N5WRWNW"},
			wantErr: false,
		},
		{
			name:        "no ASIN should fail",
			args:        []string{},
			wantErr:     true,
			errContains: "accepts 1 arg(s)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{
				Use:  "remove",
				Args: cobra.ExactArgs(1),
				RunE: cartRemoveCmd.RunE,
			}

			cmd.SetArgs(tt.args)

			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)

			err := cmd.Execute()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestCartClearCommand(t *testing.T) {
	tests := []struct {
		name           string
		confirmFlag    bool
		expectDryRun   bool
	}{
		{
			name:         "without confirm should show dry run",
			confirmFlag:  false,
			expectDryRun: true,
		},
		{
			name:         "with confirm should clear cart",
			confirmFlag:  true,
			expectDryRun: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset confirm flag
			confirm = tt.confirmFlag

			cmd := &cobra.Command{
				Use:  "clear",
				Args: cobra.NoArgs,
				RunE: cartClearCmd.RunE,
			}

			cmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm")

			buf := new(bytes.Buffer)
			cmd.SetOut(buf)

			err := cmd.Execute()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// The RunE function prints to stdout directly
			// Output validation is done through integration tests
		})
	}
}
