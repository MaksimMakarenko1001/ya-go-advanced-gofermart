package v0

import (
	"time"

	handler "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler/getUserOrders/v0"
)

type accrualItem struct {
	handler.AccrualItem
	sortTS time.Time
}

func (ai accrualItem) Convert() handler.AccrualItem {
	return handler.AccrualItem{
		Number:     ai.Number,
		Status:     ai.Status,
		Accrual:    ai.Accrual,
		UploadedAt: ai.UploadedAt,
	}
}
