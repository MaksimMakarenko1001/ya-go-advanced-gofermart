package v0

import (
	"context"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/entity"
)

type OrderRepository interface {
	OrdersListAccrualsByUserId(ctx context.Context, userId int64) (items []entity.AccrualItem, err error)
}
