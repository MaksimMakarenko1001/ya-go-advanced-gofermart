package v0

import (
	"context"
	"net/http"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/types/moneys"
)

const MethodPath = "/api/user/balance"
const Method = http.MethodGet

type HandlerFunc func(ctx context.Context, r Request) (resp *Response, err error)

type Request struct {
	UserID int64 `json:"userId"`
}

type Response struct {
	Current   moneys.Money `json:"current"`
	Withdrawn moneys.Money `json:"withdrawn"`
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

		handler.WriteJSONResult(w, resp, http.StatusOK)
	}
}
