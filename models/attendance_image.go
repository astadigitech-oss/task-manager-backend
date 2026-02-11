package models

import (
	"time"

	"gorm.io/gorm"
)

type AttendanceImage struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	AttendanceID uint           `gorm:"not null" json:"attendance_id"`
	URL          string         `gorm:"not null" json:"url"`
	CreatedAt    time.Time      `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"default:current_timestamp" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Attendance Attendance `gorm:"foreignKey:AttendanceID"`
}
