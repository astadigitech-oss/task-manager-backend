package models

import (
	"time"
)

type TaskFile struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TaskID     uint      `gorm:"not null" json:"task_id"`
	FileName   string    `gorm:"column:filename;not null" json:"filename"`   // ← ADD column:filename
	FileData   []byte    `gorm:"type:longblob;not null" json:"-"`            // ← Change type to longblob
	MimeType   string    `gorm:"column:mime_type;not null" json:"mime_type"` // ← ADD column:mime_type
	FileSize   int64     `gorm:"column:file_size;not null" json:"file_size"` // ← ADD column:file_size
	UploadedBy uint      `gorm:"not null" json:"uploaded_by"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	User       User      `gorm:"foreignKey:UploadedBy" json:"user"`
}

func (TaskFile) TableName() string {
	return "task_files"
}
