package v0

import (
	"context"
	"net/http"
	"slices"
	"time"

	handler "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler/getUserWithdrawals/v0"
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

	items, err := srv.orderRepository.OrdersListWithdrawalsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	length := len(items)
	if length == 0 {
		return &handler.Response{
			StatusCode:  http.StatusNoContent,
			Withdrawals: []handler.WithdrawalItem{},
		}, nil
	}

	withdrawals := make([]withdrawalItem, 0, length)
	for _, item := range items {
		withdrawals = append(withdrawals, withdrawalItem{
			WithdrawalItem: handler.WithdrawalItem{
				Number:      item.Order.OrderNumber,
				Sum:         moneys.New(item.Withdrawal.WithdrawalAmount),
				ProcessedAt: item.Order.CreatedAt.Format(time.RFC3339),
			},
			sortTS: item.Order.CreatedAt,
		})
	}

	slices.SortFunc(withdrawals, func(a, b withdrawalItem) int {
		return -1 * a.sortTS.Compare(b.sortTS)
	})

	resp = &handler.Response{
		StatusCode:  http.StatusOK,
		Withdrawals: make([]handler.WithdrawalItem, 0, length),
	}
	for _, withdrawal := range withdrawals {
		resp.Withdrawals = append(resp.Withdrawals, withdrawal.Convert())
	}

	return resp, nil
}
