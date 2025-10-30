package models

import "time"

type ActivityLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `json:"user_id"`
	Action     string    `json:"action"`
	TableName  string    `json:"table_name"`
	ItemID     uint      `json:"item_id"`
	Data       string    `json:"data"` // opsional - info perubahan/bio/hasil, dll
	User       User      `gorm:"foreignKey:UserID" json:"user"`
	DataBefore string    `json:"data_before"` // simpan JSON string data sebelum aksi
	DataAfter  string    `json:"data_after"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}
