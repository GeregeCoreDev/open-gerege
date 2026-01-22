package middleware_test

import (
	"io"
	"net/http/httptest"
	"testing"

	"templatev25/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultHSTSConfig(t *testing.T) {
	cfg := middleware.DefaultHSTSConfig()

	assert.True(t, cfg.Enabled)
	assert.Equal(t, 31536000, cfg.MaxAge)
	assert.True(t, cfg.IncludeSubDomains)
	assert.True(t, cfg.Preload)
}

func TestHSTS(t *testing.T) {
	tests := []struct {
		name           string
		config         middleware.HSTSConfig
		expectedHeader string
	}{
		{
			name:           "disabled",
			config:         middleware.HSTSConfig{Enabled: false},
			expectedHeader: "",
		},
		{
			name: "enabled with default config",
			config: middleware.HSTSConfig{
				Enabled:           true,
				MaxAge:            31536000,
				IncludeSubDomains: true,
				Preload:           true,
			},
			expectedHeader: "max-age=31536000; includeSubDomains; preload",
		},
		{
			name: "enabled without subdomains",
			config: middleware.HSTSConfig{
				Enabled:           true,
				MaxAge:            86400,
				IncludeSubDomains: false,
				Preload:           false,
			},
			expectedHeader: "max-age=86400",
		},
		{
			name: "enabled with subdomains only",
			config: middleware.HSTSConfig{
				Enabled:           true,
				MaxAge:            3600,
				IncludeSubDomains: true,
				Preload:           false,
			},
			expectedHeader: "max-age=3600; includeSubDomains",
		},
		{
			name: "enabled with preload only",
			config: middleware.HSTSConfig{
				Enabled:           true,
				MaxAge:            7200,
				IncludeSubDomains: false,
				Preload:           true,
			},
			expectedHeader: "max-age=7200; preload",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Use(middleware.HSTS(tt.config))
			app.Get("/", func(c *fiber.Ctx) error {
				return c.SendString("OK")
			})

			req := httptest.NewRequest("GET", "/", nil)
			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, 200, resp.StatusCode)

			header := resp.Header.Get("Strict-Transport-Security")
			assert.Equal(t, tt.expectedHeader, header)
		})
	}
}

func TestHSTS_PassesThrough(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.HSTS(middleware.DefaultHSTSConfig()))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "Hello World", string(body))
}
