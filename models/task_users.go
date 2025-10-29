package models

import "time"

type TaskUser struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TaskID     uint      `json:"task_id"`
	UserID     uint      `json:"user_id"`
	RoleInTask string    `json:"role_in_task"`
	AssignedAt time.Time `json:"assigned_at"`
	User       User      `gorm:"foreignKey:UserID" json:"user"`
	Task       Task      `gorm:"foreignKey:TaskID" json:"task"`
}

func (TaskUser) TableName() string { return "task_users" }
