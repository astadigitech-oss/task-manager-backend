// repositories/task_image_repository.go
package repositories

import (
	"project-management-backend/config"
	"project-management-backend/models"
)

type TaskImageRepository interface {
	CreateTaskImage(image *models.TaskImage) error
	GetTaskImages(taskID uint) ([]models.TaskImage, error)
	GetTaskImageByID(imageID uint) (*models.TaskImage, error)
	DeleteTaskImage(imageID uint) error
}

type taskImageRepository struct{}

func NewTaskImageRepository() TaskImageRepository {
	return &taskImageRepository{}
}

func (r *taskImageRepository) CreateTaskImage(image *models.TaskImage) error {
	return config.DB.Create(image).Error
}

func (r *taskImageRepository) GetTaskImages(taskID uint) ([]models.TaskImage, error) {
	var images []models.TaskImage
	err := config.DB.Where("task_id = ?", taskID).Find(&images).Error
	return images, err
}

func (r *taskImageRepository) GetTaskImageByID(imageID uint) (*models.TaskImage, error) {
	var image models.TaskImage
	err := config.DB.First(&image, imageID).Error
	return &image, err
}

func (r *taskImageRepository) DeleteTaskImage(imageID uint) error {
	return config.DB.Delete(&models.TaskImage{}, imageID).Error
}
