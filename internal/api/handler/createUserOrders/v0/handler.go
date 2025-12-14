package v0

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler"
)

const MethodPath = "/api/user/orders"
const Method = http.MethodPost

var MiddlewareTypeContent = handler.MiddlewareTypeContentTextPlain

type HandlerFunc func(ctx context.Context, r Request) (resp *Response, err error)

type Request struct {
	OrderNumber string `json:"orderNumber"`
	UserID      int64  `json:"userId"`
}

type Response struct {
	Status int
}

func Handle(handlerFunc HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			handler.WriteError(w, fmt.Errorf("request body not ok, %w", err))
			return
		}

		resp, err := handlerFunc(r.Context(), Request{
			UserID:      12345,
			OrderNumber: string(body),
		})
		if err != nil {
			handler.WriteError(w, err)
			return
		}

		w.Header().Set("Content-Type", handler.TypeContentTextPlain)
		w.WriteHeader(resp.Status)
	}
}
