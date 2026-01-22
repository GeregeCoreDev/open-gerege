// Package circuitbreaker provides a circuit breaker pattern implementation
//
// File: breaker.go
// Description: Circuit breaker for external API calls with retry and fallback support
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package circuitbreaker

import (
	"context"
	"errors"
	"sync"
	"time"
)

// State represents the state of the circuit breaker
type State int

const (
	StateClosed State = iota
	StateHalfOpen
	StateOpen
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateHalfOpen:
		return "half-open"
	case StateOpen:
		return "open"
	default:
		return "unknown"
	}
}

// Common errors
var (
	ErrCircuitOpen    = errors.New("circuit breaker is open")
	ErrTooManyRetries = errors.New("too many retries")
	ErrTimeout        = errors.New("operation timed out")
)

// Settings holds the configuration for a circuit breaker
type Settings struct {
	// Name is the identifier for this circuit breaker
	Name string

	// MaxRequests is the maximum number of requests allowed in half-open state
	MaxRequests uint32

	// Interval is the cyclic period of the closed state for clearing internal counts
	// If 0, counts are never cleared in closed state
	Interval time.Duration

	// Timeout is the period of the open state, after which state becomes half-open
	Timeout time.Duration

	// FailureThreshold is the number of consecutive failures needed to open the circuit
	FailureThreshold uint32

	// SuccessThreshold is the number of consecutive successes needed to close the circuit
	SuccessThreshold uint32

	// ReadyToTrip is called with a copy of Counts whenever a request fails in closed state
	// If it returns true, the circuit breaker will be opened
	ReadyToTrip func(counts Counts) bool

	// OnStateChange is called whenever the state of the circuit breaker changes
	OnStateChange func(name string, from State, to State)
}

// Counts holds the request counts
type Counts struct {
	Requests             uint32
	TotalSuccesses       uint32
	TotalFailures        uint32
	ConsecutiveSuccesses uint32
	ConsecutiveFailures  uint32
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	name          string
	maxRequests   uint32
	interval      time.Duration
	timeout       time.Duration
	readyToTrip   func(counts Counts) bool
	onStateChange func(name string, from State, to State)

	mu          sync.Mutex
	state       State
	generation  uint64
	counts      Counts
	expiry      time.Time
	lastFailure time.Time
}

// DefaultSettings returns default circuit breaker settings
func DefaultSettings(name string) Settings {
	return Settings{
		Name:             name,
		MaxRequests:      1,
		Interval:         60 * time.Second,
		Timeout:          30 * time.Second,
		FailureThreshold: 5,
		SuccessThreshold: 2,
		ReadyToTrip: func(counts Counts) bool {
			return counts.ConsecutiveFailures >= 5
		},
	}
}

// New creates a new CircuitBreaker with the given settings
func New(settings Settings) *CircuitBreaker {
	cb := &CircuitBreaker{
		name:          settings.Name,
		maxRequests:   settings.MaxRequests,
		interval:      settings.Interval,
		timeout:       settings.Timeout,
		onStateChange: settings.OnStateChange,
	}

	if settings.ReadyToTrip != nil {
		cb.readyToTrip = settings.ReadyToTrip
	} else {
		cb.readyToTrip = func(counts Counts) bool {
			return counts.ConsecutiveFailures >= settings.FailureThreshold
		}
	}

	if cb.maxRequests == 0 {
		cb.maxRequests = 1
	}
	if cb.timeout == 0 {
		cb.timeout = 30 * time.Second
	}

	cb.toNewGeneration(time.Now())
	return cb
}

// Name returns the name of the circuit breaker
func (cb *CircuitBreaker) Name() string {
	return cb.name
}

// State returns the current state of the circuit breaker
func (cb *CircuitBreaker) State() State {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	now := time.Now()
	state, _ := cb.currentState(now)
	return state
}

// Counts returns the current counts of the circuit breaker
func (cb *CircuitBreaker) Counts() Counts {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.counts
}

// Execute runs the given request if the circuit breaker permits it
func (cb *CircuitBreaker) Execute(req func() (interface{}, error)) (interface{}, error) {
	generation, err := cb.beforeRequest()
	if err != nil {
		return nil, err
	}

	defer func() {
		e := recover()
		if e != nil {
			cb.afterRequest(generation, false)
			panic(e)
		}
	}()

	result, err := req()
	cb.afterRequest(generation, err == nil)
	return result, err
}

// ExecuteWithContext runs the given request with context support
func (cb *CircuitBreaker) ExecuteWithContext(ctx context.Context, req func(context.Context) (interface{}, error)) (interface{}, error) {
	generation, err := cb.beforeRequest()
	if err != nil {
		return nil, err
	}

	defer func() {
		e := recover()
		if e != nil {
			cb.afterRequest(generation, false)
			panic(e)
		}
	}()

	result, err := req(ctx)
	cb.afterRequest(generation, err == nil)
	return result, err
}

func (cb *CircuitBreaker) beforeRequest() (uint64, error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := time.Now()
	state, generation := cb.currentState(now)

	if state == StateOpen {
		return generation, ErrCircuitOpen
	}

	if state == StateHalfOpen && cb.counts.Requests >= cb.maxRequests {
		return generation, ErrCircuitOpen
	}

	cb.counts.Requests++
	return generation, nil
}

func (cb *CircuitBreaker) afterRequest(before uint64, success bool) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := time.Now()
	state, generation := cb.currentState(now)
	if generation != before {
		return
	}

	if success {
		cb.onSuccess(state, now)
	} else {
		cb.onFailure(state, now)
	}
}

func (cb *CircuitBreaker) onSuccess(state State, now time.Time) {
	cb.counts.TotalSuccesses++
	cb.counts.ConsecutiveSuccesses++
	cb.counts.ConsecutiveFailures = 0

	switch state {
	case StateClosed:
		// nothing to do
	case StateHalfOpen:
		if cb.counts.ConsecutiveSuccesses >= cb.maxRequests {
			cb.setState(StateClosed, now)
		}
	}
}

func (cb *CircuitBreaker) onFailure(state State, now time.Time) {
	cb.counts.TotalFailures++
	cb.counts.ConsecutiveFailures++
	cb.counts.ConsecutiveSuccesses = 0
	cb.lastFailure = now

	switch state {
	case StateClosed:
		if cb.readyToTrip(cb.counts) {
			cb.setState(StateOpen, now)
		}
	case StateHalfOpen:
		cb.setState(StateOpen, now)
	}
}

func (cb *CircuitBreaker) currentState(now time.Time) (State, uint64) {
	switch cb.state {
	case StateClosed:
		if !cb.expiry.IsZero() && cb.expiry.Before(now) {
			cb.toNewGeneration(now)
		}
	case StateOpen:
		if cb.expiry.Before(now) {
			cb.setState(StateHalfOpen, now)
		}
	}
	return cb.state, cb.generation
}

func (cb *CircuitBreaker) setState(state State, now time.Time) {
	if cb.state == state {
		return
	}

	prev := cb.state
	cb.state = state
	cb.toNewGeneration(now)

	if cb.onStateChange != nil {
		cb.onStateChange(cb.name, prev, state)
	}
}

func (cb *CircuitBreaker) toNewGeneration(now time.Time) {
	cb.generation++
	cb.counts = Counts{}

	var zero time.Time
	switch cb.state {
	case StateClosed:
		if cb.interval == 0 {
			cb.expiry = zero
		} else {
			cb.expiry = now.Add(cb.interval)
		}
	case StateOpen:
		cb.expiry = now.Add(cb.timeout)
	default:
		cb.expiry = zero
	}
}

// RetryConfig holds the configuration for retry logic
type RetryConfig struct {
	MaxRetries      int
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
}

// DefaultRetryConfig returns default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:      3,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     5 * time.Second,
		Multiplier:      2.0,
	}
}

// ExecuteWithRetry runs the given request with exponential backoff retry
func ExecuteWithRetry(ctx context.Context, config RetryConfig, req func(context.Context) error) error {
	var lastErr error
	interval := config.InitialInterval

	for i := 0; i <= config.MaxRetries; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := req(ctx); err == nil {
			return nil
		} else {
			lastErr = err
		}

		if i < config.MaxRetries {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(interval):
			}

			interval = time.Duration(float64(interval) * config.Multiplier)
			if interval > config.MaxInterval {
				interval = config.MaxInterval
			}
		}
	}

	return lastErr
}
