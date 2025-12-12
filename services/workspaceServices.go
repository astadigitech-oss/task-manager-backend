package services

import (
	"errors"
	"fmt"
	"project-management-backend/config"
	"project-management-backend/models"
	"project-management-backend/repositories"
)

type WorkspaceMember struct {
	UserID uint    `json:"user_id"`
	Role   *string `json:"role_in_workspace"`
}
type WorkspaceService interface {
	CreateWorkspace(workspace *models.Workspace, user *models.User) error
	GetAllWorkspaces(user *models.User) ([]models.Workspace, error)
	UpdateWorkspace(workspace *models.Workspace, user *models.User) error
	SoftDeleteWorkspace(workspaceID uint, user *models.User) error
	DeleteWorkspace(workspaceID uint, user *models.User) error
	GetMembers(workspaceID uint, user *models.User) ([]models.WorkspaceUser, error)
	GetByID(workspaceID uint, user *models.User) (*models.Workspace, error)
	AddMembers(workspaceID uint, members []WorkspaceMember, currentUser *models.User) error
	RemoveMembers(workspaceID uint, userIDs []uint, currentUser *models.User) error
	RemoveMember(workspaceID uint, userID uint, currentUser *models.User) error
}

type workspaceService struct {
	repo        repositories.WorkspaceRepository
	projectRepo repositories.ProjectRepository
	taskRepo    repositories.TaskRepository
}

func NewWorkspaceService(r repositories.WorkspaceRepository, p repositories.ProjectRepository, t repositories.TaskRepository) WorkspaceService {
	return &workspaceService{repo: r, projectRepo: p, taskRepo: t}
}

func (s *workspaceService) CreateWorkspace(workspace *models.Workspace, user *models.User) error {
	if user.Role != "admin" {
		return errors.New("hanya admin yang boleh buat workspace")
	}
	workspace.CreatedBy = user.ID

	if err := s.repo.CreateWorkspace(workspace); err != nil {
		return err
	}

	ownerRole := "owner"
	creatorMember := &models.WorkspaceUser{
		WorkspaceID:     workspace.ID,
		UserID:          user.ID,
		RoleInWorkspace: &ownerRole,
	}

	if err := s.repo.AddMember(creatorMember); err != nil {
		s.repo.DeleteWorkspace(workspace.ID)
		return errors.New("gagal menambahkan creator sebagai member workspace")
	}

	return nil
}

func (s *workspaceService) GetAllWorkspaces(user *models.User) ([]models.Workspace, error) {
	if user.Role == "admin" {
		return s.repo.GetAllWorkspaces()
	}
	return s.repo.GetWorkspacesByUserID(user.ID)
}

func (s *workspaceService) GetByID(workspaceID uint, user *models.User) (*models.Workspace, error) {
	workspace, err := s.repo.GetByID(workspaceID)
	if err != nil {
		return nil, errors.New("workspace tidak ditemukan")
	}

	if user.Role != "admin" {
		isMember, err := s.repo.IsUserMember(workspaceID, user.ID)
		if err != nil || !isMember {
			return nil, errors.New("akses ditolak untuk workspace ini")
		}
	}

	return workspace, nil
}

func (s *workspaceService) UpdateWorkspace(workspace *models.Workspace, user *models.User) error {

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
	workspace, err := s.repo.GetByID(workspaceID)
	if err != nil {
		return errors.New("workspace tidak ditemukan")
	}

	if workspace.CreatedBy != user.ID {
		return errors.New("hanya creator workspace yang boleh soft delete")
	}

	projects, err := s.projectRepo.GetProjectsByWorkspace(workspaceID)
	if err != nil {
		return errors.New("gagal mengambil data projects")
	}

	for _, project := range projects {
		if err := s.projectRepo.SoftDeleteProject(project.ID); err != nil {
			return fmt.Errorf("gagal soft delete project %d: %w", project.ID, err)
		}

		if err := s.taskRepo.SoftDeleteAllTasksInProject(project.ID); err != nil {
			return fmt.Errorf("gagal soft delete tasks project %d: %w", project.ID, err)
		}
	}

	return s.repo.SoftDeleteWorkspace(workspaceID)
}

func (s *workspaceService) DeleteWorkspace(workspaceID uint, user *models.User) error {

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

func (s *workspaceService) AddMembers(workspaceID uint, members []WorkspaceMember, currentUser *models.User) error {
	_, err := s.repo.GetByID(workspaceID)
	if err != nil {
		return errors.New("workspace tidak ditemukan")
	}

	for _, member := range members {
		_, err := s.repo.GetUserByID(member.UserID)
		if err != nil {
			return fmt.Errorf("user %d tidak ditemukan", member.UserID)
		}

		isMember, err := s.repo.IsUserMember(workspaceID, member.UserID)
		if err != nil {
			return fmt.Errorf("gagal memvalidasi member %d", member.UserID)
		}
		if isMember {
			return fmt.Errorf("user %d sudah menjadi member di workspace ini", member.UserID)
		}

		workspaceMember := &models.WorkspaceUser{
			WorkspaceID:     workspaceID,
			UserID:          member.UserID,
			RoleInWorkspace: member.Role,
		}
		if err := s.repo.AddMember(workspaceMember); err != nil {
			return fmt.Errorf("gagal menambahkan user %d", member.UserID)
		}
	}

	return nil
}

func (s *workspaceService) GetMembers(workspaceID uint, user *models.User) ([]models.WorkspaceUser, error) {

	_, err := s.repo.GetByID(workspaceID)
	if err != nil {
		return nil, errors.New("workspace tidak ditemukan")
	}

	return s.repo.GetMembers(workspaceID)
}

func (s *workspaceService) RemoveMember(workspaceID uint, userID uint, currentUser *models.User) error {
	workspace, err := s.repo.GetByID(workspaceID)
	if err != nil {
		return errors.New("workspace tidak ditemukan")
	}

	targetMember, err := s.repo.GetWorkspaceMember(workspaceID, userID)
	if err != nil {
		return errors.New("member tidak ditemukan di workspace ini")
	}

	if currentUser.Role != "admin" {

		if workspace.CreatedBy != currentUser.ID {

			currentMember, err := s.repo.GetWorkspaceMember(workspaceID, currentUser.ID)
			if err != nil {
				return errors.New("anda bukan member workspace ini")
			}

			if *currentMember.RoleInWorkspace != "admin" {
				return errors.New("role Anda tidak cukup untuk menghapus member")
			}

			if targetMember.RoleInWorkspace != nil && *targetMember.RoleInWorkspace == "owner" {
				return errors.New("tidak bisa menghapus owner workspace")
			}
		}
	}

	if userID == currentUser.ID {
		return errors.New("tidak bisa menghapus diri sendiri")
	}

	return s.repo.RemoveMember(workspaceID, userID)
}

func (s *workspaceService) RemoveMembers(workspaceID uint, userIDs []uint, currentUser *models.User) error {
	workspace, err := s.repo.GetByID(workspaceID)
	if err != nil {
		return errors.New("workspace tidak ditemukan")
	}

	if currentUser.Role != "admin" {
		if workspace.CreatedBy != currentUser.ID {
			currentMember, err := s.repo.GetWorkspaceMember(workspaceID, currentUser.ID)
			if err != nil {
				return errors.New("anda bukan member workspace ini")
			}
			if *currentMember.RoleInWorkspace != "admin" {
				return errors.New("role Anda tidak cukup untuk menghapus member")
			}
		}
	}

	for _, userID := range userIDs {
		_, err := s.repo.GetWorkspaceMember(workspaceID, userID)
		if err != nil {
			return fmt.Errorf("user %d tidak ditemukan di workspace ini", userID)
		}

		if userID == currentUser.ID {
			return fmt.Errorf("tidak bisa menghapus diri sendiri (user_id: %d)", userID)
		}
	}

	return s.repo.RemoveMembers(workspaceID, userIDs)
}
