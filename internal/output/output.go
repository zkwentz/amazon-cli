package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

// OutputFormat represents the output format type
type OutputFormat string

const (
	FormatJSON  OutputFormat = "json"
	FormatTable OutputFormat = "table"
	FormatRaw   OutputFormat = "raw"
)

// Printer handles output formatting
type Printer struct {
	format OutputFormat
	quiet  bool
}

// NewPrinter creates a new Printer
func NewPrinter(format string, quiet bool) *Printer {
	if format == "" {
		format = "json"
	}
	return &Printer{
		format: OutputFormat(format),
		quiet:  quiet,
	}
}

// Print outputs data in the configured format
func (p *Printer) Print(data interface{}) error {
	switch p.format {
	case FormatJSON:
		return p.printJSON(data)
	case FormatTable:
		// TODO: implement table format
		return p.printJSON(data)
	case FormatRaw:
		fmt.Println(data)
		return nil
	default:
		return p.printJSON(data)
	}
}

// PrintError outputs an error in JSON format
func (p *Printer) PrintError(err error) error {
	var cliErr *models.CLIError

	// Check if it's already a CLIError
	if e, ok := err.(*models.CLIError); ok {
		cliErr = e
	} else {
		// Wrap generic errors
		cliErr = models.NewCLIError(models.ErrCodeAmazonError, err.Error(), nil)
	}

	errResp := models.ErrorResponse{
		Error: cliErr,
	}

	return p.printJSON(errResp)
}

// printJSON marshals and prints data as JSON
func (p *Printer) printJSON(data interface{}) error {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	fmt.Println(string(jsonBytes))
	return nil
}

// Exit exits with appropriate code based on error
func Exit(err error) {
	if err == nil {
		os.Exit(0)
		return
	}

	exitCode := 1 // general error by default

	if cliErr, ok := err.(*models.CLIError); ok {
		switch cliErr.Code {
		case models.ErrCodeAuthRequired, models.ErrCodeAuthExpired:
			exitCode = 3
		case models.ErrCodeNetworkError:
			exitCode = 4
		case models.ErrCodeRateLimited:
			exitCode = 5
		case models.ErrCodeNotFound:
			exitCode = 6
		case models.ErrCodeInvalidInput:
			exitCode = 2
		}
	}

	os.Exit(exitCode)
}
