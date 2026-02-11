package repositories

import (
	"project-management-backend/models"

	"gorm.io/gorm"
)

type AttendanceImageRepository struct {
	db *gorm.DB
}

func NewAttendanceImageRepository(db *gorm.DB) *AttendanceImageRepository {
	return &AttendanceImageRepository{db: db}
}

func (r *AttendanceImageRepository) Create(image *models.AttendanceImage) error {
	return r.db.Create(image).Error
}

func (r *AttendanceImageRepository) GetByAttendanceID(attendanceID uint) ([]models.AttendanceImage, error) {
	var images []models.AttendanceImage
	err := r.db.Where("attendance_id = ?", attendanceID).Find(&images).Error
	return images, err
}

func (r *AttendanceImageRepository) Delete(imageID uint) error {
	return r.db.Where("id = ?", imageID).Delete(&models.AttendanceImage{}).Error
}
