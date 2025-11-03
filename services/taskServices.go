package services

import (
	"project-management-backend/models"
)

type TaskService interface {
	CreateTask(task *models.Task, user *models.User) error
	GetAllTasks(projectID uint, user *models.User) ([]models.Task, error)
	GetByID(taskID uint, user *models.User) (*models.Task, error)
	AddMember(taskID uint, userID uint, role string, currentUser *models.User) error
	GetMembers(taskID uint, user *models.User) ([]models.TaskUser, error)
}
