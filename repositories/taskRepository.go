package repositories

import (
	"project-management-backend/config"
	"project-management-backend/models"
	"time"

	"gorm.io/gorm"
)

type TaskRepository interface {
	CreateTask(task *models.Task) error
	GetAllTasks(projectID uint) ([]models.Task, error)
	GetByID(taskID uint) (*models.Task, error)
	UpdateTask(taskID uint, updates map[string]interface{}) error
	SoftDeleteTask(taskID uint) error
	SoftDeleteAllTasksInProject(projectID uint) error
	DeleteTask(taskID uint) error // Hard delete
	AddMember(tu *models.TaskUser) error
	GetMembers(taskID uint) ([]models.TaskUser, error)
	DeleteMember(taskID uint, userID uint) error
	GetUserByID(userID uint) (*models.User, error)
	IsProjectInWorkspace(projectID uint, workspaceID uint) (bool, error)
	IsUserMember(taskID uint, userID uint) (bool, error)
	IsUserInProject(projectID uint, userID uint) (bool, error)
	GetTasksByUserID(projectID uint, userID uint) ([]models.Task, error)
	GetAllTasksByUserID(userID uint) ([]models.Task, error)
}

type taskRepository struct{}

func NewTaskRepository() TaskRepository {
	return &taskRepository{}
}

func (r *taskRepository) CreateTask(task *models.Task) error {
	return config.DB.Create(task).Error
}

func (r *taskRepository) GetAllTasks(projectID uint) ([]models.Task, error) {
	var tasks []models.Task
	err := config.DB.
		Joins("JOIN projects ON projects.id = tasks.project_id").
		Joins("JOIN workspaces ON workspaces.id = projects.workspace_id").
		Where("tasks.project_id = ? AND tasks.deleted_at IS NULL AND projects.deleted_at IS NULL AND workspaces.deleted_at IS NULL", projectID).
		Preload("Members.User").
		Preload("Images").
		Preload("Project").
		Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) GetAllTasksByUserID(userID uint) ([]models.Task, error) {
	var tasks []models.Task
	err := config.DB.
		Joins("JOIN task_users ON task_users.task_id = tasks.id").
		Joins("JOIN projects ON projects.id = tasks.project_id").
		Joins("JOIN workspaces ON workspaces.id = projects.workspace_id").
		Where("task_users.user_id = ? AND tasks.deleted_at IS NULL AND projects.deleted_at IS NULL AND workspaces.deleted_at IS NULL", userID).
		Preload("Members.User").
		Preload("Images").
		Preload("Project").
		Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) GetTasksByUserID(projectID uint, userID uint) ([]models.Task, error) {
	var tasks []models.Task
	err := config.DB.
		Joins("JOIN projects ON projects.id = tasks.project_id").
		Joins("JOIN workspaces ON workspaces.id = projects.workspace_id").
		Where("tasks.project_id = ? AND tasks.deleted_at IS NULL AND projects.deleted_at IS NULL AND workspaces.deleted_at IS NULL", projectID).
		Where("tasks.id IN (SELECT task_id FROM task_users WHERE user_id = ?)", userID).
		Preload("Members.User").
		Preload("Images").
		Preload("Project").
		Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) SoftDeleteAllTasksInProject(projectID uint) error {
	return config.DB.Model(&models.Task{}).
		Where("project_id = ?", projectID).
		Update("deleted_at", time.Now()).Error
}

func (r *taskRepository) GetProjectByID(projectID uint) (*models.Project, error) {
	var project models.Project
	err := config.DB.First(&project, projectID).Error
	return &project, err
}

func (r *taskRepository) GetByID(taskID uint) (*models.Task, error) {
	var task models.Task
	err := config.DB.
		Joins("JOIN projects ON projects.id = tasks.project_id").
		Joins("JOIN workspaces ON workspaces.id = projects.workspace_id").
		Where("tasks.id = ? AND tasks.deleted_at IS NULL AND projects.deleted_at IS NULL AND workspaces.deleted_at IS NULL", taskID).
		Preload("Members.User").
		Preload("Images").
		Preload("Project").
		Preload("Project.Workspace").
		First(&task).Error
	return &task, err
}

func (r *taskRepository) UpdateTask(taskID uint, updates map[string]interface{}) error {
	return config.DB.Model(&models.Task{}).
		Where("id = ? AND deleted_at IS NULL", taskID).
		Updates(updates).Error
}
func (r *taskRepository) SoftDeleteTask(taskID uint) error {
	return config.DB.Model(&models.Task{}).
		Where("id = ?", taskID).
		Update("deleted_at", time.Now()).Error
}

func (r *taskRepository) DeleteTask(taskID uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Delete task_users (members)
		if err := tx.Where("task_id = ?", taskID).Delete(&models.TaskUser{}).Error; err != nil {
			return err
		}

		// 2. Delete task_images
		if err := tx.Where("task_id = ?", taskID).Delete(&models.TaskImage{}).Error; err != nil {
			return err
		}

		// 3. HARD DELETE task
		if err := tx.Unscoped().Where("id = ?", taskID).Delete(&models.Task{}).Error; err != nil {
			return err
		}

		return nil
	})
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

func (r *taskRepository) DeleteMember(taskID uint, userID uint) error {
	return config.DB.Where("task_id = ? AND user_id = ?", taskID, userID).Delete(&models.TaskUser{}).Error
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
