package repositories

import (
	"project-management-backend/config"
	"project-management-backend/models"
	"time"
)

type UserRepository interface {
	GetByID(userID uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	GetAllUsers() ([]models.User, error)
	UpdateUser(user *models.User) error
	GetAllUsersPaginated(page, limit int) ([]models.User, int64, error)
	GetAllUsersWithFilters(filters UserFilters) ([]models.User, int64, error)
	GetUserByID(userID uint) (*models.User, error)
	SearchUsers(query string) ([]models.User, error)
	GetUsersWithDetails(userIDs []uint) ([]models.User, error)
	UpdateUserOnlineStatus(userID uint, isOnline bool) error
	GetOnlineUsers() ([]models.User, error)
	GetOnlineWorkspaceMembers(workspaceID uint) ([]models.User, error)
	IsUserMemberOfWorkspace(userID uint, workspaceID uint) (bool, error)
	DeleteUser(userID uint) error
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

type UserFilters struct {
	Search   string
	Role     string
	Position string
	Page     int
	Limit    int
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

func (r *userRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := config.DB.
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&users).Error
	return users, err
}

func (r *userRepository) GetAllUsersPaginated(page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	offset := (page - 1) * limit

	// Count total
	err := config.DB.Model(&models.User{}).
		Where("deleted_at IS NULL").
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated data
	err = config.DB.
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&users).Error

	return users, total, err
}

func (r *userRepository) GetAllUsersWithFilters(filters UserFilters) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := config.DB.Model(&models.User{}).
		Where("deleted_at IS NULL")

	// Apply filters
	if filters.Search != "" {
		search := "%" + filters.Search + "%"
		query = query.Where("name LIKE ? OR email LIKE ?", search, search)
	}

	if filters.Role != "" {
		query = query.Where("role = ?", filters.Role)
	}

	if filters.Position != "" {
		query = query.Where("position LIKE ?", "%"+filters.Position+"%")
	}

	// Count total with filters
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (filters.Page - 1) * filters.Limit
	if filters.Limit > 0 {
		query = query.Offset(offset).Limit(filters.Limit)
	}

	// Get data
	err = query.Order("created_at DESC").Find(&users).Error
	return users, total, err
}

func (r *userRepository) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	err := config.DB.
		Where("id = ? AND deleted_at IS NULL", userID).
		First(&user).Error
	return &user, err
}

func (r *userRepository) UpdateUser(user *models.User) error {
	return config.DB.Save(user).Error
}

func (r *userRepository) SearchUsers(queryStr string) ([]models.User, error) {
	var users []models.User

	search := "%" + queryStr + "%"
	err := config.DB.
		Where("deleted_at IS NULL AND (name LIKE ? OR email LIKE ?)", search, search).
		Order("name ASC").
		Limit(20). // Limit untuk autocomplete
		Find(&users).Error

	return users, err
}

func (r *userRepository) GetUsersWithDetails(userIDs []uint) ([]models.User, error) {
	var users []models.User

	err := config.DB.
		Where("id IN ? AND deleted_at IS NULL", userIDs).
		Order("name ASC").
		Find(&users).Error

	return users, err
}

func (r *userRepository) UpdateUserOnlineStatus(userID uint, isOnline bool) error {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return err
	}

	user.IsOnline = isOnline
	now := time.Now()
	if !isOnline {
		user.LastSeen = &now
	}

	return config.DB.Save(&user).Error
}

func (r *userRepository) GetOnlineUsers() ([]models.User, error) {
	var users []models.User
	err := config.DB.
		Preload("Workspaces").
		Preload("Projects").
		Preload("Tasks.Task").
		Where("is_online = ?", true).
		Find(&users).Error
	return users, err
}

func (r *userRepository) GetOnlineWorkspaceMembers(workspaceID uint) ([]models.User, error) {
	var memberIDs []uint
	err := config.DB.Model(&models.WorkspaceUser{}).Where("workspace_id = ?", workspaceID).Pluck("user_id", &memberIDs).Error
	if err != nil {
		return nil, err
	}

	var onlineUsers []models.User
	if len(memberIDs) > 0 {
		err = config.DB.Where("id IN ? AND is_online = ?", memberIDs, true).Find(&onlineUsers).Error
	}

	return onlineUsers, err
}

func (r *userRepository) IsUserMemberOfWorkspace(userID uint, workspaceID uint) (bool, error) {
	var membership models.WorkspaceUser
	err := config.DB.Where("user_id = ? AND workspace_id = ?", userID, workspaceID).First(&membership).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *userRepository) DeleteUser(userID uint) error {
	return config.DB.Delete(&models.User{}, userID).Error
}
