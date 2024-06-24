package repository

import (
	"errors"
	"fmt"
	"github.com/Nicolas-ggd/ch-mod/internal/db/models"
	"github.com/Nicolas-ggd/ch-mod/internal/db/models/request"
	"github.com/Nicolas-ggd/ch-mod/pkg/common"
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

func (r *UserRepository) Register(user *models.Users) (*models.Users, error) {
	err := r.DB.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.Users, error) {
	var user models.Users

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

func (r *UserRepository) ChangePassword(newPassword *request.SetPasswordRequest, userId uint) error {
	hash, err := common.HashPassword(newPassword.Password)
	if err != nil {
		return err
	}

	result := r.DB.Model(&models.Users{}).Where("id = ?", userId).Update("password", hash)
	if result.Error != nil {
		return err
	}

	return nil
}

func (r *UserRepository) FindByEmail(email string) (*[]models.Users, error) {
	var user []models.Users

	err := r.DB.Debug().Where("email LIKE ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with email: %s doesn't exist", email)
		}

		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByID(id uint) (*models.Users, error) {
	var user *models.Users

	err := r.DB.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with ID: %d doesn't exist", id)
		}

		return nil, err
	}

	return user, nil
}
