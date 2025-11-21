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
	GetAllProjects(user *models.User) ([]models.Project, error)
	UpdateProject(project *models.Project, user *models.User) error
	SoftDeleteProject(projectID uint, user *models.User) error
	DeleteProject(projectID uint, user *models.User) error
	GetByID(projectID uint, user *models.User) (*models.Project, error)
	AddMember(projectID uint, userID uint, role string, currentUser *models.User) error // role string, bukan *string
	GetMembers(projectID uint, user *models.User) ([]models.ProjectUser, error)
}

type projectService struct {
	repo          repositories.ProjectRepository
	workspaceRepo repositories.WorkspaceRepository
}

func NewProjectService(repo repositories.ProjectRepository, workspaceRepo repositories.WorkspaceRepository) ProjectService {
	return &projectService{
		repo:          repo,
		workspaceRepo: workspaceRepo,
	}
}

func (s *projectService) CreateProject(project *models.Project, user *models.User) error {
	hasWorkspaceAccess, err := s.workspaceRepo.IsUserMember(project.WorkspaceID, user.ID)
	if err != nil || !hasWorkspaceAccess {
		return errors.New("hanya member workspace yang boleh buat project")
	}

	project.CreatedBy = user.ID

	if err := s.repo.CreateProject(project); err != nil {
		return err
	}

	creatorMember := &models.ProjectUser{
		ProjectID:     project.ID,
		UserID:        user.ID,
		RoleInProject: "admin",
	}

	if err := s.repo.AddMember(creatorMember); err != nil {
		s.repo.DeleteProject(project.ID)
		return errors.New("gagal menambahkan creator sebagai member project")
	}

	return nil
}

func (s *projectService) GetAllProjects(user *models.User) ([]models.Project, error) {
	if user.Role == "admin" {
		return s.repo.GetAllProjects()
	}

	// User regular hanya lihat projects yang dia punya akses
	return s.repo.GetProjectsByUserID(user.ID)
}

func (s *projectService) GetByID(projectID uint, user *models.User) (*models.Project, error) {
	project, err := s.repo.GetByID(projectID)
	if err != nil {
		return nil, errors.New("project tidak ditemukan")
	}

	if user.Role != "Admin" {
		isProjectMember, err := s.repo.IsUserMember(projectID, user.ID)
		if err != nil {
			return nil, errors.New("harus menjadi member project atau workspace untuk mengakses project ini")
		}
		isWorkspaceMember, err := s.workspaceRepo.IsUserMember(project.WorkspaceID, user.ID)
		if err != nil || !isWorkspaceMember {
			return nil, errors.New("tidak memiliki akses ke workspace project ini")
		}

		if !isProjectMember && !isWorkspaceMember {
			return nil, errors.New("akses ditolak untuk project ini")
		}
	}

	return project, nil
}

func (s *projectService) UpdateProject(project *models.Project, user *models.User) error {
	existingProject, err := s.repo.GetByID(project.ID)
	if err != nil {
		return errors.New("project tidak ditemukan")
	}

	isWorkspaceMember, err := s.workspaceRepo.IsUserMember(existingProject.WorkspaceID, user.ID)
	if err != nil || !isWorkspaceMember {
		return errors.New("hanya member workspace yang boleh update project")
	}

	return s.repo.UpdateProject(project)
}

func (s *projectService) SoftDeleteProject(projectID uint, user *models.User) error {
	project, err := s.repo.GetByID(projectID)
	if err != nil {
		return errors.New("project tidak ditemukan")
	}

	isWorkspaceMember, err := s.workspaceRepo.IsUserMember(project.WorkspaceID, user.ID)
	if err != nil || !isWorkspaceMember {
		return errors.New("hanya member workspace yang boleh soft delete project")
	}

	return s.repo.SoftDeleteProject(projectID)
}

func (s *projectService) DeleteProject(projectID uint, user *models.User) error {
	var project models.Project
	err := config.DB.Unscoped().Where("id = ?", projectID).First(&project).Error
	if err != nil {
		return errors.New("project tidak ditemukan")
	}

	isWorkspaceMember, err := s.workspaceRepo.IsUserMember(project.WorkspaceID, user.ID)
	if err != nil || !isWorkspaceMember {
		return errors.New("hanya member workspace yang boleh hard delete project")
	}

	return s.repo.DeleteProject(projectID)
}

func (s *projectService) AddMember(projectID uint, userID uint, role string, currentUser *models.User) error {
	project, err := s.repo.GetByID(projectID)
	if err != nil {
		return errors.New("project tidak ditemukan")
	}

	isWorkspaceMember, err := s.workspaceRepo.IsUserMember(project.WorkspaceID, currentUser.ID)
	if err != nil || !isWorkspaceMember {
		return errors.New("hanya member workspace yang boleh menambah member project")
	}

	isTargetUserInWorkspace, err := s.workspaceRepo.IsUserMember(project.WorkspaceID, userID)
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
	isWorkspaceMember, _ := s.workspaceRepo.IsUserMember(project.WorkspaceID, user.ID)

	if !isProjectMember && !isWorkspaceMember {
		return nil, errors.New("akses ditolak untuk melihat members project")
	}

	return s.repo.GetMembers(projectID)
}
