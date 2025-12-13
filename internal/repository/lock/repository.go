package lock

import (
	"context"
	"time"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/db"
)

type Repository struct {
	db *db.PGConnect
}

func New(db *db.PGConnect) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) LockAcquire(ctx context.Context, key string, segment string, until time.Time, pid string) (ok bool, err error) {
	return false, nil
}

func (r *Repository) LockRelease(ctx context.Context, key string, segment string, pid string) (ok bool, err error) {
	return false, nil
}
