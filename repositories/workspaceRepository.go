package repositories

import (
	"project-management-backend/config"
	"project-management-backend/models"
)

type WorkspaceRepository struct{}

func NewWorkspaceRepository() WorkspaceRepository {
	return WorkspaceRepository{}
}

func (r WorkspaceRepository) FindByUserID(userID uint) ([]models.Workspace, error) {
	var workspaces []models.Workspace
	err := config.DB.
		Where("created_by = ?", userID).
		Preload("Creator").
		Preload("Projects").
		Find(&workspaces).Error
	return workspaces, err
}

func (r WorkspaceRepository) Create(workspace *models.Workspace) error {
	return config.DB.Create(workspace).Error
}
