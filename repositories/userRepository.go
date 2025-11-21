package repositories

import (
	"project-management-backend/config"
	"project-management-backend/models"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(userID uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetAll() ([]models.User, error)
	Update(user *models.User) error
	Delete(userID uint) error
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) Create(user *models.User) error {
	return config.DB.Create(user).Error
}

func (r *userRepository) GetByID(userID uint) (*models.User, error) {
	var user models.User
	err := config.DB.Where("id = ?", userID).First(&user).Error
	return &user, err
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err // return error asli dari GORM
	}
	return &user, nil
}

func (r *userRepository) GetAll() ([]models.User, error) {
	var users []models.User
	err := config.DB.Find(&users).Error
	return users, err
}

func (r *userRepository) Update(user *models.User) error {
	return config.DB.Save(user).Error
}

func (r *userRepository) Delete(userID uint) error {
	return config.DB.Where("id = ?", userID).Delete(&models.User{}).Error
}
