package repository

import (
	"fmt"
	"github.com/Nicolas-ggd/ch-mod/internal/db/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{
		DB: db,
	}
}

func (r *UserRepository) Register(user *models.User) (*models.User, error) {
	err := r.DB.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User

	err := r.DB.Where("email = ?", email).
		Preload("UserToken").
		First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetAuthToken(userId uint) (*models.UserToken, error) {
	var tokenModel models.UserToken

	err := r.DB.Model(models.UserToken{}).
		Where("user_id = ?", userId).
		Where("type = ?", models.Auth).
		First(&tokenModel).
		Error

	if err != nil {
		return nil, err
	}

	return &tokenModel, nil
}

func (r *UserRepository) GetUserTokenByHash(hash string) (*models.UserToken, error) {
	var userToken models.UserToken
	if err := r.DB.Where("hash = ?", hash).First(&userToken).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve user by hash %s: %v", hash, err.Error())
	}

	return &userToken, nil
}
