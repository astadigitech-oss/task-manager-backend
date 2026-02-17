package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Name           string         `json:"name"`
	Email          string         `gorm:"unique;not null" json:"email"`
	Password       string         `json:"-"`
	Role           string         `json:"role"`
	ProfileImage   *string        `json:"profile_image"`
	Position       *string        `json:"position"` // Jabatan, nullable
	IsOnline       bool           `gorm:"default:false" json:"is_online"`
	LastSeen       *time.Time     `json:"last_seen,omitempty"`
	Workspaces     []Workspace    `gorm:"many2many:workspace_users" json:"workspaces"`
	TelegramChatID *string        `json:"telegram_chat_id,omitempty"`
	Projects       []Project      `gorm:"many2many:project_users" json:"projects"`
	Tasks          []TaskUser     `gorm:"foreignKey:UserID" json:"tasks"`
	CreatedAt      time.Time      `gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // Soft delete
}
