// repositories/project_repository.go
package repositories

import (
	"project-management-backend/config"
	"project-management-backend/models"
)

type ProjectRepository interface {
	CreateProject(project *models.Project) error
	GetAllProjects(workspaceID uint) ([]models.Project, error)
	GetByID(projectID uint) (*models.Project, error)
	AddMember(pu *models.ProjectUser) error
	GetMembers(projectID uint) ([]models.ProjectUser, error)
	GetUserByID(userID uint) (*models.User, error)
	IsUserMember(projectID uint, userID uint) (bool, error)
	IsUserInWorkspace(workspaceID uint, userID uint) (bool, error)
}

type projectRepository struct{}

func NewProjectRepository() ProjectRepository {
	return &projectRepository{}
}

func (r *projectRepository) CreateProject(project *models.Project) error {
	return config.DB.Create(project).Error
}

func (r *projectRepository) GetAllProjects(workspaceID uint) ([]models.Project, error) {
	var projects []models.Project
	err := config.DB.
		Preload("Members.User").
		Preload("Tasks").
		Preload("Images").
		Where("workspace_id = ?", workspaceID).
		Find(&projects).Error
	return projects, err
}

func (r *projectRepository) GetByID(projectID uint) (*models.Project, error) {
	var project models.Project
	err := config.DB.
		Preload("Members.User").
		Preload("Tasks").
		Preload("Images").
		Preload("Workspace").
		First(&project, projectID).Error
	return &project, err
}

func (r *projectRepository) AddMember(pu *models.ProjectUser) error {
	return config.DB.Create(pu).Error
}

func (r *projectRepository) GetMembers(projectID uint) ([]models.ProjectUser, error) {
	var members []models.ProjectUser
	err := config.DB.
		Preload("User").
		Where("project_id = ?", projectID).
		Find(&members).Error
	return members, err
}

func (r *projectRepository) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	err := config.DB.First(&user, userID).Error
	return &user, err
}

func (r *projectRepository) IsUserMember(projectID uint, userID uint) (bool, error) {
	var count int64
	err := config.DB.Model(&models.ProjectUser{}).
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *projectRepository) IsUserInWorkspace(workspaceID uint, userID uint) (bool, error) {
	var count int64
	err := config.DB.Model(&models.WorkspaceUser{}).
		Where("workspace_id = ? AND user_id = ?", workspaceID, userID).
		Count(&count).Error
	return count > 0, err
}
