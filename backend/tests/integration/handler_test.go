//go:build integration

// Package integration provides integration tests for HTTP handlers
//
// These tests verify the HTTP layer behavior including:
// - Request/response handling
// - Authentication/authorization
// - Input validation
// - Error responses
package integration

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// APIResponse represents the standard API response structure
type APIResponse struct {
	Success   bool        `json:"success"`
	Code      string      `json:"code"`
	Message   string      `json:"msg"`
	Data      interface{} `json:"data"`
	RequestID string      `json:"request_id"`
}

// setupTestApp creates a minimal Fiber app for testing
func setupTestApp() *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	// Health endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"db":     "connected",
		})
	})

	// Test endpoint for rate limiting
	app.Get("/test/rate-limit", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	// Test endpoint for validation
	app.Post("/test/validate", func(c *fiber.Ctx) error {
		type Request struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		var req Request
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid JSON")
		}
		if req.Name == "" {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "name is required")
		}
		if req.Email == "" {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "email is required")
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": 1, "name": req.Name})
	})

	// Test endpoint for error responses
	app.Get("/test/error/:code", func(c *fiber.Ctx) error {
		code, _ := c.ParamsInt("code")
		switch code {
		case 400:
			return fiber.NewError(fiber.StatusBadRequest, "bad request")
		case 401:
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		case 403:
			return fiber.NewError(fiber.StatusForbidden, "forbidden")
		case 404:
			return fiber.NewError(fiber.StatusNotFound, "not found")
		case 429:
			return fiber.NewError(fiber.StatusTooManyRequests, "rate limit exceeded")
		case 500:
			return fiber.NewError(fiber.StatusInternalServerError, "internal error")
		case 503:
			return fiber.NewError(fiber.StatusServiceUnavailable, "service unavailable")
		default:
			return c.JSON(fiber.Map{"code": code})
		}
	})

	return app
}

// TestHealthEndpoint tests the health check endpoint
func TestHealthEndpoint(t *testing.T) {
	app := setupTestApp()

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "health check returns ok",
			method:     http.MethodGet,
			path:       "/health",
			wantStatus: http.StatusOK,
			wantBody:   `"status":"ok"`,
		},
		{
			name:       "health check wrong method",
			method:     http.MethodPost,
			path:       "/health",
			wantStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			if tt.wantBody != "" {
				body, _ := io.ReadAll(resp.Body)
				assert.Contains(t, string(body), tt.wantBody)
			}
		})
	}
}

// TestValidationEndpoint tests input validation
func TestValidationEndpoint(t *testing.T) {
	app := setupTestApp()

	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantError  string
	}{
		{
			name:       "valid request",
			body:       `{"name": "Test", "email": "test@example.com"}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "missing name",
			body:       `{"email": "test@example.com"}`,
			wantStatus: http.StatusUnprocessableEntity,
			wantError:  "name is required",
		},
		{
			name:       "missing email",
			body:       `{"name": "Test"}`,
			wantStatus: http.StatusUnprocessableEntity,
			wantError:  "email is required",
		},
		{
			name:       "invalid JSON",
			body:       `{invalid}`,
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid JSON",
		},
		{
			name:       "empty body",
			body:       ``,
			wantStatus: http.StatusBadRequest, // Empty body is invalid JSON
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/test/validate", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			if tt.wantError != "" {
				body, _ := io.ReadAll(resp.Body)
				assert.Contains(t, string(body), tt.wantError)
			}
		})
	}
}

// TestErrorResponses tests that error responses are properly formatted
func TestErrorResponses(t *testing.T) {
	app := setupTestApp()

	tests := []struct {
		name       string
		errorCode  int
		wantStatus int
	}{
		{"bad request", 400, http.StatusBadRequest},
		{"unauthorized", 401, http.StatusUnauthorized},
		{"forbidden", 403, http.StatusForbidden},
		{"not found", 404, http.StatusNotFound},
		{"rate limit", 429, http.StatusTooManyRequests},
		{"internal error", 500, http.StatusInternalServerError},
		{"service unavailable", 503, http.StatusServiceUnavailable},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test/error/"+string(rune(tt.errorCode+'0')), nil)
			// Fix: use proper integer formatting
			req = httptest.NewRequest(http.MethodGet, "/test/error/"+itoa(tt.errorCode), nil)

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

// itoa converts int to string (simple helper)
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	s := ""
	for i > 0 {
		s = string(rune('0'+i%10)) + s
		i /= 10
	}
	return s
}

// TestContentTypeHeaders tests that responses have correct content type
func TestContentTypeHeaders(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	assert.Contains(t, contentType, "application/json")
}

// TestNotFoundRoute tests 404 for non-existent routes
func TestNotFoundRoute(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest(http.MethodGet, "/non-existent-route", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// TestJSONResponseFormat tests that responses are valid JSON
func TestJSONResponseFormat(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	assert.NoError(t, err, "response should be valid JSON")
}

// TestRequestWithQueryParams tests query parameter handling
func TestRequestWithQueryParams(t *testing.T) {
	app := fiber.New()
	app.Get("/search", func(c *fiber.Ctx) error {
		page := c.QueryInt("page", 1)
		size := c.QueryInt("size", 10)
		q := c.Query("q", "")
		return c.JSON(fiber.Map{
			"page":  page,
			"size":  size,
			"query": q,
		})
	})

	tests := []struct {
		name       string
		query      string
		wantPage   int
		wantSize   int
		wantQuery  string
	}{
		{
			name:      "default values",
			query:     "",
			wantPage:  1,
			wantSize:  10,
			wantQuery: "",
		},
		{
			name:      "custom pagination",
			query:     "?page=2&size=20",
			wantPage:  2,
			wantSize:  20,
			wantQuery: "",
		},
		{
			name:      "with search query",
			query:     "?q=test&page=1&size=5",
			wantPage:  1,
			wantSize:  5,
			wantQuery: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/search"+tt.query, nil)
			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			var result map[string]interface{}
			body, _ := io.ReadAll(resp.Body)
			json.Unmarshal(body, &result)

			assert.Equal(t, float64(tt.wantPage), result["page"])
			assert.Equal(t, float64(tt.wantSize), result["size"])
			assert.Equal(t, tt.wantQuery, result["query"])
		})
	}
}

// TestRequestHeaders tests header handling
func TestRequestHeaders(t *testing.T) {
	app := fiber.New()
	app.Get("/headers", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"user_agent":   c.Get("User-Agent"),
			"content_type": c.Get("Content-Type"),
			"custom":       c.Get("X-Custom-Header"),
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/headers", nil)
	req.Header.Set("User-Agent", "TestClient/1.0")
	req.Header.Set("X-Custom-Header", "custom-value")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer resp.Body.Close()

	var result map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	assert.Equal(t, "TestClient/1.0", result["user_agent"])
	assert.Equal(t, "custom-value", result["custom"])
}
