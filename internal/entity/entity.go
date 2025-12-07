package entity

import (
	"time"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/types/moneys"
)

type Order struct {
	OrderNumber string    `json:"order_number"`
	OrderStatus string    `json:"order_status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserID      int64     `json:"user_id"`
}

type Accrual struct {
	AccrualStatus string       `json:"accrual_status"`
	AccrualAmount moneys.Money `json:"accrual_amount"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	OrderID       int64        `json:"order_id"`
}

type Withdrawal struct {
	WithdrawalAmount moneys.Money `json:"withdrawal_amount"`
	CreatedAt        time.Time    `json:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at"`
	OrderID          int64        `json:"order_id"`
}

type UserBalance struct {
	AccrualAmount    moneys.Money `json:"accrual_amount"`
	WithdrawalAmount moneys.Money `json:"withdrawal_amount"`
	CreatedAt        time.Time    `json:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at"`
	UserID           int64        `json:"user_id"`
}

type AccrualItem struct {
	Accrual Accrual `json:"accrual"`
	Order   Order   `json:"order"`
}

type WithdrawalItem struct {
	Withdrawal Withdrawal `json:"withdrawal"`
	Order      Order      `json:"order"`
}

type OrderUpdate struct {
	OrderNumber string    `json:"order_number"`
	OrderStatus string    `json:"order_status"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AccrualUpdate struct {
	AccrualStatus string       `json:"accrual_status"`
	AccrualAmount moneys.Money `json:"accrual_amount"`
	UpdatedAt     time.Time    `json:"updated_at"`
	OrderID       int64        `json:"order_id"`
}

type UserBalanceUpdate struct {
	AccrualAmount    moneys.Money `json:"accrual_amount"`
	WithdrawalAmount moneys.Money `json:"withdrawal_amount"`
	UpdatedAt        time.Time    `json:"updated_at"`
	UserID           int64        `json:"user_id"`
}
