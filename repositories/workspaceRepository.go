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
	err := config.DB.Find(&workspaces).Error
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
	err := config.DB.Preload("Members.User").Preload("Projects").First(&workspace, workspaceID).Error
	return &workspace, err
}
