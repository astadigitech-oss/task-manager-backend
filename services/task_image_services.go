// services/task_image_service.go
package services

import (
	"errors"
	"mime/multipart"
	"os"
	"path/filepath"
	"project-management-backend/models"
	"project-management-backend/repositories"

	"github.com/google/uuid"
)

type TaskImageService interface {
	UploadTaskImage(taskID uint, workspaceID uint, file *multipart.FileHeader, user *models.User) (*models.TaskImage, error)
	GetTaskImages(taskID uint, projectID uint, workspaceID uint, user *models.User) ([]models.TaskImage, error)
	DeleteTaskImage(imageID uint, workspaceID uint, user *models.User) error
}

type taskImageService struct {
	repo          repositories.TaskImageRepository
	taskRepo      repositories.TaskRepository
	projectRepo   repositories.ProjectRepository
	workspaceRepo repositories.WorkspaceRepository
}

func NewTaskImageService(
	repo repositories.TaskImageRepository,
	taskRepo repositories.TaskRepository,
	projectRepo repositories.ProjectRepository,
	workspaceRepo repositories.WorkspaceRepository,
) TaskImageService {
	return &taskImageService{
		repo:          repo,
		taskRepo:      taskRepo,
		projectRepo:   projectRepo,
		workspaceRepo: workspaceRepo,
	}
}

func (s *taskImageService) UploadTaskImage(taskID uint, workspaceID uint, file *multipart.FileHeader, user *models.User) (*models.TaskImage, error) {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return nil, errors.New("task tidak ditemukan")
	}

	if task.Project.WorkspaceID != workspaceID {
		return nil, errors.New("task tidak ditemukan di workspace ini")
	}

	if user.Role != "admin" {
		hasWorkspaceAccess, err := s.workspaceRepo.IsUserMember(workspaceID, user.ID)
		if err != nil || !hasWorkspaceAccess {
			return nil, errors.New("tidak memiliki akses ke workspace ini")
		}

		isTaskMember, err := s.taskRepo.IsUserMember(taskID, user.ID)
		if err != nil || !isTaskMember {
			return nil, errors.New("hanya task member yang boleh upload image")
		}
	}

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	fileType := file.Header.Get("Content-Type")
	if !allowedTypes[fileType] {
		return nil, errors.New("format file tidak didukung. Hanya JPEG, PNG, GIF, WebP")
	}

	// Validasi file size (max 5MB)
	if file.Size > 5*1024*1024 {
		return nil, errors.New("ukuran file maksimal 5MB")
	}

	// Generate unique filename
	fileExt := filepath.Ext(file.Filename)
	fileName := uuid.New().String() + fileExt

	// Create uploads directory
	uploadDir := "./uploads/tasks"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, errors.New("gagal membuat directory upload")
	}

	// Save file
	filePath := filepath.Join(uploadDir, fileName)
	if err := saveUploadedFile(file, filePath); err != nil {
		return nil, errors.New("gagal menyimpan file: " + err.Error())
	}

	// Create relative URL
	relativeURL := "/uploads/tasks/" + fileName

	// Create task image record
	taskImage := &models.TaskImage{
		TaskID:     taskID,
		URL:        relativeURL,
		UploadedBy: user.ID,
	}

	if err := s.repo.CreateTaskImage(taskImage); err != nil {
		// Rollback: delete file jika gagal save ke database
		os.Remove(filePath)
		return nil, errors.New("gagal menyimpan data image: " + err.Error())
	}

	return taskImage, nil
}

func (s *taskImageService) GetTaskImages(taskID uint, projectID uint, workspaceID uint, user *models.User) ([]models.TaskImage, error) {
	task, err := s.taskRepo.GetByID(taskID)
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
		return s.repo.GetTaskImages(taskID)
	}

	hasWorkspaceAccess, err := s.workspaceRepo.IsUserMember(workspaceID, user.ID)
	if err != nil || !hasWorkspaceAccess {
		return nil, errors.New("tidak memiliki akses ke workspace ini")
	}

	isTaskMember, err := s.taskRepo.IsUserMember(taskID, user.ID)
	if err != nil || !isTaskMember {
		return nil, errors.New("hanya task member yang boleh melihat image")
	}

	isProjectMember, err := s.taskRepo.IsUserInProject(task.ProjectID, user.ID)
	if err != nil || !isProjectMember {
		return nil, errors.New("hanya project member yang boleh melihat image")
	}

	return s.repo.GetTaskImages(taskID)
}

func (s *taskImageService) DeleteTaskImage(imageID uint, workspaceID uint, user *models.User) error {
	image, err := s.repo.GetTaskImageByID(imageID)
	if err != nil {
		return errors.New("image tidak ditemukan")
	}

	task, err := s.taskRepo.GetByID(image.TaskID)
	if err != nil {
		return errors.New("task tidak ditemukan")
	}

	if task.Project.WorkspaceID != workspaceID {
		return errors.New("task tidak ditemukan di workspace ini")
	}

	if user.Role != "admin" {
		hasWorkspaceAccess, err := s.workspaceRepo.IsUserMember(workspaceID, user.ID)
		if err != nil || !hasWorkspaceAccess {
			return errors.New("tidak memiliki akses ke workspace ini")
		}
		if image.UploadedBy != user.ID {
			return errors.New("hanya uploader yang boleh menghapus image")
		}
	}

	filePath := "." + image.URL
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		// Log error tapi tetap lanjut delete dari database
	}

	return s.repo.DeleteTaskImage(imageID)
}
