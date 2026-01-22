// Package circuitbreaker provides a circuit breaker pattern implementation
//
// File: breaker_test.go
// Description: Unit tests for circuit breaker
package circuitbreaker_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"templatev25/internal/circuitbreaker"

	"github.com/stretchr/testify/assert"
)

func TestState_String(t *testing.T) {
	tests := []struct {
		state circuitbreaker.State
		want  string
	}{
		{circuitbreaker.StateClosed, "closed"},
		{circuitbreaker.StateHalfOpen, "half-open"},
		{circuitbreaker.StateOpen, "open"},
		{circuitbreaker.State(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.state.String())
		})
	}
}

func TestDefaultSettings(t *testing.T) {
	settings := circuitbreaker.DefaultSettings("test")

	assert.Equal(t, "test", settings.Name)
	assert.Equal(t, uint32(1), settings.MaxRequests)
	assert.Equal(t, 60*time.Second, settings.Interval)
	assert.Equal(t, 30*time.Second, settings.Timeout)
	assert.Equal(t, uint32(5), settings.FailureThreshold)
	assert.Equal(t, uint32(2), settings.SuccessThreshold)
	assert.NotNil(t, settings.ReadyToTrip)
}

func TestCircuitBreaker_New(t *testing.T) {
	cb := circuitbreaker.New(circuitbreaker.DefaultSettings("test"))

	assert.Equal(t, "test", cb.Name())
	assert.Equal(t, circuitbreaker.StateClosed, cb.State())
	assert.Equal(t, uint32(0), cb.Counts().Requests)
}

func TestCircuitBreaker_Execute_Success(t *testing.T) {
	cb := circuitbreaker.New(circuitbreaker.DefaultSettings("test"))

	result, err := cb.Execute(func() (interface{}, error) {
		return "success", nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "success", result)
	assert.Equal(t, uint32(1), cb.Counts().TotalSuccesses)
	assert.Equal(t, uint32(0), cb.Counts().TotalFailures)
}

func TestCircuitBreaker_Execute_Failure(t *testing.T) {
	cb := circuitbreaker.New(circuitbreaker.DefaultSettings("test"))

	_, err := cb.Execute(func() (interface{}, error) {
		return nil, errors.New("failed")
	})

	assert.Error(t, err)
	assert.Equal(t, uint32(0), cb.Counts().TotalSuccesses)
	assert.Equal(t, uint32(1), cb.Counts().TotalFailures)
}

func TestCircuitBreaker_OpenAfterFailures(t *testing.T) {
	settings := circuitbreaker.DefaultSettings("test")
	settings.FailureThreshold = 3
	settings.ReadyToTrip = func(counts circuitbreaker.Counts) bool {
		return counts.ConsecutiveFailures >= 3
	}

	cb := circuitbreaker.New(settings)

	// Fail 3 times
	for i := 0; i < 3; i++ {
		_, _ = cb.Execute(func() (interface{}, error) {
			return nil, errors.New("failed")
		})
	}

	assert.Equal(t, circuitbreaker.StateOpen, cb.State())
}

func TestCircuitBreaker_RejectsRequestsWhenOpen(t *testing.T) {
	settings := circuitbreaker.DefaultSettings("test")
	settings.FailureThreshold = 1
	settings.ReadyToTrip = func(counts circuitbreaker.Counts) bool {
		return counts.ConsecutiveFailures >= 1
	}

	cb := circuitbreaker.New(settings)

	// Fail once to open the circuit
	_, _ = cb.Execute(func() (interface{}, error) {
		return nil, errors.New("failed")
	})

	assert.Equal(t, circuitbreaker.StateOpen, cb.State())

	// Next request should be rejected
	_, err := cb.Execute(func() (interface{}, error) {
		return "should not run", nil
	})

	assert.Equal(t, circuitbreaker.ErrCircuitOpen, err)
}

func TestCircuitBreaker_HalfOpenAfterTimeout(t *testing.T) {
	settings := circuitbreaker.DefaultSettings("test")
	settings.FailureThreshold = 1
	settings.Timeout = 50 * time.Millisecond
	settings.ReadyToTrip = func(counts circuitbreaker.Counts) bool {
		return counts.ConsecutiveFailures >= 1
	}

	cb := circuitbreaker.New(settings)

	// Fail once to open the circuit
	_, _ = cb.Execute(func() (interface{}, error) {
		return nil, errors.New("failed")
	})

	assert.Equal(t, circuitbreaker.StateOpen, cb.State())

	// Wait for timeout
	time.Sleep(60 * time.Millisecond)

	// State should transition to half-open
	assert.Equal(t, circuitbreaker.StateHalfOpen, cb.State())
}

func TestCircuitBreaker_ClosesAfterSuccessInHalfOpen(t *testing.T) {
	settings := circuitbreaker.DefaultSettings("test")
	settings.FailureThreshold = 1
	settings.Timeout = 50 * time.Millisecond
	settings.MaxRequests = 1
	settings.ReadyToTrip = func(counts circuitbreaker.Counts) bool {
		return counts.ConsecutiveFailures >= 1
	}

	cb := circuitbreaker.New(settings)

	// Fail once to open the circuit
	_, _ = cb.Execute(func() (interface{}, error) {
		return nil, errors.New("failed")
	})

	// Wait for timeout
	time.Sleep(60 * time.Millisecond)

	// Execute successful request in half-open state
	_, err := cb.Execute(func() (interface{}, error) {
		return "success", nil
	})

	assert.NoError(t, err)
	assert.Equal(t, circuitbreaker.StateClosed, cb.State())
}

func TestCircuitBreaker_ReopensAfterFailureInHalfOpen(t *testing.T) {
	settings := circuitbreaker.DefaultSettings("test")
	settings.FailureThreshold = 1
	settings.Timeout = 50 * time.Millisecond
	settings.MaxRequests = 1
	settings.ReadyToTrip = func(counts circuitbreaker.Counts) bool {
		return counts.ConsecutiveFailures >= 1
	}

	cb := circuitbreaker.New(settings)

	// Fail once to open the circuit
	_, _ = cb.Execute(func() (interface{}, error) {
		return nil, errors.New("failed")
	})

	// Wait for timeout
	time.Sleep(60 * time.Millisecond)

	// Execute failing request in half-open state
	_, _ = cb.Execute(func() (interface{}, error) {
		return nil, errors.New("failed again")
	})

	assert.Equal(t, circuitbreaker.StateOpen, cb.State())
}

func TestCircuitBreaker_ExecuteWithContext(t *testing.T) {
	cb := circuitbreaker.New(circuitbreaker.DefaultSettings("test"))

	ctx := context.Background()
	result, err := cb.ExecuteWithContext(ctx, func(ctx context.Context) (interface{}, error) {
		return "success", nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "success", result)
}

func TestCircuitBreaker_OnStateChange(t *testing.T) {
	var stateChanges []circuitbreaker.State
	settings := circuitbreaker.DefaultSettings("test")
	settings.FailureThreshold = 1
	settings.ReadyToTrip = func(counts circuitbreaker.Counts) bool {
		return counts.ConsecutiveFailures >= 1
	}
	settings.OnStateChange = func(name string, from circuitbreaker.State, to circuitbreaker.State) {
		stateChanges = append(stateChanges, to)
	}

	cb := circuitbreaker.New(settings)

	// Fail once to open the circuit
	_, _ = cb.Execute(func() (interface{}, error) {
		return nil, errors.New("failed")
	})

	assert.Contains(t, stateChanges, circuitbreaker.StateOpen)
}

func TestDefaultRetryConfig(t *testing.T) {
	config := circuitbreaker.DefaultRetryConfig()

	assert.Equal(t, 3, config.MaxRetries)
	assert.Equal(t, 100*time.Millisecond, config.InitialInterval)
	assert.Equal(t, 5*time.Second, config.MaxInterval)
	assert.Equal(t, 2.0, config.Multiplier)
}

func TestExecuteWithRetry_Success(t *testing.T) {
	config := circuitbreaker.RetryConfig{
		MaxRetries:      3,
		InitialInterval: 10 * time.Millisecond,
		MaxInterval:     100 * time.Millisecond,
		Multiplier:      2.0,
	}

	var callCount int32
	err := circuitbreaker.ExecuteWithRetry(context.Background(), config, func(ctx context.Context) error {
		atomic.AddInt32(&callCount, 1)
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, int32(1), callCount)
}

func TestExecuteWithRetry_SuccessAfterRetries(t *testing.T) {
	config := circuitbreaker.RetryConfig{
		MaxRetries:      3,
		InitialInterval: 10 * time.Millisecond,
		MaxInterval:     100 * time.Millisecond,
		Multiplier:      2.0,
	}

	var callCount int32
	err := circuitbreaker.ExecuteWithRetry(context.Background(), config, func(ctx context.Context) error {
		count := atomic.AddInt32(&callCount, 1)
		if count < 3 {
			return errors.New("temporary error")
		}
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, int32(3), callCount)
}

func TestExecuteWithRetry_AllRetriesFail(t *testing.T) {
	config := circuitbreaker.RetryConfig{
		MaxRetries:      2,
		InitialInterval: 10 * time.Millisecond,
		MaxInterval:     100 * time.Millisecond,
		Multiplier:      2.0,
	}

	var callCount int32
	err := circuitbreaker.ExecuteWithRetry(context.Background(), config, func(ctx context.Context) error {
		atomic.AddInt32(&callCount, 1)
		return errors.New("persistent error")
	})

	assert.Error(t, err)
	assert.Equal(t, "persistent error", err.Error())
	assert.Equal(t, int32(3), callCount) // 1 initial + 2 retries
}

func TestExecuteWithRetry_ContextCanceled(t *testing.T) {
	config := circuitbreaker.RetryConfig{
		MaxRetries:      10,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     1 * time.Second,
		Multiplier:      2.0,
	}

	ctx, cancel := context.WithCancel(context.Background())

	var callCount int32
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	err := circuitbreaker.ExecuteWithRetry(ctx, config, func(ctx context.Context) error {
		atomic.AddInt32(&callCount, 1)
		return errors.New("error")
	})

	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
}

func TestExecuteWithRetry_MaxIntervalRespected(t *testing.T) {
	config := circuitbreaker.RetryConfig{
		MaxRetries:      5,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     150 * time.Millisecond,
		Multiplier:      10.0, // Should quickly exceed max
	}

	start := time.Now()
	var callCount int32
	_ = circuitbreaker.ExecuteWithRetry(context.Background(), config, func(ctx context.Context) error {
		count := atomic.AddInt32(&callCount, 1)
		if count <= 3 {
			return errors.New("error")
		}
		return nil
	})

	elapsed := time.Since(start)

	// With maxInterval of 150ms and 3 retries, max time should be around 450ms + overhead
	// Without maxInterval cap it would be 100 + 1000 + 10000 = 11100ms
	assert.Less(t, elapsed, 1*time.Second)
}
