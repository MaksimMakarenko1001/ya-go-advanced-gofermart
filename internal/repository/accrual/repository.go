package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	srv "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/getAccrualInfoByOrderIds/v0"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/backoff"
)

type Repository struct {
	semaphore chan struct{}
	address   string
	backoff   *backoff.LinearBackoff
	client    *http.Client
}

func New(address string, reqTimeout time.Duration, throttlingRate uint, backoff *backoff.LinearBackoff) *Repository {
	return &Repository{
		semaphore: make(chan struct{}, throttlingRate),
		address:   address,
		backoff:   backoff,
		client: &http.Client{
			Timeout: reqTimeout,
		},
	}
}

func (r *Repository) AccrualsListInfoByOrderIds(ctx context.Context, orderIds []string) ([]srv.AccrualResponse, error) {
	var wg sync.WaitGroup

	respCh := make(chan srv.AccrualResponse, len(orderIds))
	defer close(respCh)

	for _, id := range orderIds {
		wg.Add(1)

		go func(orderId string) {
			r.semaphore <- struct{}{}
			defer wg.Done()
			defer func() { <-r.semaphore }()

			resp, err := r.sendWithBackoff(ctx, orderId)
			respCh <- srv.AccrualResponse{
				Err:     err,
				Payload: resp,
			}

		}(id)
	}
	wg.Wait()

	res := make([]srv.AccrualResponse, 0, len(orderIds))
	for item := range respCh {
		res = append(res, item)
	}

	return res, nil
}

func (r *Repository) sendWithBackoff(ctx context.Context, orderID string) (res *srv.AccrualPayload, err error) {
	fn := func(ctxBackoff context.Context) error {
		var errBackoff error
		res, errBackoff = r.send(ctxBackoff, orderID)
		return errBackoff
	}

	backoff := r.backoff.WithRetry()
	err = backoff(fn)(ctx)
	return res, err
}

func (r *Repository) send(ctx context.Context, orderID string) (res *srv.AccrualPayload, err error) {
	u := url.URL{
		Scheme: "http",
		Host:   r.address,
		Path:   "/api/orders/" + url.PathEscape(orderID),
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
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
		errResp = fmt.Errorf("%w{order-id=%v}", errAccuralNoContent, orderID)
	case http.StatusTooManyRequests:
		errResp = fmt.Errorf("%w{retry-after=%vs}", errAccuralTooManyRequests, response.Header.Get("Retry-After"))
	case http.StatusInternalServerError:
		errResp = fmt.Errorf("%w", errAccuralInternalServer)
	default:
		errResp = fmt.Errorf("accural unhandled status, %v", response.StatusCode)
	}

	_ = response.Body.Close()

	return res, errResp
}
