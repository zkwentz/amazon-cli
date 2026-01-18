package amazon

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/michaelshimeles/amazon-cli/pkg/models"
)

func TestDetectHTMLResponse_JSON(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
	}
	body := []byte(`{"status": "success"}`)

	responseType, isHTML := DetectHTMLResponse(resp, body)

	if isHTML {
		t.Error("Expected non-HTML response, got HTML")
	}
	if responseType != HTMLResponseNone {
		t.Errorf("Expected HTMLResponseNone, got %v", responseType)
	}
}

func TestDetectHTMLResponse_Captcha(t *testing.T) {
	testCases := []struct {
		name string
		body string
	}{
		{
			name: "captcha keyword",
			body: `<html><body><h1>Please solve this CAPTCHA</h1></body></html>`,
		},
		{
			name: "robot check",
			body: `<html><body><p>Robot Check - We need to verify you're not a robot</p></body></html>`,
		},
		{
			name: "automated access",
			body: `<html><body><p>Automated access detected. Please complete verification.</p></body></html>`,
		},
		{
			name: "security check",
			body: `<html><body><h1>Security Check</h1><p>Please verify your identity</p></body></html>`,
		},
		{
			name: "Amazon style message",
			body: `<html><body><p>Sorry, we just need to make sure you're not a robot</p></body></html>`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := &http.Response{
				Header: http.Header{
					"Content-Type": []string{"text/html"},
				},
			}

			responseType, isHTML := DetectHTMLResponse(resp, []byte(tc.body))

			if !isHTML {
				t.Error("Expected HTML response detection")
			}
			if responseType != HTMLResponseCaptcha {
				t.Errorf("Expected HTMLResponseCaptcha, got %v", responseType)
			}
		})
	}
}

func TestDetectHTMLResponse_Login(t *testing.T) {
	testCases := []struct {
		name string
		body string
	}{
		{
			name: "signin page",
			body: `<html><body><form action="/ap_signin"><input name="ap_email" /></form></body></html>`,
		},
		{
			name: "authentication required",
			body: `<html><body><h1>Authentication Required</h1><p>Please sign in to continue</p></body></html>`,
		},
		{
			name: "sign-in with hyphen",
			body: `<html><body><a href="/sign-in">Sign-In</a></body></html>`,
		},
		{
			name: "your account page",
			body: `<html><body><h1>Your Account</h1><p>Sign in to access your account</p></body></html>`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := &http.Response{
				Header: http.Header{
					"Content-Type": []string{"text/html; charset=UTF-8"},
				},
			}

			responseType, isHTML := DetectHTMLResponse(resp, []byte(tc.body))

			if !isHTML {
				t.Error("Expected HTML response detection")
			}
			if responseType != HTMLResponseLogin {
				t.Errorf("Expected HTMLResponseLogin, got %v", responseType)
			}
		})
	}
}

func TestDetectHTMLResponse_UnknownHTML(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{
			"Content-Type": []string{"text/html"},
		},
	}
	body := []byte(`<html><body><h1>Some Other Page</h1><p>This is an unexpected page</p></body></html>`)

	responseType, isHTML := DetectHTMLResponse(resp, body)

	if !isHTML {
		t.Error("Expected HTML response detection")
	}
	if responseType != HTMLResponseUnknown {
		t.Errorf("Expected HTMLResponseUnknown, got %v", responseType)
	}
}

func TestHandleHTMLResponse_Captcha(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": []string{"text/html"},
		},
		Request: &http.Request{
			URL: &url.URL{
				Scheme: "https",
				Host:   "www.amazon.com",
				Path:   "/cart",
			},
		},
	}
	body := []byte(`<html><body><h1>CAPTCHA Required</h1></body></html>`)

	err := HandleHTMLResponse(resp, body)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	cliErr, ok := err.(*models.CLIError)
	if !ok {
		t.Fatalf("Expected CLIError, got %T", err)
	}

	if cliErr.Code != models.ErrCodeCaptchaRequired {
		t.Errorf("Expected error code %s, got %s", models.ErrCodeCaptchaRequired, cliErr.Code)
	}

	if cliErr.Details == nil {
		t.Error("Expected details to be set")
	}
}

func TestHandleHTMLResponse_Login(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": []string{"text/html"},
		},
		Request: &http.Request{
			URL: &url.URL{
				Scheme: "https",
				Host:   "www.amazon.com",
				Path:   "/signin",
			},
		},
	}
	body := []byte(`<html><body><form action="/ap_signin"><input name="ap_email" /></form></body></html>`)

	err := HandleHTMLResponse(resp, body)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	cliErr, ok := err.(*models.CLIError)
	if !ok {
		t.Fatalf("Expected CLIError, got %T", err)
	}

	if cliErr.Code != models.ErrCodeLoginRequired {
		t.Errorf("Expected error code %s, got %s", models.ErrCodeLoginRequired, cliErr.Code)
	}
}

func TestHandleHTMLResponse_UnknownHTML(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": []string{"text/html"},
		},
		Request: &http.Request{
			URL: &url.URL{
				Scheme: "https",
				Host:   "www.amazon.com",
				Path:   "/somepage",
			},
		},
	}
	body := []byte(`<html><body><h1>Unknown Page</h1></body></html>`)

	err := HandleHTMLResponse(resp, body)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	cliErr, ok := err.(*models.CLIError)
	if !ok {
		t.Fatalf("Expected CLIError, got %T", err)
	}

	if cliErr.Code != models.ErrCodeHTMLResponse {
		t.Errorf("Expected error code %s, got %s", models.ErrCodeHTMLResponse, cliErr.Code)
	}
}

func TestHandleHTMLResponse_ValidJSON(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Request: &http.Request{
			URL: &url.URL{
				Scheme: "https",
				Host:   "www.amazon.com",
				Path:   "/api/cart",
			},
		},
	}
	body := []byte(`{"status": "success"}`)

	err := HandleHTMLResponse(resp, body)

	if err != nil {
		t.Errorf("Expected no error for valid JSON, got: %v", err)
	}
}

func TestReadAndCheckResponse_ValidJSON(t *testing.T) {
	jsonBody := []byte(`{"status": "success"}`)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: io.NopCloser(bytes.NewReader(jsonBody)),
		Request: &http.Request{
			URL: &url.URL{
				Scheme: "https",
				Host:   "www.amazon.com",
				Path:   "/api/cart",
			},
		},
	}

	body, err := ReadAndCheckResponse(resp)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if !bytes.Equal(body, jsonBody) {
		t.Errorf("Expected body %s, got %s", jsonBody, body)
	}
}

func TestReadAndCheckResponse_CaptchaHTML(t *testing.T) {
	htmlBody := []byte(`<html><body><h1>CAPTCHA Required</h1></body></html>`)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": []string{"text/html"},
		},
		Body: io.NopCloser(bytes.NewReader(htmlBody)),
		Request: &http.Request{
			URL: &url.URL{
				Scheme: "https",
				Host:   "www.amazon.com",
				Path:   "/cart",
			},
		},
	}

	body, err := ReadAndCheckResponse(resp)

	if err == nil {
		t.Fatal("Expected error for CAPTCHA page, got nil")
	}

	cliErr, ok := err.(*models.CLIError)
	if !ok {
		t.Fatalf("Expected CLIError, got %T", err)
	}

	if cliErr.Code != models.ErrCodeCaptchaRequired {
		t.Errorf("Expected error code %s, got %s", models.ErrCodeCaptchaRequired, cliErr.Code)
	}

	// Body should still be returned even on error
	if !bytes.Equal(body, htmlBody) {
		t.Errorf("Expected body to be returned even on error")
	}
}

func TestIsHTMLContentType(t *testing.T) {
	testCases := []struct {
		contentType string
		expected    bool
	}{
		{"text/html", true},
		{"text/html; charset=UTF-8", true},
		{"TEXT/HTML", true},
		{"application/json", false},
		{"application/xml", false},
		{"text/plain", false},
		{"", false},
	}

	for _, tc := range testCases {
		t.Run(tc.contentType, func(t *testing.T) {
			result := IsHTMLContentType(tc.contentType)
			if result != tc.expected {
				t.Errorf("Expected %v for content type %s, got %v", tc.expected, tc.contentType, result)
			}
		})
	}
}

func TestIsJSONContentType(t *testing.T) {
	testCases := []struct {
		contentType string
		expected    bool
	}{
		{"application/json", true},
		{"application/json; charset=UTF-8", true},
		{"APPLICATION/JSON", true},
		{"text/json", true},
		{"text/html", false},
		{"application/xml", false},
		{"text/plain", false},
		{"", false},
	}

	for _, tc := range testCases {
		t.Run(tc.contentType, func(t *testing.T) {
			result := IsJSONContentType(tc.contentType)
			if result != tc.expected {
				t.Errorf("Expected %v for content type %s, got %v", tc.expected, tc.contentType, result)
			}
		})
	}
}

func TestDetectHTMLResponse_CaseInsensitive(t *testing.T) {
	// Test that detection is case-insensitive
	resp := &http.Response{
		Header: http.Header{
			"Content-Type": []string{"text/html"},
		},
	}
	body := []byte(`<html><body><h1>PLEASE SOLVE THIS CAPTCHA</h1></body></html>`)

	responseType, isHTML := DetectHTMLResponse(resp, body)

	if !isHTML {
		t.Error("Expected HTML response detection")
	}
	if responseType != HTMLResponseCaptcha {
		t.Errorf("Expected HTMLResponseCaptcha for uppercase CAPTCHA, got %v", responseType)
	}
}

func TestHandleHTMLResponse_SnippetTruncation(t *testing.T) {
	// Test that long responses are truncated in the error details
	longBody := bytes.Repeat([]byte("a"), 1000)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": []string{"text/html"},
		},
		Body: io.NopCloser(bytes.NewReader(longBody)),
		Request: &http.Request{
			URL: &url.URL{
				Scheme: "https",
				Host:   "www.amazon.com",
				Path:   "/test",
			},
		},
	}

	err := HandleHTMLResponse(resp, longBody)

	if err == nil {
		t.Fatal("Expected error for HTML response")
	}

	cliErr, ok := err.(*models.CLIError)
	if !ok {
		t.Fatalf("Expected CLIError, got %T", err)
	}

	if cliErr.Details == nil {
		t.Fatal("Expected details to be set")
	}

	snippet, ok := cliErr.Details["snippet"].(string)
	if !ok {
		t.Fatal("Expected snippet in details")
	}

	// Should be truncated to ~500 chars + "..."
	if len(snippet) > 510 {
		t.Errorf("Expected snippet to be truncated to ~500 chars, got %d", len(snippet))
	}
}
