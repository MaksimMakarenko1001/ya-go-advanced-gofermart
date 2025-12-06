package v0

import "context"

type AccrualRepository interface {
	AccrualsGetInfoByOrders(ctx context.Context, orderNumbers []string) ([]AccrualResponse, error)
}
