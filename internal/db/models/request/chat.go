package request

import "github.com/Nicolas-ggd/ch-mod/internal/db/models"

type ChatRequest struct {
	To      uint           `json:"to"`
	From    uint           `json:"from"`
	Name    string         `json:"name"`
	Message MessageRequest `json:"message"`
}

type MessageRequest struct {
	From    uint   `json:"from"`
	Content string `json:"content"`
	ChatID  uint   `json:"chat_id"`
}

func (cr *ChatRequest) ToModel() *models.Chat {
	return &models.Chat{
		To:   cr.To,
		From: cr.From,
		Name: cr.Name,
	}
}

func (mr *MessageRequest) ToModel() *models.Message {
	return &models.Message{
		ChatID:  mr.ChatID,
		From:    mr.From,
		Content: mr.Content,
	}
}
