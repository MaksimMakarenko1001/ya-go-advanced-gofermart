package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	srv "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/getAccrualInfoByOrders/v0"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/backoff"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/types/gofermart"
)

type Repository struct {
	semaphore         chan struct{}
	address           string
	backoff           *backoff.LinearBackoff
	client            *http.Client
	retryAfterDefault time.Duration
}

func New(cfg Config, backoff *backoff.LinearBackoff) *Repository {
	return &Repository{
		semaphore:         make(chan struct{}, cfg.ThrottlingRate),
		address:           cfg.Address,
		retryAfterDefault: cfg.RetryAfter,
		backoff:           backoff,
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

func (r *Repository) AccrualsGetInfoByOrders(ctx context.Context, orderNumbers []string) ([]srv.AccrualResponse, error) {
	var wg sync.WaitGroup
	res := make([]srv.AccrualResponse, len(orderNumbers))

	for i, id := range orderNumbers {
		wg.Add(1)
		go func(idx int, orderId string) {
			r.semaphore <- struct{}{}

			defer func() { <-r.semaphore }()
			defer wg.Done()

			resp, err := r.sendWithBackoff(ctx, orderId)
			res[idx] = srv.AccrualResponse{Err: err, Payload: resp}

		}(i, id)
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
		}
	case http.StatusTooManyRequests:
		retryAfter, err := time.ParseDuration(response.Header.Get("Retry-After") + "s")
		if err != nil {
			retryAfter = r.retryAfterDefault
		}

		res = &srv.AccrualPayload{
			OrderNumber:   orderNumber,
			AccrualStatus: gofermart.AccrualStatusWaiting,
			AccrualAfter:  pkg.ToPtr(time.Now().Add(retryAfter)),
		}
	case http.StatusInternalServerError:
		errResp = fmt.Errorf("%w", errAccuralInternalServer)
	default:
		errResp = fmt.Errorf("accural unhandled http status, %v", response.StatusCode)
	}

	_ = response.Body.Close()

	return res, errResp
}
