package services

import (
	"errors"
	"project-management-backend/models"
	"project-management-backend/repositories"
)

type TaskService interface {
	CreateTask(task *models.Task, workspaceID uint, user *models.User) error
	GetAllTasks(projectID uint, workspaceID uint, user *models.User) ([]models.Task, error)
	GetByID(taskID uint, workspaceID uint, user *models.User) (*models.Task, error)
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

	isMember, err := s.repo.IsUserInProject(projectID, user.ID)
	if err != nil || !isMember {
		return nil, errors.New("hanya member project yang boleh lihat tasks")
	}

	return s.repo.GetAllTasks(projectID, workspaceID)
}

func (s *taskService) GetByID(taskID uint, workspaceID uint, user *models.User) (*models.Task, error) {
	task, err := s.repo.GetByID(taskID, workspaceID)
	if err != nil {
		return nil, errors.New("task tidak ditemukan")
	}

	isTaskMember, _ := s.repo.IsUserMember(taskID, user.ID)
	isProjectMember, _ := s.repo.IsUserInProject(task.ProjectID, user.ID)

	if !isTaskMember && !isProjectMember {
		return nil, errors.New("akses ditolak untuk task ini")
	}

	return task, nil
}

func (s *taskService) AddMember(taskID uint, workspaceID uint, userID uint, role string, currentUser *models.User) error {
	task, err := s.repo.GetByID(taskID, workspaceID)
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
	task, err := s.repo.GetByID(taskID, workspaceID)
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
