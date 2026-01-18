package amazon

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/zkwentz/amazon-cli/internal/config"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

// mockHTTPClient is a mock implementation of HTTPClient for testing
type mockHTTPClient struct {
	getFunc func(url string) (*http.Response, error)
	doFunc  func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Get(url string) (*http.Response, error) {
	if m.getFunc != nil {
		return m.getFunc(url)
	}
	return nil, fmt.Errorf("mock Get not implemented")
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.doFunc != nil {
		return m.doFunc(req)
	}
	return nil, fmt.Errorf("mock Do not implemented")
}

func TestGetReturnStatus_Success(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return a mock HTML response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			<html>
				<body>
					<div class="return-status">
						<h1>Return Status: Received</h1>
						<p>Order #123-4567890-1234567</p>
						<p>Return initiated on 2024-01-15</p>
						<p>Reason: defective</p>
					</div>
				</body>
			</html>
		`))
	}))
	defer server.Close()

	// Create client with mock HTTP client
	cfg := &config.Config{}
	mockClient := &mockHTTPClient{
		getFunc: func(url string) (*http.Response, error) {
			// Redirect all requests to test server
			testURL := server.URL + "/returns/status/R123456789"
			return http.Get(testURL)
		},
	}
	client := NewClientWithHTTPClient(cfg, mockClient)

	// Call GetReturnStatus
	returnStatus, err := client.GetReturnStatus("R123456789")

	// Verify no error
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify return data
	if returnStatus == nil {
		t.Fatal("Expected return status, got nil")
	}

	if returnStatus.ReturnID != "R123456789" {
		t.Errorf("Expected return ID R123456789, got %s", returnStatus.ReturnID)
	}

	if returnStatus.Status != "received" {
		t.Errorf("Expected status 'received', got %s", returnStatus.Status)
	}

	if returnStatus.OrderID != "123-4567890-1234567" {
		t.Errorf("Expected order ID 123-4567890-1234567, got %s", returnStatus.OrderID)
	}

	if returnStatus.Reason != "defective" {
		t.Errorf("Expected reason 'defective', got %s", returnStatus.Reason)
	}
}

func TestGetReturnStatus_EmptyReturnID(t *testing.T) {
	cfg := &config.Config{}
	client := NewClient(cfg)

	// Call with empty return ID
	_, err := client.GetReturnStatus("")

	// Verify error
	if err == nil {
		t.Fatal("Expected error for empty return ID, got nil")
	}

	cliErr, ok := err.(*models.CLIError)
	if !ok {
		t.Fatalf("Expected CLIError, got %T", err)
	}

	if cliErr.Code != models.ErrInvalidInput {
		t.Errorf("Expected error code %s, got %s", models.ErrInvalidInput, cliErr.Code)
	}
}

func TestGetReturnStatus_NotFound(t *testing.T) {
	cfg := &config.Config{}
	mockClient := &mockHTTPClient{
		getFunc: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(strings.NewReader("Not Found")),
			}, nil
		},
	}
	client := NewClientWithHTTPClient(cfg, mockClient)

	// Call GetReturnStatus
	_, err := client.GetReturnStatus("INVALID123")

	// Verify error
	if err == nil {
		t.Fatal("Expected error for not found, got nil")
	}

	cliErr, ok := err.(*models.CLIError)
	if !ok {
		t.Fatalf("Expected CLIError, got %T", err)
	}

	if cliErr.Code != models.ErrNotFound {
		t.Errorf("Expected error code %s, got %s", models.ErrNotFound, cliErr.Code)
	}
}

func TestGetReturnStatus_Unauthorized(t *testing.T) {
	cfg := &config.Config{}
	mockClient := &mockHTTPClient{
		getFunc: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       io.NopCloser(strings.NewReader("Unauthorized")),
			}, nil
		},
	}
	client := NewClientWithHTTPClient(cfg, mockClient)

	// Call GetReturnStatus
	_, err := client.GetReturnStatus("R123456789")

	// Verify error
	if err == nil {
		t.Fatal("Expected error for unauthorized, got nil")
	}

	cliErr, ok := err.(*models.CLIError)
	if !ok {
		t.Fatalf("Expected CLIError, got %T", err)
	}

	if cliErr.Code != models.ErrAuthExpired {
		t.Errorf("Expected error code %s, got %s", models.ErrAuthExpired, cliErr.Code)
	}
}

func TestGetReturnStatus_NetworkError(t *testing.T) {
	cfg := &config.Config{}
	mockClient := &mockHTTPClient{
		getFunc: func(url string) (*http.Response, error) {
			return nil, fmt.Errorf("network error: connection refused")
		},
	}
	client := NewClientWithHTTPClient(cfg, mockClient)

	// Call GetReturnStatus
	_, err := client.GetReturnStatus("R123456789")

	// Verify error
	if err == nil {
		t.Fatal("Expected network error, got nil")
	}

	cliErr, ok := err.(*models.CLIError)
	if !ok {
		t.Fatalf("Expected CLIError, got %T", err)
	}

	if cliErr.Code != models.ErrNetworkError {
		t.Errorf("Expected error code %s, got %s", models.ErrNetworkError, cliErr.Code)
	}
}

func TestGetReturnStatus_VariousStatuses(t *testing.T) {
	testCases := []struct {
		name           string
		htmlContent    string
		expectedStatus string
	}{
		{
			name:           "Initiated status",
			htmlContent:    "<html><body>Return initiated</body></html>",
			expectedStatus: "initiated",
		},
		{
			name:           "Shipped status",
			htmlContent:    "<html><body>Return shipped</body></html>",
			expectedStatus: "shipped",
		},
		{
			name:           "Received status",
			htmlContent:    "<html><body>Return received</body></html>",
			expectedStatus: "received",
		},
		{
			name:           "Refunded status",
			htmlContent:    "<html><body>Refunded to your account</body></html>",
			expectedStatus: "refunded",
		},
		{
			name:           "Pending status (default)",
			htmlContent:    "<html><body>Processing your return</body></html>",
			expectedStatus: "pending",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &config.Config{}
			mockClient := &mockHTTPClient{
				getFunc: func(url string) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(tc.htmlContent)),
					}, nil
				},
			}
			client := NewClientWithHTTPClient(cfg, mockClient)

			returnStatus, err := client.GetReturnStatus("R123456789")
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if returnStatus.Status != tc.expectedStatus {
				t.Errorf("Expected status %s, got %s", tc.expectedStatus, returnStatus.Status)
			}
		})
	}
}

func TestGetReturnLabel_Success(t *testing.T) {
	cfg := &config.Config{}
	mockClient := &mockHTTPClient{
		getFunc: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(strings.NewReader(`
					<html>
						<body>
							<div>Return label for UPS</div>
						</body>
					</html>
				`)),
			}, nil
		},
	}
	client := NewClientWithHTTPClient(cfg, mockClient)

	label, err := client.GetReturnLabel("R123456789")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if label.Carrier != "UPS" {
		t.Errorf("Expected carrier UPS, got %s", label.Carrier)
	}

	expectedURL := "https://www.amazon.com/returns/label/R123456789.pdf"
	if label.URL != expectedURL {
		t.Errorf("Expected URL %s, got %s", expectedURL, label.URL)
	}
}

func TestGetReturnLabel_EmptyReturnID(t *testing.T) {
	cfg := &config.Config{}
	client := NewClient(cfg)

	_, err := client.GetReturnLabel("")
	if err == nil {
		t.Fatal("Expected error for empty return ID, got nil")
	}

	cliErr, ok := err.(*models.CLIError)
	if !ok {
		t.Fatalf("Expected CLIError, got %T", err)
	}

	if cliErr.Code != models.ErrInvalidInput {
		t.Errorf("Expected error code %s, got %s", models.ErrInvalidInput, cliErr.Code)
	}
}
