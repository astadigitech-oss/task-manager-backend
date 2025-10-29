package services

import (
	"errors"
	"project-management-backend/models"
	"project-management-backend/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetAllUsers() ([]models.User, error)
	CreateUser(input models.User) error
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

func (s *userService) CreateUser(input models.User) error {
	// Cek email sudah dipakai belum
	existing, _ := s.repo.GetByEmail(input.Email)
	if existing != nil && existing.Email != "" {
		return errors.New("email sudah terdaftar")
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("gagal meng-hash password")
	}

	// Set default role jika kosong
	if input.Role == "" {
		input.Role = "member"
	}

	input.Password = string(hashed)

	// Simpan ke DB
	return s.repo.CreateUser(&input)
}
