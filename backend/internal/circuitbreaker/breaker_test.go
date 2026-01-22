// Package circuitbreaker provides circuit breaker pattern implementation
//
// File: breaker_test.go
// Description: Unit tests for circuit breaker
package circuitbreaker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestState_String(t *testing.T) {
	tests := []struct {
		state    State
		expected string
	}{
		{StateClosed, "closed"},
		{StateHalfOpen, "half-open"},
		{StateOpen, "open"},
		{State(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.state.String())
		})
	}
}

func TestDefaultSettings(t *testing.T) {
	settings := DefaultSettings("test-breaker")

	assert.Equal(t, "test-breaker", settings.Name)
	assert.Equal(t, uint32(1), settings.MaxRequests)
	assert.Equal(t, 60*time.Second, settings.Interval)
	assert.Equal(t, 30*time.Second, settings.Timeout)
	assert.Equal(t, uint32(5), settings.FailureThreshold)
	assert.Equal(t, uint32(2), settings.SuccessThreshold)
	assert.NotNil(t, settings.ReadyToTrip)

	// Test default ReadyToTrip function
	counts := Counts{ConsecutiveFailures: 5}
	assert.True(t, settings.ReadyToTrip(counts))
	counts.ConsecutiveFailures = 4
	assert.False(t, settings.ReadyToTrip(counts))
}

func TestNew(t *testing.T) {
	cb := New(DefaultSettings("test"))

	assert.NotNil(t, cb)
	assert.Equal(t, "test", cb.Name())
	assert.Equal(t, StateClosed, cb.State())
}

func TestNew_WithZeroValues(t *testing.T) {
	settings := Settings{
		Name:        "test",
		MaxRequests: 0, // Should default to 1
		Timeout:     0, // Should default to 30s
	}

	cb := New(settings)

	assert.Equal(t, uint32(1), cb.maxRequests)
	assert.Equal(t, 30*time.Second, cb.timeout)
}

func TestCircuitBreaker_Execute_Success(t *testing.T) {
	cb := New(DefaultSettings("test"))

	result, err := cb.Execute(func() (interface{}, error) {
		return "success", nil
	})

	require.NoError(t, err)
	assert.Equal(t, "success", result)
	assert.Equal(t, StateClosed, cb.State())

	counts := cb.Counts()
	assert.Equal(t, uint32(1), counts.TotalSuccesses)
	assert.Equal(t, uint32(0), counts.TotalFailures)
}

func TestCircuitBreaker_Execute_Failure(t *testing.T) {
	cb := New(DefaultSettings("test"))
	expectedErr := errors.New("test error")

	result, err := cb.Execute(func() (interface{}, error) {
		return nil, expectedErr
	})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)

	counts := cb.Counts()
	assert.Equal(t, uint32(0), counts.TotalSuccesses)
	assert.Equal(t, uint32(1), counts.TotalFailures)
}

func TestCircuitBreaker_OpensAfterFailures(t *testing.T) {
	settings := DefaultSettings("test")
	settings.FailureThreshold = 3
	settings.ReadyToTrip = func(counts Counts) bool {
		return counts.ConsecutiveFailures >= 3
	}
	cb := New(settings)

	// Fail 3 times
	for i := 0; i < 3; i++ {
		_, _ = cb.Execute(func() (interface{}, error) {
			return nil, errors.New("error")
		})
	}

	assert.Equal(t, StateOpen, cb.State())
}

func TestCircuitBreaker_RejectsWhenOpen(t *testing.T) {
	settings := DefaultSettings("test")
	settings.FailureThreshold = 1
	settings.ReadyToTrip = func(counts Counts) bool {
		return counts.ConsecutiveFailures >= 1
	}
	cb := New(settings)

	// Fail once to open the circuit
	_, _ = cb.Execute(func() (interface{}, error) {
		return nil, errors.New("error")
	})

	assert.Equal(t, StateOpen, cb.State())

	// Next call should be rejected
	_, err := cb.Execute(func() (interface{}, error) {
		return "should not execute", nil
	})

	assert.Equal(t, ErrCircuitOpen, err)
}

func TestCircuitBreaker_TransitionsToHalfOpen(t *testing.T) {
	settings := DefaultSettings("test")
	settings.FailureThreshold = 1
	settings.Timeout = 50 * time.Millisecond
	settings.ReadyToTrip = func(counts Counts) bool {
		return counts.ConsecutiveFailures >= 1
	}
	cb := New(settings)

	// Open the circuit
	_, _ = cb.Execute(func() (interface{}, error) {
		return nil, errors.New("error")
	})

	assert.Equal(t, StateOpen, cb.State())

	// Wait for timeout
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, StateHalfOpen, cb.State())
}

func TestCircuitBreaker_ClosesAfterSuccessInHalfOpen(t *testing.T) {
	settings := DefaultSettings("test")
	settings.FailureThreshold = 1
	settings.MaxRequests = 1
	settings.Timeout = 50 * time.Millisecond
	settings.ReadyToTrip = func(counts Counts) bool {
		return counts.ConsecutiveFailures >= 1
	}
	cb := New(settings)

	// Open the circuit
	_, _ = cb.Execute(func() (interface{}, error) {
		return nil, errors.New("error")
	})

	// Wait for half-open
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, StateHalfOpen, cb.State())

	// Success in half-open should close the circuit
	_, err := cb.Execute(func() (interface{}, error) {
		return "success", nil
	})

	require.NoError(t, err)
	assert.Equal(t, StateClosed, cb.State())
}

func TestCircuitBreaker_ExecuteWithContext(t *testing.T) {
	cb := New(DefaultSettings("test"))
	ctx := context.Background()

	result, err := cb.ExecuteWithContext(ctx, func(c context.Context) (interface{}, error) {
		return "success", nil
	})

	require.NoError(t, err)
	assert.Equal(t, "success", result)
}

func TestCircuitBreaker_OnStateChange(t *testing.T) {
	stateChanges := []struct {
		from State
		to   State
	}{}

	settings := DefaultSettings("test")
	settings.FailureThreshold = 1
	settings.ReadyToTrip = func(counts Counts) bool {
		return counts.ConsecutiveFailures >= 1
	}
	settings.OnStateChange = func(name string, from, to State) {
		stateChanges = append(stateChanges, struct {
			from State
			to   State
		}{from, to})
	}
	cb := New(settings)

	// Trigger state change to open
	_, _ = cb.Execute(func() (interface{}, error) {
		return nil, errors.New("error")
	})

	require.Len(t, stateChanges, 1)
	assert.Equal(t, StateClosed, stateChanges[0].from)
	assert.Equal(t, StateOpen, stateChanges[0].to)
}

func TestDefaultRetryConfig(t *testing.T) {
	cfg := DefaultRetryConfig()

	assert.Equal(t, 3, cfg.MaxRetries)
	assert.Equal(t, 100*time.Millisecond, cfg.InitialInterval)
	assert.Equal(t, 5*time.Second, cfg.MaxInterval)
	assert.Equal(t, 2.0, cfg.Multiplier)
}

func TestExecuteWithRetry_Success(t *testing.T) {
	ctx := context.Background()
	cfg := RetryConfig{
		MaxRetries:      3,
		InitialInterval: 10 * time.Millisecond,
		MaxInterval:     100 * time.Millisecond,
		Multiplier:      2.0,
	}

	callCount := 0
	err := ExecuteWithRetry(ctx, cfg, func(c context.Context) error {
		callCount++
		return nil
	})

	require.NoError(t, err)
	assert.Equal(t, 1, callCount) // Should succeed on first try
}

func TestExecuteWithRetry_EventualSuccess(t *testing.T) {
	ctx := context.Background()
	cfg := RetryConfig{
		MaxRetries:      5,
		InitialInterval: 10 * time.Millisecond,
		MaxInterval:     100 * time.Millisecond,
		Multiplier:      2.0,
	}

	callCount := 0
	err := ExecuteWithRetry(ctx, cfg, func(c context.Context) error {
		callCount++
		if callCount < 3 {
			return errors.New("not yet")
		}
		return nil
	})

	require.NoError(t, err)
	assert.Equal(t, 3, callCount)
}

func TestExecuteWithRetry_AllFailures(t *testing.T) {
	ctx := context.Background()
	cfg := RetryConfig{
		MaxRetries:      2,
		InitialInterval: 10 * time.Millisecond,
		MaxInterval:     100 * time.Millisecond,
		Multiplier:      2.0,
	}

	expectedErr := errors.New("always fails")
	callCount := 0
	err := ExecuteWithRetry(ctx, cfg, func(c context.Context) error {
		callCount++
		return expectedErr
	})

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, 3, callCount) // Initial + 2 retries
}

func TestExecuteWithRetry_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cfg := RetryConfig{
		MaxRetries:      10,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     1 * time.Second,
		Multiplier:      2.0,
	}

	callCount := 0
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	err := ExecuteWithRetry(ctx, cfg, func(c context.Context) error {
		callCount++
		return errors.New("fail")
	})

	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
}

func TestErrors(t *testing.T) {
	assert.Equal(t, "circuit breaker is open", ErrCircuitOpen.Error())
	assert.Equal(t, "too many retries", ErrTooManyRetries.Error())
	assert.Equal(t, "operation timed out", ErrTimeout.Error())
}
