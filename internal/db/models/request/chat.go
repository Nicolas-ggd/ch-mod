package request

import "github.com/Nicolas-ggd/ch-mod/internal/db/models"

type ChatRequest struct {
	To        uint             `json:"to"`
	From      uint             `json:"from"`
	Name      string           `json:"name"`
	Message   []MessageRequest `json:"message"`
	IsPrivate bool             `json:"is_private"`
}

type WsChatRequest struct {
	To        uint             `json:"to"`
	From      uint             `json:"from"`
	Name      string           `json:"name"`
	Message   []MessageRequest `json:"message"`
	Clients   []string         `json:"clients"`
	IsPrivate bool             `json:"is_private"`
}

type MessageRequest struct {
	From    uint   `json:"from"`
	Content string `json:"content"`
	ChatID  uint   `json:"chat_id"`
}

func (cr *WsChatRequest) ToModel() *models.Chat {
	c := &models.Chat{
		To:        cr.To,
		From:      cr.From,
		Name:      cr.Name,
		IsPrivate: cr.IsPrivate,
	}

	for _, m := range cr.Message {
		c.Message = append(c.Message, models.Message{Content: m.Content, ChatID: m.ChatID, From: m.From})
	}

	return c
}
