// Package testutils provides test utilities for integration tests
//
// File: http.go
// Description: HTTP test helpers
package testutils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// TestRequest represents a test HTTP request
type TestRequest struct {
	Method      string
	Path        string
	Body        interface{}
	Headers     map[string]string
	QueryParams map[string]string
}

// TestResponse represents a test HTTP response
type TestResponse struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
}

// NewTestApp creates a new Fiber app for testing
func NewTestApp() *fiber.App {
	return fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"code":    "ERROR",
				"message": err.Error(),
			})
		},
	})
}

// MakeRequest sends a test request to a Fiber app
func MakeRequest(t *testing.T, app *fiber.App, req TestRequest) TestResponse {
	t.Helper()

	var bodyReader io.Reader
	if req.Body != nil {
		bodyBytes, err := json.Marshal(req.Body)
		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	httpReq := httptest.NewRequest(req.Method, req.Path, bodyReader)

	// Set default Content-Type for POST/PUT/PATCH
	if req.Body != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	// Set custom headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// Add query parameters
	if len(req.QueryParams) > 0 {
		q := httpReq.URL.Query()
		for key, value := range req.QueryParams {
			q.Add(key, value)
		}
		httpReq.URL.RawQuery = q.Encode()
	}

	resp, err := app.Test(httpReq, -1)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	return TestResponse{
		StatusCode: resp.StatusCode,
		Body:       body,
		Headers:    resp.Header,
	}
}

// ParseJSON parses the response body as JSON
func (r TestResponse) ParseJSON(t *testing.T, v interface{}) {
	t.Helper()

	if err := json.Unmarshal(r.Body, v); err != nil {
		t.Fatalf("failed to parse response JSON: %v\nBody: %s", err, string(r.Body))
	}
}

// APIResponse represents a standard API response
type APIResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ParseAPIResponse parses the response body as APIResponse
func (r TestResponse) ParseAPIResponse(t *testing.T) APIResponse {
	t.Helper()

	var resp APIResponse
	r.ParseJSON(t, &resp)
	return resp
}

// AuthToken generates a mock auth token for testing
func AuthToken(userID int) string {
	// In real tests, this would generate a proper JWT
	return "test-token"
}

// AuthHeader returns an Authorization header with the given token
func AuthHeader(token string) map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + token,
	}
}

// SessionCookie returns a session cookie header
func SessionCookie(sessionID string) map[string]string {
	return map[string]string{
		"Cookie": "sid=" + sessionID,
	}
}
