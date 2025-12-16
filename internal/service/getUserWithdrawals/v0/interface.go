package v0

import (
	"context"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/entity"
)

type OrderRepository interface {
	OrdersListWithdrawalsByUserId(ctx context.Context, UserID int64) (items []entity.WithdrawalItem, err error)
}

type JwtRepository interface {
	JwtGetUserID(tokenString string) (int64, error)
}
