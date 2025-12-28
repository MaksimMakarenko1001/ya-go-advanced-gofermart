package v0

import (
	"context"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/entity"
)

type UserRepository interface {
	UsersGetByUserName(ctx context.Context, userName string) (user *entity.User, err error)
}
type HashRepository interface {
	HashSHA256(ctx context.Context, message []byte) (string, error)
}

type JwtRepository interface {
	JwtGenerateToken(userID int64) (string, error)
}
