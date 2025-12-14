package v0

import (
	"context"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/entity"
)

type OrderRepository interface {
	OrdersCreate(
		ctx context.Context,
		orderNumber string,
		order entity.Order,
		accrual entity.Accrual,
		withdrawal entity.Withdrawal,
		userBalance entity.UserBalance,
	) (resp *Response, err error)
}

type Response struct {
	Ok                    bool  `json:"ok"`
	AlreadyExists         bool  `json:"already_exists"`
	AlreadyExistsByUserId int64 `json:"already_exists_by_user_id"`
	OrderID               int64 `json:"order_id"`
}
