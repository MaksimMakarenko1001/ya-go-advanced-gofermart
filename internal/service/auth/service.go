package auth

import (
	"context"
)

type Service struct {
	jwtRepository JwtRepository
}

func New(JwtRepository JwtRepository) *Service {
	return &Service{
		jwtRepository: JwtRepository,
	}
}

func (srv *Service) ValidateToken(ctx context.Context, token string) (bool, error) {
	return srv.jwtRepository.JwtValidate(token)
}
