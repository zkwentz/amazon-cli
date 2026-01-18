package amazon

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/michaelshimeles/amazon-cli/pkg/models"
)

// HTMLResponseType represents the type of HTML response detected
type HTMLResponseType int

const (
	// HTMLResponseNone indicates no HTML response (valid JSON/expected response)
	HTMLResponseNone HTMLResponseType = iota
	// HTMLResponseCaptcha indicates a CAPTCHA page
	HTMLResponseCaptcha
	// HTMLResponseLogin indicates a login/authentication page
	HTMLResponseLogin
	// HTMLResponseUnknown indicates an unexpected HTML page
	HTMLResponseUnknown
)

// DetectHTMLResponse checks if the response is an unexpected HTML page
// It returns the type of HTML response and whether it was detected
func DetectHTMLResponse(resp *http.Response, body []byte) (HTMLResponseType, bool) {
	// Check Content-Type header
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		// Not HTML, this is expected (JSON, etc.)
		return HTMLResponseNone, false
	}

	// Convert to lowercase for case-insensitive matching
	bodyLower := strings.ToLower(string(body))

	// Check for CAPTCHA indicators
	captchaIndicators := []string{
		"captcha",
		"robot check",
		"automated access",
		"solve this puzzle",
		"verify you're not a robot",
		"security check",
		"sorry, we just need to make sure you're not a robot",
	}

	for _, indicator := range captchaIndicators {
		if strings.Contains(bodyLower, indicator) {
			return HTMLResponseCaptcha, true
		}
	}

	// Check for login/authentication indicators
	loginIndicators := []string{
		"ap_signin",
		"sign in",
		"sign-in",
		"authentication required",
		"your account",
		"ap_email",
		"ap_password",
		"signin",
	}

	for _, indicator := range loginIndicators {
		if strings.Contains(bodyLower, indicator) {
			return HTMLResponseLogin, true
		}
	}

	// If we got HTML but couldn't identify the specific type, it's unknown
	return HTMLResponseUnknown, true
}

// HandleHTMLResponse processes an unexpected HTML response and returns an appropriate error
func HandleHTMLResponse(resp *http.Response, body []byte) error {
	responseType, isHTML := DetectHTMLResponse(resp, body)

	if !isHTML {
		// Not an HTML response, this is fine
		return nil
	}

	// Extract a snippet of the response for debugging
	snippet := string(body)
	if len(snippet) > 500 {
		snippet = snippet[:500] + "..."
	}

	details := map[string]interface{}{
		"url":          resp.Request.URL.String(),
		"status_code":  resp.StatusCode,
		"content_type": resp.Header.Get("Content-Type"),
		"snippet":      snippet,
	}

	// Return appropriate error based on the HTML type
	switch responseType {
	case HTMLResponseCaptcha:
		return models.NewCaptchaRequiredError(details)
	case HTMLResponseLogin:
		return models.NewLoginRequiredError(details)
	default:
		return models.NewHTMLResponseError(details)
	}
}

// ReadAndCheckResponse reads the response body and checks for unexpected HTML responses
// It returns the body bytes and any error encountered
func ReadAndCheckResponse(resp *http.Response) ([]byte, error) {
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, models.NewNetworkError(err)
	}

	// Close the original body
	resp.Body.Close()

	// Replace the body with a new reader so it can be read again if needed
	resp.Body = io.NopCloser(bytes.NewReader(body))

	// Check for unexpected HTML responses
	if err := HandleHTMLResponse(resp, body); err != nil {
		return body, err
	}

	return body, nil
}

// IsHTMLContentType checks if the Content-Type header indicates HTML
func IsHTMLContentType(contentType string) bool {
	return strings.Contains(strings.ToLower(contentType), "text/html")
}

// IsJSONContentType checks if the Content-Type header indicates JSON
func IsJSONContentType(contentType string) bool {
	ct := strings.ToLower(contentType)
	return strings.Contains(ct, "application/json") || strings.Contains(ct, "text/json")
}
