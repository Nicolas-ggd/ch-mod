package models

import (
	"gorm.io/gorm"
	"time"
)

type Chat struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	From      uint      `json:"from" gorm:"index"`
	To        uint      `json:"to" gorm:"index"`
	Name      string    `json:"name"`
	Message   []Message `json:"message" gorm:"foreignKey:ChatID;onDelete:CASCADE"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt time.Time `json:"-" gorm:"index"`
}

type Message struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	ChatID    uint           `json:"chat_id" gorm:"index"`
	From      uint           `json:"from"`
	Content   string         `json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
