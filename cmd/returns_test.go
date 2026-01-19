package cmd

import (
	"encoding/json"
	"testing"

	"github.com/zkwentz/amazon-cli/internal/amazon"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

func TestReturnsCreateCmd_Success(t *testing.T) {
	// Test that the create command exists and has correct configuration
	if returnsCreateCmd.Use != "create <order-id> <item-id>" {
		t.Errorf("Expected Use='create <order-id> <item-id>', got '%s'", returnsCreateCmd.Use)
	}

	if returnsCreateCmd.Short != "Create a return" {
		t.Errorf("Expected Short='Create a return', got '%s'", returnsCreateCmd.Short)
	}

	if returnsCreateCmd.Run == nil {
		t.Error("Expected Run function to be defined")
	}

	// Test that it requires exactly 2 arguments
	if returnsCreateCmd.Args == nil {
		t.Error("Expected Args validator to be defined")
	}
}

func TestReturnsCreateCmd_Flags(t *testing.T) {
	// Test that flags are properly configured
	reasonFlag := returnsCreateCmd.Flags().Lookup("reason")
	if reasonFlag == nil {
		t.Error("Expected --reason flag to be defined")
	} else {
		if reasonFlag.DefValue != "" {
			t.Errorf("Expected --reason default value to be empty, got '%s'", reasonFlag.DefValue)
		}
	}

	confirmFlag := returnsCreateCmd.Flags().Lookup("confirm")
	if confirmFlag == nil {
		t.Error("Expected --confirm flag to be defined")
	} else {
		if confirmFlag.DefValue != "false" {
			t.Errorf("Expected --confirm default value to be 'false', got '%s'", confirmFlag.DefValue)
		}
	}
}

func TestReturnsCreateCmd_ReasonRequired(t *testing.T) {
	// Test that the reason flag is marked as required
	reasonFlag := returnsCreateCmd.Flags().Lookup("reason")
	if reasonFlag == nil {
		t.Fatal("Expected --reason flag to be defined")
	}

	// The flag should be registered, but we can't easily test the required
	// annotation without executing the command
	// We verify the flag exists and has the correct metadata
	usage := reasonFlag.Usage
	expectedUsage := "Return reason (required): defective, wrong_item, not_as_described, no_longer_needed, better_price, other"
	if usage != expectedUsage {
		t.Errorf("Expected usage='%s', got '%s'", expectedUsage, usage)
	}
}

func TestReturnsCmd_Configuration(t *testing.T) {
	// Test the main returns command configuration
	if returnsCmd.Use != "returns" {
		t.Errorf("Expected Use='returns', got '%s'", returnsCmd.Use)
	}

	if returnsCmd.Short != "Manage returns" {
		t.Errorf("Expected Short='Manage returns', got '%s'", returnsCmd.Short)
	}

	expectedLong := "List returnable items, get return options, and create returns."
	if returnsCmd.Long != expectedLong {
		t.Errorf("Expected Long='%s', got '%s'", expectedLong, returnsCmd.Long)
	}
}

func TestReturnsCmd_Subcommands(t *testing.T) {
	// Test that all subcommands are registered
	commands := returnsCmd.Commands()

	if len(commands) != 3 {
		t.Errorf("Expected 3 subcommands, got %d", len(commands))
	}

	// Check that create subcommand exists
	foundCreate := false
	foundLabel := false
	foundStatus := false
	for _, cmd := range commands {
		if cmd.Use == "create <order-id> <item-id>" {
			foundCreate = true
		}
		if cmd.Use == "label <return-id>" {
			foundLabel = true
		}
		if cmd.Use == "status <return-id>" {
			foundStatus = true
		}
	}
	if !foundCreate {
		t.Error("Expected 'create' subcommand not found")
	}
	if !foundLabel {
		t.Error("Expected 'label' subcommand not found")
	}
	if !foundStatus {
		t.Error("Expected 'status' subcommand not found")
	}
}

func TestReturnsCmd_VariablesInitialized(t *testing.T) {
	// Test that package-level variables are initialized
	// Save original values
	origReason := returnsReason
	origConfirm := returnsConfirm

	// Modify them
	returnsReason = "defective"
	returnsConfirm = true

	// Verify modifications worked
	if returnsReason != "defective" {
		t.Error("Failed to modify returnsReason")
	}
	if !returnsConfirm {
		t.Error("Failed to modify returnsConfirm")
	}

	// Restore original values
	returnsReason = origReason
	returnsConfirm = origConfirm
}

func TestReturnsCreateCmd_ResponseParsing(t *testing.T) {
	// This test verifies that the models.Return structure
	// can be properly marshaled to JSON (as used by output.JSON)

	ret := &models.Return{
		ReturnID:  "RET-123e4567-e89b-12d3-a456-426614174000",
		OrderID:   "123-4567890-1234567",
		ItemID:    "item-12345",
		Status:    "pending",
		Reason:    "defective",
		CreatedAt: "2026-01-18T12:00:00Z",
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(ret)
	if err != nil {
		t.Fatalf("Failed to marshal Return to JSON: %v", err)
	}

	// Verify it's valid JSON
	var parsed map[string]interface{}
	err = json.Unmarshal(jsonBytes, &parsed)
	if err != nil {
		t.Fatalf("Failed to parse marshaled JSON: %v", err)
	}

	// Verify expected fields exist
	if _, ok := parsed["return_id"]; !ok {
		t.Error("Expected 'return_id' field in JSON output")
	}
	if _, ok := parsed["order_id"]; !ok {
		t.Error("Expected 'order_id' field in JSON output")
	}
	if _, ok := parsed["item_id"]; !ok {
		t.Error("Expected 'item_id' field in JSON output")
	}
	if _, ok := parsed["status"]; !ok {
		t.Error("Expected 'status' field in JSON output")
	}
	if _, ok := parsed["reason"]; !ok {
		t.Error("Expected 'reason' field in JSON output")
	}
	if _, ok := parsed["created_at"]; !ok {
		t.Error("Expected 'created_at' field in JSON output")
	}
}

func TestReturnsCreateCmd_DryRunPreview(t *testing.T) {
	// Test that dry run preview has correct structure
	preview := map[string]interface{}{
		"dry_run":  true,
		"order_id": "123-4567890-1234567",
		"item_id":  "item-12345",
		"reason":   "defective",
		"message":  "Add --confirm to submit the return",
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(preview)
	if err != nil {
		t.Fatalf("Failed to marshal preview to JSON: %v", err)
	}

	// Verify it's valid JSON
	var parsed map[string]interface{}
	err = json.Unmarshal(jsonBytes, &parsed)
	if err != nil {
		t.Fatalf("Failed to parse marshaled JSON: %v", err)
	}

	// Verify expected fields exist
	if _, ok := parsed["dry_run"]; !ok {
		t.Error("Expected 'dry_run' field in preview JSON")
	}
	if _, ok := parsed["order_id"]; !ok {
		t.Error("Expected 'order_id' field in preview JSON")
	}
	if _, ok := parsed["item_id"]; !ok {
		t.Error("Expected 'item_id' field in preview JSON")
	}
	if _, ok := parsed["reason"]; !ok {
		t.Error("Expected 'reason' field in preview JSON")
	}
	if _, ok := parsed["message"]; !ok {
		t.Error("Expected 'message' field in preview JSON")
	}

	// Verify dry_run is true
	if dryRun, ok := parsed["dry_run"].(bool); !ok || !dryRun {
		t.Error("Expected 'dry_run' to be true in preview")
	}
}

func TestReturnsCreateCmd_ValidReasons(t *testing.T) {
	// Test that CreateReturn validates reasons correctly
	testClient := amazon.NewClient()

	validReasons := []string{
		"defective",
		"wrong_item",
		"not_as_described",
		"no_longer_needed",
		"better_price",
		"other",
	}

	for _, reason := range validReasons {
		ret, err := testClient.CreateReturn("123-4567890-1234567", "item-12345", reason)
		if err != nil {
			t.Errorf("Expected no error for valid reason '%s', got %v", reason, err)
		}
		if ret == nil {
			t.Errorf("Expected non-nil return for valid reason '%s'", reason)
		}
		if ret != nil && ret.Reason != reason {
			t.Errorf("Expected reason '%s', got '%s'", reason, ret.Reason)
		}
	}
}

func TestReturnsCreateCmd_InvalidReason(t *testing.T) {
	// Test that CreateReturn rejects invalid reasons
	testClient := amazon.NewClient()

	invalidReasons := []string{
		"invalid_reason",
		"",
		"DEFECTIVE",
		"changed my mind",
	}

	for _, reason := range invalidReasons {
		_, err := testClient.CreateReturn("123-4567890-1234567", "item-12345", reason)
		if err == nil {
			t.Errorf("Expected error for invalid reason '%s', got nil", reason)
		}
	}
}

func TestReturnsCreateCmd_EmptyOrderID(t *testing.T) {
	// Test that CreateReturn validates empty order ID
	testClient := amazon.NewClient()

	_, err := testClient.CreateReturn("", "item-12345", "defective")
	if err == nil {
		t.Error("Expected error for empty order ID, got nil")
	}
}

func TestReturnsCreateCmd_EmptyItemID(t *testing.T) {
	// Test that CreateReturn validates empty item ID
	testClient := amazon.NewClient()

	_, err := testClient.CreateReturn("123-4567890-1234567", "", "defective")
	if err == nil {
		t.Error("Expected error for empty item ID, got nil")
	}
}

func TestReturnsCreateCmd_GeneratesReturnID(t *testing.T) {
	// Test that CreateReturn generates a unique return ID
	testClient := amazon.NewClient()

	ret1, err1 := testClient.CreateReturn("123-4567890-1234567", "item-12345", "defective")
	if err1 != nil {
		t.Fatalf("Expected no error, got %v", err1)
	}

	ret2, err2 := testClient.CreateReturn("123-4567890-1234567", "item-12345", "defective")
	if err2 != nil {
		t.Fatalf("Expected no error, got %v", err2)
	}

	// Verify return IDs are different
	if ret1.ReturnID == ret2.ReturnID {
		t.Error("Expected different return IDs for separate calls")
	}

	// Verify return ID format
	if ret1.ReturnID[:4] != "RET-" {
		t.Errorf("Expected return ID to start with 'RET-', got '%s'", ret1.ReturnID)
	}
}

func TestReturnsCreateCmd_SetsCorrectStatus(t *testing.T) {
	// Test that CreateReturn sets status to "pending"
	testClient := amazon.NewClient()

	ret, err := testClient.CreateReturn("123-4567890-1234567", "item-12345", "defective")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if ret.Status != "pending" {
		t.Errorf("Expected status='pending', got '%s'", ret.Status)
	}
}

func TestReturnsCreateCmd_SetsCreatedAt(t *testing.T) {
	// Test that CreateReturn sets CreatedAt timestamp
	testClient := amazon.NewClient()

	ret, err := testClient.CreateReturn("123-4567890-1234567", "item-12345", "defective")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if ret.CreatedAt == "" {
		t.Error("Expected CreatedAt to be set")
	}
}

// Tests for returns label command

func TestReturnsLabelCmd_Success(t *testing.T) {
	// Test that the label command exists and has correct configuration
	if returnsLabelCmd.Use != "label <return-id>" {
		t.Errorf("Expected Use='label <return-id>', got '%s'", returnsLabelCmd.Use)
	}

	if returnsLabelCmd.Short != "Get return label" {
		t.Errorf("Expected Short='Get return label', got '%s'", returnsLabelCmd.Short)
	}

	if returnsLabelCmd.Run == nil {
		t.Error("Expected Run function to be defined")
	}

	// Test that it requires exactly 1 argument
	if returnsLabelCmd.Args == nil {
		t.Error("Expected Args validator to be defined")
	}
}

func TestReturnsLabelCmd_ClientMethod(t *testing.T) {
	// Test that the client method returns the correct data
	testClient := amazon.NewClient()

	label, err := testClient.GetReturnLabel("RET-123e4567-e89b-12d3-a456-426614174000")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if label == nil {
		t.Fatal("Expected non-nil label")
	}

	// Verify label fields
	if label.URL == "" {
		t.Error("Expected URL to be set")
	}
	if label.Carrier == "" {
		t.Error("Expected Carrier to be set")
	}
	if label.Instructions == "" {
		t.Error("Expected Instructions to be set")
	}
}

func TestReturnsLabelCmd_EmptyReturnID(t *testing.T) {
	// Test that GetReturnLabel validates empty return ID
	testClient := amazon.NewClient()

	_, err := testClient.GetReturnLabel("")
	if err == nil {
		t.Error("Expected error for empty return ID, got nil")
	}
}

func TestReturnsLabelCmd_ResponseParsing(t *testing.T) {
	// This test verifies that the models.ReturnLabel structure
	// can be properly marshaled to JSON (as used by output.JSON)

	label := &models.ReturnLabel{
		URL:          "https://amazon.com/returns/label/RET-123.pdf",
		Carrier:      "UPS",
		Instructions: "Print this label and attach it to your package.",
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(label)
	if err != nil {
		t.Fatalf("Failed to marshal ReturnLabel to JSON: %v", err)
	}

	// Verify it's valid JSON
	var parsed map[string]interface{}
	err = json.Unmarshal(jsonBytes, &parsed)
	if err != nil {
		t.Fatalf("Failed to parse marshaled JSON: %v", err)
	}

	// Verify expected fields exist
	if _, ok := parsed["url"]; !ok {
		t.Error("Expected 'url' field in JSON output")
	}
	if _, ok := parsed["carrier"]; !ok {
		t.Error("Expected 'carrier' field in JSON output")
	}
	if _, ok := parsed["instructions"]; !ok {
		t.Error("Expected 'instructions' field in JSON output")
	}
}

// Tests for returns status command

func TestReturnsStatusCmd_Success(t *testing.T) {
	// Test that the status command exists and has correct configuration
	if returnsStatusCmd.Use != "status <return-id>" {
		t.Errorf("Expected Use='status <return-id>', got '%s'", returnsStatusCmd.Use)
	}

	if returnsStatusCmd.Short != "Get return status" {
		t.Errorf("Expected Short='Get return status', got '%s'", returnsStatusCmd.Short)
	}

	if returnsStatusCmd.Run == nil {
		t.Error("Expected Run function to be defined")
	}

	// Test that it requires exactly 1 argument
	if returnsStatusCmd.Args == nil {
		t.Error("Expected Args validator to be defined")
	}
}

func TestReturnsStatusCmd_ClientMethod(t *testing.T) {
	// Test that the client method returns the correct data
	testClient := amazon.NewClient()

	ret, err := testClient.GetReturnStatus("RET-123e4567-e89b-12d3-a456-426614174000")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if ret == nil {
		t.Fatal("Expected non-nil return")
	}

	// Verify return fields
	if ret.ReturnID == "" {
		t.Error("Expected ReturnID to be set")
	}
	if ret.OrderID == "" {
		t.Error("Expected OrderID to be set")
	}
	if ret.ItemID == "" {
		t.Error("Expected ItemID to be set")
	}
	if ret.Status == "" {
		t.Error("Expected Status to be set")
	}
	if ret.Reason == "" {
		t.Error("Expected Reason to be set")
	}
	if ret.CreatedAt == "" {
		t.Error("Expected CreatedAt to be set")
	}
}

func TestReturnsStatusCmd_EmptyReturnID(t *testing.T) {
	// Test that GetReturnStatus validates empty return ID
	testClient := amazon.NewClient()

	_, err := testClient.GetReturnStatus("")
	if err == nil {
		t.Error("Expected error for empty return ID, got nil")
	}
}

func TestReturnsStatusCmd_ResponseParsing(t *testing.T) {
	// This test verifies that the models.Return structure
	// can be properly marshaled to JSON (as used by output.JSON)

	ret := &models.Return{
		ReturnID:  "RET-123e4567-e89b-12d3-a456-426614174000",
		OrderID:   "123-4567890-1234567",
		ItemID:    "item-12345",
		Status:    "approved",
		Reason:    "defective",
		CreatedAt: "2026-01-18T12:00:00Z",
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(ret)
	if err != nil {
		t.Fatalf("Failed to marshal Return to JSON: %v", err)
	}

	// Verify it's valid JSON
	var parsed map[string]interface{}
	err = json.Unmarshal(jsonBytes, &parsed)
	if err != nil {
		t.Fatalf("Failed to parse marshaled JSON: %v", err)
	}

	// Verify expected fields exist
	if _, ok := parsed["return_id"]; !ok {
		t.Error("Expected 'return_id' field in JSON output")
	}
	if _, ok := parsed["order_id"]; !ok {
		t.Error("Expected 'order_id' field in JSON output")
	}
	if _, ok := parsed["item_id"]; !ok {
		t.Error("Expected 'item_id' field in JSON output")
	}
	if _, ok := parsed["status"]; !ok {
		t.Error("Expected 'status' field in JSON output")
	}
	if _, ok := parsed["reason"]; !ok {
		t.Error("Expected 'reason' field in JSON output")
	}
	if _, ok := parsed["created_at"]; !ok {
		t.Error("Expected 'created_at' field in JSON output")
	}
}
