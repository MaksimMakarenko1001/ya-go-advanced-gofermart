package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	srv "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/getAccrualInfoByOrders/v0"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/backoff"
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
	respCh := make(chan srv.AccrualResponse, len(orderNumbers))
	defer close(respCh)

	for _, id := range orderNumbers {

		go func(orderId string) {
			r.semaphore <- struct{}{}
			defer func() { <-r.semaphore }()

			resp, err := r.sendWithBackoff(ctx, orderId)
			respCh <- srv.AccrualResponse{
				Err:     err,
				Payload: resp,
			}

		}(id)
	}

	res := make([]srv.AccrualResponse, 0, len(orderNumbers))
	for item := range respCh {
		res = append(res, item)
	}

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
	u := url.URL{
		Scheme: "http",
		Host:   r.address,
		Path:   "/api/orders/" + url.PathEscape(orderNumber),
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
		errResp = fmt.Errorf("%w{order-id=%v}", errAccuralNoContent, orderNumber)
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
