package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

type OutputFormat string

const (
	FormatJSON  OutputFormat = "json"
	FormatTable OutputFormat = "table"
	FormatRaw   OutputFormat = "raw"
)

type Printer struct {
	format OutputFormat
	quiet  bool
}

func NewPrinter(format string, quiet bool) *Printer {
	outputFormat := FormatJSON
	if format == string(FormatTable) {
		outputFormat = FormatTable
	} else if format == string(FormatRaw) {
		outputFormat = FormatRaw
	}

	return &Printer{
		format: outputFormat,
		quiet:  quiet,
	}
}

func (p *Printer) Print(data interface{}) error {
	switch p.format {
	case FormatJSON:
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(jsonData))
	case FormatTable, FormatRaw:
		// For now, just print as JSON. Table formatting can be added later
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(jsonData))
	}
	return nil
}

func (p *Printer) PrintError(err error) error {
	errorResponse := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    "UNKNOWN_ERROR",
			"message": err.Error(),
			"details": map[string]interface{}{},
		},
	}

	if cliErr, ok := err.(*models.CLIError); ok {
		errorResponse["error"] = map[string]interface{}{
			"code":    cliErr.Code,
			"message": cliErr.Message,
			"details": cliErr.Details,
		}
	}

	jsonData, jsonErr := json.MarshalIndent(errorResponse, "", "  ")
	if jsonErr != nil {
		return jsonErr
	}
	fmt.Fprintln(os.Stderr, string(jsonData))
	return nil
}
