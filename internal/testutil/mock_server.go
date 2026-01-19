package testutil

import (
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
)

// MockAmazonServer represents a mock HTTP server for testing Amazon-related functionality.
// It uses httptest.Server to simulate HTTP responses and allows configuring different
// fixture files to be served for different URL paths.
type MockAmazonServer struct {
	Server   *httptest.Server
	fixtures map[string]string
	mu       sync.RWMutex
}

// NewMockAmazonServer creates and returns a new MockAmazonServer instance.
// The server is started immediately and can be configured with fixture files
// using the ServeFixture method. The caller is responsible for closing the
// server when done by calling server.Server.Close().
func NewMockAmazonServer() *MockAmazonServer {
	mock := &MockAmazonServer{
		fixtures: make(map[string]string),
	}

	mock.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mock.mu.RLock()
		fixtureFile, exists := mock.fixtures[r.URL.Path]
		mock.mu.RUnlock()

		if !exists {
			http.NotFound(w, r)
			return
		}

		data, err := os.ReadFile(fixtureFile)
		if err != nil {
			http.Error(w, "Failed to read fixture file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}))

	return mock
}

// ServeFixture configures the mock server to serve the specified fixture file
// when the given URL path is requested. This allows tests to set up different
// responses for different endpoints.
//
// Example:
//
//	server := NewMockAmazonServer()
//	defer server.Server.Close()
//	server.ServeFixture("/product/B08N5WRWNW", "testdata/product.html")
func (m *MockAmazonServer) ServeFixture(path, fixtureFile string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.fixtures[path] = fixtureFile
}
