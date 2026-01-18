package cmd

import (
	"testing"
)

func TestBuyCommand(t *testing.T) {
	// Test that the buy command is registered
	cmd := rootCmd
	buyCmd := cmd.Commands()

	found := false
	for _, c := range buyCmd {
		if c.Name() == "buy" {
			found = true
			break
		}
	}

	if !found {
		t.Error("buy command not found in root command")
	}
}

func TestBuyFlags(t *testing.T) {
	// Test that all expected flags are defined
	flags := []string{"confirm", "quantity", "address-id", "payment-id"}

	for _, flagName := range flags {
		flag := buyCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("flag --%s not found", flagName)
		}
	}
}

func TestBuyCommandArgs(t *testing.T) {
	// Test that buy command requires exactly one argument (ASIN)
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no args should fail",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "one arg should succeed",
			args:    []string{"B08N5WRWNW"},
			wantErr: false,
		},
		{
			name:    "two args should fail",
			args:    []string{"B08N5WRWNW", "extra"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := buyCmd.Args(buyCmd, tt.args)
			if tt.wantErr && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
