package pg

import (
	"context"

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

func (r *Repository) OrdersGetNext(ctx context.Context, limit int) (orderNumbers []string, err error) {
	return nil, nil
}

func (r *Repository) OrdersComplete(ctx context.Context, orderNumbers []string) (err error) {
	return nil
}

func (r *Repository) OrdersFail(ctx context.Context, orderNumbers []string) (err error) {
	return nil
}
