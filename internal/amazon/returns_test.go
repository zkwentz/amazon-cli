package amazon

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetReturnLabel(t *testing.T) {
	tests := []struct {
		name           string
		returnID       string
		statusCode     int
		responseBody   string
		expectError    bool
		expectedURL    string
		expectedCarrier string
	}{
		{
			name:        "empty return ID",
			returnID:    "",
			expectError: true,
		},
		{
			name:         "successful label retrieval with JSON format",
			returnID:     "R123456789",
			statusCode:   200,
			responseBody: `{"labelUrl":"https://amazon.com/label/R123456789.pdf","carrier":"UPS"}`,
			expectError:  false,
			expectedURL:  "https://amazon.com/label/R123456789.pdf",
			expectedCarrier: "UPS",
		},
		{
			name:         "successful label retrieval with HTML format",
			returnID:     "R987654321",
			statusCode:   200,
			responseBody: `<html><body><div>UPS Return Label</div><a href="https://amazon.com/returns/label.pdf">Download Label</a></body></html>`,
			expectError:  false,
			expectedCarrier: "UPS",
		},
		{
			name:        "404 not found",
			returnID:    "INVALID123",
			statusCode:  404,
			expectError: true,
		},
		{
			name:        "500 server error",
			returnID:    "R123456789",
			statusCode:  500,
			expectError: true,
		},
		{
			name:        "401 unauthorized",
			returnID:    "R123456789",
			statusCode:  401,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip server setup for empty return ID test
			if tt.returnID == "" {
				client := NewClient(nil)
				_, err := client.GetReturnLabel(tt.returnID)
				if err == nil {
					t.Error("expected error for empty return ID but got none")
				}
				return
			}

			// Create a test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify the URL contains the return ID
				if !strings.Contains(r.URL.String(), tt.returnID) {
					t.Errorf("expected URL to contain return ID %s, got %s", tt.returnID, r.URL.String())
				}

				w.WriteHeader(tt.statusCode)
				if tt.responseBody != "" {
					w.Write([]byte(tt.responseBody))
				}
			}))
			defer server.Close()

			// Create a client with test server
			client := newTestClient(server)

			// Call GetReturnLabel
			label, err := client.GetReturnLabel(tt.returnID)

			// Check error expectation
			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Validate the label
			if label == nil {
				t.Error("expected label but got nil")
				return
			}

			// Check URL if specified
			if tt.expectedURL != "" && label.URL != tt.expectedURL {
				t.Errorf("expected URL %s, got %s", tt.expectedURL, label.URL)
			}

			// Check carrier if specified
			if tt.expectedCarrier != "" && label.Carrier != tt.expectedCarrier {
				t.Errorf("expected carrier %s, got %s", tt.expectedCarrier, label.Carrier)
			}

			// Validate that we have at least a URL
			if label.URL == "" {
				t.Error("expected non-empty URL in label")
			}

			// Validate that we have instructions
			if label.Instructions == "" {
				t.Error("expected non-empty instructions in label")
			}
		})
	}
}

func TestGetReturnLabelWithRetries(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		// Fail first two attempts with 503
		if attempts < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		// Succeed on third attempt
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><body>UPS label available</body></html>`))
	}))
	defer server.Close()

	client := newTestClient(server)

	label, err := client.GetReturnLabel("R123456789")

	if err != nil {
		t.Errorf("expected successful retry, got error: %v", err)
	}

	if label == nil {
		t.Error("expected label after retries")
	}

	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestGetReturnLabelIntegration(t *testing.T) {
	// This is a more realistic integration test
	testCases := []struct {
		name       string
		returnID   string
		setupMock  func() *httptest.Server
		validate   func(*testing.T, interface{}, error)
	}{
		{
			name:     "complete successful flow",
			returnID: "R111222333",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					response := fmt.Sprintf(`
						<html>
							<body>
								<div class="return-label">
									<h1>UPS Return Label</h1>
									<a href="https://amazon.com/returns/label/%s.pdf">Download Label</a>
									<p>Print this label and drop off at any UPS location.</p>
								</div>
							</body>
						</html>
					`, r.URL.Query().Get("returnId"))
					w.Write([]byte(response))
				}))
			},
			validate: func(t *testing.T, result interface{}, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				// Type assertion would normally be needed here
				// but we're keeping it simple for this test
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := tc.setupMock()
			defer server.Close()

			client := newTestClient(server)

			result, err := client.GetReturnLabel(tc.returnID)
			tc.validate(t, result, err)
		})
	}
}

// newTestClient creates a client that redirects all requests to the test server
func newTestClient(server *httptest.Server) *Client {
	client := NewClient(nil)

	// Replace the http client with one that redirects to test server
	testTransport := &testTransport{
		serverURL: server.URL,
		base:      http.DefaultTransport,
	}

	client.httpClient = &http.Client{
		Transport: testTransport,
	}

	return client
}

// testTransport redirects all requests to the test server
type testTransport struct {
	serverURL string
	base      http.RoundTripper
}

func (t *testTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Redirect to test server
	req.URL.Scheme = "http"
	req.URL.Host = strings.TrimPrefix(t.serverURL, "http://")
	return t.base.RoundTrip(req)
}

