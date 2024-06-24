package services

import "github.com/Nicolas-ggd/ch-mod/pkg/repository"

type Service struct {
	AuthService IAuthService
	ChatService IChatService
}

func NewService(repositories *repository.Repository) *Service {
	return &Service{
		AuthService: NewAuthService(repositories.UserRepository, repositories.TokenRepository),
		ChatService: NewChatService(repositories.ChatRepository),
	}
}
