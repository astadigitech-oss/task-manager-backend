package services

import (
	"errors"
	"project-management-backend/models"
	"project-management-backend/repositories"
)

type WorkspaceService interface {
	CreateWorkspace(workspace *models.Workspace, user *models.User) error
	GetAllWorkspaces(user *models.User) ([]models.Workspace, error)
	AddMember(workspaceID uint, userID uint, role *string) error
	GetMembers(workspaceID uint) ([]models.WorkspaceUser, error)
	GetByID(workspaceID uint) (*models.Workspace, error)
}

type workspaceService struct {
	repo repositories.WorkspaceRepository
}

func NewWorkspaceService(r repositories.WorkspaceRepository) WorkspaceService {
	return &workspaceService{repo: r}
}

func (s *workspaceService) CreateWorkspace(workspace *models.Workspace, user *models.User) error {
	if user.Role != "admin" {
		return errors.New("hanya admin yang boleh buat workspace")
	}
	workspace.CreatedBy = user.ID
	return s.repo.CreateWorkspace(workspace)
}

func (s *workspaceService) GetAllWorkspaces(user *models.User) ([]models.Workspace, error) {
	if user.Role != "admin" {
		return nil, errors.New("hanya admin yang boleh lihat semua workspace")
	}
	return s.repo.GetAllWorkspaces()
}

func (s *workspaceService) AddMember(workspaceID uint, userID uint, role *string) error {
	member := &models.WorkspaceUser{
		WorkspaceID:     workspaceID,
		UserID:          userID,
		RoleInWorkspace: role,
	}
	return s.repo.AddMember(member)
}

func (s *workspaceService) GetMembers(workspaceID uint) ([]models.WorkspaceUser, error) {
	return s.repo.GetMembers(workspaceID)
}

func (s *workspaceService) GetByID(workspaceID uint) (*models.Workspace, error) {
	return s.repo.GetByID(workspaceID)
}
