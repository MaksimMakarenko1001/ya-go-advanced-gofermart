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
	response := struct {
		Ok bool `json:"ok"`
	}{}

	err = r.db.QueryWithOneResultJSON(
		ctx,
		&response,
		"select lock.locks_acquire(_key=>$1, _segment=>$2, _until=>$3, _pid=>$4);",
		key, segment, until, pid,
	)
	if err != nil {
		return false, err
	}

	return response.Ok, nil
}

func (r *Repository) LockRelease(ctx context.Context, key string, segment string, pid string) (ok bool, err error) {
	response := struct {
		Ok bool `json:"ok"`
	}{}

	err = r.db.QueryWithOneResultJSON(
		ctx,
		&response,
		"select lock.locks_release(_key=>$1, _segment=>$2, _pid=>$3);",
		key, segment, pid,
	)
	if err != nil {
		return false, err
	}

	return response.Ok, nil
}
