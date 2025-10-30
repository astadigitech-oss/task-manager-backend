package models

import "time"

type ErrorLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `json:"user_id"`
	Action     string    `json:"action"`
	TableName  string    `json:"table_name"`
	ItemID     uint      `json:"item_id"`
	ErrorMsg   string    `json:"error_msg"`
	StackTrace string    `json:"stack_trace"`
	User       User      `gorm:"foreignKey:UserID" json:"user"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}
