package v0

import (
	"context"
	"time"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/entity"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/types/moneys"
)

type OrderRepository interface {
	OrdersGetNext(ctx context.Context, orderStatus string, limit int) (items []entity.AccrualItem, err error)
	OrdersUpdateAccruals(ctx context.Context, orders []OrderUpdate, accruals []AccrualUpdate) (err error)
}

type OrderUpdate struct {
	OrderNumber string    `json:"order_number"`
	OrderStatus string    `json:"order_status"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AccrualUpdate struct {
	AccrualStatus string       `json:"accrual_status"`
	AccrualAmount moneys.Money `json:"accrual_amount"`
	UpdatedAt     time.Time    `json:"updated_at"`
	OrderID       int64        `json:"order_id"`
}
