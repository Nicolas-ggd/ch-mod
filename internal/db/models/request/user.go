package request

import (
	"github.com/Nicolas-ggd/ch-mod/internal/db/models"
	"time"
)

type UserRegisterRequest struct {
	FullName        string `json:"full_name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	BirthDate       string `json:"birth_date" binding:"required"`
	PhoneNumber     string `json:"phone_number"`
	Gender          string `json:"gender" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (ur *UserRegisterRequest) ToModel() *models.User {
	return &models.User{
		FullName:    ur.FullName,
		Email:       ur.Email,
		BirthDate:   parseTime(ur.BirthDate),
		PhoneNumber: ur.PhoneNumber,
		Password:    ur.Password,
		Gender:      ur.Gender,
	}
}

func parseTime(datetimeStr string) time.Time {
	t, err := time.Parse(time.RFC3339, datetimeStr)
	if err != nil {
		return time.Time{}
	}

	return t
}
