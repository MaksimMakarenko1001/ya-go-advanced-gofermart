package v0

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/types/moneys"
)

const MethodPath = "/api/user/balance/withdraw"
const Method = http.MethodPost

type HandlerFunc func(ctx context.Context, r Request) (resp *Response, err error)

type Request struct {
	Order  string       `json:"order"`
	Sum    moneys.Money `json:"sum"`
	UserID int64        `json:"userId"`
}

type Response struct {
}

func Handle(handlerFunc HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handler.WriteError(w, fmt.Errorf("request body not ok, %w", err))
			return
		}

		_, err := handlerFunc(r.Context(), req)
		if err != nil {
			handler.WriteError(w, err)
			return
		}

		handler.WriteOK(w, http.StatusOK)
	}
}
