package v0

import (
	"context"
	"fmt"
	"time"

	handler "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler/postUserRegister/v0"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/entity"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg"
)

type Service struct {
	userRepository UserRepository
	hashRepository HashRepository
}

func New(userRepository UserRepository, hashRepository HashRepository) *Service {
	return &Service{
		userRepository: userRepository,
		hashRepository: hashRepository,
	}
}

func (srv *Service) Do(ctx context.Context, r handler.Request) (resp *handler.Response, err error) {
	hash, err := srv.hashRepository.Hash(ctx, []byte(r.Password))
	if err != nil {
		return nil, err
	}

	ts := time.Now()
	createResp, err := srv.userRepository.UsersCreate(ctx, r.Login, entity.User{
		UserName:     r.Login,
		PasswordHash: hash,
		CreatedAt:    ts,
		UpdatedAt:    ts,
	})
	if err != nil {
		return nil, err
	}
	if createResp.AlreadyExists {
		return nil, pkg.ErrConflict
	}
	if !createResp.Ok {
		return nil, fmt.Errorf("user not ok, username=%s", r.Login)
	}

	return &handler.Response{}, nil
}
