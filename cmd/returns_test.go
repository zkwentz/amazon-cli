package cmd

import (
	"testing"

	"github.com/zkwentz/amazon-cli/internal/amazon"
	"github.com/zkwentz/amazon-cli/internal/config"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

// TestGetReturnStatusIntegration tests the GetReturnStatus function
func TestGetReturnStatusIntegration(t *testing.T) {
	cfg := config.GetDefaultConfig()
	client := amazon.NewClient(cfg)

	// Test with a valid return ID
	result, err := client.GetReturnStatus("RET12345")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if result == nil {
		t.Errorf("Expected non-nil result")
		return
	}

	if result.ReturnID != "RET12345" {
		t.Errorf("Expected return_id=RET12345, got %s", result.ReturnID)
	}

	if result.Status == "" {
		t.Errorf("Status should not be empty")
	}
}

// TestGetReturnStatusEmptyID tests error handling for empty return ID
func TestGetReturnStatusEmptyID(t *testing.T) {
	cfg := config.GetDefaultConfig()
	client := amazon.NewClient(cfg)

	_, err := client.GetReturnStatus("")
	if err == nil {
		t.Errorf("Expected error for empty return ID")
		return
	}

	cliErr, ok := err.(*models.CLIError)
	if !ok {
		t.Errorf("Expected CLIError type")
		return
	}

	if cliErr.Code != models.ErrCodeInvalidInput {
		t.Errorf("Expected error code %s, got %s", models.ErrCodeInvalidInput, cliErr.Code)
	}
}
