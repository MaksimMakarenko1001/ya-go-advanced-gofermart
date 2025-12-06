package v0

import (
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/types/gofermart"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/types/moneys"
)

type AccrualResponse struct {
	Err     error
	Payload *AccrualPayload
}

type AccrualPayload struct {
	OrderNumber   string                      `json:"order"`
	AccrualStatus gofermart.AccrualStatusType `json:"status"`
	AccrualAmount moneys.Money                `json:"accrual"`
}
