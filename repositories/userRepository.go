package repositories

import (
	"project-management-backend/config"
	"project-management-backend/models"
)

type UserRepository interface {
	GetAllUsers() ([]models.User, error)
	GetByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) error
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := config.DB.Preload("Projects").Preload("Tasks").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(user *models.User) error {
	return config.DB.Create(user).Error
}
