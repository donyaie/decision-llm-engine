package reliability

import (
	"context"
	"fmt"
	"time"
)

// DoWithResult retries a function with bounded linear backoff.
func DoWithResult[T any](ctx context.Context, attempts int, baseDelay time.Duration, fn func() (T, error)) (T, error) {
	var zero T
	if attempts < 1 {
		attempts = 1
	}

	var lastErr error
	for attempt := 1; attempt <= attempts; attempt++ {
		result, err := fn()
		if err == nil {
			return result, nil
		}

		lastErr = err
		if attempt == attempts {
			break
		}

		wait := time.Duration(attempt) * baseDelay
		select {
		case <-ctx.Done():
			return zero, fmt.Errorf("retry aborted: %w", ctx.Err())
		case <-time.After(wait):
		}
	}

	return zero, fmt.Errorf("operation failed after %d attempt(s): %w", attempts, lastErr)
}
