package repositories

import (
	"project-management-backend/models"

	"gorm.io/gorm"
)

type TaskFileRepository struct {
	DB *gorm.DB
}

func NewTaskFileRepository(db *gorm.DB) *TaskFileRepository {
	return &TaskFileRepository{DB: db}
}

func (r *TaskFileRepository) Create(taskFile *models.TaskFile) error {
	return r.DB.Create(taskFile).Error
}

func (r *TaskFileRepository) FindByTaskID(taskID uint) ([]models.TaskFile, error) {
	var taskFiles []models.TaskFile
	// IMPORTANT: Don't load FileData when listing files (for performance)
	err := r.DB.Select("id", "task_id", "filename", "mime_type", "file_size", "uploaded_by", "created_at").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "email") // Only select needed user fields
		}).
		Where("task_id = ?", taskID).
		Find(&taskFiles).Error
	return taskFiles, err
}

func (r *TaskFileRepository) FindByID(id uint) (*models.TaskFile, error) {
	var taskFile models.TaskFile
	err := r.DB.First(&taskFile, id).Error
	return &taskFile, err
}

func (r *TaskFileRepository) GetFileData(id uint) ([]byte, string, error) {
	var taskFile models.TaskFile
	err := r.DB.Select("file_data", "mime_type", "filename").
		First(&taskFile, id).Error
	return taskFile.FileData, taskFile.MimeType, err
}

func (r *TaskFileRepository) Delete(taskFile *models.TaskFile) error {
	return r.DB.Delete(taskFile).Error
}
