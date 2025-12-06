package v0

import (
	"context"
	"errors"
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

	ctxRepo, cancel := context.WithTimeout(ctx, srv.config.AccrualTimeout)
	defer cancel()

	accrualResp, err := srv.accrualRepository.AccrualsGetInfoByOrders(ctxRepo, orderNumbers)
	if err != nil {
		return nil, err
	}

	errs := make([]error, 0, len(orderNumbers))
	response := make(map[string]AccrualPayload, len(orderNumbers))
	for _, resp := range accrualResp {
		err := resp.Err
		if err == nil && resp.Payload != nil {
			response[resp.Payload.OrderNumber] = *resp.Payload
		}
		errs = append(errs, err)
	}

	return response, errors.Join(errs...)
}
