package repository

import (
	"github.com/Nicolas-ggd/ch-mod/internal/db/models"
	"gorm.io/gorm"
)

type ChatRepository struct {
	DB *gorm.DB
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return ChatRepository{
		DB: db,
	}
}

func (r *ChatRepository) Create(model *models.Chat) (*models.Chat, error) {
	err := r.DB.Where("from = ? AND to = ?", model.From, model.To).First(model).Error
	if err == nil {
		err = r.DB.Create(model).Error
		if err != nil {
			return nil, err
		}

		return nil, err
	}

	err = r.DB.Create(model.Message).Error
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (r *ChatRepository) FindByName(name string) ([]*models.Chat, error) {
	var chat []*models.Chat
	err := r.DB.Where("name LIKE ?", name).Find(&chat).Error
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (r *ChatRepository) DeleteConversation(id uint) error {
	err := r.DB.Delete(&models.Chat{}, id).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *ChatRepository) DeleteMessage(id uint) error {
	err := r.DB.Delete(&models.Message{}, id).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *ChatRepository) UpdateMessage(model *models.Message) (*models.Message, error) {
	err := r.DB.Where("id = ? AND chat_id = ?", model.ID, model.ChatID).Updates(&model).Error
	if err != nil {
		return nil, err
	}

	return model, nil
}
