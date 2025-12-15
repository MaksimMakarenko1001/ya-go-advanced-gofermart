package v0

import (
	"context"
	"fmt"

	handler "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler/getUserBalance/v0"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/types/moneys"
)

type Service struct {
	orderRepository OrderRepository
}

func New(orderRepository OrderRepository) *Service {
	return &Service{orderRepository: orderRepository}
}

func (srv *Service) Do(ctx context.Context, r handler.Request) (resp *handler.Response, err error) {
	balance, err := srv.orderRepository.OrdersGetUserBalanceByUserId(ctx, r.UserID)
	if err != nil {
		return nil, err
	}
	if balance == nil {
		return nil, fmt.Errorf("balance not ok, user_id=%v", r.UserID)
	}

	return &handler.Response{
		Current:   moneys.New(balance.AccrualAmount - balance.WithdrawalAmount),
		Withdrawn: moneys.New(balance.WithdrawalAmount),
	}, nil
}
