package services

import (
	"project-management-backend/models"
	"project-management-backend/repositories"
)

type DashboardService struct {
	repo repositories.TaskRepository
}

func NewDashboardService(repo repositories.TaskRepository) *DashboardService {
	return &DashboardService{repo: repo}
}

func (s *DashboardService) GetAllTasksForUser(userID uint) ([]models.Task, error) {
	return s.repo.GetAllTasksByUserID(userID)
}
