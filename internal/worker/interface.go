package worker

import (
	"context"
	"time"
)

type Locker interface {
	LockAcquire(ctx context.Context, key string, segment string, until time.Time, pid string) (ok bool, err error)
	LockRelease(ctx context.Context, key string, segment string, pid string) (ok bool, err error)
}
