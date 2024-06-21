package repository

import (
	"fmt"
	"github.com/Nicolas-ggd/ch-mod/internal/db/models"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

var secret = os.Getenv("PRIVATE_SECRET")

type TokenRepository struct {
	DB *gorm.DB
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return TokenRepository{DB: db}
}

func (r *TokenRepository) CreateJWT(userId uint) (string, error) {

	now := time.Now().String()
	claims := &jwt.MapClaims{
		"iss":       "ch-mod.api",
		"ExpiresAt": 15000,
		"user": models.TokenClaim{
			UserId: userId,
			Role:   nil,
			Time:   &now,
		},
	}

	content := []byte(secret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(content)
	if err != nil {
		fmt.Printf("Error signing JWT: %v", err)
		return "", err
	}

	return ss, nil
}

func (r *TokenRepository) UpdateToken(userToken *models.UserToken) error {
	result := r.DB.Model(&models.UserToken{}).Where("user_id = ? AND type = ?", userToken.UserID, userToken.Type).Update("hash", userToken.Hash)
	if result.Error != nil {
		log.Printf("Failed to update user_token, where user_id = %v and type = %v, got error: %v", userToken.UserID, userToken.Type, result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		log.Printf("Record not found for UserID: %v and Type: %v", userToken.UserID, userToken.Type)
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *TokenRepository) CreateToken(token *models.UserToken) error {
	err := r.DB.Save(token).Error
	if err != nil {
		log.Printf("Error creating token: %v", err)
		return fmt.Errorf("failed to creating token: %v", err)
	}

	return nil
}

func (r *TokenRepository) DeleteToken(userId uint, tokenType models.Type) error {
	err := r.DB.
		Where("type = ?", tokenType).
		Where("user_id = ?", userId).
		Delete(&models.UserToken{}).Error

	if err != nil {
		log.Printf("Error deleting token: %v", err)

		return fmt.Errorf("failed to delete token: %v", err)
	}

	return nil
}
