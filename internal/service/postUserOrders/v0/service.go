package v0

import (
	"context"
	"fmt"
	"net/http"
	"time"

	handler "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler/postUserOrders/v0"
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
	if !service.IsLuhnValid(r.OrderNumber) {
		return nil, pkg.ErrUnprocessableEntity
	}

	ts := time.Now()

	createResp, err := srv.orderRepository.OrdersCreateAccrual(ctx, r.OrderNumber,
		entity.Order{
			OrderNumber: r.OrderNumber,
			OrderStatus: gofermart.OrderStatusNew.String(),
			CreatedAt:   ts,
			UpdatedAt:   ts,
			UserID:      r.UserID,
		},
		entity.Accrual{
			AccrualStatus: gofermart.AccrualStatusNew.String(),
			CreatedAt:     ts,
			UpdatedAt:     ts,
		},
		entity.UserBalance{
			CreatedAt: ts,
			UpdatedAt: ts,
			UserID:    r.UserID,
		},
	)
	if err != nil {
		return nil, err
	}
	if createResp.AlreadyExists && createResp.AlreadyExistsByUserId != r.UserID {
		return nil, pkg.ErrConflict
	}
	if createResp.AlreadyExists && createResp.AlreadyExistsByUserId == r.UserID {
		return &handler.Response{StatusCode: http.StatusOK}, nil
	}
	if !createResp.Ok {
		return nil, fmt.Errorf("order not ok, number=%s", r.OrderNumber)
	}

	return &handler.Response{StatusCode: http.StatusAccepted}, nil
}
