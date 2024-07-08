package request

import (
	"github.com/Nicolas-ggd/ch-mod/internal/db/models"
	"strconv"
)

type ChatRequest struct {
	Name      string             `json:"name"`
	Message   []MessageRequest   `json:"message"`
	Users     []ChatUsersRequest `json:"users"`
	IsPrivate bool               `json:"is_private"`
}

type WsChatRequest struct {
	Name      string           `json:"name"`
	Message   []MessageRequest `json:"message"`
	Clients   []string         `json:"clients"`
	IsPrivate bool             `json:"is_private"`
}

type MessageRequest struct {
	FromID  uint   `json:"from_id"`
	Content string `json:"content"`
	ChatID  uint   `json:"chat_id"`
}

type ChatUsersRequest struct {
	UsersID []uint `json:"users_id"`
}

func (cr *WsChatRequest) ToModel() *models.Chat {
	c := &models.Chat{
		Name:      cr.Name,
		IsPrivate: cr.IsPrivate,
	}

	for _, userId := range cr.Clients {
		id, err := strconv.ParseUint(userId, 10, 64)
		if err != nil {
			return nil
		}

		c.Users = append(c.Users, models.ChatUsers{UserId: uint(id)})
	}

	for _, messageRequest := range cr.Message {
		c.Messages = append(c.Messages, models.Message{
			Content: messageRequest.Content,
			FromID:  messageRequest.FromID,
		})
	}

	return c
}
