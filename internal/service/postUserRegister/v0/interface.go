package v0

import (
	"context"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/entity"
)

type UserRepository interface {
	UsersCreate(ctx context.Context, userName string, user entity.User) (resp *Response, err error)
}
type HashRepository interface {
	Hash(ctx context.Context, message []byte) (string, error)
}

type Response struct {
	Ok            bool  `json:"ok"`
	AlreadyExists bool  `json:"already_exists"`
	UserID        int64 `json:"user_id"`
}
