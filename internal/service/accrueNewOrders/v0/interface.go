package v0

import (
	"context"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/entity"
)

type OrderRepository interface {
	OrdersListAccrualsByOrderStatus(ctx context.Context, orderStatus string, limit int) (items []entity.AccrualItem, err error)
	OrdersUpdateAccruals(
		ctx context.Context, orders []entity.OrderUpdate, accruals []entity.AccrualUpdate, userBalances []entity.UserBalanceUpdate,
	) (err error)
}
