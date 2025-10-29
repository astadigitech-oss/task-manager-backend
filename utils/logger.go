package utils

import (
	"project-management-backend/config"
	"project-management-backend/models"
	"time"
)

// Activity Log
func Activity(userID uint, action, table string, itemID uint, data string) {
	log := models.ActivityLog{
		UserID:    userID,
		Action:    action,
		TableName: table,
		ItemID:    itemID,
		Data:      data,
		CreatedAt: time.Now(),
	}
	config.DB.Create(&log)
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
