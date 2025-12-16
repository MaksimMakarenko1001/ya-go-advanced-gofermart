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

type HandlerFunc func(ctx context.Context, r Request) (resp *Response, err error)

type Request struct {
	OrderNumber string `json:"orderNumber"`
	AccessToken string `json:"accessToken"`
}

type Response struct {
	StatusCode int
}

func Handle(handlerFunc HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			handler.WriteError(w, fmt.Errorf("request body not ok, %w", err))
			return
		}

		resp, err := handlerFunc(r.Context(), Request{
			OrderNumber: string(body),
			AccessToken: r.Header.Get(handler.HeaderAccessToken),
		})
		if err != nil {
			handler.WriteError(w, err)
			return
		}

		handler.WriteOK(w, resp.StatusCode)
	}
}
