package models

import "time"

type DailyActivityItem struct {
	ActivityTime time.Time
	User         string
	ProjectTitle string
	TaskTitle    string
	TaskPriority string
	StatusAtLog  string
	Overdue      time.Duration
}
