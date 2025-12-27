package utils

import (
	"encoding/json"
	"project-management-backend/config"
	"project-management-backend/models"
	"time"

	"gorm.io/gorm"
)

type ActivityLogger interface {
	Log(activity models.ActivityLog)
}

type activityLogger struct {
	db *gorm.DB
}

func (l *activityLogger) Log(activity models.ActivityLog) {
	l.db.Create(&activity)
}

func NewActivityLogger(db *gorm.DB) ActivityLogger {
	return &activityLogger{db: db}
}

// Activity Log
func ActivityLog(userID uint, action, table string, itemID uint, before interface{}, after interface{}) error {
	beforeJSON, _ := json.Marshal(before)
	afterJSON, _ := json.Marshal(after)

	log := models.ActivityLog{
		UserID:     userID,
		Action:     action,
		TableName:  table,
		ItemID:     itemID,
		DataBefore: string(beforeJSON),
		DataAfter:  string(afterJSON),
	}
	return config.DB.Create(&log).Error
}

// Error Log
func Error(userID uint, action, table string, itemID uint, errMsg, stack string) {
	log := models.ErrorLog{
		UserID:     userID,
		Action:     action,
		TableName:  table,
		ItemID:     itemID,
		ErrorMsg:   errMsg,
		StackTrace: stack,
		CreatedAt:  time.Now(),
	}
	config.DB.Create(&log)
}
