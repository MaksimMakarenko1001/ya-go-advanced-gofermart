package v0

import (
	"context"
	"fmt"
	"time"

	handler "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler/postUserBalanceWithdraw/v0"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/entity"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/types/gofermart"
)

type Service struct {
	orderRepository OrderRepository
}

func New(orderRepository OrderRepository) *Service {
	return &Service{orderRepository: orderRepository}
}

func (srv *Service) Do(ctx context.Context, r handler.Request) (resp *handler.Response, err error) {
	if !service.IsLuhnValid(r.Order) {
		return nil, pkg.ErrUnprocessableEntity
	}

	balance, err := srv.orderRepository.OrdersGetUserBalanceByUserId(ctx, r.UserID)
	if err != nil {
		return nil, err
	}
	if balance == nil {
		return nil, fmt.Errorf("balance not ok, user_id=%v", r.UserID)
	}

	if current := balance.AccrualAmount - balance.WithdrawalAmount; current < r.Sum.Amount() {
		return nil, pkg.ErrPaymentRequired
	}

	ts := time.Now()
	createResp, err := srv.orderRepository.OrdersCreateWithdrawal(ctx, r.Order,
		entity.Order{
			OrderNumber: r.Order,
			OrderStatus: gofermart.OrderStatusProcessed.String(),
			CreatedAt:   ts,
			UpdatedAt:   ts,
			UserID:      r.UserID,
		},
		entity.Withdrawal{
			WithdrawalAmount: r.Sum.Amount(),
			CreatedAt:        ts,
			UpdatedAt:        ts,
		},
		entity.UserBalance{
			WithdrawalAmount: r.Sum.Amount(),
			UpdatedAt:        ts,
			UserID:           r.UserID,
		},
	)
	if err != nil {
		return nil, err
	}
	if createResp.AlreadyExists {
		return nil, pkg.ErrUnprocessableEntity
	}
	if !createResp.Ok {
		return nil, fmt.Errorf("order not ok, number=%s", r.Order)
	}

	return &handler.Response{}, nil
}
