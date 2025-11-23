package backoff

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Backoff struct {
	errClassifyFunc func(err error) ErrorClassification
	maxRetries      uint16
	t0              time.Duration
}

func NewBackoff(
	errClassifyFunc func(err error) ErrorClassification,
	maxRetries uint16,
	t0 time.Duration,

) *Backoff {
	return &Backoff{
		errClassifyFunc: errClassifyFunc,
		maxRetries:      maxRetries,
		t0:              t0,
	}
}

func (b *Backoff) WithRetry() func(retried) retried {
	return func(fn retried) retried {
		return func(ctx context.Context) error {
			for attempt := range b.maxRetries + 1 {
				if err := b.doAttempt(ctx, fn, attempt, b.t0); err != nil {
					return err
				}
			}
			return nil
		}
	}
}

func (b *Backoff) doAttempt(ctx context.Context, fn retried, attempt uint16, delay time.Duration) error {
	var err error

	select {
	case <-ctx.Done():
		err = ctx.Err()
	default:
		errAttempt := fn(ctx)
		if errAttempt == nil {
			break
		}

		if b.errClassifyFunc(errAttempt) == NonRetriable {
			err = fmt.Errorf("non retriable, %w", errAttempt)
		}

		log.Printf("attempt #%d failed: %v", attempt+1, errAttempt)
		if attempt < b.maxRetries {
			log.Printf("retrying in %vs...", delay.Seconds())
			time.Sleep(delay)
		} else {
			err = fmt.Errorf("max attempts reached, %w", err)
		}
	}

	return err
}
