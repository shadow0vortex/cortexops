package retry

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

// Policy defines the parameters for exponential backoff retries.
type Policy struct {
	MaxAttempts int
	MinDelay    time.Duration
	MaxDelay    time.Duration
	Factor      float64
}

// DefaultPolicy provides a sensible default retry policy.
var DefaultPolicy = Policy{
	MaxAttempts: 5,
	MinDelay:    100 * time.Millisecond,
	MaxDelay:    2 * time.Second,
	Factor:      2.0,
}

// Do executes a function and retries it according to the policy using exponential backoff with jitter.
func Do(ctx context.Context, policy Policy, fn func() error) error {
	var err error
	delay := policy.MinDelay

	for attempt := 1; attempt <= policy.MaxAttempts; attempt++ {
		err = fn()
		if err == nil {
			return nil
		}

		// Don't sleep if it's the last attempt
		if attempt == policy.MaxAttempts {
			break
		}

		// Calculate next delay with jitter
		jitter := time.Duration(rand.Float64() * float64(delay) * 0.1) // 10% jitter
		sleep := delay + jitter

		select {
		case <-time.After(sleep):
			// Update delay for next iteration
			delay = time.Duration(float64(delay) * policy.Factor)
			if delay > policy.MaxDelay {
				delay = policy.MaxDelay
			}
		case <-ctx.Done():
			return errors.Join(err, ctx.Err())
		}
	}

	return err
}
