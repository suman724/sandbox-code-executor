package orchestration

import (
	"context"
	"errors"
	"time"
)

type RetryPolicy struct {
	MaxAttempts int
	Backoff     time.Duration
}

func DefaultRetryPolicy() RetryPolicy {
	return RetryPolicy{MaxAttempts: 3, Backoff: 500 * time.Millisecond}
}

func Retry(ctx context.Context, policy RetryPolicy, fn func(context.Context) error) error {
	if policy.MaxAttempts <= 0 {
		return errors.New("invalid retry policy")
	}
	var lastErr error
	for attempt := 1; attempt <= policy.MaxAttempts; attempt++ {
		if err := fn(ctx); err != nil {
			lastErr = err
			if attempt == policy.MaxAttempts {
				break
			}
			if err := sleep(ctx, policy.Backoff); err != nil {
				return err
			}
			continue
		}
		return nil
	}
	return lastErr
}

func sleep(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
