package services

import (
	"errors"
	"project-management-backend/models"
	"project-management-backend/repositories"
)

type UserService interface {
	GetAllUsers(currentUser *models.User) ([]models.User, error)
	GetAllUsersPaginated(page, limit int, currentUser *models.User) (PaginatedUsersResponse, error)
	GetAllUsersWithFilters(filters UserFilters, currentUser *models.User) (PaginatedUsersResponse, error)
	GetUserByID(userID uint, currentUser *models.User) (*models.User, error)
	SearchUsers(query string, currentUser *models.User) ([]models.User, error)
	GetUsersWithDetails(userIDs []uint, currentUser *models.User) ([]models.User, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

type UserFilters struct {
	Search   string `form:"search"`
	Role     string `form:"role"`
	Position string `form:"position"`
	Page     int    `form:"page" binding:"min=1"`
	Limit    int    `form:"limit" binding:"min=1,max=100"`
}

type PaginatedUsersResponse struct {
	Users      []models.User `json:"users"`
	Total      int64         `json:"total"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalPages int           `json:"total_pages"`
}

func (s *userService) GetAllUsers(currentUser *models.User) ([]models.User, error) {
	// Hanya admin yang bisa melihat semua users
	if currentUser.Role != "admin" {
		return nil, errors.New("hanya admin yang bisa mengakses semua user")
	}

	return s.repo.GetAllUsers()
}

func (s *userService) GetAllUsersPaginated(page, limit int, currentUser *models.User) (PaginatedUsersResponse, error) {
	if currentUser.Role != "admin" {
		return PaginatedUsersResponse{}, errors.New("hanya admin yang bisa mengakses semua user")
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20 // default limit
	}

	users, total, err := s.repo.GetAllUsersPaginated(page, limit)
	if err != nil {
		return PaginatedUsersResponse{}, err
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return PaginatedUsersResponse{
		Users:      users,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *userService) GetAllUsersWithFilters(filters UserFilters, currentUser *models.User) (PaginatedUsersResponse, error) {
	if currentUser.Role != "admin" {
		return PaginatedUsersResponse{}, errors.New("hanya admin yang bisa mengakses semua user")
	}

	// Set default values
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.Limit < 1 || filters.Limit > 100 {
		filters.Limit = 20
	}

	// Convert to repository filters
	repoFilters := repositories.UserFilters{
		Search:   filters.Search,
		Role:     filters.Role,
		Position: filters.Position,
		Page:     filters.Page,
		Limit:    filters.Limit,
	}

	users, total, err := s.repo.GetAllUsersWithFilters(repoFilters)
	if err != nil {
		return PaginatedUsersResponse{}, err
	}

	totalPages := int(total) / filters.Limit
	if int(total)%filters.Limit > 0 {
		totalPages++
	}

	return PaginatedUsersResponse{
		Users:      users,
		Total:      total,
		Page:       filters.Page,
		Limit:      filters.Limit,
		TotalPages: totalPages,
	}, nil
}

func (s *userService) GetUserByID(userID uint, currentUser *models.User) (*models.User, error) {
	// Admin bisa lihat semua, user biasa hanya bisa lihat dirinya sendiri
	if currentUser.Role != "admin" && currentUser.ID != userID {
		return nil, errors.New("hanya bisa melihat profil sendiri")
	}

	return s.repo.GetUserByID(userID)
}

func (s *userService) SearchUsers(query string, currentUser *models.User) ([]models.User, error) {
	// Hanya admin yang bisa search semua users
	if currentUser.Role != "admin" {
		return nil, errors.New("hanya admin yang bisa mencari semua user")
	}

	if len(query) < 2 {
		return nil, errors.New("minimal 2 karakter untuk pencarian")
	}

	return s.repo.SearchUsers(query)
}

func (s *userService) GetUsersWithDetails(userIDs []uint, currentUser *models.User) ([]models.User, error) {
	// Hanya admin yang bisa get multiple users with details
	if currentUser.Role != "admin" {
		return nil, errors.New("hanya admin yang bisa mengakses detail multiple users")
	}

	return s.repo.GetUsersWithDetails(userIDs)
}
