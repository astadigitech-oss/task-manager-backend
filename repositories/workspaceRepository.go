package repositories

import (
	"project-management-backend/config"
	"project-management-backend/models"
)

type WorkspaceRepository interface {
	CreateWorkspace(workspace *models.Workspace) error
	GetAllWorkspaces() ([]models.Workspace, error)
	AddMember(wu *models.WorkspaceUser) error
	GetMembers(workspaceID uint) ([]models.WorkspaceUser, error)
	GetByID(workspaceID uint) (*models.Workspace, error)
	GetUserByID(userID uint) (*models.User, error)
	IsUserMember(workspaceID uint, userID uint) (bool, error)
	GetWorkspaceMember(workspaceID uint, userID uint) (*models.WorkspaceUser, error)
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
		Preload("Projects"). // Tetap preload projects
		Preload("Members").
		Find(&workspaces).Error
	return workspaces, err
}

func (r *workspaceRepository) AddMember(wu *models.WorkspaceUser) error {
	return config.DB.Create(wu).Error
}

func (r *workspaceRepository) GetMembers(workspaceID uint) ([]models.WorkspaceUser, error) {
	var members []models.WorkspaceUser
	err := config.DB.Preload("User").Where("workspace_id = ?", workspaceID).Find(&members).Error
	return members, err
}

func (r *workspaceRepository) GetByID(workspaceID uint) (*models.Workspace, error) {
	var workspace models.Workspace
	err := config.DB.
		Preload("Members.User").
		Preload("Projects").
		First(&workspace, workspaceID).Error
	return &workspace, err
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
