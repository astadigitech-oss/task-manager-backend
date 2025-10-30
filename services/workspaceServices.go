package services

import (
	"errors"
	"project-management-backend/models"
	"project-management-backend/repositories"
	"project-management-backend/utils"
)

type WorkspaceService struct {
	Repo repositories.WorkspaceRepository
}

func NewWorkspaceService(repo repositories.WorkspaceRepository) WorkspaceService {
	return WorkspaceService{Repo: repo}
}

func (s WorkspaceService) GetUserWorkspaces(userID uint) ([]models.Workspace, error) {
	workspaces, err := s.Repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	return workspaces, nil
}

func (s WorkspaceService) CreateWorkspace(userID uint, role string, name string) (*models.Workspace, error) {
	if role != "admin" {
		return nil, errors.New("Hanya admin yang boleh membuat workspace")
	}

	workspace := models.Workspace{
		Name:      name,
		CreatedBy: userID,
	}

	if err := s.Repo.Create(&workspace); err != nil {
		utils.Error(userID, "CREATE_WORKSPACE", "workspaces", 400, err.Error(), "Failed to create workspace")
		return nil, err
	}
	return &workspace, nil
}
