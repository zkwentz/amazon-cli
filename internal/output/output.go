package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// OutputFormat represents the format for outputting data
type OutputFormat string

const (
	// FormatJSON outputs data as JSON
	FormatJSON OutputFormat = "json"
	// FormatTable outputs data as a human-readable table
	FormatTable OutputFormat = "table"
	// FormatRaw outputs data in raw string format
	FormatRaw OutputFormat = "raw"
)

// Printer handles formatting and printing of output data
type Printer struct {
	format OutputFormat
	quiet  bool
	writer io.Writer
}

// NewPrinter creates a new Printer with the specified format and quiet mode
func NewPrinter(format string, quiet bool) *Printer {
	// Parse format string to OutputFormat
	outputFormat := FormatJSON // default to JSON
	switch format {
	case "json":
		outputFormat = FormatJSON
	case "table":
		outputFormat = FormatTable
	case "raw":
		outputFormat = FormatRaw
	default:
		outputFormat = FormatJSON
	}

	return &Printer{
		format: outputFormat,
		quiet:  quiet,
		writer: os.Stdout,
	}
}

// Print outputs the given data according to the printer's format
func (p *Printer) Print(data interface{}) error {
	if p.quiet {
		// In quiet mode, suppress non-essential output
		return nil
	}

	switch p.format {
	case FormatJSON:
		return p.printJSON(data)
	case FormatTable:
		return p.printTable(data)
	case FormatRaw:
		return p.printRaw(data)
	default:
		return p.printJSON(data)
	}
}

// printJSON marshals data to JSON and prints it
func (p *Printer) printJSON(data interface{}) error {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	_, err = fmt.Fprintln(p.writer, string(jsonBytes))
	return err
}

// printTable formats data as a human-readable table
func (p *Printer) printTable(data interface{}) error {
	// TODO: Implement table formatting using tablewriter package
	// For now, fall back to JSON
	return p.printJSON(data)
}

// printRaw prints the raw string representation of data
func (p *Printer) printRaw(data interface{}) error {
	_, err := fmt.Fprintln(p.writer, data)
	return err
}

// PrintError outputs an error in the appropriate format
func (p *Printer) PrintError(err error) error {
	errorOutput := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    "UNKNOWN_ERROR",
			"message": err.Error(),
			"details": map[string]interface{}{},
		},
	}

	// Always output errors, even in quiet mode
	jsonBytes, jsonErr := json.MarshalIndent(errorOutput, "", "  ")
	if jsonErr != nil {
		return fmt.Errorf("failed to marshal error JSON: %w", jsonErr)
	}
	_, writeErr := fmt.Fprintln(p.writer, string(jsonBytes))
	return writeErr
}

// SetWriter sets the output writer (useful for testing)
func (p *Printer) SetWriter(w io.Writer) {
	p.writer = w
}
