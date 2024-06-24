package models

import (
	"fmt"
	"github.com/Nicolas-ggd/ch-mod/pkg/common"
	"gorm.io/gorm"
	"time"
)

type Users struct {
	ID          uint           `json:"id" gorm:"primary_key"`
	FullName    string         `json:"full_name"`
	Email       string         `json:"email" gorm:"unique"`
	BirthDate   time.Time      `json:"birth_date"`
	PhoneNumber string         `json:"phone_number"`
	Gender      string         `json:"gender"`
	Password    string         `json:"-"`
	UserToken   []*UserToken   `json:"-" gorm:"foreignKey:UserID;constraint:onDelete:CASCADE"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index" `
}

// BeforeCreate is GORM built-in hook, it's used to generate user unique password before record is created in database
func (u *Users) BeforeCreate(tx *gorm.DB) (err error) {
	hash, err := common.HashPassword(u.Password)
	if err != nil {
		return fmt.Errorf("failed to generate password hash: %v", err)
	}
	u.Password = hash

	return nil
}
