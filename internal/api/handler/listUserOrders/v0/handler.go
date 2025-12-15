package v0

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/types/gofermart"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/types/moneys"
)

const MethodPath = "/api/user/orders"
const Method = http.MethodGet

type HandlerFunc func(ctx context.Context, r Request) (resp *Response, err error)

type Request struct {
	UserID int64 `json:"userId"`
}

type Response struct {
	StatusCode int
	Accruals   []AccrualItem
}

type AccrualItem struct {
	Number     string                    `json:"number"`
	Status     gofermart.OrderStatusType `json:"status"`
	Accrual    *moneys.Money             `json:"accrual,omitempty"`
	UploadedAt string                    `json:"uploaded_at"`
}

func Handle(handlerFunc HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := handlerFunc(r.Context(), Request{
			UserID: 12345,
		})
		if err != nil {
			handler.WriteError(w, err)
			return
		}

		res, err := json.Marshal(resp.Accruals)
		if err != nil {
			handler.WriteError(w, fmt.Errorf("accruals not ok, %w", err))
			return
		}

		handler.WriteJSONResult(w, res, resp.StatusCode)
	}
}
