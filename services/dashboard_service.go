package services

import (
	"project-management-backend/models"
	"project-management-backend/repositories"
	"time"
)

type DashboardService struct {
	repo repositories.TaskRepository
}

func NewDashboardService(repo repositories.TaskRepository) *DashboardService {
	return &DashboardService{repo: repo}
}

func (s *DashboardService) GetAllTasksForUser(userID uint) ([]models.Task, error) {
	tasks, err := s.repo.GetAllTasksByUserID(userID)
	if err != nil {
		return nil, err
	}
	var filteredTasks []models.Task
	now := time.Now()
	threeDaysFromNow := now.AddDate(0, 0, 3)

	for _, task := range tasks {
		if task.Status != "Done" && task.DueDate.Before(threeDaysFromNow) || task.DueDate.Equal(now) {
			filteredTasks = append(filteredTasks, task)
		}
	}

	return filteredTasks, nil
}
