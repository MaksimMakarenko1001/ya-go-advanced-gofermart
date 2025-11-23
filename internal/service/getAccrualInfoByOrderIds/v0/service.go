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

func (srv *Service) Do(ctx context.Context, orderIds []string) (map[string]AccrualPayload, error) {
	if len(orderIds) == 0 {
		return nil, nil
	}

	ctxRepo, cancel := context.WithTimeout(ctx, srv.config.AccrualTimeout)
	defer cancel()

	accrualResp, err := srv.accrualRepository.AccrualsListInfoByOrderIds(ctxRepo, orderIds)
	if err != nil {
		return nil, err
	}

	errs := make([]error, 0, len(orderIds))
	response := make(map[string]AccrualPayload, len(orderIds))
	for _, resp := range accrualResp {
		err := resp.Err
		if err == nil && resp.Payload != nil {
			response[resp.Payload.Order] = *resp.Payload
		}
		errs = append(errs, err)
	}

	return response, errors.Join(errs...)
}
