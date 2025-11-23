package backoff

import (
	"context"
	"time"
)

type LinearBackoff struct {
	backoff *Backoff
	dt      time.Duration
}

func NewLinearBackoff(backoff *Backoff, dt time.Duration) *LinearBackoff {
	return &LinearBackoff{
		backoff: backoff,
		dt:      dt,
	}
}

func (lb *LinearBackoff) WithRetry() func(retried) retried {
	return func(fn retried) retried {
		return func(ctx context.Context) error {
			delay := lb.backoff.t0
			for attempt := range lb.backoff.maxRetries + 1 {
				if err := lb.backoff.doAttempt(ctx, fn, attempt, delay); err != nil {
					return err
				}
				delay += lb.dt
			}
			return nil
		}
	}
}
