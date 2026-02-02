package services

import (
	"errors"
	"io"
	"mime/multipart"
	"project-management-backend/models"
	"project-management-backend/repositories"
)

type TaskFileService struct {
	repo *repositories.TaskFileRepository
}

func NewTaskFileService(repo *repositories.TaskFileRepository) *TaskFileService {
	return &TaskFileService{repo: repo}
}

func (s *TaskFileService) UploadFile(taskID uint, userID uint, file *multipart.FileHeader) (*models.TaskFile, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	fileData, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	taskFile := &models.TaskFile{
		TaskID:     taskID,
		FileName:   file.Filename,
		FileData:   fileData,
		MimeType:   file.Header.Get("Content-Type"),
		FileSize:   file.Size,
		UploadedBy: userID,
	}

	err = s.repo.Create(taskFile)
	if err != nil {
		return nil, err
	}

	return taskFile, nil
}

func (s *TaskFileService) GetFilesByTaskID(taskID uint) ([]models.TaskFile, error) {
	return s.repo.FindByTaskID(taskID)
}

func (s *TaskFileService) DownloadFile(fileID uint) ([]byte, string, string, error) {
	fileData, mimeType, err := s.repo.GetFileData(fileID)
	if err != nil {
		return nil, "", "", err
	}

	taskFile, err := s.repo.FindByID(fileID)
	if err != nil {
		return nil, "", "", err
	}

	return fileData, mimeType, taskFile.FileName, nil
}

func (s *TaskFileService) DeleteFile(fileID uint, userID uint) error {
	taskFile, err := s.repo.FindByID(fileID)
	if err != nil {
		return err
	}

	return s.repo.Delete(taskFile)
}

// Optional: Add file size limit validation
func (s *TaskFileService) ValidateFile(file *multipart.FileHeader) error {
	// Max 10MB file size
	const maxFileSize = 10 * 1024 * 1024 // 10MB

	if file.Size > maxFileSize {
		return errors.New("file size exceeds 10MB limit")
	}

	// Optional: Validate file type
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
