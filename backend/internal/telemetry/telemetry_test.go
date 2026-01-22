// Package telemetry provides OpenTelemetry observability setup
//
// File: telemetry_test.go
// Description: Unit tests for telemetry package
package telemetry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultTracerConfig(t *testing.T) {
	cfg := DefaultTracerConfig()

	assert.True(t, cfg.Enabled)
	assert.Equal(t, "localhost:4317", cfg.Endpoint)
	assert.True(t, cfg.Insecure)
	assert.Equal(t, 1.0, cfg.SampleRate)
	assert.False(t, cfg.UseStdout)
}

func TestInitTracer_Disabled(t *testing.T) {
	ctx := context.Background()
	cfg := TracerConfig{
		Enabled: false,
	}

	shutdown, err := InitTracer(ctx, cfg, "test-service", "1.0.0")

	require.NoError(t, err)
	require.NotNil(t, shutdown)

	// Shutdown should be a no-op
	err = shutdown(ctx)
	assert.NoError(t, err)
}

func TestInitTracer_StdoutExporter(t *testing.T) {
	ctx := context.Background()
	cfg := TracerConfig{
		Enabled:    true,
		UseStdout:  true,
		SampleRate: 1.0,
	}

	shutdown, err := InitTracer(ctx, cfg, "test-service", "1.0.0")

	require.NoError(t, err)
	require.NotNil(t, shutdown)

	// Cleanup
	err = shutdown(ctx)
	assert.NoError(t, err)
}

func TestInitTracer_SamplerConfigurations(t *testing.T) {
	tests := []struct {
		name       string
		sampleRate float64
	}{
		{"always sample", 1.0},
		{"never sample", 0.0},
		{"ratio sample 50%", 0.5},
		{"ratio sample 10%", 0.1},
		{"above 1.0 should always sample", 2.0},
		{"below 0.0 should never sample", -1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			cfg := TracerConfig{
				Enabled:    true,
				UseStdout:  true,
				SampleRate: tt.sampleRate,
			}

			shutdown, err := InitTracer(ctx, cfg, "test-service", "1.0.0")
			require.NoError(t, err)
			defer shutdown(ctx)
		})
	}
}

func TestTracer(t *testing.T) {
	// Tracer should return a tracer from the global provider
	tracer := Tracer("test-tracer")
	assert.NotNil(t, tracer)
}

func TestTracerConfig_Fields(t *testing.T) {
	cfg := TracerConfig{
		Enabled:    true,
		Endpoint:   "custom:4317",
		Insecure:   false,
		SampleRate: 0.5,
		UseStdout:  true,
	}

	assert.True(t, cfg.Enabled)
	assert.Equal(t, "custom:4317", cfg.Endpoint)
	assert.False(t, cfg.Insecure)
	assert.Equal(t, 0.5, cfg.SampleRate)
	assert.True(t, cfg.UseStdout)
}
