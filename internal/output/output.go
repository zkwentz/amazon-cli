package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/michaelshimeles/amazon-cli/pkg/models"
)

// OutputFormat represents the output format type
type OutputFormat string

const (
	// FormatJSON represents JSON output format
	FormatJSON OutputFormat = "json"
	// FormatTable represents table output format
	FormatTable OutputFormat = "table"
	// FormatRaw represents raw output format
	FormatRaw OutputFormat = "raw"
)

// Printer handles output formatting and printing
type Printer struct {
	format  OutputFormat
	quiet   bool
	verbose bool
	writer  io.Writer
}

// NewPrinter creates a new output printer with the specified format and options
func NewPrinter(format string, quiet bool, verbose bool) *Printer {
	return &Printer{
		format:  OutputFormat(format),
		quiet:   quiet,
		verbose: verbose,
		writer:  os.Stdout,
	}
}

// NewPrinterWithWriter creates a new output printer with a custom writer (useful for testing)
func NewPrinterWithWriter(format string, quiet bool, verbose bool, writer io.Writer) *Printer {
	return &Printer{
		format:  OutputFormat(format),
		quiet:   quiet,
		verbose: verbose,
		writer:  writer,
	}
}

// Print outputs data in the configured format
func (p *Printer) Print(data interface{}) error {
	if p.quiet {
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

// PrintError outputs an error in the configured format with verbose debug info if enabled
func (p *Printer) PrintError(err error) error {
	if err == nil {
		return nil
	}

	// Convert to CLIError if possible
	var cliErr *models.CLIError
	var ok bool

	if cliErr, ok = err.(*models.CLIError); !ok {
		// Not a CLIError, wrap it as a general error
		cliErr = models.NewCLIErrorWithCause(models.ErrAmazonError, err.Error(), err)
	}

	// Build the error response
	errDetail := &models.ErrorDetail{
		Code:    cliErr.Code,
		Message: cliErr.Message,
		Details: cliErr.Details,
	}

	// Add debug info if verbose flag is set
	if p.verbose {
		debugInfo := &models.DebugInfo{}

		if cliErr.Cause() != nil {
			debugInfo.Cause = cliErr.Cause().Error()
		}

		if cliErr.StackTrace() != "" {
			debugInfo.StackTrace = cliErr.StackTrace()
		}

		// Only add debug info if there's something to show
		if debugInfo.Cause != "" || debugInfo.StackTrace != "" {
			errDetail.Debug = debugInfo
		}
	}

	errResponse := &models.ErrorResponse{
		Error: errDetail,
	}

	// Print to stderr only if using default stdout writer
	originalWriter := p.writer
	if p.writer == os.Stdout {
		p.writer = os.Stderr
		defer func() { p.writer = originalWriter }()
	}

	return p.Print(errResponse)
}

// printJSON outputs data as formatted JSON
func (p *Printer) printJSON(data interface{}) error {
	encoder := json.NewEncoder(p.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// printTable outputs data as a formatted table
// This is a placeholder implementation that will be expanded later
func (p *Printer) printTable(data interface{}) error {
	// TODO: Implement table formatting using a library like tablewriter
	// For now, fall back to JSON
	return p.printJSON(data)
}

// printRaw outputs data as a raw string
func (p *Printer) printRaw(data interface{}) error {
	_, err := fmt.Fprintf(p.writer, "%v\n", data)
	return err
}

// GetVerbose returns the verbose flag setting
func (p *Printer) GetVerbose() bool {
	return p.verbose
}

// GetQuiet returns the quiet flag setting
func (p *Printer) GetQuiet() bool {
	return p.quiet
}
