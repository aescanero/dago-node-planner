package llm

import (
	"context"
	"fmt"
	"time"

	"github.com/aescanero/dago-node-planner/internal/config"
)

// Retrier provides retry logic with exponential backoff.
type Retrier struct {
	maxAttempts  int
	initialDelay time.Duration
	maxDelay     time.Duration
	multiplier   float64
}

// NewRetrier creates a new retrier with the given configuration.
func NewRetrier(cfg config.RetryConfig) *Retrier {
	return &Retrier{
		maxAttempts:  cfg.MaxAttempts,
		initialDelay: cfg.InitialDelay,
		maxDelay:     cfg.MaxDelay,
		multiplier:   cfg.Multiplier,
	}
}

// Do executes a function with retry logic.
func (r *Retrier) Do(ctx context.Context, fn func() error) error {
	var lastErr error
	delay := r.initialDelay

	for attempt := 1; attempt <= r.maxAttempts; attempt++ {
		// Execute function
		err := fn()
		if err == nil {
			return nil // Success
		}

		lastErr = err

		// Check if we should retry
		if !r.shouldRetry(err) {
			return err
		}

		// Don't wait after the last attempt
		if attempt >= r.maxAttempts {
			break
		}

		// Wait before retry
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			// Calculate next delay with exponential backoff
			delay = time.Duration(float64(delay) * r.multiplier)
			if delay > r.maxDelay {
				delay = r.maxDelay
			}
		}
	}

	return fmt.Errorf("max retry attempts (%d) exceeded: %w", r.maxAttempts, lastErr)
}

// shouldRetry determines if an error is retryable.
func (r *Retrier) shouldRetry(err error) bool {
	// For now, retry all errors
	// In a production system, you'd check for specific error types
	// (e.g., rate limits, temporary network errors)
	return true
}
