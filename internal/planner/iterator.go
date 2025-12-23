package planner

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// Iterator manages iterative refinement of graph generation.
type Iterator struct {
	maxIterations int
	logger        *zap.Logger
}

// NewIterator creates a new iterator.
func NewIterator(maxIterations int, logger *zap.Logger) *Iterator {
	return &Iterator{
		maxIterations: maxIterations,
		logger:        logger,
	}
}

// IterateFunc is a function that performs a single iteration.
// It should return nil on success, or an error to trigger another iteration.
type IterateFunc func(ctx context.Context, attempt int) error

// Iterate performs iterative refinement up to maxIterations.
func (it *Iterator) Iterate(ctx context.Context, fn IterateFunc) error {
	var lastErr error

	for attempt := 1; attempt <= it.maxIterations; attempt++ {
		it.logger.Debug("starting iteration",
			zap.Int("attempt", attempt),
			zap.Int("max_iterations", it.maxIterations),
		)

		// Check context
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Execute iteration
		err := fn(ctx, attempt)
		if err == nil {
			it.logger.Debug("iteration succeeded",
				zap.Int("attempt", attempt),
			)
			return nil // Success!
		}

		lastErr = err
		it.logger.Debug("iteration failed, will retry",
			zap.Int("attempt", attempt),
			zap.Error(err),
		)

		// Check if we've exhausted attempts
		if attempt >= it.maxIterations {
			break
		}
	}

	return fmt.Errorf("max iterations (%d) exceeded: %w", it.maxIterations, lastErr)
}
