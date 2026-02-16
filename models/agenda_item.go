package models

import "time"

type AgendaItem struct {
	ProjectTitle string     `json:"project_title"`
	TaskTitle    string     `json:"task_title"`
	MemberName   string     `json:"member_name"`
	Status       string     `json:"status"`
	Kondisi      string     `json:"kondisi"`
	StartDate    time.Time  `json:"start_date"`
	DueDate      time.Time  `json:"due_date"`
	Notes        string     `json:"notes"`
	WorkDuration string     `json:"work_duration"`
	FinishedAt   *time.Time `json:"finished_at"`
}
