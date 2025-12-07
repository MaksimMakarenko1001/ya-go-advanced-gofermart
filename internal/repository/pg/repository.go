package pg

import (
	"context"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/db"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/entity"
)

type Repository struct {
	db *db.PGConnect
}

func New(db *db.PGConnect) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) OrdersListAccrualsByOrderStatus(ctx context.Context, orderStatus string, limit int) (items []entity.AccrualItem, err error) {
	return nil, nil
}

func (r *Repository) OrdersUpdateAccruals(
	ctx context.Context, orders []entity.OrderUpdate, accruals []entity.AccrualUpdate, userBalances []entity.UserBalanceUpdate,
) (err error) {
	return nil
}
