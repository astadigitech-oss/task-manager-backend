package models

import "time"

type ProjectUser struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ProjectID     uint      `json:"project_id"`
	UserID        uint      `json:"user_id"`
	RoleInProject string    `json:"role_in_project"`
	User          User      `gorm:"foreignKey:UserID" json:"user"`
	Project       Project   `gorm:"foreignKey:ProjectID" json:"project"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

func (ProjectUser) TableName() string { return "project_users" }
