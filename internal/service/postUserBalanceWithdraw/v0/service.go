package v0

import (
	"context"
	"fmt"
	"strconv"
	"time"

	handler "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler/postUserBalanceWithdraw/v0"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/entity"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/types/gofermart"
)

const lockKey = "postUserBalanceWithdraw"

type Service struct {
	config          Config
	orderRepository OrderRepository
	jwtRepository   JwtRepository
	locker          Locker
}

func New(config Config, orderRepository OrderRepository, jwtRepository JwtRepository, locker Locker) *Service {
	return &Service{
		config:          config,
		orderRepository: orderRepository,
		jwtRepository:   jwtRepository,
		locker:          locker,
	}
}

func (srv *Service) Do(ctx context.Context, r handler.Request) (resp *handler.Response, err error) {
	if !service.IsLuhnValid(r.Order) {
		return nil, pkg.ErrUnprocessableEntity
	}

	userID, err := srv.jwtRepository.JwtGetUserID(r.AccessToken)
	if err != nil {
		return nil, err
	}

	ts := time.Now()
	userIDStr := strconv.FormatInt(userID, 10)

	ok, err := srv.locker.LockAcquire(ctx, lockKey, userIDStr, ts.Add(srv.config.LockInterval), "")
	if err != nil {
		return nil, fmt.Errorf("error to acquire lock: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("failed to acquire lock")
	}
	defer func() {
		srv.locker.LockRelease(ctx, lockKey, userIDStr, "")
	}()

	balance, err := srv.orderRepository.OrdersGetUserBalanceByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if balance == nil {
		return nil, fmt.Errorf("balance not ok, user_id=%v", userID)
	}

	if current := balance.AccrualAmount - balance.WithdrawalAmount; current < r.Sum.Amount() {
		return nil, pkg.ErrPaymentRequired
	}

	createResp, err := srv.orderRepository.OrdersCreateWithdrawal(ctx, r.Order,
		entity.Order{
			OrderNumber: r.Order,
			OrderStatus: gofermart.OrderStatusProcessed.String(),
			CreatedAt:   ts,
			UpdatedAt:   ts,
			UserID:      userID,
		},
		entity.Withdrawal{
			WithdrawalAmount: r.Sum.Amount(),
			CreatedAt:        ts,
			UpdatedAt:        ts,
		},
		entity.UserBalance{
			WithdrawalAmount: r.Sum.Amount(),
			UpdatedAt:        ts,
			UserID:           userID,
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
