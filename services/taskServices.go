package services

import (
	"errors"
	"project-management-backend/config"
	"project-management-backend/models"
	"project-management-backend/repositories"

	"gorm.io/gorm"
)

type TaskService interface {
	CreateTask(task *models.Task, workspaceID uint, user *models.User) error
	GetAllTasks(projectID uint, workspaceID uint, user *models.User) ([]models.Task, error)
	GetByID(taskID uint, workspaceID uint, user *models.User) (*models.Task, error)
	UpdateTask(task *models.Task, workspaceID uint, user *models.User) error
	SoftDeleteTask(taskID uint, workspaceID uint, user *models.User) error
	DeleteTask(taskID uint, workspaceID uint, user *models.User) error
	AddMember(taskID uint, workspaceID uint, userID uint, role string, currentUser *models.User) error
	GetMembers(taskID uint, workspaceID uint, user *models.User) ([]models.TaskUser, error)
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

	isMember, err := s.repo.IsUserInProject(task.ProjectID, user.ID)
	if err != nil || !isMember {
		return errors.New("hanya member project yang boleh buat task")
	}

	if task.Status == "" {
		task.Status = "todo"
	}
	if task.Priority == "" {
		task.Priority = "medium"
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

	isProjectAdmin, err := s.isProjectAdmin(projectID, user.ID)
	if err != nil {
		return nil, err
	}
	if isProjectAdmin {
		return s.repo.GetAllTasks(projectID)
	} else {
		return s.repo.GetTasksByUserID(projectID, user.ID)
	}
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

func (s *taskService) UpdateTask(task *models.Task, workspaceID uint, user *models.User) error {
	existingTask, err := s.GetByID(task.ID, workspaceID, user)
	if err != nil {
		return errors.New("task tidak ditemukan")
	}

	isProjectMember, err := s.repo.IsUserInProject(existingTask.ProjectID, user.ID)
	if err != nil || !isProjectMember {
		return errors.New("hanya member project yang boleh update task")
	}

	return s.repo.UpdateTask(task)
}

func (s *taskService) SoftDeleteTask(taskID uint, workspaceID uint, user *models.User) error {
	task, err := s.GetByID(taskID, workspaceID, user)
	if err != nil {
		return errors.New("task tidak ditemukan")
	}

	isProjectMember, err := s.repo.IsUserInProject(task.ProjectID, user.ID)
	if err != nil || !isProjectMember {
		return errors.New("hanya member project yang boleh soft delete task")
	}

	return s.repo.SoftDeleteTask(taskID)
}

func (s *taskService) DeleteTask(taskID uint, workspaceID uint, user *models.User) error {
	task, err := s.GetByID(taskID, workspaceID, user)
	if err != nil {
		return errors.New("task tidak ditemukan")
	}

	isProjectMember, err := s.repo.IsUserInProject(task.ProjectID, user.ID)
	if err != nil || !isProjectMember {
		return errors.New("hanya member project yang boleh hard delete task")
	}

	return s.repo.DeleteTask(taskID)
}

func (s *taskService) AddMember(taskID uint, workspaceID uint, userID uint, role string, currentUser *models.User) error {
	task, err := s.repo.GetByID(taskID)
	if err != nil {
		return errors.New("task tidak ditemukan")
	}

	isProjectMember, err := s.repo.IsUserInProject(task.ProjectID, currentUser.ID)
	if err != nil || !isProjectMember {
		return errors.New("hanya member project yang boleh menambah member task")
	}

	isTargetUserInProject, err := s.repo.IsUserInProject(task.ProjectID, userID)
	if err != nil || !isTargetUserInProject {
		return errors.New("user harus menjadi member project terlebih dahulu")
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

func (s *taskService) GetMembers(taskID uint, workspaceID uint, user *models.User) ([]models.TaskUser, error) {
	task, err := s.repo.GetByID(taskID)
	if err != nil {
		return nil, errors.New("task tidak ditemukan")
	}

	isTaskMember, _ := s.repo.IsUserMember(taskID, user.ID)
	isProjectMember, _ := s.repo.IsUserInProject(task.ProjectID, user.ID)

	if !isTaskMember && !isProjectMember {
		return nil, errors.New("akses ditolak untuk melihat members task")
	}

	return s.repo.GetMembers(taskID)
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
