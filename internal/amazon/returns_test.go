package amazon

import (
	"testing"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

func TestGetReturnableItems(t *testing.T) {
	client := NewClient()

	response, err := client.GetReturnableItems()
	if err != nil {
		t.Fatalf("GetReturnableItems failed: %v", err)
	}

	if response == nil {
		t.Fatal("Response is nil")
	}

	if response.TotalCount == 0 {
		t.Error("Expected non-zero total count")
	}

	if len(response.Items) != response.TotalCount {
		t.Errorf("Expected %d items, got %d", response.TotalCount, len(response.Items))
	}

	// Verify structure of first item
	if len(response.Items) > 0 {
		item := response.Items[0]
		if item.OrderID == "" {
			t.Error("Expected non-empty OrderID")
		}
		if item.ItemID == "" {
			t.Error("Expected non-empty ItemID")
		}
		if item.ASIN == "" {
			t.Error("Expected non-empty ASIN")
		}
		if item.Title == "" {
			t.Error("Expected non-empty Title")
		}
		if item.Price <= 0 {
			t.Error("Expected positive price")
		}
		if item.PurchaseDate == "" {
			t.Error("Expected non-empty PurchaseDate")
		}
		if item.ReturnWindow == "" {
			t.Error("Expected non-empty ReturnWindow")
		}
	}
}

func TestGetReturnOptions(t *testing.T) {
	client := NewClient()

	_, err := client.GetReturnOptions("test-order", "test-item")
	if err == nil {
		t.Error("Expected error for unimplemented method")
	}

	cliErr, ok := err.(*models.CLIError)
	if !ok {
		t.Error("Expected CLIError type")
	}

	if cliErr.Code != models.ErrCodeAmazonError {
		t.Errorf("Expected error code %s, got %s", models.ErrCodeAmazonError, cliErr.Code)
	}
}

func TestCreateReturn(t *testing.T) {
	client := NewClient()

	_, err := client.CreateReturn("test-order", "test-item", "defective")
	if err == nil {
		t.Error("Expected error for unimplemented method")
	}

	cliErr, ok := err.(*models.CLIError)
	if !ok {
		t.Error("Expected CLIError type")
	}

	if cliErr.Code != models.ErrCodeAmazonError {
		t.Errorf("Expected error code %s, got %s", models.ErrCodeAmazonError, cliErr.Code)
	}
}

func TestGetReturnLabel(t *testing.T) {
	client := NewClient()

	_, err := client.GetReturnLabel("test-return")
	if err == nil {
		t.Error("Expected error for unimplemented method")
	}

	cliErr, ok := err.(*models.CLIError)
	if !ok {
		t.Error("Expected CLIError type")
	}

	if cliErr.Code != models.ErrCodeAmazonError {
		t.Errorf("Expected error code %s, got %s", models.ErrCodeAmazonError, cliErr.Code)
	}
}

func TestGetReturnStatus(t *testing.T) {
	client := NewClient()

	_, err := client.GetReturnStatus("test-return")
	if err == nil {
		t.Error("Expected error for unimplemented method")
	}

	cliErr, ok := err.(*models.CLIError)
	if !ok {
		t.Error("Expected CLIError type")
	}

	if cliErr.Code != models.ErrCodeAmazonError {
		t.Errorf("Expected error code %s, got %s", models.ErrCodeAmazonError, cliErr.Code)
	}
}
