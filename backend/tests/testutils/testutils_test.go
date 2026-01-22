// Package testutils provides test utilities for integration tests
//
// File: testutils_test.go
// Description: Tests for test utilities
package testutils

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// ============================================================
// ASSERTIONS TESTS
// ============================================================

func TestAssertEqual(t *testing.T) {
	tests := []struct {
		name     string
		expected interface{}
		actual   interface{}
		wantFail bool
	}{
		{
			name:     "equal integers",
			expected: 1,
			actual:   1,
			wantFail: false,
		},
		{
			name:     "equal strings",
			expected: "hello",
			actual:   "hello",
			wantFail: false,
		},
		{
			name:     "equal slices",
			expected: []int{1, 2, 3},
			actual:   []int{1, 2, 3},
			wantFail: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Using a mock testing.T to avoid actual test failures
			mockT := &testing.T{}
			AssertEqual(mockT, tt.expected, tt.actual)
			if mockT.Failed() != tt.wantFail {
				t.Errorf("AssertEqual() failed = %v, wantFail %v", mockT.Failed(), tt.wantFail)
			}
		})
	}
}

func TestAssertNotEqual(t *testing.T) {
	tests := []struct {
		name     string
		expected interface{}
		actual   interface{}
		wantFail bool
	}{
		{
			name:     "different integers",
			expected: 1,
			actual:   2,
			wantFail: false,
		},
		{
			name:     "different strings",
			expected: "hello",
			actual:   "world",
			wantFail: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := &testing.T{}
			AssertNotEqual(mockT, tt.expected, tt.actual)
			if mockT.Failed() != tt.wantFail {
				t.Errorf("AssertNotEqual() failed = %v, wantFail %v", mockT.Failed(), tt.wantFail)
			}
		})
	}
}

func TestAssertNoError(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		mockT := &testing.T{}
		AssertNoError(mockT, nil)
		if mockT.Failed() {
			t.Error("AssertNoError() should not fail for nil error")
		}
	})
}

func TestAssertError(t *testing.T) {
	t.Run("non-nil error", func(t *testing.T) {
		mockT := &testing.T{}
		AssertError(mockT, errors.New("test error"))
		if mockT.Failed() {
			t.Error("AssertError() should not fail for non-nil error")
		}
	})
}

func TestAssertErrorContains(t *testing.T) {
	t.Run("error contains message", func(t *testing.T) {
		mockT := &testing.T{}
		AssertErrorContains(mockT, errors.New("connection refused"), "refused")
		if mockT.Failed() {
			t.Error("AssertErrorContains() should not fail when error contains message")
		}
	})
}

func TestAssertTrue(t *testing.T) {
	t.Run("true value", func(t *testing.T) {
		mockT := &testing.T{}
		AssertTrue(mockT, true)
		if mockT.Failed() {
			t.Error("AssertTrue() should not fail for true")
		}
	})
}

func TestAssertFalse(t *testing.T) {
	t.Run("false value", func(t *testing.T) {
		mockT := &testing.T{}
		AssertFalse(mockT, false)
		if mockT.Failed() {
			t.Error("AssertFalse() should not fail for false")
		}
	})
}

func TestAssertLen(t *testing.T) {
	tests := []struct {
		name     string
		obj      interface{}
		expected int
		wantFail bool
	}{
		{
			name:     "slice with correct length",
			obj:      []int{1, 2, 3},
			expected: 3,
			wantFail: false,
		},
		{
			name:     "string with correct length",
			obj:      "hello",
			expected: 5,
			wantFail: false,
		},
		{
			name:     "empty slice",
			obj:      []string{},
			expected: 0,
			wantFail: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := &testing.T{}
			AssertLen(mockT, tt.obj, tt.expected)
			if mockT.Failed() != tt.wantFail {
				t.Errorf("AssertLen() failed = %v, wantFail %v", mockT.Failed(), tt.wantFail)
			}
		})
	}
}

func TestAssertContains(t *testing.T) {
	t.Run("string contains substring", func(t *testing.T) {
		mockT := &testing.T{}
		AssertContains(mockT, "hello world", "world")
		if mockT.Failed() {
			t.Error("AssertContains() should not fail when string contains substring")
		}
	})
}

func TestAssertNotContains(t *testing.T) {
	t.Run("string does not contain substring", func(t *testing.T) {
		mockT := &testing.T{}
		AssertNotContains(mockT, "hello world", "foo")
		if mockT.Failed() {
			t.Error("AssertNotContains() should not fail when string does not contain substring")
		}
	})
}

func TestAssertStatusCode(t *testing.T) {
	t.Run("matching status codes", func(t *testing.T) {
		mockT := &testing.T{}
		AssertStatusCode(mockT, 200, 200)
		if mockT.Failed() {
			t.Error("AssertStatusCode() should not fail for matching codes")
		}
	})
}

func TestAssertJSONEqual(t *testing.T) {
	t.Run("equal JSON objects", func(t *testing.T) {
		mockT := &testing.T{}
		AssertJSONEqual(mockT, `{"a":1,"b":2}`, `{"b":2,"a":1}`)
		if mockT.Failed() {
			t.Error("AssertJSONEqual() should not fail for equal JSON")
		}
	})
}

func TestAssertGreater(t *testing.T) {
	t.Run("greater value", func(t *testing.T) {
		mockT := &testing.T{}
		AssertGreater(mockT, 10, 5)
		if mockT.Failed() {
			t.Error("AssertGreater() should not fail when value is greater")
		}
	})
}

func TestAssertLess(t *testing.T) {
	t.Run("less value", func(t *testing.T) {
		mockT := &testing.T{}
		AssertLess(mockT, 5, 10)
		if mockT.Failed() {
			t.Error("AssertLess() should not fail when value is less")
		}
	})
}

func TestAssertEmpty(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		mockT := &testing.T{}
		AssertEmpty(mockT, []int{})
		if mockT.Failed() {
			t.Error("AssertEmpty() should not fail for empty slice")
		}
	})

	t.Run("empty string", func(t *testing.T) {
		mockT := &testing.T{}
		AssertEmpty(mockT, "")
		if mockT.Failed() {
			t.Error("AssertEmpty() should not fail for empty string")
		}
	})
}

func TestAssertNotEmpty(t *testing.T) {
	t.Run("non-empty slice", func(t *testing.T) {
		mockT := &testing.T{}
		AssertNotEmpty(mockT, []int{1, 2, 3})
		if mockT.Failed() {
			t.Error("AssertNotEmpty() should not fail for non-empty slice")
		}
	})

	t.Run("non-empty string", func(t *testing.T) {
		mockT := &testing.T{}
		AssertNotEmpty(mockT, "hello")
		if mockT.Failed() {
			t.Error("AssertNotEmpty() should not fail for non-empty string")
		}
	})
}

// ============================================================
// HTTP TESTS
// ============================================================

func TestNewTestApp(t *testing.T) {
	app := NewTestApp()

	if app == nil {
		t.Fatal("NewTestApp() returned nil")
	}

	// Test that the app handles errors correctly
	app.Get("/error", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusBadRequest, "test error")
	})

	resp := MakeRequest(t, app, TestRequest{
		Method: "GET",
		Path:   "/error",
	})

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestMakeRequest(t *testing.T) {
	t.Run("GET request", func(t *testing.T) {
		app := NewTestApp()
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "success"})
		})

		resp := MakeRequest(t, app, TestRequest{
			Method: "GET",
			Path:   "/test",
		})

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST request with body", func(t *testing.T) {
		app := NewTestApp()
		app.Post("/data", func(c *fiber.Ctx) error {
			var body map[string]interface{}
			if err := c.BodyParser(&body); err != nil {
				return fiber.NewError(fiber.StatusBadRequest, err.Error())
			}
			return c.JSON(body)
		})

		resp := MakeRequest(t, app, TestRequest{
			Method: "POST",
			Path:   "/data",
			Body:   map[string]string{"key": "value"},
		})

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("request with headers", func(t *testing.T) {
		app := NewTestApp()
		app.Get("/header", func(c *fiber.Ctx) error {
			return c.SendString(c.Get("X-Custom-Header"))
		})

		resp := MakeRequest(t, app, TestRequest{
			Method: "GET",
			Path:   "/header",
			Headers: map[string]string{
				"X-Custom-Header": "test-value",
			},
		})

		if string(resp.Body) != "test-value" {
			t.Errorf("expected body 'test-value', got '%s'", string(resp.Body))
		}
	})

	t.Run("request with query params in path", func(t *testing.T) {
		app := NewTestApp()
		app.Get("/query", func(c *fiber.Ctx) error {
			return c.SendString(c.Query("param"))
		})

		// Query params can be added directly to the path
		resp := MakeRequest(t, app, TestRequest{
			Method: "GET",
			Path:   "/query?param=test-value",
		})

		if string(resp.Body) != "test-value" {
			t.Errorf("expected body 'test-value', got '%s'", string(resp.Body))
		}
	})
}

func TestTestResponse_ParseJSON(t *testing.T) {
	resp := TestResponse{
		StatusCode: 200,
		Body:       []byte(`{"message":"hello"}`),
	}

	var result map[string]string
	resp.ParseJSON(t, &result)

	if result["message"] != "hello" {
		t.Errorf("expected message 'hello', got '%s'", result["message"])
	}
}

func TestTestResponse_ParseAPIResponse(t *testing.T) {
	resp := TestResponse{
		StatusCode: 200,
		Body:       []byte(`{"code":"SUCCESS","message":"OK","data":{"id":1}}`),
	}

	apiResp := resp.ParseAPIResponse(t)

	if apiResp.Code != "SUCCESS" {
		t.Errorf("expected code 'SUCCESS', got '%s'", apiResp.Code)
	}
	if apiResp.Message != "OK" {
		t.Errorf("expected message 'OK', got '%s'", apiResp.Message)
	}
}

func TestAuthToken(t *testing.T) {
	token := AuthToken(1)
	if token == "" {
		t.Error("AuthToken() returned empty string")
	}
}

func TestAuthHeader(t *testing.T) {
	headers := AuthHeader("test-token")

	if headers["Authorization"] != "Bearer test-token" {
		t.Errorf("expected 'Bearer test-token', got '%s'", headers["Authorization"])
	}
}

func TestSessionCookie(t *testing.T) {
	headers := SessionCookie("session123")

	if headers["Cookie"] != "sid=session123" {
		t.Errorf("expected 'sid=session123', got '%s'", headers["Cookie"])
	}
}

// ============================================================
// DB TESTS (Unit tests only, no actual DB connection)
// ============================================================

func TestDefaultTestDSN(t *testing.T) {
	dsn := DefaultTestDSN()

	if dsn == "" {
		t.Error("DefaultTestDSN() returned empty string")
	}

	// Should contain postgres
	if !contains(dsn, "postgres") {
		t.Errorf("expected DSN to contain 'postgres', got '%s'", dsn)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
