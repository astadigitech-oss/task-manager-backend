// services/project_service.go
package services

import (
	"errors"
	"project-management-backend/models"
	"project-management-backend/repositories"
)

type ProjectService interface {
	CreateProject(project *models.Project, user *models.User) error
	GetAllProjects(workspaceID uint, user *models.User) ([]models.Project, error)
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
