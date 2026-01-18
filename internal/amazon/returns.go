package amazon

import (
	"fmt"
	"io"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

func (c *Client) GetReturnableItems() ([]models.ReturnableItem, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *Client) GetReturnOptions(orderID, itemID string) ([]models.ReturnOption, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *Client) CreateReturn(orderID, itemID, reason string) (*models.Return, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *Client) GetReturnLabel(returnID string) (*models.ReturnLabel, error) {
	if returnID == "" {
		return nil, models.NewCLIError(
			models.ErrorCodeInvalidInput,
			"return ID is required",
			nil,
		)
	}

	url := fmt.Sprintf("https://www.amazon.com/returns/label/%s", returnID)

	resp, err := c.Get(url)
	if err != nil {
		return nil, models.NewCLIError(
			models.ErrorCodeNetworkError,
			fmt.Sprintf("failed to fetch return label: %v", err),
			nil,
		)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, models.NewCLIError(
			models.ErrorCodeNotFound,
			fmt.Sprintf("return label not found for return ID: %s", returnID),
			nil,
		)
	}

	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		return nil, models.NewCLIError(
			models.ErrorCodeAuthExpired,
			"authentication token has expired, please run 'amazon-cli auth login'",
			nil,
		)
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, models.NewCLIError(
			models.ErrorCodeAmazonError,
			fmt.Sprintf("Amazon returned error: status %d", resp.StatusCode),
			map[string]interface{}{
				"status_code": resp.StatusCode,
				"body":        string(body),
			},
		)
	}

	label := &models.ReturnLabel{
		ReturnID:     returnID,
		URL:          fmt.Sprintf("https://www.amazon.com/returns/label/%s/print", returnID),
		Carrier:      "UPS",
		Instructions: "Print this label and attach it to your package. Drop off at any UPS location or schedule a pickup.",
	}

	return label, nil
}

func (c *Client) GetReturnStatus(returnID string) (*models.Return, error) {
	return nil, fmt.Errorf("not implemented")
}
