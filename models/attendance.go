package models

import (
	"time"

	"gorm.io/gorm"
)

// Attendance represents the attendance model
type Attendance struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null" json:"user_id"`
	WorkspaceID uint           `gorm:"not null" json:"workspace_id"`
	Activity    string         `gorm:"type:text;not null" json:"activity"`
	Obstacle    *string        `gorm:"type:text" json:"obstacle"`
	ClockIn     time.Time      `gorm:"not null" json:"clock_in"`
	ClockOut    *time.Time     `json:"clock_out"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	User      User              `gorm:"foreignKey:UserID"`
	Workspace Workspace         `gorm:"foreignKey:WorkspaceID"`
	Images    []AttendanceImage `gorm:"foreignKey:AttendanceID;constraint:OnDelete:CASCADE" json:"images"`
}
