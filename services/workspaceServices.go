package services

import (
	"errors"
	"project-management-backend/config"
	"project-management-backend/models"
	"project-management-backend/repositories"
)

type WorkspaceService interface {
	CreateWorkspace(workspace *models.Workspace, user *models.User) error
	GetAllWorkspaces(user *models.User) ([]models.Workspace, error)
	UpdateWorkspace(workspace *models.Workspace, user *models.User) error
	SoftDeleteWorkspace(workspaceID uint, user *models.User) error
	DeleteWorkspace(workspaceID uint, user *models.User) error
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

func (s *workspaceService) GetByID(workspaceID uint, user *models.User) (*models.Workspace, error) {
	if user.Role != "admin" {
		return nil, errors.New("hanya admin yang boleh melihat detail workspace")
	}
	return s.repo.GetByID(workspaceID)
}

func (s *workspaceService) UpdateWorkspace(workspace *models.Workspace, user *models.User) error {
	if user.Role != "admin" {
		return errors.New("hanya admin yang boleh update workspace")
	}

	existingWorkspace, err := s.repo.GetByID(workspace.ID)
	if err != nil {
		return errors.New("workspace tidak ditemukan")
	}

	if existingWorkspace.CreatedBy != user.ID {
		return errors.New("hanya creator workspace yang boleh mengupdate")
	}

	return s.repo.UpdateWorkspace(workspace)
}

func (s *workspaceService) SoftDeleteWorkspace(workspaceID uint, user *models.User) error {
	if user.Role != "admin" {
		return errors.New("hanya admin yang boleh soft delete workspace")
	}

	workspace, err := s.repo.GetByID(workspaceID)
	if err != nil {
		return errors.New("workspace tidak ditemukan")
	}

	if workspace.CreatedBy != user.ID {
		return errors.New("hanya creator workspace yang boleh soft delete")
	}

	return s.repo.SoftDeleteWorkspace(workspaceID)
}

func (s *workspaceService) DeleteWorkspace(workspaceID uint, user *models.User) error {
	if user.Role != "admin" {
		return errors.New("hanya admin yang boleh hard delete workspace")
	}

	var workspace models.Workspace
	err := config.DB.Unscoped().Where("id = ?", workspaceID).First(&workspace).Error
	if err != nil {
		return errors.New("workspace tidak ditemukan")
	}

	if workspace.CreatedBy != user.ID {
		return errors.New("hanya creator workspace yang boleh hard delete")
	}

	return s.repo.DeleteWorkspace(workspaceID)
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
