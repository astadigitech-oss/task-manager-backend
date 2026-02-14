package models

import (
	"time"

	"gorm.io/gorm"
)

type TaskStatusLog struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TaskID    uint           `json:"task_id"`
	Status    string         `json:"status"`
	ClockIn   time.Time      `json:"clock_in"`
	ClockOut  *time.Time     `json:"clock_out"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Task Task `gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE" json:"task"`
}
