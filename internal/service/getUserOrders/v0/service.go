package v0

import (
	"context"
	"net/http"
	"slices"
	"time"

	handler "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler/getUserOrders/v0"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/types/gofermart"
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

	items, err := srv.orderRepository.OrdersListAccrualsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	length := len(items)
	if length == 0 {
		return &handler.Response{
			StatusCode: http.StatusNoContent,
			Accruals:   []handler.AccrualItem{},
		}, nil
	}

	accruals := make([]accrualItem, 0, length)
	for _, item := range items {
		accrual := accrualItem{
			AccrualItem: handler.AccrualItem{
				Number:     item.Order.OrderNumber,
				Status:     gofermart.OrderStatusType(item.Order.OrderStatus),
				UploadedAt: item.Order.CreatedAt.Format(time.RFC3339),
			},
			sortTS: item.Order.CreatedAt,
		}

		if item.Accrual.AccruedAt != nil {
			accrual.Accrual = pkg.ToPtr(moneys.New(item.Accrual.AccrualAmount))
		}
		accruals = append(accruals, accrual)
	}

	slices.SortFunc(accruals, func(a, b accrualItem) int {
		return -1 * a.sortTS.Compare(b.sortTS)
	})

	resp = &handler.Response{
		StatusCode: http.StatusOK,
		Accruals:   make([]handler.AccrualItem, 0, length),
	}
	for _, accrual := range accruals {
		resp.Accruals = append(resp.Accruals, accrual.Convert())
	}

	return resp, nil
}
