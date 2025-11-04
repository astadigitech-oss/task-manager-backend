package repositories

import (
	"project-management-backend/config"
	"project-management-backend/models"
)

type TaskRepository interface {
	CreateTask(task *models.Task) error
	GetAllTasks(projectID uint, workspaceID uint) ([]models.Task, error)
	GetByID(taskID uint, workspaceID uint) (*models.Task, error)
	AddMember(tu *models.TaskUser) error
	GetMembers(taskID uint) ([]models.TaskUser, error)
	GetUserByID(userID uint) (*models.User, error)
	IsProjectInWorkspace(projectID uint, workspaceID uint) (bool, error)
	IsUserMember(taskID uint, userID uint) (bool, error)
	IsUserInProject(projectID uint, userID uint) (bool, error)
}

type taskRepository struct{}

func NewTaskRepository() TaskRepository {
	return &taskRepository{}
}

func (r *taskRepository) CreateTask(task *models.Task) error {
	return config.DB.Create(task).Error
}

func (r *taskRepository) GetAllTasks(projectID uint, workspaceID uint) ([]models.Task, error) {
	var tasks []models.Task
	err := config.DB.
		Preload("Members.User").
		Preload("Images").
		Joins("JOIN projects ON projects.id = tasks.project_id").
		Where("tasks.project_id = ? AND projects.workspace_id = ?", projectID, workspaceID).
		Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) GetByID(taskID uint, workspaceID uint) (*models.Task, error) {
	var task models.Task
	err := config.DB.
		Preload("Members.User").
		Preload("Images").
		Preload("Project").
		Joins("JOIN projects ON projects.id = tasks.project_id").
		Where("tasks.id = ? AND projects.workspace_id = ?", taskID, workspaceID).
		First(&task).Error
	return &task, err
}

func (r *taskRepository) AddMember(tu *models.TaskUser) error {
	return config.DB.Create(tu).Error
}

func (r *taskRepository) GetMembers(taskID uint) ([]models.TaskUser, error) {
	var members []models.TaskUser
	err := config.DB.
		Preload("User").
		Where("task_id = ?", taskID).
		Find(&members).Error
	return members, err
}

func (r *taskRepository) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	err := config.DB.First(&user, userID).Error
	return &user, err
}

func (r *taskRepository) IsProjectInWorkspace(projectID uint, workspaceID uint) (bool, error) {
	var count int64
	err := config.DB.Model(&models.Project{}).
		Where("id = ? AND workspace_id = ?", projectID, workspaceID).
		Count(&count).Error
	return count > 0, err
}

func (r *taskRepository) IsUserMember(taskID uint, userID uint) (bool, error) {
	var count int64
	err := config.DB.Model(&models.TaskUser{}).
		Where("task_id = ? AND user_id = ?", taskID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *taskRepository) IsUserInProject(projectID uint, userID uint) (bool, error) {
	var count int64
	err := config.DB.Model(&models.ProjectUser{}).
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Count(&count).Error
	return count > 0, err
}
