package v0

import "context"

type AccrualRepository interface {
	AccrualsListInfoByOrderIds(ctx context.Context, orderIds []string) ([]AccrualResponse, error)
}
