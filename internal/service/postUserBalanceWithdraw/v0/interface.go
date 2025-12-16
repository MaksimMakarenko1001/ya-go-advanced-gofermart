package v0

import (
	"context"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/entity"
)

type OrderRepository interface {
	OrdersGetUserBalanceByUserID(ctx context.Context, userID int64) (balance *entity.UserBalance, err error)

	OrdersCreateWithdrawal(
		ctx context.Context,
		orderNumber string,
		order entity.Order,
		withdrawal entity.Withdrawal,
		userBalance entity.UserBalance,
	) (resp *Response, err error)
}

type JwtRepository interface {
	JwtGetUserID(tokenString string) (int64, error)
}

type Response struct {
	Ok            bool  `json:"ok"`
	AlreadyExists bool  `json:"already_exists"`
	OrderID       int64 `json:"order_id"`
}
