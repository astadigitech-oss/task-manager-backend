// services/project_image_service.go
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

type ProjectImageService interface {
	UploadProjectImage(projectID uint, file *multipart.FileHeader, userID uint) (*models.ProjectImage, error)
	GetProjectImages(projectID uint, userID uint) ([]models.ProjectImage, error)
	DeleteProjectImage(imageID uint, userID uint) error
}

type projectImageService struct {
	repo          repositories.ProjectImageRepository
	projectRepo   repositories.ProjectRepository
	workspaceRepo repositories.WorkspaceRepository
	userRepo      repositories.UserRepository
}

func NewProjectImageService(repo repositories.ProjectImageRepository, projectRepo repositories.ProjectRepository, workspaceRepo repositories.WorkspaceRepository, userRepo repositories.UserRepository) ProjectImageService {
	return &projectImageService{repo: repo, projectRepo: projectRepo, workspaceRepo: workspaceRepo, userRepo: userRepo}
}

func (s *projectImageService) UploadProjectImage(projectID uint, file *multipart.FileHeader, userID uint) (*models.ProjectImage, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, errors.New("project tidak ditemukan atau workspace sudah dihapus")
	}

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	if user.Role != "admin" {
		hasWorkspaceAccess, err := s.workspaceRepo.IsUserMember(project.WorkspaceID, userID)
		if err != nil || !hasWorkspaceAccess {
			return nil, errors.New("tidak memiliki akses ke workspace project ini")
		}

		isMember, err := s.projectRepo.IsUserMember(projectID, userID)
		if err != nil || !isMember {
			return nil, errors.New("hanya member project yang boleh upload image")
		}
	}

	// Validasi file type
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

	// Create uploads directory if not exists
	uploadDir := "./uploads/projects"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, errors.New("gagal membuat directory upload")
	}

	// Save file
	filePath := filepath.Join(uploadDir, fileName)
	if err := saveUploadedFile(file, filePath); err != nil {
		return nil, errors.New("gagal menyimpan file: " + err.Error())
	}

	// Create relative URL untuk database
	relativeURL := "/uploads/projects/" + fileName

	// Create project image record
	projectImage := &models.ProjectImage{
		ProjectID:  projectID,
		URL:        relativeURL,
		UploadedBy: userID,
	}

	if err := s.repo.CreateProjectImage(projectImage); err != nil {
		os.Remove(filePath)
		return nil, errors.New("gagal menyimpan data image: " + err.Error())
	}

	return projectImage, nil
}

func (s *projectImageService) GetProjectImages(projectID uint, userID uint) ([]models.ProjectImage, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, errors.New("project tidak ditemukan atau workspace sudah dihapus")
	}

	user, err := s.projectRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	if user.Role == "admin" {
		return s.repo.GetProjectImages(projectID)
	}

	hasWorkspaceAccess, err := s.workspaceRepo.IsUserMember(project.WorkspaceID, userID)
	if err != nil || !hasWorkspaceAccess {
		return nil, errors.New("tidak memiliki akses ke workspace project ini")
	}

	isMember, err := s.projectRepo.IsUserMember(projectID, userID)
	if err != nil || !isMember {
		return nil, errors.New("hanya member project yang boleh melihat images")
	}

	return s.repo.GetProjectImages(projectID)
}

func (s *projectImageService) DeleteProjectImage(imageID uint, userID uint) error {
	image, err := s.repo.GetProjectImageByID(imageID)
	if err != nil {
		return errors.New("image tidak ditemukan")
	}

	project, err := s.projectRepo.GetByID(image.ProjectID)
	if err != nil {
		return errors.New("project tidak ditemukan atau workspace sudah dihapus")
	}

	hasWorkspaceAccess, err := s.workspaceRepo.IsUserMember(project.WorkspaceID, userID)
	if err != nil || !hasWorkspaceAccess {
		return errors.New("tidak memiliki akses ke workspace project ini")
	}

	isProjectMember, err := s.projectRepo.IsUserMember(image.ProjectID, userID)
	if err != nil || !isProjectMember {
		return errors.New("hanya member project yang boleh menghapus image")
	}

	filePath := "." + image.URL // Karena URL relative
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		// Log error tapi tetap lanjut delete dari database
		// return errors.New("gagal menghapus file: " + err.Error())
	}

	return s.repo.DeleteProjectImage(imageID)
}

// Helper function untuk save file
func saveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Create destination file
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// Copy file
	_, err = out.ReadFrom(src)
	return err
}
