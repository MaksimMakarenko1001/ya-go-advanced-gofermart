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
	items, err := srv.orderRepository.OrdersListAccrualsByOrderStatus(ctx, gofermart.OrderStatusNew, srv.config.Limit)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}

	ts := time.Now()
	items = pkg.FilterSlice(items, func(x entity.AccrualItem) bool { return pkg.Deref(x.Accrual.AccrueAfter).Before(ts) })

	accrualResp, err := srv.getAccrualInfoByOrders.Do(ctx, pkg.SelectSlice(items, func(x entity.AccrualItem) string { return x.Order.OrderNumber }))
	if err != nil {
		return err
	}

	orders := make([]entity.Order, 0, len(items))
	accruals := make([]entity.Accrual, 0, len(items))
	userAccrualMap := make(map[int64]moneys.Money, len(items))

	for _, item := range items {
		var accruedAt *time.Time
		var orderStatus gofermart.OrderStatusType
		payload := accrualResp[item.Order.OrderNumber]

		switch payload.AccrualStatus {
		case gofermart.AccrualStatusWaiting:
			orderStatus = gofermart.OrderStatusNew
		case gofermart.AccrualStatusInvalid, gofermart.AccrualStatusNone:
			orderStatus = gofermart.OrderStatusInvalid
		case gofermart.AccrualStatusRegistered, gofermart.AccrualStatusProcessing:
			orderStatus = gofermart.OrderStatusProcessing
		case gofermart.AccrualStatusProcessed:
			accruedAt = &ts
			orderStatus = gofermart.OrderStatusProcessed
		default:
			orderStatus = gofermart.OrderStatusNone
		}

		orders = append(orders, entity.Order{
			OrderNumber: item.Order.OrderNumber,
			OrderStatus: orderStatus.String(),
			UpdatedAt:   ts,
		})
		accruals = append(accruals, entity.Accrual{
			AccrualStatus: payload.AccrualStatus.String(),
			AccrualAmount: payload.AccrualAmount.Amount(),
			UpdatedAt:     ts,
			OrderID:       item.Accrual.OrderID,
			AccruedAt:     accruedAt,
		})

		userAccrual, ok := userAccrualMap[item.Order.UserID]
		if !ok {
			userAccrual = moneys.Money{}
		}
		userAccrual = userAccrual.Add(payload.AccrualAmount)
		userAccrualMap[item.Order.UserID] = userAccrual
	}

	userBalanceUpdates := make([]entity.UserBalance, 0, len(userAccrualMap))
	for userID, accrual := range userAccrualMap {
		userBalanceUpdates = append(userBalanceUpdates, entity.UserBalance{
			AccrualAmount: accrual.Amount(),
			UpdatedAt:     ts,
			UserID:        userID,
		})
	}

	_, err = srv.orderRepository.OrdersUpdateAccruals(ctx, orders, accruals, userBalanceUpdates)
	return err
}
