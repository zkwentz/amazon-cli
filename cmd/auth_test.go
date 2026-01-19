package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestAuthLoginCmd(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Create a new command instance for testing
	cmd := &cobra.Command{
		Use: "login",
		Run: authLoginCmd.Run,
	}

	// Execute the command
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Command execution failed: %v", err)
	}

	// Restore stdout and read captured output
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Parse the JSON output
	var result map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &result)
	if err != nil {
		t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, buf.String())
	}

	// Verify the status field
	status, ok := result["status"].(string)
	if !ok {
		t.Fatal("status field is missing or not a string")
	}
	if status != "login_required" {
		t.Errorf("Expected status 'login_required', got '%s'", status)
	}

	// Verify the message field
	message, ok := result["message"].(string)
	if !ok {
		t.Fatal("message field is missing or not a string")
	}
	expectedMessage := "Browser-based login not yet implemented"
	if message != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, message)
	}
}

// setupTestViper creates a fresh viper instance for testing
func setupTestViper(t *testing.T) {
	viper.Reset()
	// Use a unique config name to avoid conflicts
	tmpDir := t.TempDir()
	viper.SetConfigFile(tmpDir + "/test-config.json")
}

func TestAuthStatusCmd_NoToken(t *testing.T) {
	setupTestViper(t)

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the command
	authStatusCmd.Run(authStatusCmd, []string{})

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Parse JSON output
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, buf.String())
	}

	// Verify authenticated is false
	authenticated, ok := result["authenticated"].(bool)
	if !ok {
		t.Fatal("authenticated field is missing or not a boolean")
	}
	if authenticated {
		t.Error("Expected authenticated to be false when no token exists")
	}

	// Verify message exists
	if _, ok := result["message"]; !ok {
		t.Error("Expected message field to be present when not authenticated")
	}
}

func TestAuthStatusCmd_ValidToken(t *testing.T) {
	setupTestViper(t)

	// Set up a valid token
	expiresAt := time.Now().Add(1 * time.Hour)
	viper.Set("auth.access_token", "valid_token")
	viper.Set("auth.expires_at", expiresAt.Format(time.RFC3339))

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the command
	authStatusCmd.Run(authStatusCmd, []string{})

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Parse JSON output
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, buf.String())
	}

	// Verify authenticated is true
	authenticated, ok := result["authenticated"].(bool)
	if !ok {
		t.Fatal("authenticated field is missing or not a boolean")
	}
	if !authenticated {
		t.Error("Expected authenticated to be true with valid token")
	}

	// Verify expires_at exists and matches
	expiresAtStr, ok := result["expires_at"].(string)
	if !ok {
		t.Fatal("expires_at field is missing or not a string")
	}
	if expiresAtStr == "" {
		t.Error("Expected expires_at to not be empty")
	}

	// Verify expires_in_seconds exists and is positive
	expiresInSeconds, ok := result["expires_in_seconds"].(float64)
	if !ok {
		t.Fatal("expires_in_seconds field is missing or not a number")
	}
	if expiresInSeconds <= 0 {
		t.Error("Expected expires_in_seconds to be positive for valid token")
	}

	// Verify message does not exist for authenticated status
	if _, ok := result["message"]; ok {
		t.Error("Expected no message field for authenticated status")
	}
}

func TestAuthStatusCmd_ExpiredToken(t *testing.T) {
	setupTestViper(t)

	// Set up an expired token
	expiresAt := time.Now().Add(-1 * time.Hour)
	viper.Set("auth.access_token", "expired_token")
	viper.Set("auth.expires_at", expiresAt.Format(time.RFC3339))

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the command
	authStatusCmd.Run(authStatusCmd, []string{})

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Parse JSON output
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, buf.String())
	}

	// Verify authenticated is false
	authenticated, ok := result["authenticated"].(bool)
	if !ok {
		t.Fatal("authenticated field is missing or not a boolean")
	}
	if authenticated {
		t.Error("Expected authenticated to be false with expired token")
	}

	// Verify expired flag exists and is true
	expired, ok := result["expired"].(bool)
	if !ok {
		t.Fatal("expired field is missing or not a boolean")
	}
	if !expired {
		t.Error("Expected expired to be true for expired token")
	}

	// Verify expires_at exists
	expiresAtStr, ok := result["expires_at"].(string)
	if !ok {
		t.Fatal("expires_at field is missing or not a string")
	}
	if expiresAtStr == "" {
		t.Error("Expected expires_at to not be empty")
	}

	// Verify message exists
	if _, ok := result["message"]; !ok {
		t.Error("Expected message field to be present for expired token")
	}
}

func TestAuthStatusCmd_InvalidExpiryFormat(t *testing.T) {
	setupTestViper(t)

	// Set up a token with invalid expiry format
	viper.Set("auth.access_token", "token_with_invalid_expiry")
	viper.Set("auth.expires_at", "invalid-date-format")

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the command
	authStatusCmd.Run(authStatusCmd, []string{})

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Parse JSON output
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, buf.String())
	}

	// Verify authenticated is false
	authenticated, ok := result["authenticated"].(bool)
	if !ok {
		t.Fatal("authenticated field is missing or not a boolean")
	}
	if authenticated {
		t.Error("Expected authenticated to be false with invalid expiry format")
	}

	// Verify message exists
	message, ok := result["message"].(string)
	if !ok {
		t.Fatal("message field is missing or not a string")
	}
	if message == "" {
		t.Error("Expected message to explain the invalid token expiry")
	}
}

func TestAuthStatusCmd_ExpiresInSecondsCalculation(t *testing.T) {
	setupTestViper(t)

	// Set up a token that expires in exactly 3600 seconds (1 hour)
	expiresAt := time.Now().Add(1 * time.Hour)
	viper.Set("auth.access_token", "valid_token")
	viper.Set("auth.expires_at", expiresAt.Format(time.RFC3339))

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the command
	authStatusCmd.Run(authStatusCmd, []string{})

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Parse JSON output
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, buf.String())
	}

	// Verify expires_in_seconds is approximately 3600 (allowing for test execution time)
	expiresInSeconds, ok := result["expires_in_seconds"].(float64)
	if !ok {
		t.Fatal("expires_in_seconds field is missing or not a number")
	}

	// Allow for 5 seconds of test execution time
	if expiresInSeconds < 3595 || expiresInSeconds > 3600 {
		t.Errorf("Expected expires_in_seconds to be approximately 3600, got %f", expiresInSeconds)
	}
}

func TestAuthStatusCmd_ShortLivedToken(t *testing.T) {
	setupTestViper(t)

	// Set up a token that expires in 30 seconds
	expiresAt := time.Now().Add(30 * time.Second)
	viper.Set("auth.access_token", "short_lived_token")
	viper.Set("auth.expires_at", expiresAt.Format(time.RFC3339))

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the command
	authStatusCmd.Run(authStatusCmd, []string{})

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Parse JSON output
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, buf.String())
	}

	// Verify authenticated is true
	authenticated, ok := result["authenticated"].(bool)
	if !ok {
		t.Fatal("authenticated field is missing or not a boolean")
	}
	if !authenticated {
		t.Error("Expected authenticated to be true with short-lived but valid token")
	}

	// Verify expires_in_seconds is less than or equal to 30
	expiresInSeconds, ok := result["expires_in_seconds"].(float64)
	if !ok {
		t.Fatal("expires_in_seconds field is missing or not a number")
	}
	if expiresInSeconds > 30 {
		t.Errorf("Expected expires_in_seconds to be at most 30, got %f", expiresInSeconds)
	}
	if expiresInSeconds < 0 {
		t.Errorf("Expected expires_in_seconds to be positive, got %f", expiresInSeconds)
	}
}
