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

func (r *ChatRepository) FindByUser(fromId, userId uint) (*[]models.Chat, error) {
	var chats []models.Chat

	err := r.DB.
		Joins("JOIN chat_users cu1 ON chats.id = cu1.chat_id").
		Joins("JOIN chat_users cu2 ON chats.id = cu2.chat_id").
		Where("cu1.user_id = ?", userId).
		Where("cu2.user_id = ?", fromId).
		Preload("Users.User").
		Preload("Messages").
		Find(&chats).Error
	if err != nil {
		return nil, err
	}

	return &chats, nil
}

func (r *ChatRepository) UserConversations(userId uint) (*[]models.Chat, error) {
	var model []models.Chat

	err := r.DB.
		Debug().
		Table("chats").
		Select("DISTINCT chats.*").
		Joins("JOIN chat_users cu ON cu.chat_id = chats.id").
		Joins("JOIN messages ON messages.chat_id = chats.id").
		Where("cu.user_id != ?", userId).
		Preload("Users.User").
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("messages.id ASC")
		}).
		Find(&model).Error
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (r *ChatRepository) Create(model *models.Chat) (*models.Chat, error) {
	if model.IsPrivate {
		var existingChat models.Chat
		err := r.DB.
			Where("name = ?", model.Name).
			Preload("Users.User").
			Preload("Messages").
			First(&existingChat).Error
		if err == nil {
			for _, chat := range model.Messages {
				chat.ChatID = existingChat.ID
				err = r.DB.Create(&chat).Error
				if err != nil {
					return nil, err
				}
			}
			return &existingChat, err
		}
		if existingChat.ID > 0 {
			return &existingChat, nil
		}
	}

	// Create the chat record in the database
	err := r.DB.Create(model).Error
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (r *ChatRepository) FindByName(name string) ([]*models.Chat, error) {
	var chat []*models.Chat
	err := r.DB.Where("name LIKE ?", "%"+name+"%").Find(&chat).Error
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
