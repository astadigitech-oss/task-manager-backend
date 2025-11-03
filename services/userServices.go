package services

import (
	"errors"
	"project-management-backend/models"
	"project-management-backend/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetAllUsers() ([]models.User, error)
	CreateUser(user *models.User) error
	GetUsersByRole(role string) ([]models.User, error)
	GetUsersByEmail(email string) (*models.User, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(r repositories.UserRepository) UserService {
	return &userService{repo: r}
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAllUsers()
}

func (s *userService) GetUsersByEmail(email string) (*models.User, error) {
	return s.repo.GetByEmail(email)
}

func (s *userService) GetUsersByRole(role string) ([]models.User, error) {
	return s.repo.GetByRole(role)
}

func (s *userService) CreateUser(user *models.User) error {
	existing, _ := s.repo.GetByEmail(user.Email)
	if existing != nil && existing.Email != "" {
		return errors.New("email sudah terdaftar")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("gagal meng-hash password")
	}

	user.Password = string(hashed)

	if user.Role == "" {
		user.Role = "member"
	}

	return s.repo.CreateUser(user)
}
