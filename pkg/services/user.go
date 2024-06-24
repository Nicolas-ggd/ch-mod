package services

import (
	"github.com/Nicolas-ggd/ch-mod/internal/db/models"
	"github.com/Nicolas-ggd/ch-mod/pkg/repository"
)

type IUserService interface {
	FindByEmail(email string) (*[]models.Users, error)
	FindByID(id uint) (*models.Users, error)
}

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) IUserService {
	return &UserService{
		userRepository: repository,
	}
}

func (us *UserService) FindByEmail(email string) (*[]models.Users, error) {
	user, err := us.userRepository.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) FindByID(id uint) (*models.Users, error) {
	user, err := us.userRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
