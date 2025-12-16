package v0

import (
	"context"
	"fmt"

	handler "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler/getUserBalance/v0"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/types/moneys"
)

type Service struct {
	orderRepository OrderRepository
	jwtRepository   JwtRepository
}

func New(orderRepository OrderRepository, jwtRepository JwtRepository) *Service {
	return &Service{
		orderRepository: orderRepository,
		jwtRepository:   jwtRepository,
	}
}

func (srv *Service) Do(ctx context.Context, r handler.Request) (resp *handler.Response, err error) {
	userID, err := srv.jwtRepository.JwtGetUserID(r.AccessToken)
	if err != nil {
		return nil, err
	}

	balance, err := srv.orderRepository.OrdersGetUserBalanceByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if balance == nil {
		return nil, fmt.Errorf("balance not ok, user_id=%v", userID)
	}

	return &handler.Response{
		Current:   moneys.New(balance.AccrualAmount - balance.WithdrawalAmount),
		Withdrawn: moneys.New(balance.WithdrawalAmount),
	}, nil
}
