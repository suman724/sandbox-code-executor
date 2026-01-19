package orchestration

import "time"

type RetryPolicy struct {
	MaxAttempts int
	Backoff     time.Duration
}

func DefaultRetryPolicy() RetryPolicy {
	return RetryPolicy{MaxAttempts: 3, Backoff: 500 * time.Millisecond}
}
