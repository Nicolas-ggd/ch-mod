package models

import (
	"gorm.io/gorm"
	"time"
)

type Chat struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name"`
	Messages  []Message      `json:"messages" gorm:"foreignKey:ChatID;onDelete:CASCADE"`
	Users     []ChatUsers    `json:"users" gorm:"foreignKey:ChatID;onDelete:CASCADE"`
	IsPrivate bool           `json:"is_private"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Message struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	ChatID    uint           `json:"chat_id" gorm:"index"`
	FromID    uint           `json:"from_id" gorm:"index"`
	From      Users          `json:"-" gorm:"foreignKey:FromID;references:ID;onDelete:CASCADE"`
	Content   string         `json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type ChatUsers struct {
	UserId uint  `json:"user_id" gorm:"uniqueIndex:idx_unique_chat_users;primaryKey;autoIncrement:false;"`
	User   Users `json:"user" gorm:"foreignKey:UserId;references:ID;onDelete:CASCADE"`
	ChatID uint  `json:"chat_id" gorm:"uniqueIndex:idx_unique_chat_users;primaryKey;autoIncrement:false;"`
	Chat   Chat  `json:"chat" gorm:"foreignKey:ChatID;references:ID;onDelete:CASCADE"`
}
