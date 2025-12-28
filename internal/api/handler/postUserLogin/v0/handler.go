package v0

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler"
)

const MethodPath = "/api/user/login"
const Method = http.MethodPost

type HandlerFunc func(ctx context.Context, r Request) (resp *Response, err error)

type Request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Response struct {
	Token string
}

func Handle(handlerFunc HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handler.WriteError(w, fmt.Errorf("request body not ok, %w", err))
			return
		}

		resp, err := handlerFunc(r.Context(), req)
		if err != nil {
			handler.WriteError(w, err)
			return
		}

		w.Header().Set("Authorization", fmt.Sprintf("Bearer %v", resp.Token))
		handler.WriteOK(w, http.StatusOK)
	}
}
