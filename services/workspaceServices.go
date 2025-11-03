package services

import (
	"errors"
	"project-management-backend/models"
	"project-management-backend/repositories"
)

type WorkspaceService interface {
	CreateWorkspace(workspace *models.Workspace, user *models.User) error
	GetAllWorkspaces(user *models.User) ([]models.Workspace, error)
	AddMember(workspaceID uint, userID uint, role *string, user *models.User) error
	GetMembers(workspaceID uint, user *models.User) ([]models.WorkspaceUser, error)
	GetByID(workspaceID uint, user *models.User) (*models.Workspace, error)
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

func (s *workspaceService) AddMember(workspaceID uint, userID uint, role *string, user *models.User) error {
	if user.Role != "admin" {
		return errors.New("hanya admin yang boleh menambah member")
	}

	_, err := s.repo.GetByID(workspaceID)
	if err != nil {
		return errors.New("workspace tidak ditemukan")
	}

	_, err = s.repo.GetUserByID(userID)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	isMember, err := s.repo.IsUserMember(workspaceID, userID)
	if err != nil {
		return errors.New("gagal memvalidasi member")
	}
	if isMember {
		return errors.New("user sudah menjadi member di workspace ini")
	}

	member := &models.WorkspaceUser{
		WorkspaceID:     workspaceID,
		UserID:          userID,
		RoleInWorkspace: role,
	}
	return s.repo.AddMember(member)
}

func (s *workspaceService) GetMembers(workspaceID uint, user *models.User) ([]models.WorkspaceUser, error) {
	if user.Role != "admin" {
		return nil, errors.New("hanya admin yang boleh melihat members")
	}

	_, err := s.repo.GetByID(workspaceID)
	if err != nil {
		return nil, errors.New("workspace tidak ditemukan")
	}

	return s.repo.GetMembers(workspaceID)
}

func (s *workspaceService) GetByID(workspaceID uint, user *models.User) (*models.Workspace, error) {
	if user.Role != "admin" {
		return nil, errors.New("hanya admin yang boleh melihat detail workspace")
	}
	return s.repo.GetByID(workspaceID)
}
