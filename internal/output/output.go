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

// Printer handles output formatting and printing
type Printer struct {
	format OutputFormat
	quiet  bool
}

// NewPrinter creates a new Printer with the specified format and quiet mode
func NewPrinter(format string, quiet bool) *Printer {
	outputFormat := FormatJSON
	switch format {
	case "table":
		outputFormat = FormatTable
	case "raw":
		outputFormat = FormatRaw
	case "json":
		outputFormat = FormatJSON
	}

	return &Printer{
		format: outputFormat,
		quiet:  quiet,
	}
}

// Print outputs data in the configured format
func (p *Printer) Print(data interface{}) error {
	if p.quiet {
		return nil
	}

	switch p.format {
	case FormatJSON:
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(jsonData))
	case FormatTable:
		// TODO: Implement table formatting using tablewriter
		fmt.Fprintf(os.Stdout, "%+v\n", data)
	case FormatRaw:
		fmt.Fprintf(os.Stdout, "%v\n", data)
	}

	return nil
}

// PrintError formats and outputs an error as JSON
// Returns the original error for chaining
func (p *Printer) PrintError(err error) error {
	if err == nil {
		return nil
	}

	// Check if it's already a CLIError
	cliErr, ok := err.(*models.CLIError)
	if !ok {
		// Wrap regular errors as generic CLIErrors
		cliErr = &models.CLIError{
			Code:    models.ErrCodeAmazonError,
			Message: err.Error(),
			Details: make(map[string]interface{}),
		}
	}

	// Create the error response structure
	errorResponse := map[string]interface{}{
		"error": cliErr,
	}

	// Marshal to JSON with indentation
	jsonData, marshalErr := json.MarshalIndent(errorResponse, "", "  ")
	if marshalErr != nil {
		// Fallback to simple error output if JSON marshaling fails
		fmt.Fprintf(os.Stderr, `{"error": {"code": "INTERNAL_ERROR", "message": "Failed to format error: %s"}}`, marshalErr.Error())
		fmt.Fprintln(os.Stderr)
		return err
	}

	// Print to stderr
	fmt.Fprintln(os.Stderr, string(jsonData))

	return err
}
