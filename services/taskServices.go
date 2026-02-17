package services

import (
	"errors"
	"fmt"
	"project-management-backend/config"
	"project-management-backend/models"
	"project-management-backend/repositories"
	"project-management-backend/utils"
	"strings"
	"time"

	"gorm.io/gorm"
)

type TaskService interface {
	CreateTask(task *models.Task, workspaceID uint, user *models.User) error
	GetAllTasks(projectID uint, workspaceID uint, user *models.User) ([]models.Task, error)
	GetByID(taskID uint, workspaceID uint, user *models.User) (*models.Task, error)
	UpdateTask(taskID uint, updates map[string]interface{}, workspaceID uint, user *models.User) error
	SoftDeleteTask(taskID uint, workspaceID uint, user *models.User) error
	DeleteTask(taskID uint, workspaceID uint, user *models.User) error
	AddMember(taskID uint, projectID uint, workspaceID uint, userID uint, role string, currentUser *models.User) error
	GetMembers(taskID uint, projectID uint, workspaceID uint, user *models.User) ([]models.TaskUser, error)
	DeleteMember(taskID uint, projectID uint, workspaceID uint, userID uint, currentUser *models.User) error
}
type taskService struct {
	repo            repositories.TaskRepository
	userRepo        repositories.UserRepository
	activityLogger  utils.ActivityLogger
	taskStatusLog   repositories.TaskStatusLogRepository
	telegramService TelegramService
}

func NewTaskService(repo repositories.TaskRepository, userRepo repositories.UserRepository, taskStatusLogRepo repositories.TaskStatusLogRepository, activityLogger utils.ActivityLogger, telegramService TelegramService) TaskService {
	return &taskService{
		repo:            repo,
		userRepo:        userRepo,
		activityLogger:  activityLogger,
		taskStatusLog:   taskStatusLogRepo,
		telegramService: telegramService,
	}

}

func (s *taskService) CreateTask(task *models.Task, workspaceID uint, user *models.User) error {
	isProjectInWorkspace, err := s.repo.IsProjectInWorkspace(task.ProjectID, workspaceID)
	if err != nil || !isProjectInWorkspace {
		return errors.New("project tidak ditemukan di workspace ini")
	}

	if task.Status == "" {
		task.Status = "On board"
	}
	if task.Priority == "" {
		task.Priority = "Normal"
	}

	err = s.repo.CreateTask(task)
	if err != nil {
		return err
	}

	activity := models.ActivityLog{
		UserID:    user.ID,
		Action:    fmt.Sprintf("User created task '%s' with status '%s'", task.Title, task.Status),
		TableName: "tasks",
		ItemID:    task.ID,
	}
	s.activityLogger.Log(activity)

	taskStatusLog := &models.TaskStatusLog{
		TaskID:   task.ID,
		Status:   task.Status,
		ClockIn:  task.StartDate,
		ClockOut: nil,
	}

	return s.taskStatusLog.Create(taskStatusLog)
}

func (s *taskService) GetAllTasks(projectID uint, workspaceID uint, user *models.User) ([]models.Task, error) {
	isProjectInWorkspace, err := s.repo.IsProjectInWorkspace(projectID, workspaceID)
	if err != nil || !isProjectInWorkspace {
		return nil, errors.New("project tidak ditemukan di workspace ini")
	}

	if user.Role == "admin" {
		return s.repo.GetAllTasks(projectID)
	}

	isProjectMember, err := s.repo.IsUserInProject(projectID, user.ID)
	if err != nil || !isProjectMember {
		return nil, errors.New("hanya member project yang boleh lihat tasks")
	}

	return s.repo.GetTasksByUserID(projectID, user.ID)
}

func (s *taskService) GetByID(taskID uint, workspaceID uint, user *models.User) (*models.Task, error) {
	task, err := s.repo.GetByID(taskID)
	if err != nil {
		return nil, errors.New("task tidak ditemukan")
	}

	if task.Project.WorkspaceID != workspaceID {
		return nil, errors.New("task tidak ditemukan di workspace ini")
	}

	if user.Role != "admin" {
		isTaskMember, _ := s.repo.IsUserMember(taskID, user.ID)
		isProjectMember, _ := s.repo.IsUserInProject(task.ProjectID, user.ID)

		if !isTaskMember && !isProjectMember {
			return nil, errors.New("akses ditolak untuk task ini")
		}
	}

	return task, nil
}

func (s *taskService) UpdateTask(taskID uint, updates map[string]interface{}, workspaceID uint, user *models.User) error {
	existingTask, err := s.repo.GetByID(taskID)
	if err != nil {
		return errors.New("task tidak ditemukan")
	}

	if existingTask.Project.WorkspaceID != workspaceID {
		return errors.New("task tidak ditemukan di workspace ini")
	}

	finalUpdates := updates
	isProjectAdminOrAdmin := user.Role == "admin"

	if !isProjectAdminOrAdmin {
		isPAdmin, err := s.isProjectAdmin(existingTask.ProjectID, user.ID)
		if err != nil {
			return errors.New("gagal memvalidasi admin project")
		}
		isProjectAdminOrAdmin = isPAdmin
	}

	if !isProjectAdminOrAdmin {
		isMember, err := s.repo.IsUserMember(taskID, user.ID)
		if err != nil {
			return errors.New("gagal memvalidasi member task")
		}
		if !isMember {
			return errors.New("anda bukan member dari task ini, tidak bisa mengupdate")
		}

		allowedUpdates := make(map[string]interface{})
		for key, value := range updates {
			if key == "status" || key == "notes" {
				allowedUpdates[key] = value
			} else {
				return errors.New("anda hanya diizinkan untuk mengupdate status dan notes")
			}
		}
		finalUpdates = allowedUpdates
	}

	if len(finalUpdates) == 0 {
		return errors.New("tidak ada field yang diizinkan untuk diupdate")
	}

	if newStatus, ok := finalUpdates["status"].(string); ok && newStatus != existingTask.Status {
		activity := models.ActivityLog{
			UserID:    user.ID,
			Action:    fmt.Sprintf("User changed status of task '%s' from '%s' to '%s'", existingTask.Title, existingTask.Status, newStatus),
			TableName: "tasks",
			ItemID:    taskID,
		}
		s.activityLogger.Log(activity)

		now := time.Now()
		lastLog, err := s.taskStatusLog.FindLastLog(taskID)
		if err != nil {
			return err
		}
		if lastLog != nil {
			err = s.taskStatusLog.UpdateClockOut(lastLog.ID, now)
			if err != nil {
				return err
			}
		}

		newLog := &models.TaskStatusLog{
			TaskID:  taskID,
			Status:  newStatus,
			ClockIn: now,
		}
		err = s.taskStatusLog.Create(newLog)
		if err != nil {
			return err
		}
	}

	if status, ok := finalUpdates["status"]; ok {
		if statusStr, isString := status.(string); isString {
			newStatusLower := strings.ToLower(statusStr)
			oldStatusLower := strings.ToLower(existingTask.Status)

			if newStatusLower == "done" {
				if oldStatusLower != "done" {
					now := time.Now()
					finalUpdates["finished_at"] = &now
				}
			} else if oldStatusLower == "done" && newStatusLower != "done" {
				finalUpdates["finished_at"] = nil
			}
		}
	}

	if time.Now().After(existingTask.DueDate) {
		duration := time.Since(existingTask.DueDate)
		finalUpdates["overdue_duration"] = duration
	} else {
		finalUpdates["overdue_duration"] = time.Duration(0)
	}

	return s.repo.UpdateTask(taskID, finalUpdates)
}

func (s *taskService) SoftDeleteTask(taskID uint, workspaceID uint, user *models.User) error {
	_, err := s.GetByID(taskID, workspaceID, user)
	if err != nil {
		return errors.New("task tidak ditemukan")
	}

	return s.repo.SoftDeleteTask(taskID)
}

func (s *taskService) DeleteTask(taskID uint, workspaceID uint, user *models.User) error {
	_, err := s.GetByID(taskID, workspaceID, user)
	if err != nil {
		return errors.New("task tidak ditemukan")
	}

	return s.repo.DeleteTask(taskID)
}

func (s *taskService) AddMember(taskID uint, projectID uint, workspaceID uint, userID uint, role string, currentUser *models.User) error {
	task, err := s.repo.GetByID(taskID)
	if err != nil {
		return errors.New("task tidak ditemukan")
	}

	if task.Project.WorkspaceID != workspaceID {
		return errors.New("task tidak ditemukan di workspace ini")
	}

	if task.ProjectID != projectID {
		return errors.New("task tidak ditemukan di project ini")
	}

	isMember, err := s.repo.IsUserMember(taskID, userID)
	if err != nil {
		return errors.New("gagal memvalidasi member")
	}
	if isMember {
		return errors.New("user sudah menjadi member di task ini")
	}

	member := &models.TaskUser{
		TaskID:     taskID,
		UserID:     userID,
		RoleInTask: role,
		AssignedAt: task.CreatedAt,
	}
	if err := s.repo.AddMember(member); err != nil {
		return err
	}

	assignedUser, err := s.userRepo.GetByID(userID)
	if err == nil && assignedUser.TelegramChatID != nil && *assignedUser.TelegramChatID != "" {
		message := fmt.Sprintf("Anda telah ditugaskan untuk task baru: %s", task.Title)
		go s.telegramService.SendNotification(*assignedUser.TelegramChatID, message)
	}

	return nil
}

func (s *taskService) GetMembers(taskID uint, projectID uint, workspaceID uint, user *models.User) ([]models.TaskUser, error) {
	task, err := s.repo.GetByID(taskID)
	if err != nil {
		return nil, errors.New("task tidak ditemukan")
	}

	if task.Project.WorkspaceID != workspaceID {
		return nil, errors.New("task tidak ditemukan di workspace ini")
	}

	if task.ProjectID != projectID {
		return nil, errors.New("task tidak ditemukan di project ini")
	}

	if user.Role == "admin" {
		return s.repo.GetMembers(taskID)
	}

	if user.Role != "admin" {
		isTaskMember, err := s.repo.IsUserMember(taskID, user.ID)
		if err != nil {
			return nil, errors.New("gagal memvalidasi member task")
		}
		if !isTaskMember {
			return nil, errors.New("anda bukan member dari task ini")
		}

		isProjectMember, err := s.repo.IsUserInProject(task.ProjectID, user.ID)
		if err != nil || !isProjectMember {
			return nil, errors.New("hanya member project yang boleh melihat members task")
		}

		isWorkspaceMember, err := s.repo.IsProjectInWorkspace(task.ProjectID, workspaceID)
		if err != nil || !isWorkspaceMember {
			return nil, errors.New("task tidak ditemukan di workspace ini")
		}

	}
	return s.repo.GetMembers(taskID)
}

func (s *taskService) DeleteMember(taskID uint, projectID uint, workspaceID uint, userID uint, currentUser *models.User) error {
	task, err := s.repo.GetByID(taskID)
	if err != nil {
		return errors.New("task tidak ditemukan")
	}

	if task.Project.WorkspaceID != workspaceID {
		return errors.New("task tidak ditemukan di workspace ini")
	}

	if task.ProjectID != projectID {
		return errors.New("task tidak ditemukan di project ini")
	}

	return s.repo.DeleteMember(taskID, userID)

}

func (s *taskService) isProjectAdmin(projectID uint, userID uint) (bool, error) {
	var projectUser models.ProjectUser
	err := config.DB.
		Where("project_id = ? AND user_id = ?", projectID, userID).
		First(&projectUser).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return projectUser.RoleInProject == "admin", nil
}
