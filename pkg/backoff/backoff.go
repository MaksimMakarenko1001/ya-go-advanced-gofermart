package backoff

import (
	"context"
	"fmt"
	"time"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/logger"
)

type Backoff struct {
	logger          *logger.ZapLogger
	errClassifyFunc func(err error) ErrorClassification
	maxRetries      uint16
	t0              time.Duration
}

func NewBackoff(
	logger *logger.ZapLogger,
	errClassifyFunc func(err error) ErrorClassification,
	maxRetries uint16,
	t0 time.Duration,
) *Backoff {
	return &Backoff{
		logger:          logger,
		errClassifyFunc: errClassifyFunc,
		maxRetries:      maxRetries,
		t0:              t0,
	}
}

func (b *Backoff) WithRetry() func(retried) retried {
	return func(fn retried) retried {
		return func(ctx context.Context) error {
			for attempt := uint16(0); ; attempt++ {
				ok, err := b.doAttempt(ctx, fn, attempt, b.t0)
				if err != nil {
					return err
				}
				if ok {
					return nil
				}
			}
		}
	}
}

func (b *Backoff) doAttempt(ctx context.Context, fn retried, attempt uint16, delay time.Duration) (ok bool, err error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	err = fn(ctx)
	if err == nil {
		return true, nil
	}

	b.logger.Debugf("backoff", "attempt #%d failed: %v", attempt+1, err)
	if b.errClassifyFunc(err) == NonRetriable {
		return false, fmt.Errorf("non retriable, %w", err)
	}
	if b.maxRetries-attempt < 1 {
		return false, fmt.Errorf("max attempts reached, %w", err)
	}

	b.logger.Debugf("backoff", "retrying in %vs...", delay.Seconds())
	time.Sleep(delay)

	return false, nil
}
