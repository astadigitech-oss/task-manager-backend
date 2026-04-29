package repositories

import (
	"project-management-backend/config"
	"project-management-backend/models"
	"time"

	"gorm.io/gorm"
)

type TaskStatusLogRepository interface {
	Create(log *models.TaskStatusLog) error
	FindLastLog(taskID uint) (*models.TaskStatusLog, error)
	UpdateClockOut(logID uint, clockOut time.Time) error
	GetLogsByTaskID(taskID uint) ([]models.TaskStatusLog, error)
}

type taskStatusLogRepository struct{}

func NewTaskStatusLogRepository() TaskStatusLogRepository {
	return &taskStatusLogRepository{}
}

func (r *taskStatusLogRepository) Create(log *models.TaskStatusLog) error {
	return config.DB.Create(log).Error
}

func (r *taskStatusLogRepository) FindLastLog(taskID uint) (*models.TaskStatusLog, error) {
	var log models.TaskStatusLog
	err := config.DB.Where("task_id = ?", taskID).Order("created_at desc").First(&log).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &log, nil
}

func (r *taskStatusLogRepository) UpdateClockOut(logID uint, clockOut time.Time) error {
	var log models.TaskStatusLog
	if err := config.DB.First(&log, logID).Error; err != nil {
		return err
	}
	if log.ClockIn.IsZero() {
		return config.DB.Model(&models.TaskStatusLog{}).Where("id = ?", logID).Update("clock_out", clockOut).Error
	}
	duration := clockOut.Sub(log.ClockIn).Milliseconds()
	return config.DB.Model(&models.TaskStatusLog{}).Where("id = ?", logID).Updates(map[string]interface{}{
		"clock_out": clockOut,
		"duration":  duration,
	}).Error
}

func (r *taskStatusLogRepository) GetLogsByTaskID(taskID uint) ([]models.TaskStatusLog, error) {
	var logs []models.TaskStatusLog
	err := config.DB.Where("task_id = ?", taskID).Order("created_at asc").Find(&logs).Error
	return logs, err
}
