package main

import (
	"os"

	"github.com/michaelshimeles/amazon-cli/cmd"
	"github.com/michaelshimeles/amazon-cli/internal/output"
	"github.com/michaelshimeles/amazon-cli/pkg/models"
)

func main() {
	if err := cmd.Execute(); err != nil {
		// Create printer with verbose flag to show debug info on errors
		printer := output.NewPrinter(
			cmd.GetOutput(),
			cmd.GetQuiet(),
			cmd.GetVerbose(),
		)

		// Print error with verbose debug info if enabled
		if printErr := printer.PrintError(err); printErr != nil {
			// Fallback if printing fails
			os.Stderr.WriteString(err.Error() + "\n")
		}

		// Exit with appropriate code based on error type
		exitCode := getExitCode(err)
		os.Exit(exitCode)
	}
}

// getExitCode maps error codes to exit codes as defined in the PRD
func getExitCode(err error) int {
	if cliErr, ok := err.(*models.CLIError); ok {
		switch cliErr.Code {
		case models.ErrInvalidInput:
			return 2 // Invalid arguments
		case models.ErrAuthRequired, models.ErrAuthExpired:
			return 3 // Authentication error
		case models.ErrNetworkError:
			return 4 // Network error
		case models.ErrRateLimited:
			return 5 // Rate limited
		case models.ErrNotFound:
			return 6 // Not found
		default:
			return 1 // General error
		}
	}
	return 1 // General error for non-CLIError errors
}
