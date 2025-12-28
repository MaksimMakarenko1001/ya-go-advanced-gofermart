package v0

import (
	"context"
	"errors"
	"log"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/types/gofermart"
)

type Service struct {
	config            Config
	accrualRepository AccrualRepository
}

func New(cfg Config, accrualRepo AccrualRepository) *Service {
	return &Service{
		config:            cfg,
		accrualRepository: accrualRepo,
	}
}

func (srv *Service) Do(ctx context.Context, orderNumbers []string) (map[string]AccrualPayload, error) {
	if len(orderNumbers) == 0 {
		return nil, nil
	}

	ctxRepo, cancel := context.WithTimeout(ctx, srv.config.Timeout)
	defer cancel()

	accrualResp, err := srv.accrualRepository.AccrualsGetInfoByOrders(ctxRepo, orderNumbers)
	if err != nil {
		return nil, err
	}
	if len(orderNumbers) != len(accrualResp) {
		log.Printf("accrual response count mismatch, expected=%d, actual=%d", len(orderNumbers), len(accrualResp))
	}

	errs := make([]error, 0, len(orderNumbers))
	res := make(map[string]AccrualPayload, len(orderNumbers))
	for _, resp := range accrualResp {
		err := resp.Err
		if err == nil && resp.Payload != nil {
			res[resp.Payload.OrderNumber] = *resp.Payload
		}
		errs = append(errs, err)
	}

	for _, number := range orderNumbers {
		if _, exists := res[number]; !exists {
			res[number] = AccrualPayload{
				OrderNumber:   number,
				AccrualStatus: gofermart.AccrualStatusNone,
			}
		}
	}

	return res, errors.Join(errs...)
}
