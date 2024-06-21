package services

import "github.com/Nicolas-ggd/ch-mod/pkg/repository"

type Service struct {
	AuthService IAuthService
}

func NewService(repositories *repository.Repository) *Service {
	return &Service{
		AuthService: NewAuthService(repositories.UserRepository, repositories.TokenRepository),
	}
}
