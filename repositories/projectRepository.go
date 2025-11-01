package repositories

import (
	"project-management-backend/config"
	"project-management-backend/models"
)

type ProjectRepository interface {
	CreateProject(project *models.Project) error
	GetAllProject() ([]models.Project, error)
	AddMember(pu *models.ProjectUser) error
	GetMembers(projectID uint) ([]models.ProjectUser, error)
	GetByID(projectID uint) (*models.Project, error)
}

type projectRepository struct{}

func NewProjectRepository() ProjectRepository {
	return &projectRepository{}
}

func (r *projectRepository) CreateProject(project *models.Project) error {
	return config.DB.Create(project).Error
}

func (r *projectRepository) GetAllProject() ([]models.Project, error) {
	var projects []models.Project
	err := config.DB.Find(&projects).Error
	return projects, err
}

func (r *projectRepository) AddMember(pu *models.ProjectUser) error {
	return config.DB.Create(pu).Error
}

func (r *projectRepository) GetMembers(projectID uint) ([]models.ProjectUser, error) {
	var members []models.ProjectUser
	err := config.DB.Preload("User").Where("project_id = ?", projectID).Find(&members).Error
	return members, err
}

func (r *projectRepository) GetByID(projectID uint) (*models.Project, error) {
	var project models.Project
	err := config.DB.Preload("Members.User").First(&project, projectID).Error
	return &project, err
}
