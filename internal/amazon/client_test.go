package amazon

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/zkwentz/amazon-cli/internal/config"
)

func TestPostForm(t *testing.T) {
	receivedMethod := ""
	receivedContentType := ""
	receivedBody := ""

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedContentType = r.Header.Get("Content-Type")
		bodyBytes, _ := io.ReadAll(r.Body)
		receivedBody = string(bodyBytes)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	cfg := config.DefaultConfig()
	cfg.RateLimiting.MinDelayMs = 0

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	formData := url.Values{}
	formData.Set("username", "testuser")
	formData.Set("password", "testpass")

	resp, err := client.PostForm(server.URL, formData)
	if err != nil {
		t.Fatalf("PostForm failed: %v", err)
	}
	defer resp.Body.Close()

	if receivedMethod != "POST" {
		t.Errorf("Expected POST method, got %s", receivedMethod)
	}

	if receivedContentType != "application/x-www-form-urlencoded" {
		t.Errorf("Expected Content-Type 'application/x-www-form-urlencoded', got %s", receivedContentType)
	}

	expectedBody := formData.Encode()
	if receivedBody != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, receivedBody)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if string(body) != "OK" {
		t.Errorf("Expected body 'OK', got %s", string(body))
	}
}

func TestPostFormWithRetry(t *testing.T) {
	attemptCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		if attemptCount < 3 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success"))
	}))
	defer server.Close()

	cfg := config.DefaultConfig()
	cfg.RateLimiting.MinDelayMs = 0
	cfg.RateLimiting.MaxRetries = 3

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	formData := url.Values{}
	formData.Set("test", "value")

	resp, err := client.PostForm(server.URL, formData)
	if err != nil {
		t.Fatalf("PostForm failed: %v", err)
	}
	defer resp.Body.Close()

	if attemptCount != 3 {
		t.Errorf("Expected 3 attempts, got %d", attemptCount)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestGet(t *testing.T) {
	receivedMethod := ""

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("GET OK"))
	}))
	defer server.Close()

	cfg := config.DefaultConfig()
	cfg.RateLimiting.MinDelayMs = 0

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	defer resp.Body.Close()

	if receivedMethod != "GET" {
		t.Errorf("Expected GET method, got %s", receivedMethod)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestUserAgentRotation(t *testing.T) {
	receivedUAs := make(map[string]bool)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ua := r.Header.Get("User-Agent")
		receivedUAs[ua] = true
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := config.DefaultConfig()
	cfg.RateLimiting.MinDelayMs = 0

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	for i := 0; i < 5; i++ {
		resp, err := client.Get(server.URL)
		if err != nil {
			t.Fatalf("Get failed on attempt %d: %v", i, err)
		}
		resp.Body.Close()
	}

	if len(receivedUAs) < 2 {
		t.Errorf("Expected user agent rotation, got only %d unique UAs", len(receivedUAs))
	}
}
