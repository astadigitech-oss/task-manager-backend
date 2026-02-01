package services

import (
	"errors"
	"project-management-backend/config"
	"project-management-backend/models"
	"project-management-backend/repositories"
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
	repo repositories.TaskRepository
}

func NewTaskService(repo repositories.TaskRepository) TaskService {
	return &taskService{repo: repo}
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

	return s.repo.CreateTask(task)
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

	if user.Role != "admin" {
		if existingTask.Project.WorkspaceID != workspaceID {
			return errors.New("task tidak ditemukan di workspace ini")
		}

		isPAdmin, err := s.isProjectAdmin(existingTask.ProjectID, user.ID)
		if err != nil {
			return errors.New("gagal memvalidasi admin project")
		}

		if !isPAdmin {
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

			if len(allowedUpdates) == 0 {
				return errors.New("tidak ada field yang diizinkan untuk diupdate")
			}

			return s.repo.UpdateTask(taskID, allowedUpdates)
		}
	}

	if status, ok := updates["status"]; ok {
		if status == "done" {
			now := time.Now()
			updates["finished_at"] = &now

			if now.After(existingTask.DueDate) {
				duration := now.Sub(existingTask.DueDate)
				updates["overdue_duration"] = duration
			} else {
				updates["overdue_duration"] = time.Duration(0)
			}
		} else if existingTask.Status == "done" && status != "done" {
			updates["finished_at"] = nil
			updates["overdue_duration"] = time.Duration(0)
		}
	}

	return s.repo.UpdateTask(taskID, updates)
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
	return s.repo.AddMember(member)
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
