package v0

import (
	"time"

	handler "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler/getUserWithdrawals/v0"
)

type withdrawalItem struct {
	handler.WithdrawalItem
	sortTS time.Time
}

func (wi withdrawalItem) Convert() handler.WithdrawalItem {
	return handler.WithdrawalItem{
		Order:       wi.Order,
		Sum:         wi.Sum,
		ProcessedAt: wi.ProcessedAt,
	}
}
