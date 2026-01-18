package amazon

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

func TestGetAddresses_Success_JSON(t *testing.T) {
	// Create mock addresses
	mockAddresses := []models.Address{
		{
			ID:      "addr_1",
			Name:    "John Doe",
			Street:  "123 Main St",
			City:    "Seattle",
			State:   "WA",
			Zip:     "98101",
			Country: "US",
			Default: true,
		},
		{
			ID:      "addr_2",
			Name:    "Jane Doe",
			Street:  "456 Oak Ave",
			City:    "Portland",
			State:   "OR",
			Zip:     "97201",
			Country: "US",
			Default: false,
		},
	}

	// Create test server that returns JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		if r.URL.Path != "/a/addresses" {
			t.Errorf("Expected path /a/addresses, got %s", r.URL.Path)
		}

		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		// Check headers
		userAgent := r.Header.Get("User-Agent")
		if userAgent == "" {
			t.Error("Expected User-Agent header to be set")
		}

		// Return mock addresses as JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockAddresses)
	}))
	defer server.Close()

	// Create client with test server URL
	client := &Client{
		httpClient: server.Client(),
		baseURL:    server.URL,
	}

	// Test GetAddresses
	addresses, err := client.GetAddresses()
	if err != nil {
		t.Fatalf("GetAddresses() returned error: %v", err)
	}

	// Verify results
	if len(addresses) != len(mockAddresses) {
		t.Errorf("Expected %d addresses, got %d", len(mockAddresses), len(addresses))
	}

	for i, addr := range addresses {
		expected := mockAddresses[i]
		if addr.ID != expected.ID {
			t.Errorf("Address[%d].ID = %s, want %s", i, addr.ID, expected.ID)
		}
		if addr.Name != expected.Name {
			t.Errorf("Address[%d].Name = %s, want %s", i, addr.Name, expected.Name)
		}
		if addr.Street != expected.Street {
			t.Errorf("Address[%d].Street = %s, want %s", i, addr.Street, expected.Street)
		}
		if addr.City != expected.City {
			t.Errorf("Address[%d].City = %s, want %s", i, addr.City, expected.City)
		}
		if addr.State != expected.State {
			t.Errorf("Address[%d].State = %s, want %s", i, addr.State, expected.State)
		}
		if addr.Zip != expected.Zip {
			t.Errorf("Address[%d].Zip = %s, want %s", i, addr.Zip, expected.Zip)
		}
		if addr.Country != expected.Country {
			t.Errorf("Address[%d].Country = %s, want %s", i, addr.Country, expected.Country)
		}
		if addr.Default != expected.Default {
			t.Errorf("Address[%d].Default = %v, want %v", i, addr.Default, expected.Default)
		}
	}
}

func TestGetAddresses_Success_HTML(t *testing.T) {
	// Create test server that returns HTML
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><body>Address book page</body></html>"))
	}))
	defer server.Close()

	client := &Client{
		httpClient: server.Client(),
		baseURL:    server.URL,
	}

	// Test GetAddresses with HTML response
	// Currently returns empty list as HTML parsing is not implemented
	addresses, err := client.GetAddresses()
	if err != nil {
		t.Fatalf("GetAddresses() returned error: %v", err)
	}

	// For now, HTML responses return empty list
	if len(addresses) != 0 {
		t.Errorf("Expected empty addresses for HTML response, got %d", len(addresses))
	}
}

func TestGetAddresses_EmptyList(t *testing.T) {
	// Create test server that returns empty JSON array
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]models.Address{})
	}))
	defer server.Close()

	client := &Client{
		httpClient: server.Client(),
		baseURL:    server.URL,
	}

	addresses, err := client.GetAddresses()
	if err != nil {
		t.Fatalf("GetAddresses() returned error: %v", err)
	}

	if len(addresses) != 0 {
		t.Errorf("Expected 0 addresses, got %d", len(addresses))
	}
}

func TestGetAddresses_HTTPError(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
	}{
		{"Unauthorized", http.StatusUnauthorized},
		{"Forbidden", http.StatusForbidden},
		{"NotFound", http.StatusNotFound},
		{"InternalServerError", http.StatusInternalServerError},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				w.Write([]byte("Error message"))
			}))
			defer server.Close()

			client := &Client{
				httpClient: server.Client(),
				baseURL:    server.URL,
			}

			addresses, err := client.GetAddresses()
			if err == nil {
				t.Error("Expected error, got nil")
			}

			if addresses != nil {
				t.Errorf("Expected nil addresses on error, got %v", addresses)
			}
		})
	}
}

func TestGetAddresses_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := &Client{
		httpClient: server.Client(),
		baseURL:    server.URL,
	}

	addresses, err := client.GetAddresses()
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}

	if addresses != nil {
		t.Errorf("Expected nil addresses on error, got %v", addresses)
	}
}

func TestGetAddresses_NetworkError(t *testing.T) {
	// Use invalid URL to trigger network error
	client := &Client{
		httpClient: &http.Client{},
		baseURL:    "http://invalid-host-that-does-not-exist:99999",
	}

	addresses, err := client.GetAddresses()
	if err == nil {
		t.Error("Expected network error, got nil")
	}

	if addresses != nil {
		t.Errorf("Expected nil addresses on error, got %v", addresses)
	}
}

func TestNewClient(t *testing.T) {
	client := NewClient()

	if client == nil {
		t.Fatal("NewClient() returned nil")
	}

	if client.httpClient == nil {
		t.Error("Expected httpClient to be initialized")
	}

	if client.baseURL != "https://www.amazon.com" {
		t.Errorf("Expected baseURL to be https://www.amazon.com, got %s", client.baseURL)
	}
}
