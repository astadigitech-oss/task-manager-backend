package repositories

import (
	"project-management-backend/config"
	"project-management-backend/models"
	"time"

	"gorm.io/gorm"
)

type WorkspaceRepository interface {
	CreateWorkspace(workspace *models.Workspace) error
	GetAllWorkspaces() ([]models.Workspace, error)
	UpdateWorkspace(workspace *models.Workspace) error
	SoftDeleteWorkspace(workspaceID uint) error
	DeleteWorkspace(workspaceID uint) error // Hard delete
	AddMember(wu *models.WorkspaceUser) error
	GetMembers(workspaceID uint) ([]models.WorkspaceUser, error)
	GetByID(workspaceID uint) (*models.Workspace, error)
	GetUserByID(userID uint) (*models.User, error)
	IsUserMember(workspaceID uint, userID uint) (bool, error)
	GetWorkspaceMember(workspaceID uint, userID uint) (*models.WorkspaceUser, error)
	GetWorkspacesByUserID(userID uint) ([]models.Workspace, error)
	RemoveMembers(workspaceID uint, userIDs []uint) error
	RemoveMember(workspaceID uint, userID uint) error
}

type workspaceRepository struct{}

func NewWorkspaceRepository() WorkspaceRepository {
	return &workspaceRepository{}
}

func (r *workspaceRepository) CreateWorkspace(workspace *models.Workspace) error {
	return config.DB.Create(workspace).Error
}

func (r *workspaceRepository) GetAllWorkspaces() ([]models.Workspace, error) {
	var workspaces []models.Workspace
	err := config.DB.
		Where("deleted_at IS NULL").
		Preload("Projects").
		Preload("Members").
		Find(&workspaces).Error
	return workspaces, err
}

func (r *workspaceRepository) GetWorkspacesByUserID(userID uint) ([]models.Workspace, error) {
	var workspaces []models.Workspace
	err := config.DB.
		Select("DISTINCT workspaces.*").
		Joins("JOIN workspace_users ON workspace_users.workspace_id = workspaces.id").
		Where("workspace_users.user_id = ? AND workspaces.deleted_at IS NULL", userID).
		Preload("Projects", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL")
		}).
		Preload("Members").
		Find(&workspaces).Error
	return workspaces, err
}

func (r *workspaceRepository) GetByID(workspaceID uint) (*models.Workspace, error) {
	var workspace models.Workspace
	err := config.DB.
		Where("id = ? AND deleted_at IS NULL", workspaceID).
		Preload("Members.User").
		Preload("Projects").
		First(&workspace, workspaceID).Error
	return &workspace, err
}

func (r *workspaceRepository) UpdateWorkspace(workspace *models.Workspace) error {
	return config.DB.Model(&models.Workspace{}).
		Where("id = ? AND deleted_at IS NULL", workspace.ID).
		Updates(map[string]interface{}{
			"name":        workspace.Name,
			"description": workspace.Description,
			"updated_at":  time.Now(),
		}).Error
}

func (r *workspaceRepository) SoftDeleteWorkspace(workspaceID uint) error {
	return config.DB.Model(&models.Workspace{}).
		Where("id = ?", workspaceID).
		Update("deleted_at", time.Now()).Error
}

func (r *workspaceRepository) DeleteWorkspace(workspaceID uint) error {
	return config.DB.Unscoped().Where("id = ?", workspaceID).Delete(&models.Workspace{}).Error
}

func (r *workspaceRepository) AddMember(wu *models.WorkspaceUser) error {
	return config.DB.Create(wu).Error
}

func (r *workspaceRepository) GetMembers(workspaceID uint) ([]models.WorkspaceUser, error) {
	var members []models.WorkspaceUser
	err := config.DB.Preload("User").Where("workspace_id = ?", workspaceID).Find(&members).Error
	return members, err
}

func (r *workspaceRepository) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	err := config.DB.First(&user, userID).Error
	return &user, err
}

func (r *workspaceRepository) IsUserMember(workspaceID uint, userID uint) (bool, error) {
	var count int64
	err := config.DB.Model(&models.WorkspaceUser{}).
		Where("workspace_id = ? AND user_id = ?", workspaceID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *workspaceRepository) GetWorkspaceMember(workspaceID uint, userID uint) (*models.WorkspaceUser, error) {
	var workspaceUser models.WorkspaceUser
	err := config.DB.Where("workspace_id = ? AND user_id = ?", workspaceID, userID).First(&workspaceUser).Error
	if err != nil {
		return nil, err
	}
	return &workspaceUser, nil
}

func (r *workspaceRepository) RemoveMember(workspaceID uint, userID uint) error {
	return config.DB.
		Where("workspace_id = ? AND user_id = ?", workspaceID, userID).
		Delete(&models.WorkspaceUser{}).Error
}

func (r *workspaceRepository) RemoveMembers(workspaceID uint, userIDs []uint) error {
	return config.DB.
		Where("workspace_id = ? AND user_id IN ?", workspaceID, userIDs).
		Delete(&models.WorkspaceUser{}).Error
}
