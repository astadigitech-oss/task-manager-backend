package repositories

import (
	"project-management-backend/config"
	"project-management-backend/models"
)

type ProjectImageRepository interface {
	CreateProjectImage(image *models.ProjectImage) error
	GetProjectImages(projectID uint) ([]models.ProjectImage, error)
	GetProjectImageByID(imageID uint) (*models.ProjectImage, error)
	DeleteProjectImage(imageID uint) error
	IsUserProjectMember(projectID uint, userID uint) (bool, error)
}

type projectImageRepository struct{}

func NewProjectImageRepository() ProjectImageRepository {
	return &projectImageRepository{}
}

func (r *projectImageRepository) CreateProjectImage(image *models.ProjectImage) error {
	return config.DB.Create(image).Error
}

func (r *projectImageRepository) GetProjectImages(projectID uint) ([]models.ProjectImage, error) {
	var images []models.ProjectImage
	err := config.DB.
		Where("project_id = ?", projectID).
		Order("created_at DESC").
		Find(&images).Error
	return images, err
}

func (r *projectImageRepository) GetProjectImageByID(imageID uint) (*models.ProjectImage, error) {
	var image models.ProjectImage
	err := config.DB.
		Preload("Project").
		First(&image, imageID).Error
	return &image, err
}

func (r *projectImageRepository) DeleteProjectImage(imageID uint) error {
	return config.DB.Where("id = ?", imageID).Delete(&models.ProjectImage{}).Error
}

func (r *projectImageRepository) IsUserProjectMember(projectID uint, userID uint) (bool, error) {
	var count int64
	err := config.DB.Model(&models.ProjectUser{}).
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Count(&count).Error
	return count > 0, err
}
