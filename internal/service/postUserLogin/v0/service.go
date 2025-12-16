package v0

import (
	"context"

	handler "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler/postUserLogin/v0"
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
	user, err := srv.userRepository.UsersGetByUserName(ctx, r.Login)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, pkg.ErrUnauthorized
	}

	hash, err := srv.hashRepository.Hash(ctx, []byte(r.Password))
	if err != nil {
		return nil, err
	}

	if user.PasswordHash != hash {
		return nil, pkg.ErrUnauthorized
	}

	return &handler.Response{}, nil
}
