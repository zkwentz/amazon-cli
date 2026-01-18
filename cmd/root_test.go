package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/michaelshimeles/amazon-cli/pkg/models"
	"github.com/spf13/cobra"
)

func TestExecute_Success(t *testing.T) {
	// Save and restore original rootCmd
	originalRootCmd := rootCmd
	defer func() { rootCmd = originalRootCmd }()

	// Create a test command that succeeds
	rootCmd = &cobra.Command{
		Use: "test",
		Run: func(cmd *cobra.Command, args []string) {
			// Success - do nothing
		},
	}

	exitCode := Execute()
	if exitCode != models.ExitSuccess {
		t.Errorf("Execute() = %d, want %d", exitCode, models.ExitSuccess)
	}
}

func TestExecute_GeneralError(t *testing.T) {
	// Save and restore original rootCmd
	originalRootCmd := rootCmd
	defer func() { rootCmd = originalRootCmd }()

	// Create a test command that returns a general error
	rootCmd = &cobra.Command{
		Use: "test",
		RunE: func(cmd *cobra.Command, args []string) error {
			return models.NewAmazonError("Test error")
		},
	}

	// Suppress error output
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true

	exitCode := Execute()
	if exitCode != models.ExitGeneralError {
		t.Errorf("Execute() = %d, want %d", exitCode, models.ExitGeneralError)
	}
}

func TestExecute_AuthError(t *testing.T) {
	// Save and restore original rootCmd
	originalRootCmd := rootCmd
	defer func() { rootCmd = originalRootCmd }()

	// Create a test command that returns an auth error
	rootCmd = &cobra.Command{
		Use: "test",
		RunE: func(cmd *cobra.Command, args []string) error {
			return models.NewAuthRequiredError("Authentication required")
		},
	}

	// Suppress error output
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true

	exitCode := Execute()
	if exitCode != models.ExitAuthError {
		t.Errorf("Execute() = %d, want %d", exitCode, models.ExitAuthError)
	}
}

func TestExecute_NetworkError(t *testing.T) {
	// Save and restore original rootCmd
	originalRootCmd := rootCmd
	defer func() { rootCmd = originalRootCmd }()

	// Create a test command that returns a network error
	rootCmd = &cobra.Command{
		Use: "test",
		RunE: func(cmd *cobra.Command, args []string) error {
			return models.NewNetworkError("Network failed")
		},
	}

	// Suppress error output
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true

	exitCode := Execute()
	if exitCode != models.ExitNetworkError {
		t.Errorf("Execute() = %d, want %d", exitCode, models.ExitNetworkError)
	}
}

func TestExecute_NotFoundError(t *testing.T) {
	// Save and restore original rootCmd
	originalRootCmd := rootCmd
	defer func() { rootCmd = originalRootCmd }()

	// Create a test command that returns a not found error
	rootCmd = &cobra.Command{
		Use: "test",
		RunE: func(cmd *cobra.Command, args []string) error {
			return models.NewNotFoundError("Order not found")
		},
	}

	// Suppress error output
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true

	exitCode := Execute()
	if exitCode != models.ExitNotFound {
		t.Errorf("Execute() = %d, want %d", exitCode, models.ExitNotFound)
	}
}

func TestExecute_RateLimitedError(t *testing.T) {
	// Save and restore original rootCmd
	originalRootCmd := rootCmd
	defer func() { rootCmd = originalRootCmd }()

	// Create a test command that returns a rate limited error
	rootCmd = &cobra.Command{
		Use: "test",
		RunE: func(cmd *cobra.Command, args []string) error {
			return models.NewRateLimitedError("Too many requests")
		},
	}

	// Suppress error output
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true

	exitCode := Execute()
	if exitCode != models.ExitRateLimited {
		t.Errorf("Execute() = %d, want %d", exitCode, models.ExitRateLimited)
	}
}

func TestExecute_InvalidInputError(t *testing.T) {
	// Save and restore original rootCmd
	originalRootCmd := rootCmd
	defer func() { rootCmd = originalRootCmd }()

	// Create a test command that returns an invalid input error
	rootCmd = &cobra.Command{
		Use: "test",
		RunE: func(cmd *cobra.Command, args []string) error {
			return models.NewInvalidInputError("Invalid ASIN")
		},
	}

	// Suppress error output
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true

	exitCode := Execute()
	if exitCode != models.ExitInvalidArgs {
		t.Errorf("Execute() = %d, want %d", exitCode, models.ExitInvalidArgs)
	}
}

func TestExecute_HelpFlag(t *testing.T) {
	// Save and restore original state
	originalRootCmd := rootCmd
	originalArgs := os.Args
	defer func() {
		rootCmd = originalRootCmd
		os.Args = originalArgs
	}()

	// Create a simple test command
	rootCmd = &cobra.Command{
		Use:   "test",
		Short: "Test command",
		Run: func(cmd *cobra.Command, args []string) {
			// Success
		},
	}

	// Capture stdout to suppress help output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Set args to trigger help
	os.Args = []string{"test", "--help"}

	// Execute should return success for help
	exitCode := Execute()

	// Restore stdout
	w.Close()
	os.Stdout = old
	io.Copy(io.Discard, r)

	if exitCode != models.ExitSuccess {
		t.Errorf("Execute() with --help = %d, want %d", exitCode, models.ExitSuccess)
	}
}

func TestExecute_InvalidFlag(t *testing.T) {
	// Save and restore original state
	originalRootCmd := rootCmd
	originalArgs := os.Args
	defer func() {
		rootCmd = originalRootCmd
		os.Args = originalArgs
	}()

	// Create a simple test command
	rootCmd = &cobra.Command{
		Use:   "test",
		Short: "Test command",
		Run: func(cmd *cobra.Command, args []string) {
			// Success
		},
	}

	// Suppress error output
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true

	// Capture stderr
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Set args with invalid flag
	os.Args = []string{"test", "--invalid-flag"}

	exitCode := Execute()

	// Restore stderr
	w.Close()
	os.Stderr = old

	// Read and discard captured output
	buf := new(bytes.Buffer)
	io.Copy(buf, r)

	// Invalid flag should return general error (Cobra's default behavior)
	if exitCode != models.ExitGeneralError {
		t.Errorf("Execute() with invalid flag = %d, want %d", exitCode, models.ExitGeneralError)
	}
}
