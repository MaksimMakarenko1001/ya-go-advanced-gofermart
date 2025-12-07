package v0

import (
	"context"
	"time"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/entity"
	v0 "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/getAccrualInfoByOrders/v0"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/types/gofermart"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/types/moneys"
)

type Service struct {
	config                 Config
	orderRepository        OrderRepository
	getAccrualInfoByOrders *v0.Service
}

func New(config Config, orderRepository OrderRepository, getAccrualInfoByOrders *v0.Service) *Service {
	return &Service{
		config:                 config,
		orderRepository:        orderRepository,
		getAccrualInfoByOrders: getAccrualInfoByOrders,
	}
}

func (srv *Service) Do(ctx context.Context) (err error) {
	items, err := srv.orderRepository.OrdersListAccrualsByOrderStatus(ctx, gofermart.OrderStatusNew.String(), srv.config.Limit)
	if err != nil {
		return err
	}

	accrualResp, err := srv.getAccrualInfoByOrders.Do(ctx, pkg.Select(items, func(x entity.AccrualItem) string { return x.Order.OrderNumber }))
	if err != nil {
		return err
	}

	orderUpdates := make([]entity.OrderUpdate, 0, len(accrualResp))
	accrualUpdates := make([]entity.AccrualUpdate, 0, len(accrualResp))
	userBalanceMap := make(map[int64]entity.UserBalanceUpdate, len(accrualResp))
	for _, item := range items {
		payload, ok := accrualResp[item.Order.OrderNumber]
		if !ok {
			continue
		}

		var orderStatus gofermart.OrderStatusType
		switch payload.AccrualStatus {
		case gofermart.AccrualStatusInvalid:
			orderStatus = gofermart.OrderStatusInvalid
		case gofermart.AccrualStatusProcessed:
			orderStatus = gofermart.OrderStatusProcessed
		case gofermart.AccrualStatusRegistered, gofermart.AccrualStatusProcessing:
			orderStatus = gofermart.OrderStatusProcessing
		default:
			orderStatus = gofermart.OrderStatusNew
		}

		orderUpdates = append(orderUpdates, entity.OrderUpdate{
			OrderNumber: item.Order.OrderNumber,
			OrderStatus: orderStatus.String(),
			UpdatedAt:   time.Now(),
		})
		accrualUpdates = append(accrualUpdates, entity.AccrualUpdate{
			AccrualStatus: payload.AccrualStatus.String(),
			AccrualAmount: payload.AccrualAmount,
			UpdatedAt:     time.Now(),
			OrderID:       item.Accrual.OrderID,
		})

		userBalance, ok := userBalanceMap[item.Order.UserID]
		if !ok {
			userBalance = entity.UserBalanceUpdate{
				AccrualAmount:    moneys.New(0),
				WithdrawalAmount: moneys.New(0),
				UpdatedAt:        time.Now(),
				UserID:           item.Order.UserID,
			}
		}
		userBalance.AccrualAmount = userBalance.AccrualAmount.Add(payload.AccrualAmount)
		userBalanceMap[item.Order.UserID] = userBalance

	}

	return srv.orderRepository.OrdersUpdateAccruals(ctx, orderUpdates, accrualUpdates, pkg.ValuesToList(userBalanceMap))

}
