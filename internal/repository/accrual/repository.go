package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	srv "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/getAccrualInfoByOrders/v0"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/backoff"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/types/gofermart"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/types/moneys"
)

type Repository struct {
	semaphore chan struct{}
	address   string
	backoff   *backoff.LinearBackoff
	client    *http.Client
}

func New(cfg Config, backoff *backoff.LinearBackoff) *Repository {
	return &Repository{
		semaphore: make(chan struct{}, cfg.ThrottlingRate),
		address:   cfg.Address,
		backoff:   backoff,
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

func (r *Repository) AccrualsGetInfoByOrders(ctx context.Context, orderNumbers []string) ([]srv.AccrualResponse, error) {
	var wg sync.WaitGroup
	res := make([]srv.AccrualResponse, 0, len(orderNumbers))

	for _, id := range orderNumbers {
		wg.Add(1)
		go func(orderId string) {
			r.semaphore <- struct{}{}

			defer func() { <-r.semaphore }()
			defer wg.Done()

			resp, err := r.sendWithBackoff(ctx, orderId)
			res = append(res, srv.AccrualResponse{
				Err:     err,
				Payload: resp,
			})

		}(id)
	}
	wg.Wait()

	return res, nil
}

func (r *Repository) sendWithBackoff(ctx context.Context, orderNumber string) (res *srv.AccrualPayload, err error) {
	fn := func(ctxBackoff context.Context) error {
		var errBackoff error
		res, errBackoff = r.send(ctxBackoff, orderNumber)
		return errBackoff
	}

	backoff := r.backoff.WithRetry()
	err = backoff(fn)(ctx)
	return res, err
}

func (r *Repository) send(ctx context.Context, orderNumber string) (res *srv.AccrualPayload, err error) {
	path := "/api/orders/" + url.PathEscape(orderNumber)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, r.address+path, nil)
	if err != nil {
		return nil, fmt.Errorf("accural request not ok, %w", err)
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := r.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("accural http not ok, %w", err)
	}

	var errResp error
	switch response.StatusCode {
	case http.StatusOK:
		if err := json.NewDecoder(response.Body).Decode(&res); err != nil {
			errResp = fmt.Errorf("accural response not ok, %w", err)
		}
	case http.StatusNoContent:
		res = &srv.AccrualPayload{
			OrderNumber:   orderNumber,
			AccrualStatus: gofermart.AccrualStatusNone,
			AccrualAmount: moneys.Money{},
		}
	case http.StatusTooManyRequests:
		errResp = fmt.Errorf("%w{retry-after=%vs}", errAccuralTooManyRequests, response.Header.Get("Retry-After"))
	case http.StatusInternalServerError:
		errResp = fmt.Errorf("%w", errAccuralInternalServer)
	default:
		errResp = fmt.Errorf("accural unhandled http status, %v", response.StatusCode)
	}

	_ = response.Body.Close()

	return res, errResp
}
