package v0

import (
	"context"

	handler "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler/postUserLogin/v0"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg"
)

type Service struct {
	userRepository UserRepository
	hashRepository HashRepository
	jwtRepository  JwtRepository
}

func New(userRepository UserRepository, hashRepository HashRepository, jwtRepository JwtRepository) *Service {
	return &Service{
		userRepository: userRepository,
		hashRepository: hashRepository,
		jwtRepository:  jwtRepository,
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

	hash, err := srv.hashRepository.HashSHA256(ctx, []byte(r.Password))
	if err != nil {
		return nil, err
	}

	if user.PasswordHash != hash {
		return nil, pkg.ErrUnauthorized
	}

	token, err := srv.jwtRepository.JwtGenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &handler.Response{Token: token}, nil
}
