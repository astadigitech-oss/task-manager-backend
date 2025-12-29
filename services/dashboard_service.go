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

func (s *DashboardService) GetAllTaskForAdmin() ([]models.Task, error) {
	tasks, err := s.repo.GetAllTasksForAdmin()
	if err != nil {
		return nil, err
	}

	var overdueTasks []models.Task
	now := time.Now()

	for i := range tasks {
		task := &tasks[i]
		isOverdue := task.DueDate.Before(now)

		if isOverdue {
			if task.Status == "Done" && task.FinishedAt != nil && task.FinishedAt.After(task.DueDate) {
				task.OverdueDuration = task.FinishedAt.Sub(task.DueDate)
				overdueTasks = append(overdueTasks, *task)
			} else if task.Status != "Done" {
				task.OverdueDuration = now.Sub(task.DueDate)
				overdueTasks = append(overdueTasks, *task)
			}
		}
	}
	return overdueTasks, nil
}
