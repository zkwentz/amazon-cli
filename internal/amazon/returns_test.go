package amazon

import (
	"net/http"
	"net/http/httptest"
	"net/http/cookiejar"
	"net/url"
	"testing"
	"time"

	"github.com/zkwentz/amazon-cli/internal/config"
	"github.com/zkwentz/amazon-cli/internal/ratelimit"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

func newTestClient(cfg *config.Config, testServer *httptest.Server) *Client {
	jar, _ := cookiejar.New(nil)

	serverURL, _ := url.Parse(testServer.URL)
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Jar:     jar,
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return serverURL, nil
			},
		},
	}

	return &Client{
		httpClient:  httpClient,
		rateLimiter: ratelimit.NewRateLimiter(cfg.RateLimiting),
		config:      cfg,
		userAgents:  []string{"test-agent"},
		currentUA:   0,
	}
}

func TestGetReturnLabel(t *testing.T) {
	tests := []struct {
		name           string
		returnID       string
		statusCode     int
		wantErr        bool
		wantErrCode    string
		checkResult    bool
	}{
		{
			name:        "valid return ID",
			returnID:    "RET123456",
			statusCode:  200,
			wantErr:     false,
			checkResult: true,
		},
		{
			name:        "empty return ID",
			returnID:    "",
			statusCode:  0,
			wantErr:     true,
			wantErrCode: models.ErrorCodeInvalidInput,
		},
		{
			name:        "return not found",
			returnID:    "INVALID123",
			statusCode:  404,
			wantErr:     true,
			wantErrCode: models.ErrorCodeNotFound,
		},
		{
			name:        "authentication expired",
			returnID:    "RET123456",
			statusCode:  401,
			wantErr:     true,
			wantErrCode: models.ErrorCodeAuthExpired,
		},
		{
			name:        "forbidden",
			returnID:    "RET123456",
			statusCode:  403,
			wantErr:     true,
			wantErrCode: models.ErrorCodeAuthExpired,
		},
		{
			name:        "server error",
			returnID:    "RET123456",
			statusCode:  500,
			wantErr:     true,
			wantErrCode: models.ErrorCodeAmazonError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.statusCode == 0 {
				cfg := &config.Config{
					RateLimiting: config.RateLimitConfig{
						MinDelayMs: 0,
						MaxDelayMs: 0,
						MaxRetries: 1,
					},
				}
				client := NewClient(cfg)
				result, err := client.GetReturnLabel(tt.returnID)

				if !tt.wantErr {
					t.Errorf("expected error but got none")
					return
				}

				if err == nil {
					t.Errorf("expected error but got none")
					return
				}

				cliErr, ok := err.(*models.CLIError)
				if !ok {
					t.Errorf("expected CLIError but got %T", err)
					return
				}

				if cliErr.Code != tt.wantErrCode {
					t.Errorf("expected error code %s but got %s", tt.wantErrCode, cliErr.Code)
				}

				if result != nil {
					t.Errorf("expected nil result but got %v", result)
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte("mock response"))
			}))
			defer server.Close()

			cfg := &config.Config{
				RateLimiting: config.RateLimitConfig{
					MinDelayMs: 0,
					MaxDelayMs: 0,
					MaxRetries: 1,
				},
			}

			jar, _ := cookiejar.New(nil)
			client := &Client{
				httpClient: &http.Client{
					Timeout: 30 * time.Second,
					Jar:     jar,
					Transport: &testRoundTripper{server: server},
				},
				rateLimiter: ratelimit.NewRateLimiter(cfg.RateLimiting),
				config:      cfg,
				userAgents:  []string{"test-agent"},
				currentUA:   0,
			}

			result, err := client.GetReturnLabel(tt.returnID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}

				cliErr, ok := err.(*models.CLIError)
				if !ok {
					t.Errorf("expected CLIError but got %T", err)
					return
				}

				if cliErr.Code != tt.wantErrCode {
					t.Errorf("expected error code %s but got %s", tt.wantErrCode, cliErr.Code)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}

				if tt.checkResult {
					if result == nil {
						t.Errorf("expected result but got nil")
						return
					}

					if result.ReturnID != tt.returnID {
						t.Errorf("expected return ID %s but got %s", tt.returnID, result.ReturnID)
					}

					if result.URL == "" {
						t.Errorf("expected URL but got empty string")
					}

					if result.Carrier == "" {
						t.Errorf("expected Carrier but got empty string")
					}

					if result.Instructions == "" {
						t.Errorf("expected Instructions but got empty string")
					}
				}
			}
		})
	}
}

func TestGetReturnLabelURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("mock label"))
	}))
	defer server.Close()

	cfg := &config.Config{
		RateLimiting: config.RateLimitConfig{
			MinDelayMs: 0,
			MaxDelayMs: 0,
			MaxRetries: 1,
		},
	}

	jar, _ := cookiejar.New(nil)
	client := &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Jar:     jar,
			Transport: &testRoundTripper{server: server},
		},
		rateLimiter: ratelimit.NewRateLimiter(cfg.RateLimiting),
		config:      cfg,
		userAgents:  []string{"test-agent"},
		currentUA:   0,
	}

	returnID := "RET123456"
	result, err := client.GetReturnLabel(returnID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedURL := "https://www.amazon.com/returns/label/RET123456/print"
	if result.URL != expectedURL {
		t.Errorf("expected URL %s but got %s", expectedURL, result.URL)
	}
}

type testRoundTripper struct {
	server *httptest.Server
}

func (t *testRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.server.Client().Get(t.server.URL)
}
