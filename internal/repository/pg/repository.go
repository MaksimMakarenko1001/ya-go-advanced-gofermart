package pg

import (
	"context"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/db"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/entity"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/types/gofermart"
)

type Repository struct {
	db *db.PGConnect
}

func New(db *db.PGConnect) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) OrdersListAccrualsByOrderStatus(ctx context.Context, status gofermart.OrderStatusType, limit int) (items []entity.AccrualItem, err error) {
	err = r.db.QueryWithOneResultJSON(
		ctx,
		&items,
		"select orders.orders_list_accruals_by_order_status(_order_status=>$1, _limit=>$2);",
		status.String(), limit,
	)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *Repository) OrdersUpdateAccruals(
	ctx context.Context, orders []entity.OrderUpdate, accruals []entity.AccrualUpdate, userBalances []entity.UserBalanceUpdate,
) (orderUpdatedNumbers []string, err error) {
	err = r.db.QueryWithOneResultJSON(
		ctx,
		&orderUpdatedNumbers,
		"select orders.orders_update_accruals(_accruals=>$1, _orders=>$2, _user_balances=>$3);",
		accruals, orders, userBalances,
	)
	if err != nil {
		return nil, err
	}

	return orderUpdatedNumbers, nil
}
