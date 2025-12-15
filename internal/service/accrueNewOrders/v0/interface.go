package v0

import (
	"context"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/entity"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/types/gofermart"
)

type OrderRepository interface {
	OrdersListAccrualsByOrderStatus(ctx context.Context, status gofermart.OrderStatusType, limit int) (items []entity.AccrualItem, err error)
	OrdersUpdateAccruals(
		ctx context.Context, orders []entity.Order, accruals []entity.Accrual, userBalances []entity.UserBalance,
	) (orderUpdatedNumbers []string, err error)
}
