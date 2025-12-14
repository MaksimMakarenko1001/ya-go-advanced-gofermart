package v0

import (
	"context"
	"fmt"
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
	items, err := srv.orderRepository.OrdersListAccrualsByOrderStatus(ctx, gofermart.OrderStatusNew, srv.config.Limit)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}

	accrualResp, err := srv.getAccrualInfoByOrders.Do(ctx, pkg.Select(items, func(x entity.AccrualItem) string { return x.Order.OrderNumber }))
	if err != nil {
		return err
	}

	ts := time.Now()
	orderUpdates := make([]entity.OrderUpdate, 0, len(items))
	accrualUpdates := make([]entity.AccrualUpdate, 0, len(items))
	userAccrualMap := make(map[int64]moneys.Money, len(items))
	orderUpdateHit := make(map[string]struct{}, len(items))
	for _, item := range items {
		payload, ok := accrualResp[item.Order.OrderNumber]
		if !ok {
			return fmt.Errorf("order accrual not ok, %s", item.Order.OrderNumber)
		}

		orderUpdateHit[item.Order.OrderNumber] = struct{}{}

		var orderStatus gofermart.OrderStatusType
		switch payload.AccrualStatus {
		case gofermart.AccrualStatusInvalid, gofermart.AccrualStatusNone:
			orderStatus = gofermart.OrderStatusInvalid
		case gofermart.AccrualStatusProcessed:
			orderStatus = gofermart.OrderStatusProcessed
		case gofermart.AccrualStatusRegistered, gofermart.AccrualStatusProcessing:
			orderStatus = gofermart.OrderStatusProcessing
		default:
			orderStatus = gofermart.OrderStatusNone
		}

		orderUpdates = append(orderUpdates, entity.OrderUpdate{
			OrderNumber: payload.OrderNumber,
			OrderStatus: orderStatus.String(),
			UpdatedAt:   ts,
		})
		accrualUpdates = append(accrualUpdates, entity.AccrualUpdate{
			AccrualStatus: payload.AccrualStatus.String(),
			AccrualAmount: payload.AccrualAmount.Amount(),
			UpdatedAt:     ts,
			OrderID:       item.Accrual.OrderID,
		})

		userAccrual, ok := userAccrualMap[item.Order.UserID]
		if !ok {
			userAccrual = moneys.Money{}
		}
		userAccrual = userAccrual.Add(payload.AccrualAmount)
		userAccrualMap[item.Order.UserID] = userAccrual
	}

	userBalanceUpdates := make([]entity.UserBalanceUpdate, 0, len(userAccrualMap))
	for userId, accrual := range userAccrualMap {
		userBalanceUpdates = append(userBalanceUpdates, entity.UserBalanceUpdate{
			AccrualAmount:    accrual.Amount(),
			WithdrawalAmount: 0,
			UpdatedAt:        ts,
			UserID:           userId,
		})
	}

	orderUpdatedNumbers, err := srv.orderRepository.OrdersUpdateAccruals(ctx, orderUpdates, accrualUpdates, userBalanceUpdates)
	if err != nil {
		return err
	}

	for _, number := range orderUpdatedNumbers {
		delete(orderUpdateHit, number)
	}
	if len(orderUpdateHit) > 0 {
		return fmt.Errorf("failed to update order numbers, %v", pkg.KeysToList(orderUpdateHit))
	}

	return nil
}
