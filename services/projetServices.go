// services/project_service.go
package services

import (
	"errors"
	"project-management-backend/config"
	"project-management-backend/models"
	"project-management-backend/repositories"
)

type ProjectService interface {
	CreateProject(project *models.Project, user *models.User) error
	GetAllProjects(workspaceID uint, user *models.User) ([]models.Project, error)
	UpdateProject(project *models.Project, user *models.User) error
	SoftDeleteProject(projectID uint, user *models.User) error
	DeleteProject(projectID uint, user *models.User) error
	GetByID(projectID uint, user *models.User) (*models.Project, error)
	AddMember(projectID uint, userID uint, role string, currentUser *models.User) error // role string, bukan *string
	GetMembers(projectID uint, user *models.User) ([]models.ProjectUser, error)
}

type projectService struct {
	repo repositories.ProjectRepository
}

func NewProjectService(repo repositories.ProjectRepository) ProjectService {
	return &projectService{repo: repo}
}

func (s *projectService) CreateProject(project *models.Project, user *models.User) error {
	isMember, err := s.repo.IsUserInWorkspace(project.WorkspaceID, user.ID)
	if err != nil || !isMember {
		return errors.New("hanya member workspace yang boleh buat project")
	}

	project.CreatedBy = user.ID
	return s.repo.CreateProject(project)
}

func (s *projectService) GetAllProjects(workspaceID uint, user *models.User) ([]models.Project, error) {
	isMember, err := s.repo.IsUserInWorkspace(workspaceID, user.ID)
	if err != nil || !isMember {
		return nil, errors.New("hanya member workspace yang boleh lihat projects")
	}

	return s.repo.GetAllProjects(workspaceID)
}

func (s *projectService) GetByID(projectID uint, user *models.User) (*models.Project, error) {
	project, err := s.repo.GetByID(projectID)
	if err != nil {
		return nil, errors.New("project tidak ditemukan")
	}

	isProjectMember, _ := s.repo.IsUserMember(projectID, user.ID)
	isWorkspaceMember, _ := s.repo.IsUserInWorkspace(project.WorkspaceID, user.ID)

	if !isProjectMember && !isWorkspaceMember {
		return nil, errors.New("akses ditolak untuk project ini")
	}

	return project, nil
}

func (s *projectService) UpdateProject(project *models.Project, user *models.User) error {
	// Cek apakah project exists
	existingProject, err := s.repo.GetByID(project.ID)
	if err != nil {
		return errors.New("project tidak ditemukan")
	}

	// Cek apakah user adalah member di workspace project
	isWorkspaceMember, err := s.repo.IsUserInWorkspace(existingProject.WorkspaceID, user.ID)
	if err != nil || !isWorkspaceMember {
		return errors.New("hanya member workspace yang boleh update project")
	}

	// Validasi: hanya creator yang bisa update
	if existingProject.CreatedBy != user.ID {
		return errors.New("hanya creator project yang boleh mengupdate")
	}

	return s.repo.UpdateProject(project)
}

func (s *projectService) SoftDeleteProject(projectID uint, user *models.User) error {
	// Cek apakah project exists
	project, err := s.repo.GetByID(projectID)
	if err != nil {
		return errors.New("project tidak ditemukan")
	}

	// Cek apakah user adalah member di workspace project
	isWorkspaceMember, err := s.repo.IsUserInWorkspace(project.WorkspaceID, user.ID)
	if err != nil || !isWorkspaceMember {
		return errors.New("hanya member workspace yang boleh soft delete project")
	}

	// Validasi: hanya creator yang bisa soft delete
	if project.CreatedBy != user.ID {
		return errors.New("hanya creator project yang boleh soft delete")
	}

	return s.repo.SoftDeleteProject(projectID)
}

func (s *projectService) DeleteProject(projectID uint, user *models.User) error {
	// Cek apakah project exists (termasuk yang sudah di soft delete)
	var project models.Project
	err := config.DB.Unscoped().Where("id = ?", projectID).First(&project).Error
	if err != nil {
		return errors.New("project tidak ditemukan")
	}

	// Cek apakah user adalah member di workspace project
	isWorkspaceMember, err := s.repo.IsUserInWorkspace(project.WorkspaceID, user.ID)
	if err != nil || !isWorkspaceMember {
		return errors.New("hanya member workspace yang boleh hard delete project")
	}

	// Validasi: hanya creator yang bisa hard delete
	if project.CreatedBy != user.ID {
		return errors.New("hanya creator project yang boleh hard delete")
	}

	return s.repo.DeleteProject(projectID)
}

func (s *projectService) AddMember(projectID uint, userID uint, role string, currentUser *models.User) error {
	project, err := s.repo.GetByID(projectID)
	if err != nil {
		return errors.New("project tidak ditemukan")
	}

	isWorkspaceMember, err := s.repo.IsUserInWorkspace(project.WorkspaceID, currentUser.ID)
	if err != nil || !isWorkspaceMember {
		return errors.New("hanya member workspace yang boleh menambah member project")
	}

	isTargetUserInWorkspace, err := s.repo.IsUserInWorkspace(project.WorkspaceID, userID)
	if err != nil || !isTargetUserInWorkspace {
		return errors.New("user harus menjadi member workspace terlebih dahulu")
	}

	isMember, err := s.repo.IsUserMember(projectID, userID)
	if err != nil {
		return errors.New("gagal memvalidasi member")
	}
	if isMember {
		return errors.New("user sudah menjadi member di project ini")
	}

	member := &models.ProjectUser{
		ProjectID:     projectID,
		UserID:        userID,
		RoleInProject: role,
	}
	return s.repo.AddMember(member)
}

func (s *projectService) GetMembers(projectID uint, user *models.User) ([]models.ProjectUser, error) {
	project, err := s.repo.GetByID(projectID)
	if err != nil {
		return nil, errors.New("project tidak ditemukan")
	}

	isProjectMember, _ := s.repo.IsUserMember(projectID, user.ID)
	isWorkspaceMember, _ := s.repo.IsUserInWorkspace(project.WorkspaceID, user.ID)

	if !isProjectMember && !isWorkspaceMember {
		return nil, errors.New("akses ditolak untuk melihat members project")
	}

	return s.repo.GetMembers(projectID)
}
