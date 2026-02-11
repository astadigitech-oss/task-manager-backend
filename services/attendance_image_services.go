package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"project-management-backend/models"
	"project-management-backend/repositories"
)

const (
	maxSize  = 2 << 20 // 2 MB
	imageDir = "./uploads/attendance/"
)

type AttendanceImageService struct {
	repo *repositories.AttendanceImageRepository
}

func NewAttendanceImageService(repo *repositories.AttendanceImageRepository) *AttendanceImageService {
	return &AttendanceImageService{repo: repo}
}

func (s *AttendanceImageService) UploadImage(attendanceID uint, file *multipart.FileHeader) (*models.AttendanceImage, error) {
	if file.Size > maxSize {
		return nil, fmt.Errorf("file size exceeds the limit of 2MB")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}

	fileName := fmt.Sprintf("%d_%d%s", attendanceID, time.Now().UnixNano(), ext)
	filePath := filepath.Join(imageDir, fileName)

	if err := os.MkdirAll(imageDir, os.ModePerm); err != nil {
		return nil, err
	}

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, err
	}

	imageURL := "/uploads/attendance/" + fileName

	image := &models.AttendanceImage{
		AttendanceID: attendanceID,
		URL:          imageURL,
	}

	if err := s.repo.Create(image); err != nil {
		return nil, err
	}

	return image, nil
}

func (s *AttendanceImageService) GetImagesByAttendanceID(attendanceID uint) ([]models.AttendanceImage, error) {
	return s.repo.GetByAttendanceID(attendanceID)
}

func (s *AttendanceImageService) DeleteImage(imageID uint) error {
	return s.repo.Delete(imageID)
}
