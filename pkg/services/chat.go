package services

import (
	"github.com/Nicolas-ggd/ch-mod/internal/db/models"
	"github.com/Nicolas-ggd/ch-mod/internal/db/models/request"
	"github.com/Nicolas-ggd/ch-mod/pkg/repository"
)

type IChatService interface {
	Create(model *request.WsChatRequest) (*models.Chat, error)
	FindByUser(fromId, userId uint) (*[]models.Chat, error)
	UserConversations(userId uint) (*[]models.Chat, error)
}

type ChatService struct {
	chatRepository repository.ChatRepository
}

func NewChatService(chatRepository repository.ChatRepository) IChatService {
	return &ChatService{chatRepository: chatRepository}
}

func (cs *ChatService) Create(model *request.WsChatRequest) (*models.Chat, error) {
	chat, err := cs.chatRepository.Create(model.ToModel())
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (cs *ChatService) FindByUser(fromId, userId uint) (*[]models.Chat, error) {
	model, err := cs.chatRepository.FindByUser(fromId, userId)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (cs *ChatService) UserConversations(userId uint) (*[]models.Chat, error) {
	model, err := cs.chatRepository.UserConversations(userId)
	if err != nil {
		return nil, err
	}

	return model, nil
}
