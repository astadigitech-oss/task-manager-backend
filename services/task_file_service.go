package services

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"project-management-backend/models"
	"project-management-backend/repositories"

	"github.com/google/uuid"
)

type TaskFileService struct {
	repo        *repositories.TaskFileRepository
	taskRepo    repositories.TaskRepository
	projectRepo repositories.ProjectRepository
}

func NewTaskFileService(repo *repositories.TaskFileRepository, taskRepo repositories.TaskRepository, projectRepo repositories.ProjectRepository) *TaskFileService {
	return &TaskFileService{
		repo:        repo,
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
	}
}

func (s *TaskFileService) UploadFile(workspaceID, projectID, taskID, userID uint, file *multipart.FileHeader) (*models.TaskFile, error) {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return nil, errors.New("task not found")
	}
	if task.ProjectID != projectID {
		return nil, errors.New("task not found in project")
	}

	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, errors.New("project not found")
	}

	if project.WorkspaceID != workspaceID {
		return nil, errors.New("project not found in workspace")
	}

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	fileExt := filepath.Ext(file.Filename)
	fileName := uuid.New().String() + fileExt

	uploadDir := "./uploads/tasks/files"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, errors.New("gagal membuat directory upload")
	}

	filePath, err := filepath.Abs(filepath.Join(uploadDir, fileName))
	if err != nil {
		return nil, errors.New("gagal membuat absolute path")
	}

	out, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		return nil, err
	}

	taskFile := &models.TaskFile{
		TaskID:     taskID,
		FileName:   file.Filename,
		URL:        filePath,
		MimeType:   file.Header.Get("Content-Type"),
		FileSize:   file.Size,
		UploadedBy: userID,
	}

	err = s.repo.Create(taskFile)
	if err != nil {
		os.Remove(filePath)
		return nil, err
	}

	return taskFile, nil
}

func (s *TaskFileService) GetFilesByTaskID(taskID, projectID, workspaceID uint) ([]models.TaskFile, error) {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return nil, errors.New("task not found")
	}
	if task.ProjectID != projectID {
		return nil, errors.New("task not found in project")
	}

	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, errors.New("project not found")
	}

	if project.WorkspaceID != workspaceID {
		return nil, errors.New("project not found in workspace")
	}
	return s.repo.FindByTaskID(taskID)
}

func (s *TaskFileService) DownloadFile(fileID uint) ([]byte, string, string, error) {
	taskFile, err := s.repo.FindByID(fileID)
	if err != nil {
		return nil, "", "", err
	}

	fileData, err := os.ReadFile(taskFile.URL)
	if err != nil {
		return nil, "", "", err
	}

	return fileData, taskFile.MimeType, taskFile.FileName, nil
}

func (s *TaskFileService) DeleteFile(fileID uint, userID uint) error {
	taskFile, err := s.repo.FindByID(fileID)
	if err != nil {
		return err
	}

	if err := os.Remove(taskFile.URL); err != nil && !os.IsNotExist(err) {
	}

	return s.repo.Delete(taskFile)
}

func (s *TaskFileService) ValidateFile(file *multipart.FileHeader) error {
	const maxFileSize = 10 * 1024 * 1024

	if file.Size > maxFileSize {
		return errors.New("file size exceeds 10MB limit")
	}

	allowedTypes := map[string]bool{
		"application/pdf":    true,
		"image/jpeg":         true,
		"image/png":          true,
		"image/gif":          true,
		"text/plain":         true,
		"application/msword": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	}

	contentType := file.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return errors.New("file type not allowed")
	}

	return nil
}
