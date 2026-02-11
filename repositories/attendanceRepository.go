package repositories

import (
	"time"

	"project-management-backend/models"

	"gorm.io/gorm"
)

type AttendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) *AttendanceRepository {
	return &AttendanceRepository{db: db}
}

func (r *AttendanceRepository) GetByID(attendanceID uint) (*models.Attendance, error) {
	var attendance models.Attendance
	err := r.db.Preload("Images").First(&attendance, "id = ?", attendanceID).Error
	return &attendance, err
}

func (r *AttendanceRepository) Create(attendance *models.Attendance) error {
	return r.db.Create(attendance).Error
}

func (r *AttendanceRepository) GetAttendanceByUserIDAndDateRange(userID uint, start, end time.Time) (models.Attendance, error) {
	var attendance models.Attendance
	err := r.db.Where("user_id = ? AND clock_in >= ? AND clock_in < ?", userID, start, end).First(&attendance).Error
	return attendance, err
}

func (r *AttendanceRepository) GetAttendancesByWorkspaceIDAndDateRange(workspaceID uint, start, end time.Time) ([]models.Attendance, error) {
	var attendances []models.Attendance
	err := r.db.Preload("User").Where("workspace_id = ? AND clock_in >= ? AND clock_in < ?", workspaceID, start, end).Find(&attendances).Error
	return attendances, err
}
