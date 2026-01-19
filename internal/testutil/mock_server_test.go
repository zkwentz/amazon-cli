package testutil

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestNewMockAmazonServer(t *testing.T) {
	server := NewMockAmazonServer()
	defer server.Server.Close()

	if server.Server == nil {
		t.Error("NewMockAmazonServer returned nil Server")
	}

	if server.fixtures == nil {
		t.Error("NewMockAmazonServer did not initialize fixtures map")
	}

	if server.Server.URL == "" {
		t.Error("Server URL is empty")
	}
}

func TestServeFixture_Success(t *testing.T) {
	server := NewMockAmazonServer()
	defer server.Server.Close()

	// Create a temporary fixture file
	tmpDir := t.TempDir()
	fixtureFile := filepath.Join(tmpDir, "test.html")
	expectedContent := "<html><body>Test Content</body></html>"
	err := os.WriteFile(fixtureFile, []byte(expectedContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test fixture: %v", err)
	}

	// Configure the server to serve the fixture
	server.ServeFixture("/test-path", fixtureFile)

	// Make a request to the server
	resp, err := http.Get(server.Server.URL + "/test-path")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Verify response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if string(body) != expectedContent {
		t.Errorf("Expected body %q, got %q", expectedContent, string(body))
	}

	// Verify Content-Type header
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Errorf("Expected Content-Type to contain 'text/html', got %q", contentType)
	}
}

func TestServeFixture_NotFound(t *testing.T) {
	server := NewMockAmazonServer()
	defer server.Server.Close()

	// Request a path that hasn't been configured
	resp, err := http.Get(server.Server.URL + "/non-existent-path")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}
}

func TestServeFixture_InvalidFile(t *testing.T) {
	server := NewMockAmazonServer()
	defer server.Server.Close()

	// Configure with a non-existent file
	server.ServeFixture("/invalid", "/path/to/non/existent/file.html")

	resp, err := http.Get(server.Server.URL + "/invalid")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	bodyStr := string(body)
	if !strings.Contains(bodyStr, "Failed to read fixture file") {
		t.Errorf("Expected error message about fixture file, got %q", bodyStr)
	}
}

func TestServeFixture_MultiplePaths(t *testing.T) {
	server := NewMockAmazonServer()
	defer server.Server.Close()

	tmpDir := t.TempDir()

	// Create multiple fixture files
	fixture1 := filepath.Join(tmpDir, "fixture1.html")
	content1 := "<html><body>Fixture 1</body></html>"
	err := os.WriteFile(fixture1, []byte(content1), 0644)
	if err != nil {
		t.Fatalf("Failed to create fixture1: %v", err)
	}

	fixture2 := filepath.Join(tmpDir, "fixture2.html")
	content2 := "<html><body>Fixture 2</body></html>"
	err = os.WriteFile(fixture2, []byte(content2), 0644)
	if err != nil {
		t.Fatalf("Failed to create fixture2: %v", err)
	}

	// Configure multiple paths
	server.ServeFixture("/path1", fixture1)
	server.ServeFixture("/path2", fixture2)

	// Test first path
	resp1, err := http.Get(server.Server.URL + "/path1")
	if err != nil {
		t.Fatalf("Request to /path1 failed: %v", err)
	}
	defer resp1.Body.Close()

	body1, _ := io.ReadAll(resp1.Body)
	if string(body1) != content1 {
		t.Errorf("Expected /path1 to return %q, got %q", content1, string(body1))
	}

	// Test second path
	resp2, err := http.Get(server.Server.URL + "/path2")
	if err != nil {
		t.Fatalf("Request to /path2 failed: %v", err)
	}
	defer resp2.Body.Close()

	body2, _ := io.ReadAll(resp2.Body)
	if string(body2) != content2 {
		t.Errorf("Expected /path2 to return %q, got %q", content2, string(body2))
	}
}

func TestServeFixture_OverwritePath(t *testing.T) {
	server := NewMockAmazonServer()
	defer server.Server.Close()

	tmpDir := t.TempDir()

	// Create two fixture files
	fixture1 := filepath.Join(tmpDir, "fixture1.html")
	content1 := "<html><body>First Content</body></html>"
	err := os.WriteFile(fixture1, []byte(content1), 0644)
	if err != nil {
		t.Fatalf("Failed to create fixture1: %v", err)
	}

	fixture2 := filepath.Join(tmpDir, "fixture2.html")
	content2 := "<html><body>Second Content</body></html>"
	err = os.WriteFile(fixture2, []byte(content2), 0644)
	if err != nil {
		t.Fatalf("Failed to create fixture2: %v", err)
	}

	// Configure path with first fixture
	server.ServeFixture("/test", fixture1)

	// Verify first fixture is served
	resp1, err := http.Get(server.Server.URL + "/test")
	if err != nil {
		t.Fatalf("First request failed: %v", err)
	}
	defer resp1.Body.Close()
	body1, _ := io.ReadAll(resp1.Body)
	if string(body1) != content1 {
		t.Errorf("Expected first content, got %q", string(body1))
	}

	// Overwrite with second fixture
	server.ServeFixture("/test", fixture2)

	// Verify second fixture is now served
	resp2, err := http.Get(server.Server.URL + "/test")
	if err != nil {
		t.Fatalf("Second request failed: %v", err)
	}
	defer resp2.Body.Close()
	body2, _ := io.ReadAll(resp2.Body)
	if string(body2) != content2 {
		t.Errorf("Expected second content, got %q", string(body2))
	}
}

func TestServeFixture_ConcurrentRequests(t *testing.T) {
	server := NewMockAmazonServer()
	defer server.Server.Close()

	tmpDir := t.TempDir()
	fixtureFile := filepath.Join(tmpDir, "concurrent.html")
	content := "<html><body>Concurrent Test</body></html>"
	err := os.WriteFile(fixtureFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create fixture: %v", err)
	}

	server.ServeFixture("/concurrent", fixtureFile)

	// Make concurrent requests
	numRequests := 10
	var wg sync.WaitGroup
	errors := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			resp, err := http.Get(server.Server.URL + "/concurrent")
			if err != nil {
				errors <- err
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				errors <- err
				return
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				errors <- err
				return
			}

			if string(body) != content {
				errors <- err
			}
		}()
	}

	wg.Wait()
	close(errors)

	// Check for any errors
	for err := range errors {
		if err != nil {
			t.Errorf("Concurrent request error: %v", err)
		}
	}
}

func TestServeFixture_WithQueryParams(t *testing.T) {
	server := NewMockAmazonServer()
	defer server.Server.Close()

	tmpDir := t.TempDir()
	fixtureFile := filepath.Join(tmpDir, "query.html")
	content := "<html><body>Query Test</body></html>"
	err := os.WriteFile(fixtureFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create fixture: %v", err)
	}

	// Configure path without query params
	server.ServeFixture("/search", fixtureFile)

	// Request with query params - should match by path only
	resp, err := http.Get(server.Server.URL + "/search?q=test&page=1")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != content {
		t.Errorf("Expected content %q, got %q", content, string(body))
	}
}

func TestServeFixture_RealFixtures(t *testing.T) {
	server := NewMockAmazonServer()
	defer server.Server.Close()

	// Test with real fixture files if they exist
	sampleFixture := "testdata/sample.html"
	if _, err := os.Stat(sampleFixture); err == nil {
		server.ServeFixture("/product", sampleFixture)

		resp, err := http.Get(server.Server.URL + "/product")
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read body: %v", err)
		}

		// Verify it contains expected HTML content
		bodyStr := string(body)
		if !strings.Contains(bodyStr, "Sample Product") {
			t.Errorf("Expected body to contain 'Sample Product', got %q", bodyStr)
		}
	}

	searchFixture := "testdata/search_results.html"
	if _, err := os.Stat(searchFixture); err == nil {
		server.ServeFixture("/search", searchFixture)

		resp, err := http.Get(server.Server.URL + "/search")
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read body: %v", err)
		}

		// Verify it contains expected HTML content
		bodyStr := string(body)
		if !strings.Contains(bodyStr, "Search Results") {
			t.Errorf("Expected body to contain 'Search Results', got %q", bodyStr)
		}
	}
}
