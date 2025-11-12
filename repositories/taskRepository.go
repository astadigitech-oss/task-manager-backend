package repositories

import (
	"project-management-backend/config"
	"project-management-backend/models"
	"time"

	"gorm.io/gorm"
)

type TaskRepository interface {
	CreateTask(task *models.Task) error
	GetAllTasks(projectID uint, workspaceID uint) ([]models.Task, error)
	GetByID(taskID uint, workspaceID uint) (*models.Task, error)
	UpdateTask(task *models.Task) error
	SoftDeleteTask(taskID uint) error
	DeleteTask(taskID uint) error // Hard delete
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
		Where("tasks.project_id = ? AND tasks.deleted_at IS NULL", projectID). // Specify table
		Preload("Members.User").
		Preload("Images").
		Preload("Project", func(db *gorm.DB) *gorm.DB {
			return db.Where("projects.deleted_at IS NULL") // Filter project yang tidak di soft delete
		}).
		Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) GetByID(taskID uint, workspaceID uint) (*models.Task, error) {
	var task models.Task
	err := config.DB.
		Where("tasks.id = ? AND tasks.deleted_at IS NULL", taskID). // Specify table
		Preload("Members.User").
		Preload("Images").
		Preload("Project", func(db *gorm.DB) *gorm.DB {
			return db.Where("projects.deleted_at IS NULL") // Filter project yang tidak di soft delete
		}).
		Preload("Project.Workspace", func(db *gorm.DB) *gorm.DB {
			return db.Where("workspaces.deleted_at IS NULL") // Filter workspace yang tidak di soft delete
		}).
		First(&task).Error
	return &task, err
}

func (r *taskRepository) UpdateTask(task *models.Task) error {
	return config.DB.Model(&models.Task{}).
		Where("id = ? AND deleted_at IS NULL", task.ID).
		Updates(map[string]interface{}{
			"title":       task.Title,
			"description": task.Description,
			"status":      task.Status,
			"priority":    task.Priority,
			"start_date":  task.StartDate,
			"due_date":    task.DueDate,
			"notes":       task.Notes,
			"updated_at":  time.Now(),
		}).Error
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
