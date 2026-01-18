package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/michaelshimeles/amazon-cli/pkg/models"
)

// PrintJSON prints data as formatted JSON to stdout
func PrintJSON(data interface{}) error {
	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		// If we can't marshal, output an error JSON instead
		PrintError(models.NewCLIError(models.ErrCodeInternalError,
			fmt.Sprintf("Failed to marshal output: %v", err)))
		return err
	}
	fmt.Println(string(output))
	return nil
}

// PrintError prints a CLIError as JSON to stderr and ensures program exits with appropriate code
func PrintError(err error) {
	var cliErr *models.CLIError

	// Check if it's already a CLIError
	if e, ok := err.(*models.CLIError); ok {
		cliErr = e
	} else {
		// Wrap regular errors as INTERNAL_ERROR
		cliErr = models.NewCLIError(models.ErrCodeInternalError, err.Error())
	}

	// Print error JSON to stderr
	fmt.Fprintln(os.Stderr, cliErr.ToJSON())
}

// HandleError is a utility function to handle errors uniformly
// It prints the error and returns the appropriate exit code
func HandleError(err error) int {
	if err == nil {
		return 0
	}

	PrintError(err)

	// Map error codes to exit codes as per PRD
	if cliErr, ok := err.(*models.CLIError); ok {
		switch cliErr.Code {
		case models.ErrCodeAuthRequired, models.ErrCodeAuthExpired:
			return 3 // Authentication error
		case models.ErrCodeNetworkError:
			return 4 // Network error
		case models.ErrCodeRateLimited:
			return 5 // Rate limited
		case models.ErrCodeNotFound:
			return 6 // Not found
		case models.ErrCodeInvalidInput:
			return 2 // Invalid arguments
		default:
			return 1 // General error
		}
	}

	return 1 // Default to general error
}

// WrapPanic recovers from panics and converts them to JSON errors
// This should be used as a defer in main execution paths
func WrapPanic() {
	if r := recover(); r != nil {
		var err *models.CLIError

		// Convert panic to error
		switch v := r.(type) {
		case error:
			err = models.NewCLIError(models.ErrCodeInternalError, v.Error())
		case string:
			err = models.NewCLIError(models.ErrCodeInternalError, v)
		default:
			err = models.NewCLIError(models.ErrCodeInternalError, fmt.Sprintf("Unexpected panic: %v", r))
		}

		// Print as JSON error
		fmt.Fprintln(os.Stderr, err.ToJSON())
		os.Exit(1)
	}
}
