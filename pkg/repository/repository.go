package repository

import "gorm.io/gorm"

type Repository struct {
	UserRepository  UserRepository
	TokenRepository TokenRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		UserRepository:  NewUserRepository(db),
		TokenRepository: NewTokenRepository(db),
	}
}
