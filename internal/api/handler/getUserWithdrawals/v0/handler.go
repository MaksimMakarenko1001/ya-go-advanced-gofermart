package v0

import (
	"context"
	"net/http"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/types/moneys"
)

const MethodPath = "/api/user/withdrawals"
const Method = http.MethodGet

type HandlerFunc func(ctx context.Context, r Request) (resp *Response, err error)

type Request struct {
	AccessToken string `json:"accessToken"`
}

type Response struct {
	StatusCode  int
	Withdrawals []WithdrawalItem
}

type WithdrawalItem struct {
	Number      string       `json:"number"`
	Sum         moneys.Money `json:"sum"`
	ProcessedAt string       `json:"processed_at"`
}

func Handle(handlerFunc HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := handlerFunc(r.Context(), Request{
			AccessToken: r.Header.Get(handler.HeaderAccessToken),
		})
		if err != nil {
			handler.WriteError(w, err)
			return
		}

		handler.WriteJSONResult(w, resp.Withdrawals, resp.StatusCode)
	}
}
