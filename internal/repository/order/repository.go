package order

import (
	"context"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/db"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/entity"
	postUserBalanceWithdraw "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/postUserBalanceWithdraw/v0"
	postUserOrders "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/postUserOrders/v0"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/types/gofermart"
)

type Repository struct {
	db *db.PGConnect
}

func New(db *db.PGConnect) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) OrdersCreateAccrual(
	ctx context.Context,
	orderNumber string,
	order entity.Order,
	accrual entity.Accrual,
) (resp *postUserOrders.Response, err error) {
	err = r.db.QueryWithOneResultJSON(
		ctx,
		&resp,
		"select orders.orders_create_accrual(_order_number=>$1, _order=>$2, _accrual=>$3);",
		orderNumber, order, accrual,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *Repository) OrdersCreateWithdrawal(
	ctx context.Context,
	orderNumber string,
	order entity.Order,
	withdrawal entity.Withdrawal,
	userBalance entity.UserBalance,
) (resp *postUserBalanceWithdraw.Response, err error) {
	err = r.db.QueryWithOneResultJSON(
		ctx,
		&resp,
		"select orders.orders_create_withdrawal(_order_number=>$1, _order=>$2, _withdrawal=>$3, _user_balance=>$4);",
		orderNumber, order, withdrawal, userBalance,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *Repository) OrdersListAccrualsByOrderStatus(ctx context.Context, status gofermart.OrderStatusType, limit int) (items []entity.AccrualItem, err error) {
	err = r.db.QueryWithOneResultJSON(
		ctx,
		&items,
		"select orders.orders_list_accruals_by_order_status(_order_status=>$1, _limit=>$2);",
		status.String(), limit,
	)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *Repository) OrdersUpdateAccruals(
	ctx context.Context, orders []entity.Order, accruals []entity.Accrual, userBalances []entity.UserBalance,
) (orderUpdatedNumbers []string, err error) {
	err = r.db.QueryWithOneResultJSON(
		ctx,
		&orderUpdatedNumbers,
		"select orders.orders_update_accruals(_accruals=>$1, _orders=>$2, _user_balances=>$3);",
		accruals, orders, userBalances,
	)
	if err != nil {
		return nil, err
	}

	return orderUpdatedNumbers, nil
}

func (r *Repository) OrdersListAccrualsByUserId(ctx context.Context, userId int64) (items []entity.AccrualItem, err error) {
	err = r.db.QueryWithOneResultJSON(
		ctx,
		&items,
		"select orders.orders_list_accruals_by_user_id(_user_id=>$1);",
		userId,
	)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *Repository) OrdersGetUserBalanceByUserId(ctx context.Context, userId int64) (balance *entity.UserBalance, err error) {
	err = r.db.QueryWithOneResultJSON(
		ctx,
		&balance,
		"select orders.orders_get_user_balance_by_user_id(_user_id=>$1);",
		userId,
	)
	if err != nil {
		return nil, err
	}

	return balance, nil
}
func (r *Repository) OrdersListWithdrawalsByUserId(ctx context.Context, userId int64) (items []entity.WithdrawalItem, err error) {
	err = r.db.QueryWithOneResultJSON(
		ctx,
		&items,
		"select orders.orders_list_withdrawals_by_user_id(_user_id=>$1);",
		userId,
	)
	if err != nil {
		return nil, err
	}

	return items, nil
}
