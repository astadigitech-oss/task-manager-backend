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
	CreatedAt  time.Time `json:"created_at"`
	User       User      `gorm:"foreignKey:UserID" json:"user"`
}
